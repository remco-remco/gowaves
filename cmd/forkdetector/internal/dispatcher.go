package internal

import (
	"bufio"
	"github.com/pkg/errors"
	"github.com/wavesplatform/gowaves/pkg/proto"
	"go.uber.org/zap"
	"io"
	"net"
	"os"
	"strings"
	"sync"
	"syscall"
	"time"
)

const (
	reconnectionInterval    = time.Second
	dialingTimeout          = 30 * time.Second
	writeTimeout            = 30 * time.Second
	readTimeout             = 30 * time.Second
	defaultApplication      = "waves"
	closedConnectionMessage = "use of closed network connection"
)

type dispatcher struct {
	interrupt            <-chan struct{}
	registry             *PublicAddressRegistry
	connections          <-chan net.Conn
	addresses            chan []net.TCPAddr
	declaredAddressBytes []byte
	name                 string
	nonce                uint64
	scheme               byte
	mu                   sync.Mutex
	handlers             map[string]*handler
	versions             []proto.Version
}

func NewDispatcher(interrupt <-chan struct{}, registry *PublicAddressRegistry, connections <-chan net.Conn, announcement, name string, nonce uint64, scheme byte, versions []proto.Version) (*dispatcher, error) {
	if connections == nil {
		return nil, errors.New("invalid connections channel")
	}

	dab := make([]byte, 0)
	if announcement != "" {
		var err error
		da, err := proto.NewPeerAddressFromString(announcement)
		if err != nil {
			return nil, errors.Wrap(err, "invalid declared address")
		}
		dab, err = da.MarshalBinary()
		if err != nil {
			return nil, errors.Wrap(err, "invalid declared address")
		}
	}
	return &dispatcher{
		interrupt:            interrupt,
		registry:             registry,
		connections:          connections,
		declaredAddressBytes: dab,
		name:                 name,
		nonce:                nonce,
		scheme:               scheme,
		versions:             versions,
		mu:                   sync.Mutex{},
		handlers:             make(map[string]*handler),
	}, nil
}

func (d *dispatcher) Start() <-chan struct{} {
	zap.S().Debug("Starting dispatcher...")
	done := make(chan struct{})
	reconnectTicker := time.NewTicker(reconnectionInterval)
	go func() {
		for {
			select {
			case <-d.interrupt:
				zap.S().Debug("Shutting down dispatcher...")
				close(done)
				return
			case <-reconnectTicker.C:
				pas, err := d.registry.FeasibleAddresses()
				if err != nil {
					zap.S().Warnf("Failed to pickup peers to connect: %v", err)
					continue
				}
				for _, pa := range pas {
					go d.dial(pa)
				}
			case conn := <-d.connections:
				zap.S().Debugf("New incoming connection to handle %s -> %s", conn.RemoteAddr().String(), conn.LocalAddr().String())
				go d.handleIncoming(conn)
			case addresses := <-d.getAddressesChanLocked():
				go func() {
					n, err := d.registry.RegisterNewAddresses(addresses)
					if err != nil {
						zap.S().Warnf("Failed to add new addresses: %v", err)
					}
					if n > 0 {
						zap.S().Debugf("%d new public addresses were registered", n)
					}
				}()
			}
		}
	}()
	return done
}

func (d *dispatcher) dial(pa PublicAddress) {
	select {
	case <-d.interrupt:
		return
	default:
		conn, err := net.DialTimeout("tcp", pa.address.String(), dialingTimeout)
		if err != nil {
			err = d.registry.Discard(&pa)
			if err != nil {
				zap.S().Warnf("Failed to discard address '%s': %v", pa, err)
				return
			}
			zap.S().Infof("Public address '%s' was discarded due to failed network connection", pa)
			return
		}
		rqh := d.handshake(pa.version)
		zap.S().Debugf("Trying to handshake with '%s' with version %s", conn.RemoteAddr(), pa.version)
		err = conn.SetWriteDeadline(time.Now().Add(writeTimeout))
		if err != nil {
			zap.S().Warnf("Failed to set write timeout: %v", err)
			err := d.registry.Connected(&pa)
			if err != nil {
				zap.S().Warnf("Failed to update public address's state: %v", err)
			}
			return
		}
		_, err = rqh.WriteTo(conn)
		if IsConnectionClosed(err) {
			zap.S().Warnf("Connection to '%s' was closed during sending handshake: %v", conn.RemoteAddr(), err)
			err := d.registry.Connected(&pa)
			if err != nil {
				zap.S().Warnf("Failed to update public address's state: %v", err)
			}
			return
		}
		select {
		case <-d.interrupt:
			err = conn.Close()
			if err != nil {
				zap.S().Warnf("Failed to close connection with '%s': %v", conn.RemoteAddr(), err)
			}
			return
		default:
		}
		err = conn.SetReadDeadline(time.Now().Add(readTimeout))
		if err != nil {
			zap.S().Warnf("Failed to set read timeout: %v", err)
			err := d.registry.Connected(&pa)
			if err != nil {
				zap.S().Warnf("Failed to update public address's state: %v", err)
			}
			return
		}
		var rph proto.Handshake
		_, err = rph.ReadFrom(conn)
		if err != nil {
			if IsConnectionClosed(err) {
				zap.S().Warnf("Connection to '%s' was closed during receiving handshake: %v", conn.RemoteAddr(), err)
				err := d.registry.Connected(&pa)
				if err != nil {
					zap.S().Warnf("Failed to update public address's state: %v", err)
				}
				return
			}
			zap.S().Warnf("Failed to read handshake from node '%s': %v", conn.RemoteAddr(), err)
			err := d.registry.Hostile(&pa)
			if err != nil {
				zap.S().Warnf("Failed to update public address's state: %v", err)
			}
			return
		}
		if rph.AppName[len(rph.AppName)-1] != d.scheme {
			err = d.registry.Hostile(&pa)
			zap.S().Debugf("Node '%s' has different blockchain scheme: %s", conn.RemoteAddr(), rph.AppName)
			if err != nil {
				zap.S().Warnf("Failed to update public address's state: %v", err)
			}
			return
		}
		err = d.registry.Greeted(&pa, rph.Version)
		if err != nil {
			zap.S().Warnf("Failed to update public address's state: %v", err)
			return
		}
		d.register(NewHandler(d.interrupt, conn, d.addresses))
		zap.S().Infof("Successful handshake with '%s'", conn.RemoteAddr())
	}
}

func (d *dispatcher) handshake(v proto.Version) *proto.Handshake {
	sb := strings.Builder{}
	sb.WriteString(defaultApplication)
	sb.WriteByte(d.scheme)
	return &proto.Handshake{
		AppName:           sb.String(),
		Version:           v,
		NodeName:          d.name,
		NodeNonce:         d.nonce,
		DeclaredAddrBytes: d.declaredAddressBytes,
		Timestamp:         proto.NewTimestampFromTime(time.Now()),
	}
}

func (d *dispatcher) handleIncoming(conn net.Conn) {
	zap.S().Debugf("New incoming connection from '%s'", conn.RemoteAddr().String())
	var in proto.Handshake
	r := bufio.NewReader(conn)
	_, err := in.ReadFrom(r)
	if err != nil {
		zap.S().Warnf("Failed to receive handshake from '%s': %v", conn.RemoteAddr(), err)
		return
	}
	if in.AppName[len(in.AppName)-1] != d.scheme {
		zap.S().Debugf("Incoming connection from the node '%s' with different blockchain scheme: %s", conn.RemoteAddr(), in.AppName)
		err := conn.Close()
		if err != nil {
			zap.S().Warnf("Failed to close connection with '%s'", conn.RemoteAddr())
		}
		return
	}
	out := d.handshake(in.Version)
	_, err = out.WriteTo(conn)
	if err != nil {
		zap.S().Warnf("Failed to send handshake from '%s': %v", conn.RemoteAddr(), err)
		return
	}
	if da, err := in.DeclaredAddress(); err == nil {
		a := net.TCPAddr{IP: da.Addr, Port: int(da.Port)}
		ok, err := d.registry.RegisterNewAddress(a, in.Version)
		if err != nil {
			zap.S().Warnf("Failed to register received declared address '%s': %v", da.String(), err)
		}
		if ok {
			zap.S().Infof("New public address '%s' was registered", da.String())
		}
	}
	h := NewHandler(d.interrupt, conn, d.addresses)
	d.register(h)
}

func (d *dispatcher) register(h *handler) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, ok := d.handlers[h.conn.RemoteAddr().String()]; !ok {
		d.handlers[h.conn.RemoteAddr().String()] = h
	}
}

func (d *dispatcher) getAddressesChan() chan<- []net.TCPAddr {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.getAddressesChanLocked()
}

func (d *dispatcher) getAddressesChanLocked() chan []net.TCPAddr {
	if d.addresses == nil {
		d.addresses = make(chan []net.TCPAddr)
	}
	return d.addresses
}

func IsConnectionClosed(err error) bool {
	if err == nil {
		return false
	}
	if err == io.EOF {
		return true
	}
	if opErr, ok := err.(*net.OpError); ok {
		if sysErr, ok := opErr.Err.(*os.SyscallError); ok {
			switch sysErr.Err {
			case syscall.ECONNRESET:
				return true
			case syscall.ECONNABORTED:
				return true
			case syscall.ECONNREFUSED:
				return true
			default:
			}
		}
		if strings.Contains(opErr.Err.Error(), closedConnectionMessage) {
			return true
		}
	}
	if strings.Contains(err.Error(), closedConnectionMessage) {
		return true
	}
	return false
}

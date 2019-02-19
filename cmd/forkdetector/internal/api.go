package internal

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// Logger is a middleware that logs the start and end of each request, along
// with some useful data about what was requested, what the response status was,
// and how long it took to return.
func Logger(l *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
			defer func() {
				l.Info("Served",
					zap.String("proto", r.Proto),
					zap.String("path", r.URL.Path),
					zap.Duration("lat", time.Since(t1)),
					zap.Int("status", ww.Status()),
					zap.Int("size", ww.BytesWritten()),
					zap.String("reqId", middleware.GetReqID(r.Context())))
			}()

			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}

type APIConfig struct {
	On   bool
	Bind string
}

type status struct {
	ShortForksCount     int `json:"short_forks_count"`
	LongForksCount      int `json:"long_forks_count"`
	KnowNodesCount      int `json:"know_nodes_count"`
	ConnectedNodesCount int `json:"connected_nodes_count"`
}

type api struct {
	interrupt <-chan struct{}
	log       *zap.SugaredLogger
}

func StartForkDetectorAPI(interrupt <-chan struct{}, logger *zap.Logger, cfg APIConfig) <-chan struct{} {
	done := make(chan struct{})
	if !cfg.On {
		close(done)
		return done
	}
	a := api{interrupt: interrupt, log: logger.Sugar()}
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(Logger(logger))
	r.Use(middleware.Recoverer)
	r.Use(middleware.SetHeader("Content-Type", "application/json"))
	r.Use(middleware.DefaultCompress)
	r.Mount("/api", a.routes())
	apiServer := &http.Server{Addr: cfg.Bind, Handler: r}
	go func() {
		err := apiServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			a.log.Fatalf("Failed to start API: %v", err)
			return
		}
	}()
	go func() {
		for {
			select {
			case <-a.interrupt:
				a.log.Info("Shutting down API...")
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				err := apiServer.Shutdown(ctx)
				if err != nil {
					a.log.Errorf("Failed to shutdown API server: %v", err)
				}
				cancel()
				close(done)
				return
			}
		}
	}()
	return done
}

func (a *api) routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/status", a.status)
	r.Get("/peers", a.peers)
	r.Get("/forks", a.forks)
	r.Get("/node/{address}", a.node)
	r.Get("/height/{height:\\d+}", a.height)
	r.Get("/block/{id:[a-km-zA-HJ-NP-Z1-9]+}", a.block)
	return r
}

func (a *api) status(w http.ResponseWriter, r *http.Request) {
	//h, err := a.Storage.Height()
	//if err != nil {
	//	http.Error(w, fmt.Sprintf("Failed to complete request: %s", err.Error()), http.StatusInternalServerError)
	//	return
	//}
	//blockID, err := a.Storage.BlockID(h)
	//if err != nil {
	//	http.Error(w, fmt.Sprintf("Failed to complete request: %s", err.Error()), http.StatusInternalServerError)
	//	return
	//}
	//s := status{CurrentHeight: h, LastBlockID: blockID}
	//err = json.NewEncoder(w).Encode(s)
	//if err != nil {
	//	http.Error(w, fmt.Sprintf("Failed to marshal status to JSON: %s", err.Error()), http.StatusInternalServerError)
	//	return
	//}
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (a *api) peers(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (a *api) forks(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (a *api) node(w http.ResponseWriter, r *http.Request) {
	//addr := chi.URLParam(r, "address")
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (a *api) height(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (a *api) block(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}
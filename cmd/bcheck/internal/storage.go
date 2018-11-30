package internal

import (
	"encoding/binary"
	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

var (
	heightKey           = []byte("height")
	defaultReadOptions  = &opt.ReadOptions{}
	defaultWriteOptions = &opt.WriteOptions{}
	iv                  = [...]byte{0, 0, 0, 0}
)

type Storage struct {
	Path string
	db   *leveldb.DB
}

func (s *Storage) Open() error {
	o := &opt.Options{}
	db, err := leveldb.OpenFile(s.Path, o)
	if err != nil {
		return errors.Wrap(err, "failed to open Storage")
	}
	s.db = db
	return nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) HasID(id []byte) (bool, error) {
	return s.db.Has(id, defaultReadOptions)
}

func (s *Storage) PutID(id []byte, h int) error {
	v := iv[:]
	binary.BigEndian.PutUint32(v, uint32(h))
	return s.db.Put(id, v, defaultWriteOptions)
}

func (s *Storage) GetID(id []byte) (int, error) {
	b, err := s.db.Get(id, defaultReadOptions)
	if err != nil {
		return 0, err
	}
	h := binary.BigEndian.Uint32(b)
	return int(h), nil
}

func (s *Storage) PutHeight(h int) error {
	v := iv[:]
	binary.BigEndian.PutUint32(v, uint32(h))
	return s.db.Put(heightKey, v, defaultWriteOptions)
}

func (s *Storage) GetHeight() (int, error) {
	b, err := s.db.Get(heightKey, defaultReadOptions)
	if err != nil {
		return 0, err
	}
	h := binary.BigEndian.Uint32(b)
	return int(h), nil
}

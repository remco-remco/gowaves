package db

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/wavesplatform/gowaves/pkg/crypto"
	"github.com/wavesplatform/gowaves/pkg/proto"
)

const (
	heightSuffix         = "h"
	blocksAtHeightPrefix = "blocks_at_height_"
)

type blocksAtHeight struct {
	Signatures []crypto.Signature
}

var ErrBlockOrphaned = errors.New("block orphaned")

type WavesDB struct {
	genesis crypto.Signature
	db      *leveldb.DB
}

func (w *WavesDB) GetRaw(key []byte) ([]byte, error) {
	return w.db.Get(key, nil)
}

func (w *WavesDB) PutRaw(key, value []byte) error {
	return w.db.Put(key, value, nil)
}

func (w *WavesDB) DelRaw(key []byte) error {
	return w.db.Delete(key, nil)
}

func (w *WavesDB) Has(block crypto.Signature) (bool, error) {
	return w.db.Has(block[:], nil)
}

func blocksAtHeightKey(height uint64) []byte {
	key := blocksAtHeightPrefix + strconv.FormatUint(height, 10)
	return []byte(key)
}

func (w *WavesDB) putBlockAtHeight(block *proto.Block) error {
	blocksAtHeightKey := blocksAtHeightKey(block.Height)

	var blocks blocksAtHeight
	bytes, err := w.db.Get([]byte(blocksAtHeightKey), nil)
	if err != nil {
		b := blocksAtHeight{Signatures: []crypto.Signature{block.BlockSignature}}
		bytes, err = json.Marshal(b)
		if err != nil {
			return err
		}

		return w.db.Put(blocksAtHeightKey, bytes, nil)
	}

	err = json.Unmarshal(bytes, &blocks)
	if err != nil {
		return err
	}

	blocks.Signatures = append(blocks.Signatures, block.BlockSignature)
	bytes, err = json.Marshal(blocks)

	return w.db.Put(blocksAtHeightKey, bytes, nil)
}

func (w *WavesDB) GetBlocksAtHeight(height uint64) ([]crypto.Signature, error) {
	key := blocksAtHeightKey(height)

	bytes, err := w.db.Get(key, nil)
	if err != nil {
		return nil, err
	}

	var blocks blocksAtHeight
	if err = json.Unmarshal(bytes, &blocks); err != nil {
		return nil, err
	}

	return blocks.Signatures, nil
}

func (w *WavesDB) Put(block *proto.Block) error {
	var height uint64

	switch {
	case block.BlockSignature == w.genesis:
		height = 1
	default:
		parentHeight := append(block.Parent[:], []byte(heightSuffix)...)
		has, err := w.db.Has(parentHeight, nil)
		if err != nil {
			return err
		}
		if !has {
			return ErrBlockOrphaned
		}
		heightBytes, err := w.db.Get(parentHeight, nil)
		if err != nil {
			return err
		}
		height = binary.BigEndian.Uint64(heightBytes)
		height++
	}
	bytes, err := block.MarshalBinary()
	if err != nil {
		return err
	}
	if err = w.db.Put(block.BlockSignature[:], bytes, nil); err != nil {
		return err
	}

	heightKey := append(block.BlockSignature[:], []byte(heightSuffix)...)
	heightValue := make([]byte, 8)
	binary.BigEndian.PutUint64(heightValue, height)
	if err = w.db.Put(heightKey, heightValue, nil); err != nil {
		return err
	}
	block.Height = height
	return w.putBlockAtHeight(block)
}

func (w *WavesDB) Get(block crypto.Signature) (*proto.Block, error) {
	bytes, err := w.db.Get(block[:], nil)
	if err != nil {
		return nil, err
	}
	var res proto.Block
	if err = res.UnmarshalBinary(bytes); err != nil {
		return nil, err
	}

	heightKey := append(block[:], []byte(heightSuffix)...)
	heightBytes, err := w.db.Get(heightKey, nil)
	if err != nil {
		return nil, err
	}
	height := binary.BigEndian.Uint64(heightBytes)
	res.Height = height

	return &res, nil
}

func NewDB(path string, genesis crypto.Signature) (*WavesDB, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}

	return &WavesDB{db: db, genesis: genesis}, nil
}

package util

import (
	"io/ioutil"
	"os"
)

type DataDir string

func NewTemporary() (DataDir, error) {
	dir, err := ioutil.TempDir("", "datatemp")
	return DataDir(dir), err
}

func (a DataDir) Clear() error {
	return os.RemoveAll(string(a))
}

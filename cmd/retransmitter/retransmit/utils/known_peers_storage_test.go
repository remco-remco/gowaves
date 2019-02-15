package utils

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	fs := afero.NewOsFs()
	bts := []byte("hello")

	s, err := NewFileBasedStorage(fs, "./known_peers.json")
	require.NoError(t, err)

	err = s.Save(bts)
	require.NoError(t, err)

	ret, err := s.Read()
	require.NoError(t, err)
	assert.Equal(t, bts, ret)

	s.Close()
	ret, err = s.Read()
	assert.Contains(t, err.Error(), "closed")
}

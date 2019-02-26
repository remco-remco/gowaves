package util

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDataDir(t *testing.T) {
	d, err := NewTemporary()
	require.NoError(t, err)

	err = d.Clear()
	require.NoError(t, err)
}

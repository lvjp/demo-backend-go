package password

import (
	"crypto/rand"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHash(t *testing.T) {
	password, err := io.ReadAll(io.LimitReader(rand.Reader, 128))
	require.NoError(t, err)

	storedPassword, err := Hash(password)
	require.NoError(t, err)

	require.NotEqual(t, []byte(storedPassword), password)

	isSame, err := IsSame(password, storedPassword)
	require.NoError(t, err)
	require.True(t, isSame)
}

package password

import (
	"bytes"
	"crypto/rand"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHash(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		password := make([]byte, 128)
		_, err := io.ReadFull(rand.Reader, password)
		require.NoError(t, err)

		storedPassword, err := Hash(password, rand.Reader)
		require.NoError(t, err)

		require.NotEqual(t, []byte(storedPassword), password)

		isSame, err := IsSame(password, storedPassword)
		require.NoError(t, err)
		require.True(t, isSame)
	})

	t.Run("bad encoded hash", func(t *testing.T) {
		isSame, err := IsSame([]byte("password"), "------")
		require.Error(t, err, "salt decoding error")
		require.False(t, isSame)
	})

	t.Run("bad argon2 version", func(t *testing.T) {
		isSame, err := IsSame([]byte("password"), "$argon2id$v=0$m=65536,t=2,p=1$gZiV/M1gPc22ElAH/Jh1Hw$CWOrkoo7oJBQ/iyh7uJ0LO2aLEfrHwTWllSAxT0zRno")
		require.Error(t, err, "encoded hash decoding error")
		require.False(t, isSame)
	})

	t.Run("entropy source error", func(t *testing.T) {
		hash, err := Hash([]byte("password"), bytes.NewReader(nil))
		require.Error(t, err, "salt generation error")
		require.Empty(t, hash)
	})
}

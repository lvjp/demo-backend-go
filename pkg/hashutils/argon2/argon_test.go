package argon2

import (
	"crypto/rand"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
	inputSalt, err := io.ReadAll(io.LimitReader(rand.Reader, 16))
	require.NoError(t, err)

	inputPassword, err := io.ReadAll(io.LimitReader(rand.Reader, 64))
	require.NoError(t, err)

	inputParams := ArgonParams{
		Memory:      64 * 1024,
		Time:        2,
		Parallelism: 2,
		KeyLength:   32,
	}

	phc := Hash(inputPassword, inputSalt, inputParams)
	storedHash := phc.String()

	t.Log("Stored password:", storedHash)

	actualPHC, actualParams, err := DecodePHCString(storedHash)
	require.NoError(t, err)

	require.Equal(
		t,
		phc,
		actualPHC,
	)

	require.Equal(
		t,
		&inputParams,
		actualParams,
	)

	rehashed := Hash(inputPassword, actualPHC.Salt, *actualParams).String()
	require.Equal(t, storedHash, rehashed)

	badHash := Hash(append(inputPassword, []byte("pa$$w0rd")...), actualPHC.Salt, *actualParams).String()
	require.NotEqual(t, storedHash, badHash)
}

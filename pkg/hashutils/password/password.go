package password

import (
	"crypto/rand"
	"crypto/subtle"
	"fmt"
	"io"

	"go.lvjp.me/demo-backend-go/pkg/hashutils/argon2"
)

// saltSize of 16Bytes/128bits of randomness
const saltSize = 32

func Hash(password []byte) (string, error) {
	salt, err := io.ReadAll(io.LimitReader(rand.Reader, saltSize))
	if err != nil {
		return "", fmt.Errorf("salt generation error: %w", err)
	}

	phc := argon2.Hash(password, salt, argon2.DefaultArgonParams)

	return phc.String(), nil
}

func IsSame(password []byte, encodedHash string) (bool, error) {
	phc, params, err := argon2.DecodePHCString(encodedHash)
	if err != nil {
		return false, fmt.Errorf("encoded hash decoding error: %w", err)
	}

	hashed := argon2.Hash(password, phc.Salt, *params)

	return subtle.ConstantTimeCompare(phc.Hash, hashed.Hash) == 1, nil
}

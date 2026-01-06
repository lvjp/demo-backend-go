package password

import (
	"crypto/subtle"
	"fmt"
	"io"

	"go.lvjp.me/demo-backend-go/pkg/hashutils/argon2"
	"go.lvjp.me/demo-backend-go/pkg/hashutils/phc"
)

// saltSize of 16Bytes/128bits of randomness
const saltSize = 32

func Hash(password []byte, randomSource io.Reader) (string, error) {
	salt := make([]byte, saltSize)
	if _, err := io.ReadFull(randomSource, salt); err != nil {
		return "", fmt.Errorf("salt generation error: %v", err)
	}

	phc := argon2.Hash(password, salt, argon2.DefaultParams)

	return phc.String(), nil
}

func IsSame(password []byte, hash string) (bool, error) {
	phc, err := phc.NewString(hash)
	if err != nil {
		return false, fmt.Errorf("PHC decode error: %v", err)
	}

	params, err := argon2.NewParams(phc)
	if err != nil {
		return false, fmt.Errorf("encoded hash decoding error: %v", err)
	}

	hashed := argon2.Hash(password, phc.Salt, *params)

	return subtle.ConstantTimeCompare(phc.Hash, hashed.Hash) == 1, nil
}

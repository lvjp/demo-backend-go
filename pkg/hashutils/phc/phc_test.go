package phc

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewString(t *testing.T) {
	salt, err := base64.RawStdEncoding.DecodeString("gZiV/M1gPc22ElAH/Jh1Hw")
	require.NoError(t, err)

	hash, err := base64.RawStdEncoding.DecodeString("CWOrkoo7oJBQ/iyh7uJ0LO2aLEfrHwTWllSAxT0zRno")
	require.NoError(t, err)

	testCases := []struct {
		name string
		raw  string
		phc  *String
	}{
		{
			name: "full",
			raw:  "$argon2id$v=19$m=65536,t=2,p=1$gZiV/M1gPc22ElAH/Jh1Hw$CWOrkoo7oJBQ/iyh7uJ0LO2aLEfrHwTWllSAxT0zRno",
			phc: &String{
				ID:      "argon2id",
				Version: "19",
				Params: []Parameter{
					{Name: "m", Value: "65536"},
					{Name: "t", Value: "2"},
					{Name: "p", Value: "1"},
				},
				Salt: salt,
				Hash: hash,
			},
		},
		{
			name: "id",
			raw:  "$argon2id",
			phc: &String{
				ID: "argon2id",
			},
		},
		{
			name: "version",
			raw:  "$argon2id$v=19",
			phc: &String{
				ID:      "argon2id",
				Version: "19",
			},
		},
		{
			name: "salt",
			raw:  "$argon2id$gZiV/M1gPc22ElAH/Jh1Hw",
			phc: &String{
				ID:   "argon2id",
				Salt: salt,
			},
		},
		{
			name: "hash",
			raw:  "$argon2id$gZiV/M1gPc22ElAH/Jh1Hw$CWOrkoo7oJBQ/iyh7uJ0LO2aLEfrHwTWllSAxT0zRno",
			phc: &String{
				ID:   "argon2id",
				Salt: salt,
				Hash: hash,
			},
		},
		{
			name: "params",
			raw:  "$argon2id$m=65536,t=2,p=1",
			phc: &String{
				ID: "argon2id",
				Params: []Parameter{
					{Name: "m", Value: "65536"},
					{Name: "t", Value: "2"},
					{Name: "p", Value: "1"},
				},
			},
		},
		{
			name: "param named v first",
			raw:  "$argon2id$v=19,m=65536,t=2,p=1",
			phc: &String{
				ID: "argon2id",
				Params: []Parameter{
					{Name: "v", Value: "19"},
					{Name: "m", Value: "65536"},
					{Name: "t", Value: "2"},
					{Name: "p", Value: "1"},
				},
			},
		},
		{
			name: "param named v only not version",
			raw:  "$argon2id$v=aparam",
			phc: &String{
				ID: "argon2id",
				Params: []Parameter{
					{Name: "v", Value: "aparam"},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := NewString(tc.raw)
			require.NoError(t, err)

			require.Equal(t, tc.phc, actual)
			require.Equal(t, tc.raw, actual.String())
		})
	}
}

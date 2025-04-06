package argon2

import (
	"encoding/base64"
	"slices"
	"strconv"
	"testing"

	"go.lvjp.me/demo-backend-go/pkg/hashutils/phc"

	"github.com/stretchr/testify/require"
)

func TestHash(t *testing.T) {
	testCases := []struct {
		params   Params
		password string
		salt     string
		hashed   string
	}{
		{Params{Time: 1, Memory: 1 << 16, Parallelism: 1}, "password", "somesalt", "$argon2id$v=19$m=65536,t=1,p=1$c29tZXNhbHQ$9qWtwbpyPd3vm1rB1GThgPzZ3/ydHL92zKL+15XZypg"},
		{Params{Time: 2, Memory: 1 << 16, Parallelism: 1}, "differentpassword", "somesalt", "$argon2id$v=19$m=65536,t=2,p=1$c29tZXNhbHQ$C4TWUs9rDEvq7w3+J4umqA32aWKB1+DSiRuBfYxFj94"},
		{Params{Time: 2, Memory: 1 << 16, Parallelism: 1}, "password", "diffsalt", "$argon2id$v=19$m=65536,t=2,p=1$ZGlmZnNhbHQ$vfMrBczELrFdWP0ZsfhWsRPaHppYdP3MVEMIVlqoFBw"},
		{Params{Time: 2, Memory: 1 << 16, Parallelism: 1}, "password", "somesalt", "$argon2id$v=19$m=65536,t=2,p=1$c29tZXNhbHQ$CTFhFdXPJO1aFaMaO6Mm5c8y7cJHAph8ArZWb2GRPPc"},
		{Params{Time: 2, Memory: 1 << 18, Parallelism: 1}, "password", "somesalt", "$argon2id$v=19$m=262144,t=2,p=1$c29tZXNhbHQ$eP4eyR+zqlZX1y5xCFTkw9m5GYx0L5YWwvCFvtlbLow"},
		{Params{Time: 2, Memory: 1 << 8, Parallelism: 1}, "password", "somesalt", "$argon2id$v=19$m=256,t=2,p=1$c29tZXNhbHQ$nf65EOgLrQMR/uIPnA4rEsF5h7TKyQwu9U1bMCHGi/4"},
		{Params{Time: 2, Memory: 1 << 8, Parallelism: 2}, "password", "somesalt", "$argon2id$v=19$m=256,t=2,p=2$c29tZXNhbHQ$bQk8UB/VmZZF4Oo79iDXuL5/0ttZwg2f/5U52iv1cDc"},
		{Params{Time: 4, Memory: 1 << 16, Parallelism: 1}, "password", "somesalt", "$argon2id$v=19$m=65536,t=4,p=1$c29tZXNhbHQ$kCXUjmjvc5XMqQedpMTsOv+zyJEf5PhtGiUghW9jFyw"},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			tc.params.KeyLength = 32

			actual := Hash([]byte(tc.password), []byte(tc.salt), tc.params)
			require.Equal(t, tc.hashed, actual.String())
		})
	}
}

func TestNewParams(t *testing.T) {
	validDecoded := phc.String{
		ID:      "argon2id",
		Version: "19",
		Params: []phc.Parameter{
			{Name: "m", Value: "65536"},
			{Name: "t", Value: "1"},
			{Name: "p", Value: "1"},
		},
		Salt: []byte("somesalt"),
		Hash: func() []byte {
			validHash, err := base64.RawStdEncoding.DecodeString("9qWtwbpyPd3vm1rB1GThgPzZ3/ydHL92zKL+15XZypg")
			require.NoError(t, err)
			return validHash
		}(),
	}

	testCases := []struct {
		err    string
		phc    phc.String
		params Params
	}{
		{
			phc: validDecoded,
			params: Params{
				Memory:      65536,
				Time:        1,
				Parallelism: 1,
				KeyLength:   32,
			},
		},
		{
			err: "unsupported hashing function",
			phc: func() phc.String {
				validCopy := validDecoded
				validCopy.ID = "dummy"
				return validCopy
			}(),
		},
		{
			err: "unsupported argon2id version",
			phc: func() phc.String {
				validCopy := validDecoded
				validCopy.Version = "----"
				return validCopy
			}(),
		},
		{
			err: "invalid parameter count: 0",
			phc: func() phc.String {
				validCopy := validDecoded
				validCopy.Params = nil
				return validCopy
			}(),
		},
		{
			err: "invalid parameter count: 2",
			phc: func() phc.String {
				validCopy := validDecoded
				validCopy.Params = validCopy.Params[1:]
				return validCopy
			}(),
		},
		{
			err: "parameters should be in the order: m, t, p",
			phc: func() phc.String {
				validCopy := validDecoded
				validCopy.Params = slices.Clone(validCopy.Params)
				validCopy.Params[0], validCopy.Params[1] = validCopy.Params[1], validCopy.Params[0]
				return validCopy
			}(),
		},
		{
			err: "memory parameter decode error",
			phc: func() phc.String {
				validCopy := validDecoded
				validCopy.Params = slices.Clone(validCopy.Params)
				validCopy.Params[0].Value = "invalid"
				return validCopy
			}(),
		},
		{
			err: "time parameter decode error",
			phc: func() phc.String {
				validCopy := validDecoded
				validCopy.Params = slices.Clone(validCopy.Params)
				validCopy.Params[1].Value = "invalid"
				return validCopy
			}(),
		},
		{
			err: "parallelims parameter decode error",
			phc: func() phc.String {
				validCopy := validDecoded
				validCopy.Params = slices.Clone(validCopy.Params)
				validCopy.Params[2].Value = "invalid"
				return validCopy
			}(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.err, func(t *testing.T) {
			actual, err := NewParams(&tc.phc)
			if tc.err != "" {
				require.ErrorContains(t, err, tc.err)
				require.Nil(t, actual)
				return
			}

			require.NoError(t, err)
		})
	}
}

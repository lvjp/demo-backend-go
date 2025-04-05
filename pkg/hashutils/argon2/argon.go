package argon2

import (
	"crypto/subtle"
	"errors"
	"fmt"
	"math"
	"strconv"

	"go.lvjp.me/demo-backend-go/pkg/hashutils/phc"

	"golang.org/x/crypto/argon2"
)

const ID = "argon2id"

const owaspRecommendedMemory = 46 * 1024
const owaspRecommendedTime = 1
const owaspRecommendedParallelism = 1

// keyLength of 32B*8 = 256bits is long enough for our use case
const keyLength = 32

var DefaultArgonParams = ArgonParams{
	Memory:      owaspRecommendedMemory,
	Time:        owaspRecommendedTime,
	Parallelism: owaspRecommendedParallelism,
	KeyLength:   keyLength,
}

type ArgonParams struct {
	Memory      uint32
	Time        uint32
	Parallelism uint8
	KeyLength   uint32
}

func Hash(password, salt []byte, params ArgonParams) *phc.String {
	copyedSalt := make([]byte, len(salt))
	subtle.ConstantTimeCopy(1, copyedSalt, salt)

	hash := argon2.IDKey(
		password,
		copyedSalt,
		params.Time,
		params.Memory,
		params.Parallelism,
		params.KeyLength,
	)

	return &phc.String{
		ID:      ID,
		Version: strconv.Itoa(argon2.Version),
		Params: []phc.Parameter{
			{Name: "m", Value: strconv.FormatUint(uint64(params.Memory), 10)},
			{Name: "t", Value: strconv.FormatUint(uint64(params.Time), 10)},
			{Name: "p", Value: strconv.FormatUint(uint64(params.Parallelism), 10)},
		},
		Hash: hash,
		Salt: copyedSalt,
	}
}

// argon2idParametersCount is the number of parameters that are expected in the PHC string.
// There are 3 parameters: m (memory), t (time), and p (parallelism).
const argon2idParametersCount = 3

func DecodePHCString(rawPhc string) (*phc.String, *ArgonParams, error) {
	phc, err := phc.NewString(rawPhc)
	if err != nil {
		return nil, nil, fmt.Errorf("PHC decode error: %w", err)
	}

	if phc.ID != ID {
		return nil, nil, errors.New("unsupported hashing function: " + phc.ID)
	}

	if phc.Version != strconv.Itoa(argon2.Version) {
		return nil, nil, errors.New("unsupported argon2id version: " + phc.Version)
	}

	if len(phc.Params) != argon2idParametersCount {
		return nil, nil, errors.New("invalid parameter count: " + strconv.Itoa(len(phc.Params)))
	}

	if phc.Params[0].Name != "m" || phc.Params[1].Name != "t" || phc.Params[2].Name != "p" {
		return nil, nil, errors.New("parameters should be in the order: m, t, p")
	}

	memory, err := strconv.ParseUint(phc.Params[0].Value, 10, 32)
	if err != nil {
		return nil, nil, fmt.Errorf("memory parameter decode error: %w", err)
	}

	time, err := strconv.ParseUint(phc.Params[1].Value, 10, 32)
	if err != nil {
		return nil, nil, fmt.Errorf("time parameter decode error: %w", err)
	}

	parallelism, err := strconv.ParseUint(phc.Params[2].Value, 10, 8)
	if err != nil {
		return nil, nil, fmt.Errorf("parallelims parameter decode error: %w", err)
	}

	// uint32 overflow check
	if len(phc.Hash) > math.MaxUint32 {
		return nil, nil, fmt.Errorf("hash is too long: %d", len(phc.Hash))
	}

	params := &ArgonParams{
		KeyLength:   uint32(len(phc.Hash)), // #nosec G115 checked
		Memory:      uint32(memory),
		Time:        uint32(time),
		Parallelism: uint8(parallelism),
	}

	return phc, params, nil
}

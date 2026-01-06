package phc

import (
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// String is a modelisation of the Password Hashing Competition string format.
// Specifications can be found at: https://github.com/P-H-C/phc-string-format
type String struct {
	ID      string
	Version string
	Params  []Parameter
	Salt    []byte
	Hash    []byte
}

type Parameter struct {
	Name  string
	Value string
}

var phcFormat = regexp.MustCompile(
	`^` +
		`\$([a-z0-9-]{1,32})` +
		`(?:\$v=([0-9]+))?` +
		`(?:\$([a-z0-9-]+=[a-zA-Z0-9/+.-]+(?:,[a-z0-9-]+=[a-zA-Z0-9/+.-]+)*))?` +
		`(?:\$([a-zA-Z0-9/+.-]+)(?:\$([a-zA-Z0-9/+.-]+))?)?` +
		`$`,
)

func NewString(raw string) (*String, error) {
	submatches := phcFormat.FindStringSubmatch(raw)
	if submatches == nil {
		return nil, errors.New("not valid PHC string format")
	}

	phc := &String{
		ID:      submatches[1],
		Version: submatches[2],
	}

	if len(submatches[3]) > 0 {
		split := strings.Split(submatches[3], ",")
		phc.Params = make([]Parameter, len(split))
		for i, param := range split {
			kv := strings.Split(param, "=")
			phc.Params[i].Name = kv[0]
			phc.Params[i].Value = kv[1]
		}
	}

	if len(submatches[4]) > 0 {
		var err error
		phc.Salt, err = base64.RawStdEncoding.DecodeString(submatches[4])
		if err != nil {
			return nil, fmt.Errorf("salt decoding error: %v", err)
		}
		if len(submatches[5]) > 0 {
			phc.Hash, err = base64.RawStdEncoding.DecodeString(submatches[5])
			if err != nil {
				return nil, fmt.Errorf("hash decoding error: %v", err)
			}
		}
	}

	return phc, nil
}

func (phc *String) String() string {
	var builder strings.Builder
	builder.WriteRune('$')
	builder.WriteString(phc.ID)

	if phc.Version != "" {
		builder.WriteString("$v=")
		builder.WriteString(phc.Version)
	}

	if len(phc.Params) > 0 {
		builder.WriteRune('$')
		for i, param := range phc.Params {
			if i > 0 {
				builder.WriteRune(',')
			}
			builder.WriteString(param.Name)
			builder.WriteRune('=')
			builder.WriteString(param.Value)
		}
	}

	if len(phc.Salt) > 0 {
		builder.WriteRune('$')
		builder.WriteString(base64.RawStdEncoding.EncodeToString(phc.Salt))

		if len(phc.Hash) > 0 {
			builder.WriteRune('$')
			builder.WriteString(base64.RawStdEncoding.EncodeToString(phc.Hash))
		}
	}

	return builder.String()
}

package headers

import (
	"errors"
	"fmt"

	"github.com/i9si-sistemas/stringx"
)

type Headers map[string]string

func New() Headers {
	return make(Headers)
}

func (h Headers) Get(key string) (value string, ok bool) {
	value, ok = h[key]
	return
}

const rn = "\r\n"

func (h Headers) Parse(data []byte) (read int, done bool, err error) {
	for {
		idx := stringx.New(data[read:]).IndexOf(rn)
		if idx < 0 {
			break
		}
		if idx == 0 {
			done = true
			break
		}
		key, value, err := parseHeader(data[read : read+idx])
		if err != nil {
			return 0, false, err
		}
		read += idx + len(rn)
		h[key] = value
	}
	read += len(rn) // Account for the final \r\n
	return
}

func parseHeader(fieldLine []byte) (key, value string, err error) {
	emptyStr := stringx.Empty.String()
	parts := stringx.New(fieldLine).SplitN(":", 2)
	if len(parts) != 2 {
		return emptyStr, emptyStr, fmt.Errorf("invalid header line: %s", fieldLine)
	}
	space := stringx.Space.String()
	removeSpace := func(s string) stringx.String { return stringx.New(s).Replace(space, emptyStr) }
	key = parts[0]
	if stringx.New(key).HasSuffix(space) {
		return emptyStr, emptyStr, errors.New("invalid header key")
	}
	value = removeSpace(parts[1]).String()
	if key == emptyStr || value == emptyStr {
		return emptyStr, emptyStr, fmt.Errorf("invalid header line: %s", fieldLine)
	}
	return key, value, nil
}

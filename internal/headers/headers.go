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
	value, ok = h[toLowerCase(key)]
	return
}

func (h Headers) Set(key, value string) {
	isEmpty := func(s string) bool { return len(s) == 0 }
	if isEmpty(key) {
		return
	}
	if isEmpty(value) {
		delete(h, toLowerCase(key))
		return
	}
	if old, ok := h[toLowerCase(key)]; ok {
		h.Replace(key, fmt.Sprintf("%s, %s", old, value))
		return
	}

	h.Replace(key, value)
}

func (h Headers) Replace(key, value string) {
	h[toLowerCase(key)] = value
}

func toLowerCase(s string) string {
	return stringx.New(s).ToLowerCase().String()
}

const rn = "\r\n"

func (h Headers) Parse(data []byte) (read int, done bool, err error) {
	for {
		idx := stringx.New(data[read:]).IndexOf(rn)
		if idx < 0 {
			break
		}

		if idx == 0 {
			read += len(rn)
			done = true
			break
		}

		line := data[read : read+idx]
		key, value, err := parseHeader(line)
		if err != nil {
			return 0, false, err
		}
		if !isToken(key) {
			return 0, false, fmt.Errorf("invalid header key: %s", key)
		}

		h.Set(key, value)
		read += idx + len(rn)
	}
	return
}

func parseHeader(fieldLine []byte) (key, value string, err error) {
	emptyStr := stringx.Empty.String()
	parts := stringx.New(string(fieldLine)).SplitN(":", 2)
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

func isToken(s string) (found bool) {
	isAlphaNum := func(ch byte) bool {
		return (ch >= 'A' && ch <= 'Z') ||
			(ch >= 'a' && ch <= 'z') ||
			(ch >= '0' && ch <= '9')
	}
	isSymbol := func(ch byte) bool {
		return ch == '!' || ch == '#' || ch == '$' ||
			ch == '%' || ch == '&' || ch == '\'' ||
			ch == '*' || ch == '+' || ch == '-' ||
			ch == '.' || ch == '^' || ch == '_' ||
			ch == '`' || ch == '|' || ch == '~'
	}

	isValid := func(s string) bool {
		for i := 0; i < len(s); i++ {
			ch := s[i]
			if !isAlphaNum(ch) && !isSymbol(ch) {
				return false
			}
		}
		return true
	}
	return isValid(s)
}

package request

import (
	"errors"
	"io"

	"github.com/gabrielluizsf/tcp_to_http/internal/headers"
	"github.com/gabrielluizsf/tcp_to_http/internal/version"
	"github.com/i9si-sistemas/stringx"
)

var (
	HTTP_VERSION              = version.HTTP
	SEPARATOR                 = "\r\n"
	ErrBodyShorterThanCL      = errors.New("body shorter than Content-Length")
	ErrMalformedRequestLine   = errors.New("malformed request-line")
	ErrMalformedHTTPVersion   = errors.New("malformed http version")
	ErrRequestInErrorState    = errors.New("request in error state")
	ErrUnsupportedHTTPVersion = errors.New("unsupported http version")
)

func NewFromReader(r io.Reader) (req *Request, err error) {
	req = newRequest()
	buf := make([]byte, 1024)
	bufLen := 0
	for !req.done() {
		n, err := r.Read(buf[bufLen:])
		if err != nil {
			return nil, err
		}
		bufLen += n
		readN, err := req.parse(buf[:bufLen])
		if err != nil {
			return nil, err
		}
		copy(buf, buf[readN:bufLen])
		bufLen -= readN
	}
	return req, nil
}

func newRequest() *Request {
	return &Request{
		state:   StateInit,
		Headers: headers.New(),
		Params:  make(Params),
	}
}

type parserState string

const (
	StateInit    parserState = "init"
	StateDone    parserState = "done"
	StateHeaders parserState = "headers"
	StateBody    parserState = "body"
	StateError   parserState = "error"
)

func parseRequestLine(
	b []byte,
) (
	line *RequestLine, index int, err error,
) {
	requestLineStr := stringx.New(string(b))
	separatorIdx := requestLineStr.IndexOf(SEPARATOR)
	if separatorIdx == -1 {
		return nil, 0, nil
	}
	startLine := b[:separatorIdx]
	read := separatorIdx + len(SEPARATOR)
	parts := stringx.New(string(startLine)).Split(stringx.Space.String())
	if len(parts) != 3 {
		return nil, 0, ErrMalformedRequestLine
	}
	httpVersion := parts[2]
	httpVersionParts := stringx.New(httpVersion).Split("/")
	isValidHTTPVersion := len(httpVersionParts) != 2 || httpVersionParts[0] != "HTTP"
	if isValidHTTPVersion {
		return nil, 0, ErrMalformedHTTPVersion
	}
	rl := &RequestLine{
		Method:  parts[0],
		Target:  parts[1],
		Version: httpVersionParts[1],
	}
	if httpVersion != HTTP_VERSION {
		return nil, 0, ErrUnsupportedHTTPVersion
	}
	return rl, read, nil
}

// Request represents the HTTP request
type Request struct {
	Line    RequestLine
	Params  Params
	Headers headers.Headers
	Body    []byte
	state   parserState
}

func (req *Request) parse(data []byte) (int, error) {
	read := 0

parseLoop:
	for {
		switch req.state {
		case StateError:
			return 0, ErrRequestInErrorState

		case StateInit:
			rl, n, err := parseRequestLine(data[read:])
			if err != nil {
				req.state = StateError
				return 0, err
			}
			if n == 0 {
				break parseLoop
			}
			req.Line = *rl
			read += n
			req.state = StateHeaders

		case StateHeaders:
			n, done, err := req.Headers.Parse(data[read:])
			if err != nil {
				req.state = StateError
				return 0, err
			}
			if n == 0 {
				break parseLoop
			}
			read += n
			if done {
				req.state = StateBody
			}

		case StateBody:
			contentLength := req.getIntHeader("Content-Length", 0)
			if contentLength == 0 {
				req.state = StateDone
				break parseLoop
			}

			toRead := min(contentLength-len(req.Body), len(data)-read)
			req.Body = append(req.Body, data[read:read+toRead]...)
			read += toRead

			if len(req.Body) == contentLength {
				req.state = StateDone
			} else if toRead == 0 {
				break parseLoop
			}

		case StateDone:
			break parseLoop
		}
	}

	if req.state == StateBody && len(req.Body) < req.getIntHeader("Content-Length", 0) && len(data) == 0 {
		req.state = StateError
		return read, ErrBodyShorterThanCL
	}

	return read, nil
}

func (req *Request) error() bool {
	return req.state == StateError
}

func (req *Request) done() bool {
	return req.state == StateDone || req.error()
}

func (req *Request) getIntHeader(key string, defaultValue int) int {
	value, ok := req.Headers.Get(key)
	if !ok {
		return defaultValue
	}
	result, err := stringx.NewParser(value).Int()
	if err != nil {
		return defaultValue
	}
	return int(result)
}

type RequestLine struct {
	Method  string
	Target  string
	Version string
}

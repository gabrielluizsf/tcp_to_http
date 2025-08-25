package server

import (
	"regexp"

	"github.com/gabrielluizsf/tcp_to_http/internal/request"
	"github.com/gabrielluizsf/tcp_to_http/internal/response"
	"github.com/i9si-sistemas/nine"
	"github.com/i9si-sistemas/stringx"
)

type Method string

const (
	GET    Method = "GET"
	POST   Method = "POST"
	PUT    Method = "PUT"
	PATCH  Method = "PATCH"
	DELETE Method = "DELETE"
)

func MethodMiddleware(method Method, endpoint string, handler Handler) Handler {
	return func(w *response.Writer, req *request.Request) {
		isEqual := stringx.New(req.Line.Method).Equal(string(method))
		if !isEqual {
			statusCode := response.StatusMethodNotAllowed
			statusText := statusCode.String()
			msg := nine.JSON{
				"error": statusText,
			}
			msgBytes, _  := msg.Bytes()
			headers := response.GetDefaultHeaders(len(msgBytes))
			headers.Replace("Content-Type", "application/json")
			w.WriteStatusLine(response.StatusMethodNotAllowed)
			w.WriteHeaders(headers)
			w.WriteBody(msgBytes)
			return
		}
		req.Params = req.Params.Reset()
		req.Params.Set(req.Line.Target, endpoint)
		handler(w, req)
	}
}

func (s *Server) Get(endpoint string, handler Handler) {
	s.handlers = append(s.handlers, route{
		endpoint: endpoint,
		matcher:  buildMatcher(endpoint),
		handler:  MethodMiddleware(GET, endpoint, handler),
	})
}

func (s *Server) Post(endpoint string, handler Handler) {
	s.handlers = append(s.handlers, route{
		endpoint: endpoint,
		matcher:  buildMatcher(endpoint),
		handler:  MethodMiddleware(POST, endpoint, handler),
	})
}

func (s *Server) Put(endpoint string, handler Handler) {
	s.handlers = append(s.handlers, route{
		endpoint: endpoint,
		matcher:  buildMatcher(endpoint),
		handler:  MethodMiddleware(PUT, endpoint, handler),
	})
}

func (s *Server) Patch(endpoint string, handler Handler) {
	s.handlers = append(s.handlers, route{
		endpoint: endpoint,
		matcher:  buildMatcher(endpoint),
		handler:  MethodMiddleware(PATCH, endpoint, handler),
	})
}

func (s *Server) Delete(endpoint string, handler Handler) {
	s.handlers = append(s.handlers, route{
		endpoint: endpoint,
		matcher:  buildMatcher(endpoint),
		handler:  MethodMiddleware(DELETE, endpoint, handler),
	})
}

func buildMatcher(endpoint string) func(string) bool {
	re := regexp.MustCompile(`:([a-zA-Z0-9_]+)`)
	pattern := re.ReplaceAllString(endpoint, `(?P<$1>[^/]+)`)
	re = regexp.MustCompile(`\{([a-zA-Z0-9_]+)\}`)
	pattern = re.ReplaceAllString(pattern, `(?P<$1>[^/]+)`)

	reg := regexp.MustCompile("^" + pattern + "$")

	return func(target string) bool {
		return reg.MatchString(target)
	}
}

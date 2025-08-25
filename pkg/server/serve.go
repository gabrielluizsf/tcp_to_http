package server

import (
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/gabrielluizsf/tcp_to_http/pkg/request"
	"github.com/gabrielluizsf/tcp_to_http/pkg/response"
	"github.com/i9si-sistemas/stringx"
)

type Server struct {
	Addr     string
	Listener net.Listener
	handlers []route
	closed   bool
}

type route struct {
	endpoint string
	matcher  func(string) bool
	handler  Handler
}

var ErrInvalidPort = errors.New("invalid port")

func New(port string) (*Server, error) {
	v, err := stringx.NewParser(port).Int()
	if err != nil {
		return nil, errors.Join(ErrInvalidPort, err)
	}
	return serve(int(v))
}

func serve(port int) (*Server, error) {
	server := &Server{
		Addr:     fmt.Sprintf(":%d", port),
		handlers: make([]route, 0),
	}
	if err := runServer(server); err != nil {
		return nil, err
	}
	return server, nil
}

func runServer(s *Server) (err error) {
	s.Listener, err = net.Listen("tcp", s.Addr)
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	go s.listen()
	return nil
}

func (s *Server) Close() error {
	s.closed = true
	return nil
}

func (s *Server) listen() {
	for {
		conn, err := s.Listener.Accept()
		if s.closed {
			return
		}
		if err != nil {
			return
		}
		defer conn.Close()
		go s.handle(conn)
	}
}

func (s *Server) handle(conn io.ReadWriteCloser) {
	defer func() {
		if err := conn.Close(); err != nil {
			fmt.Printf("Error closing connection: %v\n", err)
		}
	}()
	responseWriter := response.NewWriter(conn)
	req, err := request.NewFromReader(conn)
	if err != nil {
		responseWriter.WriteStatusLine(response.StatusBadRequest)
		responseWriter.WriteHeaders(response.GetDefaultHeaders(0))
		return
	}
	for _, r := range s.handlers {
		if r.matcher(req.Line.Target) {
			r.handler(responseWriter, req)
			return
		}
	}

}

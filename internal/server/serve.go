package server

import (
	"fmt"
	"io"
	"net"

	"github.com/gabrielluizsf/tcp_to_http/internal/request"
	"github.com/gabrielluizsf/tcp_to_http/internal/response"
)

type Server struct {
	Addr     string
	Listener net.Listener
	handlers map[string]Handler
	closed   bool
}

func Serve(port int) (*Server, error) {
	server := &Server{
		Addr:     fmt.Sprintf(":%d", port),
		handlers: make(map[string]Handler),
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
	s.handlers[req.Line.Target](responseWriter, req)
}

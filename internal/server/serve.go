package server

import (
	"fmt"
	"io"
	"net"
)

type Server struct {
	Addr     string
	Listener net.Listener
	closed   bool
}

func Serve(port int) (*Server, error) {
	server := &Server{
		Addr: fmt.Sprintf(":%d", port),
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
	output := []byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 13\r\n\r\nHello World!")
	conn.Write(output)
	if err := conn.Close(); err != nil {
		fmt.Printf("Error closing connection: %v\n", err)
	}
}

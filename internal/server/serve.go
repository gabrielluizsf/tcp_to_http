package server

import (
	"bytes"
	"fmt"
	"io"
	"net"

	"github.com/gabrielluizsf/tcp_to_http/internal/request"
	"github.com/gabrielluizsf/tcp_to_http/internal/response"
)

type Server struct {
	Addr     string
	Listener net.Listener
	handler  Handler
	closed   bool
}

func Serve(port int, h Handler) (*Server, error) {
	server := &Server{
		Addr:    fmt.Sprintf(":%d", port),
		handler: h,
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
	headers := response.GetDefaultHeaders(0)
	req, err := request.NewFromReader(conn)
	if err != nil {
		response.WriteStatusLine(conn, response.StatusBadRequest)
		response.WriteHeaders(conn, headers)
		return
	}
	writer := bytes.NewBuffer([]byte{})
	handlerError := s.handler(writer, req)
	body := writer.Bytes()
	status := response.StatusOK
	if handlerError != nil {
		body = []byte(handlerError.Message)
		status = handlerError.StatusCode
	}
	contentLen := len(body)
	headers.Replace("Content-Length", fmt.Sprint(contentLen))
	response.WriteStatusLine(conn, status)
	response.WriteHeaders(conn, headers)
	if _, err := conn.Write(body); err != nil {
		fmt.Printf("Error writing response body: %v\n", err)
	}
}

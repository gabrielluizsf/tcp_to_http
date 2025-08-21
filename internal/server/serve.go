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
	if handlerError != nil {
		headers.Replace("Content-Length", fmt.Sprint(len(handlerError.Message)))
		response.WriteStatusLine(conn, handlerError.StatusCode)
		response.WriteHeaders(conn, headers)
		conn.Write([]byte(handlerError.Message))
		return
	}
	body := writer.Bytes()
	contentLen := len(body)
	headers.Replace("Content-Length", fmt.Sprint(contentLen))
	response.WriteStatusLine(conn, response.StatusOK)
	response.WriteHeaders(conn, headers)
	if _, err := conn.Write(body); err != nil {
		fmt.Printf("Error writing response body: %v\n", err)
	}
}

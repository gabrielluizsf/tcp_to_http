package server

import (
	"errors"

	internalServer "github.com/gabrielluizsf/tcp_to_http/internal/server"
	"github.com/i9si-sistemas/stringx"
)

type Server struct {
	*internalServer.Server
}

func New(port string) (*Server, error) {
	portInt, err := stringx.NewParser(port).Int()
	if err != nil {
		return nil, errors.Join(errors.New("invalid port"), err)
	}
	sv, err := internalServer.Serve(int(portInt))
	if err != nil {
		return nil, err
	}
	return &Server{sv}, nil
}

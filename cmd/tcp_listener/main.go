package main

import (
	"fmt"
	"io"
	"net"

	"github.com/gabrielluizsf/tcp_to_http/pkg/request"
)

var (
	port     = "42069"
	listener net.Listener
)

func init() {
	address := fmt.Sprintf(":%s", port)
	l, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}
	listener = l
}

func main() {
	readFromNetConn()
}

func readFromNetConn() {
	for {
		conn, err := netConnReader()
		if err != nil {
			panic(err)
		}
		r, err := request.NewFromReader(conn)
		if err != nil {
			panic(err)
		}
		fmt.Println("Request Line:")
		fmt.Printf("Method: %s\n", r.Line.Method)
		fmt.Printf("Target: %s\n", r.Line.Target)
		fmt.Printf("Version: %s\n", r.Line.Version)
		fmt.Println("Headers:")
		for key, value := range r.Headers {
			fmt.Printf("%s: %s\n", key, value)
		}
		if len(r.Body) > 0 {
			fmt.Println("Body:")
			fmt.Printf("%s\n", r.Body)
		}

	}
}

func netConnReader() (io.ReadCloser, error) {
	return listener.Accept()
}

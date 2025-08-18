package main

import (
	"fmt"
	"io"
	"net"

	"github.com/gabrielluizsf/tcp_to_http/pkg/channel"
	"github.com/i9si-sistemas/stringx"
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
		for line := range getLinesChannel(conn) {
			fmt.Println(stringx.Space.String(), line)
		}
	}
}

func netConnReader() (io.ReadCloser, error) {
	return listener.Accept()
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	return channel.Lines{}.Get(f)
}

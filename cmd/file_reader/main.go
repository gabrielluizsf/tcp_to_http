package main

import (
	"fmt"
	"io"
	"os"

	"github.com/gabrielluizsf/tcp_to_http/pkg/channel"
	"github.com/i9si-sistemas/stringx"
)

func main() {
	readFromFile()
}

func readFromFile() {
	f, err := fileReader()
	if err != nil {
		panic(err)
	}
	for line := range getLinesChannel(f) {
		fmt.Println(stringx.Space.String(),line)
	}
}

func fileReader() (io.ReadCloser, error) {
	return os.Open("messages.txt")
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	return channel.Lines{}.Get(f)
}

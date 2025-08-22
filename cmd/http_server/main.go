package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gabrielluizsf/tcp_to_http/internal/request"
	"github.com/gabrielluizsf/tcp_to_http/internal/response"
	"github.com/gabrielluizsf/tcp_to_http/internal/server"
	"github.com/i9si-sistemas/stringx"
)

const port = 42069

func main() {
	server, err := server.Serve(port, func(res *response.Writer, req *request.Request) {
		endpoint := stringx.New(req.Line.Target)
		headers := response.GetDefaultHeaders(0)
		contentType, body := respond200()
		statusCode := response.StatusOK
		switch {
		case endpoint.Includes("/yourproblem"):
			statusCode = response.StatusBadRequest
			contentType, body = respond400()
		case endpoint.Includes("/myproblem"):
			statusCode = response.StatusInternalServerError
			contentType, body = respond500()
		case endpoint.HasPrefix("/httpbin/stream"):
			r, err := http.Get("https://httpbin.org/" + endpoint.String()[len("/httpbin/"):])
			if err != nil {
				contentType, body = respond500()
				statusCode = response.StatusInternalServerError
			} else {
				res.WriteStatusLine(response.StatusOK)
				headers.Delete("Content-Length")
				headers.Set("transfer-encoding", "chunked")
				headers.Set("Content-Type", "text/plain")
				res.WriteHeaders(headers)
				for {
					data := make([]byte, 32)
					n, err := r.Body.Read(data)
					if err != nil {
						break
					}
					res.WriteBody(fmt.Appendf([]byte{}, "%x\r\n", n))
					res.WriteBody(data[:n])
					res.WriteBody([]byte("\r\n"))
				}
				res.WriteBody([]byte("0\r\n\r\n"))
				return
			}
		}
		res.WriteStatusLine(statusCode)
		headers.Replace("Content-Length", fmt.Sprint(len(body)))
		headers.Replace("Content-Type", contentType)
		res.WriteHeaders(headers)
		res.WriteBody(body)
	})
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

const HTML_CONTENT_TYPE = "text/html; charset=utf-8"

func respond400() (contentType string, data []byte) {
	contentType = HTML_CONTENT_TYPE
	data = []byte(`<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`)
	return
}

func respond500() (contentType string, data []byte) {
	contentType = HTML_CONTENT_TYPE
	data = []byte(`<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`)
	return
}

func respond200() (contentType string, data []byte) {
	contentType = HTML_CONTENT_TYPE
	data = []byte(`<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`)
	return
}

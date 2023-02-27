package main

import (
	"bytes"
	"fmt"
	"golang.org/x/net/websocket"
	"io"
	"log"
	"net/http"
	"os"
	"syscall"
)

var c = make(chan []byte, 10)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./dist")))
	http.Handle("/ws", websocket.Handler(msgHandler))
	http.HandleFunc("/logs", receiveLogs)

	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func msgHandler(ws *websocket.Conn) {
	defer func(ws *websocket.Conn) {
		_ = ws.Close()
	}(ws)

	for {
		select {
		case msg := <-c:
			buf := bytes.NewBuffer(msg)
			err := websocket.Message.Send(ws, buf.String())
			_, _ = fmt.Fprintf(os.Stdout, "received from channel: \n%s\n", buf.String())
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "error on sending: %s\n", err)
				if err == syscall.EPIPE {
					continue
				}
			}
		}
	}
}

func receiveLogs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(r.Body)

	b, err := io.ReadAll(r.Body)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error on reading body: \n%s\n", err)
	}

	c <- b
	_, _ = fmt.Fprintf(os.Stdout, "send to channel: \n%s\n", string(b))
	w.WriteHeader(http.StatusOK)
}

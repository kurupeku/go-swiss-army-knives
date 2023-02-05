package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/websocket"
)

var c = make(chan string, 10)

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

	err := websocket.Message.Send(ws, "connect successfully!")
	if err != nil {
		log.Fatalln(err)
	}

	for {
		select {
		case msg := <-c:
			err := websocket.Message.Send(ws, msg)
			if err != nil {
				_, _ = fmt.Fprintln(os.Stderr, err)
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

	scanner := bufio.NewScanner(r.Body)
	for scanner.Scan() {
		c <- scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintln(os.Stderr, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

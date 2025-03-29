package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"

	"golang.org/x/net/websocket"
)

type Server struct {
	conns map[*websocket.Conn]bool
	mu    sync.Mutex
}

func NewServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}
func (s *Server) handleHtmx(ws *websocket.Conn) {
	fmt.Println("Connection establishing ", ws.RemoteAddr())

	s.conns[ws] = true
	s.readHtmx(ws)
}

// func handleObs(ws *websocket.Conn) {

// }

func (s *Server) readHtmx(ws *websocket.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("read err;", err)
			continue
		}
		msg := buf[:n]
		fmt.Println(string(msg))
		s.mu.Lock()
		defer s.mu.Unlock()
		ws.Write([]byte("RECV"))
		s.mu.Unlock()
	}
}

func main() {
	// concurrency for writes
	//chHtmx := make(chan int)
	//chObs := make(chan int)
	server := NewServer()
	http.Handle("/hx-obs", websocket.Handler(server.handleHtmx))
	http.ListenAndServe(":8901", nil)

}

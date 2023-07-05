package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

type Server struct {
	conns map [*websocket.Conn] bool

}

func newServer() * Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}

func (s*Server) handleWSOrderbook(ws*websocket.Conn) {
	fmt.Println("new incoming connection from client to orderbook feed",ws.RemoteAddr())
	for{
		payload:=fmt.Sprintf("orderbook data = >%d\n",time.Now().UnixNano())
		ws.Write([]byte(payload))
		time.Sleep(time.Second*2)
	}
}


func (s*Server) handleWS (ws * websocket.Conn){
	fmt.Println("new incoming connection from client",ws.RemoteAddr())
	s.conns[ws]=true
	s.readLoop(ws)
}

func (s*Server) readLoop(ws*websocket.Conn) {

	buffer:=make([]byte,1024)

	for {
		n,error:=ws.Read(buffer)
		if error!=nil {
			if error==io.EOF{
				break;
			}
			fmt.Println("read error",error)
			continue
		}
		message:=buffer[:n]
		s.broadcast(message)
	}

}

func (s*Server) broadcast (b[] byte) {
for ws:=range(s.conns) {
	go func (ws*websocket.Conn)  {
		if _,error:=ws.Write(b); error!=nil {
			fmt.Println("write error :", error)
		}
	}(ws)
}

}



func main () {

server:=newServer()
http.Handle("/ws",websocket.Handler(server.handleWS))
http.ListenAndServe(":3000", nil)




}
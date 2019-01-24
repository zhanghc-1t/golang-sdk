package main

import (
	lib "golang-sdk/lesson/telnet-chat/lib"
	"log"
	"net"
	"time"
)

func main() {
	listener, err := net.Listen("tcp", ":9999")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	server, err := lib.NewRoom(5, time.Hour*1, listener)
	defer server.Stop()

	server.Start()
}

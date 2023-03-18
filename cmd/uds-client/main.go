package main

import (
	"github.com/quells/hello-uds/internal/env"
	"log"
	"net"
)

func main() {
	addr, err := net.ResolveUnixAddr("unix", env.Config.SocketFile)
	if err != nil {
		log.Println(err)
		return
	}

	conn, err := net.DialUnix("unix", nil, addr)
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		if cErr := conn.Close(); cErr != nil {
			log.Println(cErr)
		}
	}()

	_, err = conn.Write([]byte("hello"))
	if err != nil {
		log.Println(err)
		return
	}

	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("got %d bytes: %v", n, buf[:n])
}

package main

import (
	"errors"
	"github.com/quells/exit"
	"github.com/quells/hello-uds/internal/env"
	"log"
	"net"
	"os"
	"sync"
)

func main() {
	addr, err := net.ResolveUnixAddr("unix", env.Config.SocketFile)
	if err != nil {
		log.Println(err)
		return
	}

	ln, err := net.ListenUnix("unix", addr)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("listening at %s", env.Config.SocketFile)
	defer func() {
		if rErr := os.Remove(env.Config.SocketFile); rErr != nil {
			pErr := rErr.(*os.PathError)
			if !os.IsNotExist(pErr.Err) {
				log.Printf("failed to remove socket file: %v", pErr.Err)
			}
		}
	}()

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go serve(wg, ln)

	sig := <-exit.Notification()
	log.Printf("got %s signal", sig)

	err = ln.Close()
	if err != nil {
		log.Printf("failed to close listener: %v", err)
	}

	log.Println("waiting for goroutines")
	wg.Wait()
}

func serve(wg *sync.WaitGroup, ln *net.UnixListener) {
	defer wg.Done()

	for {
		conn, err := ln.Accept()
		if err != nil {
			if !errors.Is(err, net.ErrClosed) {
				log.Printf("failed to accept connection: %v", err)
			}
			return
		}

		wg.Add(1)
		go handle(wg, conn)
	}
}

func handle(wg *sync.WaitGroup, conn net.Conn) {
	defer wg.Done()
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("failed to close connection: %v", err)
		}
	}()

	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		log.Printf("failed to read from connection: %v", err)
		return
	}

	log.Printf("got %v", buf[:n])

	_, err = conn.Write([]byte("OK\n"))
	if err != nil {
		log.Printf("failed to write to connection: %v", err)
	}
}

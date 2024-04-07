package main

import (
	"fmt"
	"hermes/internal/data"
	"hermes/internal/hermes"
	"net"
	"sync"
	"time"
)

const (
	reqs = 10_000
)

func main() {
	var wg sync.WaitGroup

	h := hermes.New(":3333")
	go h.Listen()

	time.Sleep(50 * time.Millisecond) // waits server

	started := time.Now()
	for i := 0; i < reqs; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			send(*data.New("SET", fmt.Sprintf("foo=bar-%d", i)))
		}()
	}
	wg.Wait()

	fmt.Println("SET benchmark:")
	fmt.Println("Req/s:", reqs/time.Since(started).Seconds())

	started = time.Now()
	for i := 0; i < reqs; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			send(*data.New("GET", "foo"))
		}()
	}
	wg.Wait()

	fmt.Println("GET benchmark:")
	fmt.Println("Req/s:", reqs/time.Since(started).Seconds())
}

func send(d data.Data) *data.Data {
	conn, err := net.Dial("tcp", ":3333")
	if err != nil {
		fmt.Println(err)
	}

	defer conn.Close()

	conn.Write([]byte(d.Decode()))

	buff := make([]byte, 1024)
	n, _ := conn.Read(buff)

	resp, _ := data.FromStr(string(buff[:n]))
	return resp
}

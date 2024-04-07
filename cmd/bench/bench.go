package main

import (
	"fmt"
	"hermes/internal/data"
	"net"
	"sync"
	"time"
)

const (
	total = 100_000
	conc  = 100
)

func main() {
	var wg sync.WaitGroup

	started := time.Now()

	for i := 0; i < conc; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			for i := 0; i < total/conc; i++ {
				send(*data.New("SET", fmt.Sprintf("foo=bar-%d", i)))
			}
		}()
	}

	wg.Wait()

	fmt.Println("SET benchmark:")
	fmt.Println("Req/s:", total/time.Since(started).Seconds())

	started = time.Now()

	for i := 0; i < conc; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			for i := 0; i < total/conc; i++ {
				send(*data.New("GET", fmt.Sprintf("foo-%d", i)))
			}
		}()
	}

	wg.Wait()

	fmt.Println("GET benchmark:")
	fmt.Println("Req/s:", total/time.Since(started).Seconds())
}

func send(d data.Data) *data.Data {
	conn, err := net.Dial("tcp", ":3333")
	if err != nil {
		fmt.Println(err)
	}

	defer conn.Close()

	conn.Write([]byte(d.Decode()))

	buff := make([]byte, 128)
	n, _ := conn.Read(buff)

	resp, _ := data.FromStr(string(buff[:n]))
	return resp
}

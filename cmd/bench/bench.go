package main

import (
	"fmt"
	"hermes/internal/data"
	"hermes/internal/pool"
	"log"
	"net"
	"sync"
	"time"
)

const (
	total = 10_000
)

var tcpPool pool.Pool = *pool.New(100, func() (net.Conn, error) {
	return net.Dial("tcp", ":3333")
})

func main() {
	var wg sync.WaitGroup

	started := time.Now()

	for i := 0; i < total; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			data := send(*data.New("SET", fmt.Sprintf("foo-%d=bar-%d", i, i)))
			if data.Content != fmt.Sprintf("foo-%d=bar-%d", i, i) {
				log.Panic("ops set")
			}
		}()
	}
	wg.Wait()

	fmt.Println("SET benchmark:")
	fmt.Println("Req/s:", total/time.Since(started).Seconds())

	started = time.Now()

	for i := 0; i < total; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			data := send(*data.New("GET", fmt.Sprintf("foo-%d", i)))
			if data.Content != fmt.Sprintf("foo-%d=bar-%d", i, i) {
				log.Panic("ops get")
			}
		}()
	}
	wg.Wait()

	fmt.Println("GET benchmark:")
	fmt.Println("Req/s:", total/time.Since(started).Seconds())
}

func send(d data.Data) *data.Data {
	// conn, err := net.Dial("tcp", ":3333")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// defer conn.Close()

	conn := tcpPool.Get()
	defer tcpPool.Release(conn)

	conn.Write([]byte(d.Decode()))

	buff := make([]byte, 128)
	n, _ := conn.Read(buff)

	resp, _ := data.FromStr(string(buff[:n]))
	return resp
}

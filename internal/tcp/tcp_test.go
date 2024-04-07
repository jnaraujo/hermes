package tcp

import (
	"hermes/internal/data"
	"net"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	go func() {
		s := New(":8080")

		s.AddHandler(func(d *data.Data) *data.Data {
			return d
		})

		err := s.Listen()
		if err != nil {
			t.Error("Listen() failed")
		}
	}()

	time.Sleep(10 * time.Millisecond) // wait for server to start

	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		t.Error("Dial failed")
	}
	defer conn.Close()

	conn.Write([]byte("GET\nkey"))

	buff := make([]byte, 1024)
	n, err := conn.Read(buff)
	if err != nil {
		t.Error("Read failed")
	}

	if string(buff[:n]) != "GET\nkey" {
		t.Error("Data not equal")
	}

	conn.Write([]byte("SET\nkey=val"))
	buff = make([]byte, 1024)
	n, err = conn.Read(buff)
	if err != nil {
		t.Error("Read failed")
	}

	if string(buff[:n]) != "SET\nkey=val" {
		t.Error("Data not equal")
	}
}

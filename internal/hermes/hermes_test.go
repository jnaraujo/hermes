package hermes

import (
	"fmt"
	"hermes/internal/data"
	"net"
	"testing"
	"time"
)

func TestHermes(t *testing.T) {
	hermes := New(":3333")
	go hermes.Listen()

	time.Sleep(10 * time.Millisecond) // waits server

	if data := send(*data.New("GET", "foo")); data.Content != "foo=" {
		t.Error("Data not equal")
	}

	if data := send(*data.New("SET", "foo=bar")); data.Content != "foo=bar" {
		t.Error("Data not equal")
	}

	if data := send(*data.New("GET", "foo")); data.Content != "foo=bar" {
		t.Error("Data not equal")
	}

	if data := send(*data.New("SET", "foo=baz")); data.Content != "foo=baz" {
		t.Error("Data not equal")
	}

	if data := send(*data.New("GET", "foo")); data.Content != "foo=baz" {
		t.Error("Data not equal")
	}
}

func TestConcurrentHermes(t *testing.T) {
	hermes := New(":3333")
	go hermes.Listen()

	time.Sleep(10 * time.Millisecond) // waits server

	for i := 0; i < 100; i++ {
		go func() {
			sent := fmt.Sprintf("foo=bar-%d", i)
			if data := send(*data.New("SET", sent)); data.Content != sent {
				t.Error("Data not equal")
			}
		}()
	}
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

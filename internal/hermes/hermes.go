package hermes

import (
	"fmt"
	"hermes/internal/data"
	"hermes/internal/tcp"
	"strings"
	"sync"
)

type Hermes struct {
	addr string

	storage map[string]string
	mu      sync.RWMutex
}

func New(addr string) *Hermes {
	return &Hermes{
		addr:    addr,
		storage: make(map[string]string),
	}
}

func (h *Hermes) Listen() {
	server := tcp.New(h.addr)
	server.AddHandler(func(d *data.Data) *data.Data {
		switch d.Cmd {
		case "GET":
			h.mu.RLock()
			defer h.mu.RUnlock()

			return data.New("GET", fmt.Sprintf("%s=%s", d.Content, h.storage[d.Content]))
		case "SET":
			h.mu.Lock()
			defer h.mu.Unlock()

			key_val := strings.Split(d.Content, "=")
			if len(key_val) != 2 {
				return data.New("ERROR", "invalid data")
			}
			h.storage[key_val[0]] = key_val[1]
			return data.New("SET", fmt.Sprintf("%s=%s", key_val[0], key_val[1]))
		default:
			return data.New("ERROR", "invalid command")
		}
	})

	server.Listen()
}

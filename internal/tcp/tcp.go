package tcp

import (
	"errors"
	"fmt"
	"hermes/internal/data"
	"io"
	"net"
)

type Handler func(d *data.Data) *data.Data

type Server struct {
	addr    string
	handler Handler
}

func New(addr string) *Server {
	return &Server{
		addr: addr,
	}
}

func (s *Server) AddHandler(handler Handler) {
	s.handler = handler
}

func (s *Server) Listen() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	buff := make([]byte, 32)

	defer conn.Close()

	for {
		n, err := conn.Read(buff)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}

			fmt.Println("Error reading:", err)
			return
		}

		raw := string(buff[:n])
		data, err := data.FromStr(raw)

		if err != nil {
			fmt.Println("Error decoding:", err)
			return
		}

		resp := s.handler(data)
		conn.Write([]byte(resp.Decode()))
	}
}

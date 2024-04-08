package pool

import "net"

type Pool struct {
	pool chan net.Conn
}

func New(n int, factory func() (net.Conn, error)) *Pool {
	pool := Pool{
		pool: make(chan net.Conn, n),
	}

	for i := 0; i < n; i++ {
		conn, err := factory()
		if err != nil {
			println("Conn Error: ", err)
			i--
			continue
		}
		pool.pool <- conn
	}

	return &pool
}

func (p *Pool) Get() net.Conn {
	return <-p.pool
}

func (p *Pool) Release(conn net.Conn) {
	select {
	case p.pool <- conn:
	default:
		_ = conn.Close()
	}
}

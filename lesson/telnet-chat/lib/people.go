package lib

import (
	"net"
	"time"
)

type People struct {
	conn     net.Conn
	name     string
	created  time.Time
	messaged time.Time
	message  chan *Message
}

func (p *People) Close() {
	p.conn.Close()
}

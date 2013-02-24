// common basestation / driver code
package hermes

import (
	"net"
	"encoding/gob"
	"time"
)

const (
	PORT = 1337
	Rate = (100 * time.Millisecond)
)

type Keys struct {
	UP bool
	DOWN bool
	LEFT bool
	RIGHT bool
	GEAR1 bool
	GEAR2 bool
	GEAR3 bool
}

type Conn struct {
	conn net.Conn
	enc *gob.Encoder
	dec *gob.Decoder
}

type Packet struct {
	Left byte
	Right byte
	Ack bool
}

func New(n net.Conn) *Conn {
	enc := gob.NewEncoder(n)
	dec := gob.NewDecoder(n)
	return &Conn{
		conn: n,
		enc: enc,
		dec: dec,
	}
}

func (c *Conn) Send(p Packet) error {
	err := c.enc.Encode(&p)
	return err
}

func (c *Conn) Recv() (Packet, error) {
	p := Packet{}
	err := c.dec.Decode(&p)
	return p, err
}


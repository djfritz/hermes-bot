package main

import (
	"fmt"
	"goserial"
	"os"
	"os/signal"
	"hermes"
	"net"
	"time"
)

var (
	packets chan hermes.Packet
	done chan bool
)

func main() {
	s, err := goserial.Open("/dev/ttyUSB0", 9600)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	packets = make(chan hermes.Packet)
	done = make(chan bool)

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)
	go func() {
		<-sig
		s.Close()
		os.Exit(0)
	}()

	laddr := fmt.Sprintf(":%v", hermes.PORT)
	ln, err := net.Listen("tcp", laddr)
	if err != nil {
		fmt.Println(err)
		s.Close()
		os.Exit(-1)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			s.Close()
			os.Exit(-1)
		}
		h := hermes.New(conn)
		go connHandler(h)

		CONN_LOOP:
		for {
			select {
			case packet := <-packets:
				if int(packet.Left) == 0 {
					s.Write([]byte{byte(0)})
				} else {
					data := []byte{packet.Left, packet.Right}
					s.Write(data)
				}
			case <-time.After(4 * hermes.Rate):
				fmt.Println("loss of signal!")
				s.Write([]byte{byte(0)})
			case <-done:
				break CONN_LOOP
			}
		}
	}
}

func connHandler(h *hermes.Conn) {
	for {
		packet, err := h.Recv()
		if err != nil {
			fmt.Println(err)
			break
		}

		packets <- packet

		p := hermes.Packet{
			Ack: true,
		}
		err = h.Send(p)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
	done <- true
}

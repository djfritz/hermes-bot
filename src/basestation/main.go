package main

import (
	"fmt"
	"goserial"
	"os"
	"os/signal"
	"net"
	"io"
)

type packet struct {
	left byte
	right byte
}

var incoming chan packet

func main() {
	s, err := goserial.Open("/dev/ttyUSB0", 9600)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)
	go func() {
		<-sig
		s.Close()
		os.Exit(0)
	}()

	incoming = make(chan packet, 1024)

	go listener()

	for {
		packet := <-incoming
		//fmt.Println("incoming", packet)

		if int(packet.left) == 0 {
			s.Write([]byte{byte(0)})
			continue
		}

		data := []byte{packet.left, packet.right}
		s.Write(data)
	}
}

func listener() {
	ln, err := net.Listen("tcp", ":8081")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
	
		for {	
			// handle only one connection
			data := make([]byte, 2)
			n, err := conn.Read(data)
			if err != nil {
				if err == io.EOF {
					continue
				}
				fmt.Println(err)
				break
			}
			//fmt.Println("got", n, data)
			if n != 2 {
				incoming <- packet{
					left: 0,
					right: 0,
				}
				continue
			}
			incoming <- packet{
				left: data[0],
				right: data[1],
			}
		}
	}
}

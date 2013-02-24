package main

import (
	"os/signal"
	"os"
	"time"
	"fmt"
	"hermes"
	"net"
)

const (
	LMIN = 1
	LMAX = 127
	LNEUTRAL = 64

	RMIN = 128
	RMAX = 255
	RNEUTRAL = 192
)

var (
	sig = make(chan os.Signal, 1024)
)

func main() {
	signal.Notify(sig, os.Interrupt)

	if len(os.Args) != 2 {
		fmt.Println("usage: driver <host>")
		os.Exit(-1)
	}

	addr := fmt.Sprintf("%v:%v", os.Args[1], hermes.PORT)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	h := hermes.New(conn)

	err = kbd_init()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	ncurses_init()
	defer ncurses_endwin()

	// i can get feedback from the bot (heartbeat, bytes transferred, led color)
	// i can also get feedback from the keyboard, which is polled in an event loop

	// the keyboard input is polled, and sent to the bot as a command packet every N units of time
	// the screen updates every M units of time, and should be faster than N

	// on network feedback, i'll update what the charts, output for the screen should be, and the screen 
	// redraw will just consume that data when needed. I don't want lag with ill-consumed channels, etc.

	//go RedrawHandler()
	go EventLoop(h)

	<-sig
}

func EventLoop(h *hermes.Conn) {
	for {
		left, right, err := GetMotorValues()
		if err != nil {
			fmt.Println(err)
			sig <- os.Interrupt
		}
		s := fmt.Sprintf("%d %d\n", left, right)
		ncurses_printw(s)
		ncurses_move(0,0)
		ncurses_refresh()

		packet := hermes.Packet{
			Left: left,
			Right: right,
		}
		err = h.Send(packet)
		if err != nil {
			fmt.Println(err)
			sig <- os.Interrupt
		}

		// get ack
		packet, err = h.Recv()
		if err != nil {
			fmt.Println(err)
			sig <- os.Interrupt
		}

		time.Sleep(hermes.Rate)
	}
}

func GetMotorValues() (byte, byte, error) {
	var left byte
	var right byte
	keys, err := GetKeys()
	ncurses_printw(fmt.Sprintf("%v", keys))
	ncurses_move(1,0)
	if err != nil {
		return 0, 0, err
	}

	switch {
	case keys.UP && keys.LEFT:
		left = LNEUTRAL
		right = RMAX
	case keys.UP && keys.RIGHT:
		left = LMAX
		right = RNEUTRAL
	case keys.DOWN && keys.LEFT:
		left = LNEUTRAL
		right = RMIN
	case keys.DOWN && keys.RIGHT:
		left = LMIN
		right = RNEUTRAL
	case keys.LEFT:
		left = LMIN
		right = RMAX
	case keys.RIGHT:
		left = LMAX
		right = RMIN
	case keys.UP:
		left = LMAX
		right = RMAX
	case keys.DOWN:
		left = LMIN
		right = RMIN
	default:
		left = 0
		right = 0
	}

	return left, right, nil
}

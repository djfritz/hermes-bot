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
	LMIN_3 = 1
	LMAX_3 = 127
	LMIN_2 = 32
	LMAX_2 = 96
	LMIN_1 = 54
	LMAX_1 = 74
	LNEUTRAL = 64

	RMIN_3 = 128
	RMAX_3 = 255
	RMIN_2 = 159
	RMAX_2 = 224
	RMIN_1 = 182
	RMAX_1 = 202
	RNEUTRAL = 192
)

var (
	sig = make(chan os.Signal, 1024)

	LMIN byte
	LMAX byte

	RMIN byte
	RMAX byte
	GEAR = 1
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
		ncurses_move(1,0)
		ncurses_printw(fmt.Sprintf("GEAR %v", GEAR))
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
	if err != nil {
		return 0, 0, err
	}

	// check for gear shifting
	switch {
	case keys.GEAR1:
		GEAR = 1
		LMIN = LMIN_1
		LMAX = LMAX_1
		RMIN = RMIN_1
		RMAX = RMAX_1
	case keys.GEAR2:
		GEAR = 2
		LMIN = LMIN_2
		LMAX = LMAX_2
		RMIN = RMIN_2
		RMAX = RMAX_2
	case keys.GEAR3:
		GEAR = 3
		LMIN = LMIN_3
		LMAX = LMAX_3
		RMIN = RMIN_3
		RMAX = RMAX_3
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

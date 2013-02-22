package main

import (
	"os/signal"
	"os"
	"time"
	"fmt"
)

const (
	heartbeat_rate = (100 * time.Millisecond)
)

func main() {
	sig := make(chan os.Signal, 1024)
	signal.Notify(sig, os.Interrupt)

	ncurses_init()
	defer ncurses_endwin()

	err := kbd_init()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	// i can get feedback from the bot (heartbeat, bytes transferred, led color)
	// i can also get feedback from the keyboard, which is polled in an event loop

	// the keyboard input is polled, and sent to the bot as a command packet every N units of time
	// the screen updates every M units of time, and should be faster than N

	// on network feedback, i'll update what the charts, output for the screen should be, and the screen 
	// redraw will just consume that data when needed. I don't want lag with ill-consumed channels, etc.

	//go RedrawHandler()
	go EventLoop()

	<-sig
}

func EventLoop() {
	for {
		k, err := GetKeys()
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		s := fmt.Sprintf("%v\n", k)
		ncurses_printw(s)
		ncurses_move(0,0)
		ncurses_refresh()
		time.Sleep(heartbeat_rate)
	}
}

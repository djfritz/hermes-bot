package main

import (
	"os/signal"
	"os"
	"time"
)

const (
	heartbeat_rate = (1 * time.Second)
)

func main() {
	sig := make(chan os.Signal, 1024)
	signal.Notify(sig, os.Interrupt)

	ncurses_init()
	defer ncurses_endwin()

	// i can get feedback from the bot (heartbeat, bytes transferred, led color)
	// i can also get feedback from the keyboard, which is polled in an event loop

	// the keyboard input is polled, and sent to the bot as a command packet every N units of time
	// the screen updates every M units of time, and should be faster than N

	// on network feedback, i'll update what the charts, output for the screen should be, and the screen 
	// redraw will just consume that data when needed. I don't want lag with ill-consumed channels, etc.

	//go RedrawHandler()
	go KeyboardPoller()

	<-sig
}

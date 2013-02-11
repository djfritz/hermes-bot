package main

import (
	"time"
	"fmt"
	"os"
	"kbdmap"
)

func KeyboardPoller() {
	f, err := os.Open("/dev/input/event3")
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	for {
		km, err := kbdmap.GetKBDMap(f)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		if kbdmap.IsPressed(km, kbdmap.KEY_W) {
			ncurses_printw("W\n")
		}
		if kbdmap.IsPressed(km, kbdmap.KEY_A) {
			ncurses_printw("A\n")
		}
		if kbdmap.IsPressed(km, kbdmap.KEY_S) {
			ncurses_printw("S\n")
		}
		if kbdmap.IsPressed(km, kbdmap.KEY_D) {
			ncurses_printw("D\n")
		}
		ncurses_refresh()
		time.Sleep(heartbeat_rate)
	}
}

package main

import (
	"os"
	"kbdmap"
	"fmt"
	"time"
)

func main() {
	f, err := os.Open("/dev/input/event3")
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	for {
		s, err := kbdmap.GetKBDMap(f)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}

		if kbdmap.IsPressed(s, kbdmap.KEY_W) {
			fmt.Println("W")
		}

		if kbdmap.IsPressed(s, kbdmap.KEY_A) {
			fmt.Println("A")
		}

		if kbdmap.IsPressed(s, kbdmap.KEY_S) {
			fmt.Println("S")
		}

		if kbdmap.IsPressed(s, kbdmap.KEY_D) {
			fmt.Println("D")
		}

		fmt.Println("---")

		time.Sleep(1 * time.Second)
	}
}

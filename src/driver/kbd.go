package main

import (
	"os"
	"kbdmap"
)

var (
	f_kbd *os.File
)

type keys struct {
	UP bool
	DOWN bool
	LEFT bool
	RIGHT bool
}

func kbd_init() error {
	var err error
	f_kbd, err = os.Open("/dev/input/event3")
	return err
}

func GetKeys() (keys, error) {
	k := keys{}

	km, err := kbdmap.GetKBDMap(f_kbd)
	if err != nil {
		return k, err
	}

	if kbdmap.IsPressed(km, kbdmap.KEY_W) {
		k.UP = true
	}
	if kbdmap.IsPressed(km, kbdmap.KEY_S) {
		k.DOWN = true
	}
	if kbdmap.IsPressed(km, kbdmap.KEY_A) {
		k.LEFT = true
	}
	if kbdmap.IsPressed(km, kbdmap.KEY_D) {
		k.RIGHT = true
	}

	return k, nil
}

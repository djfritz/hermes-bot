package main

import (
	"os"
	"kbdmap"
	"hermes"
)

var (
	f_kbd *os.File
)

func kbd_init() error {
	var err error
	f_kbd, err = os.Open("/dev/input/event0")
	return err
}

func GetKeys() (hermes.Keys, error) {
	k := hermes.Keys{}

	km, err := kbdmap.GetKBDMap(f_kbd)
	if err != nil {
		return k, err
	}

	if kbdmap.IsPressed(km, kbdmap.KEY_DOWN) {
		k.UP = true
	}
	if kbdmap.IsPressed(km, kbdmap.KEY_UP) {
		k.DOWN = true
	}
	if kbdmap.IsPressed(km, kbdmap.KEY_LEFT) {
		k.LEFT = true
	}
	if kbdmap.IsPressed(km, kbdmap.KEY_RIGHT) {
		k.RIGHT = true
	}
	if kbdmap.IsPressed(km, kbdmap.KEY_1) {
		k.GEAR1 = true
	}
	if kbdmap.IsPressed(km, kbdmap.KEY_2) {
		k.GEAR2 = true
	}
	if kbdmap.IsPressed(km, kbdmap.KEY_3) {
		k.GEAR3 = true
	}

	return k, nil
}

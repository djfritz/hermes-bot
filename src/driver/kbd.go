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
	f_kbd, err = os.Open("/dev/input/event3")
	return err
}

func GetKeys() (hermes.Keys, error) {
	k := hermes.Keys{}

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

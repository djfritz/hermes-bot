/*

simple shiftbrite test routine

34 - DI
35 - LI
56 - EI
59 - CI

*/

package main

import (
	"os"
	"time"
	"fmt"
)

func main() {
	// export the gpio pins
	f, err := os.OpenFile("/sys/class/gpio/export", os.O_WRONLY | os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}	
	f.WriteString("59") // internal gpio 59 is gpio pin 28 on the header
	f.Close()

	f, err = os.OpenFile("/sys/class/gpio/export", os.O_WRONLY | os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}	
	f.WriteString("56") // internal gpio 59 is gpio pin 28 on the header
	f.Close()
	
	f, err = os.OpenFile("/sys/class/gpio/export", os.O_WRONLY | os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}	
	f.WriteString("35") // internal gpio 59 is gpio pin 28 on the header
	f.Close()

	f, err = os.OpenFile("/sys/class/gpio/export", os.O_WRONLY | os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}	
	f.WriteString("34") // internal gpio 59 is gpio pin 28 on the header
	f.Close()

	// set it as an output
	f, err = os.OpenFile("/sys/class/gpio/gpio59/direction", os.O_WRONLY | os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	f.WriteString("out")
	f.Close()

	f, err = os.OpenFile("/sys/class/gpio/gpio56/direction", os.O_WRONLY | os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	f.WriteString("out")
	f.Close()

	f, err = os.OpenFile("/sys/class/gpio/gpio35/direction", os.O_WRONLY | os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	f.WriteString("out")
	f.Close()

	f, err = os.OpenFile("/sys/class/gpio/gpio34/direction", os.O_WRONLY | os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	f.WriteString("out")
	f.Close()

	// get file handles to the four signals
	DI, err := os.OpenFile("/sys/class/gpio/gpio34/value", os.O_RDWR, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	LI, err := os.OpenFile("/sys/class/gpio/gpio35/value", os.O_RDWR, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	EI, err := os.OpenFile("/sys/class/gpio/gpio56/value", os.O_RDWR, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	CI, err := os.OpenFile("/sys/class/gpio/gpio59/value", os.O_RDWR, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}

	EI.WriteString("0")
	EI.Sync()
	EI.Close()

	var red, green, blue int

	for {
		red = 0x3ff

		WritePacket(CI,DI,LI, red, green, blue)
		time.Sleep(500 * time.Millisecond)

		red = 0
	
		WritePacket(CI,DI,LI, red, green, blue)
		time.Sleep(500 * time.Millisecond)
	}
}

func WritePacket(C, D, L *os.File, red, green, blue int) {
	// write blue, which should be 00, then blue, red, green

	WriteBit(C, D, 0)
	WriteBit(C, D, 0)

	WriteBit(C, D, blue >> 9)
	WriteBit(C, D, blue >> 8)
	WriteBit(C, D, blue >> 7)
	WriteBit(C, D, blue >> 6)
	WriteBit(C, D, blue >> 5)
	WriteBit(C, D, blue >> 4)
	WriteBit(C, D, blue >> 3)
	WriteBit(C, D, blue >> 2)
	WriteBit(C, D, blue >> 1)
	WriteBit(C, D, blue)

	WriteBit(C, D, red >> 9)
	WriteBit(C, D, red >> 8)
	WriteBit(C, D, red >> 7)
	WriteBit(C, D, red >> 6)
	WriteBit(C, D, red >> 5)
	WriteBit(C, D, red >> 4)
	WriteBit(C, D, red >> 3)
	WriteBit(C, D, red >> 2)
	WriteBit(C, D, red >> 1)
	WriteBit(C, D, red)

	WriteBit(C, D, green >> 9)
	WriteBit(C, D, green >> 8)
	WriteBit(C, D, green >> 7)
	WriteBit(C, D, green >> 6)
	WriteBit(C, D, green >> 5)
	WriteBit(C, D, green >> 4)
	WriteBit(C, D, green >> 3)
	WriteBit(C, D, green >> 2)
	WriteBit(C, D, green >> 1)
	WriteBit(C, D, green)

	time.Sleep(1 * time.Microsecond)
	L.WriteString("1")
	L.Sync()
	time.Sleep(1 * time.Microsecond)
	L.WriteString("0")
	L.Sync()
	time.Sleep(1 * time.Microsecond)
}
	

func WriteBit(C *os.File, D *os.File, bit int) {
	bit &= 0x1
	if bit == 0 {
		D.WriteString("0")
	} else {
		D.WriteString("1")
	}
	C.WriteString("1")
	D.Sync()
	C.Sync()
	time.Sleep(1 * time.Microsecond)
	C.WriteString("0")
	C.Sync()
	time.Sleep(1 * time.Microsecond)
}

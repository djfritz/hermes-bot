/*

a simple test to output a 1KHz signal on gpio pin 28

*/

package main

import (
	"os"
	"time"
	"fmt"
)

func main() {
	// export the gpio pin
	f, err := os.OpenFile("/sys/class/gpio/export", os.O_WRONLY | os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}	
	f.WriteString("59") // internal gpio 59 is gpio pin 28 on the header
	f.Close()

	// set it as an output
	f, err = os.OpenFile("/sys/class/gpio/gpio59/direction", os.O_WRONLY | os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	f.WriteString("out")
	f.Close()

	// spin on it
	f, err = os.OpenFile("/sys/class/gpio/gpio59/value", os.O_RDWR, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}

	v := "0"
	for {
		f.WriteString(v)
		f.Sync()
		if v == "0" {
			v = "1"
		} else {
			v = "0"
		}
		time.Sleep(10 * time.Millisecond)
	}
}

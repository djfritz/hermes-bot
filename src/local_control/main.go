package main

import (
	"fmt"
	"goserial"
	"os"
	"os/signal"
	"strconv"
)

func main() {
	s, err := goserial.Open("/dev/ttyUSB0", 9600)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)
	go func() {
		<-sig
		s.Close()
		os.Exit(0)
	}()

	for {
		var input string
		var input2 string
		fmt.Scanln(&input, &input2)

		a, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if a == 0 {
			s.Write([]byte{byte(0)})
			continue
		}

		b, err := strconv.Atoi(input2)
		if err != nil {
			fmt.Println(err)
			continue
		}

		data := []byte{byte(a), byte(b)}
		s.Write(data)
	}
}

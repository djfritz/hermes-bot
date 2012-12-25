package goserial

// #include <termios.h>
// #include <unistd.h>
import "C"

import (
	"syscall"
	"fmt"
)

type Serial struct {
	fd int
}

func Open(path string, baud int) (Serial, error) {
	f, err := syscall.Open(path, syscall.O_NOCTTY | syscall.O_RDWR | syscall.O_NDELAY, 0666)
	ret := Serial{f}
	if err != nil {
		return ret, err
	}

	r, _, e := syscall.Syscall(syscall.SYS_FCNTL, uintptr(f), uintptr(syscall.F_SETFL), uintptr(0))
	if r != 0 || e != 0 {
		syscall.Close(f)
		return ret, fmt.Errorf("disable blocking: %v %v", r, e)
	}

	var st C.struct_termios
	_, err = C.tcgetattr(C.int(f), &st)
	if err != nil {
		syscall.Close(f)
		return ret, err
	}

	var speed C.speed_t
	switch baud {
	case 9600:
		speed = C.B9600
	default:
		syscall.Close(f)
		return ret, fmt.Errorf("bad baud rate: %v", baud)
	}

	_, err = C.cfsetispeed(&st, speed)
	if err != nil {
		syscall.Close(f)
		return ret, err
	}

	_, err = C.cfsetospeed(&st, speed)
	if err != nil {
		syscall.Close(f)
		return ret, err
	}

	st.c_cflag |= (C.CLOCAL | C.CREAD)
	st.c_lflag &= ^C.tcflag_t(C.ICANON | C.ECHO | C.ECHOE | C.ISIG)
	st.c_oflag &= ^C.tcflag_t(C.OPOST)

	_, err = C.tcsetattr(C.int(f), C.TCSANOW, &st)
	if err != nil {
		syscall.Close(f)
		return ret, err
	}

	return ret, nil
}

func (s *Serial) Read(p []byte) (n int, err error) {
	return syscall.Read(s.fd, p)
}

func (s *Serial) Write(p []byte) (n int, err error) {
	return syscall.Write(s.fd, p)
}

func (s *Serial) Close() {
	syscall.Close(s.fd)
}

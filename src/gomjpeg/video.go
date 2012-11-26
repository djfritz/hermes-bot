package main

// #cgo LDFLAGS: -lv4l2
// #include <stdlib.h>
// #include "video.h"
import "C"

import (
	"unsafe"
)


func video_start(path string) {
	p := C.CString(path)
	C.video_start(p)
	C.free(unsafe.Pointer(p))
}

func grab_image() []byte {
	p := C.grab_image()
	//fwrite(buffers[buf.index].start, buf.bytesused, 1, fout);
	ret := C.GoBytes(p.data, p.length)
	return ret
}

func video_stop() {
	C.video_stop()
}

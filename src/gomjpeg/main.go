package main

import (
	"net/http"
	"fmt"
	"mime/multipart"
	"net/textproto"
	"time"
	"os/signal"
	"os"
	"flag"
)

var (
	images = make(chan []byte)

	f_dev = flag.String("input", "/dev/video0", "input device")
)

func main() {
	flag.Parse()

	sig := make(chan os.Signal, 1024)
	signal.Notify(sig, os.Interrupt)
	go func() {
		<-sig
		teardown()
	}()

	video_start(*f_dev)

	go imageFeeder()

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
	teardown()
}

func imageFeeder() {
	for {
		images <- grab_image()
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	mh := make(textproto.MIMEHeader)
	mh.Set("Content-Type", "image/jpeg")

	m := multipart.NewWriter(w)

	h := w.Header()
	boundary := m.Boundary()
	h.Set("Content-type", fmt.Sprintf("multipart/x-mixed-replace; boundary=%s", boundary))

	// for each jpeg image
	for {
		image := <-images

		mh.Set("Content-length", fmt.Sprintf("%d", len(image)))
		fm, err := m.CreatePart(mh)
		if err != nil {
			fmt.Println(err)
			return
		}
		_, err = fm.Write(image)
		if err != nil {
			break
		}
		time.Sleep(33 * time.Millisecond)
	}
}

func teardown() {
	video_stop()
	os.Exit(0)
}

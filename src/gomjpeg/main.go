package main

import (
	"net/http"
	"fmt"
	"mime/multipart"
//	"bytes"
//	"io"
	"net/textproto"
//	"os"
	"time"
)

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	video_start("/dev/video0")

	mh := make(textproto.MIMEHeader)
	mh.Set("Content-Type", "image/jpeg")

	m := multipart.NewWriter(w)

	h := w.Header()
	boundary := m.Boundary()
	h.Set("Content-type", fmt.Sprintf("multipart/x-mixed-replace; boundary=%s", boundary))

	// for each jpeg image
	for {
//		filename := fmt.Sprintf("out%03d.jpg", i)
//		j, err := os.Open(filename)
//		if err != nil {
//			fmt.Println(err)
//			return
//		}
//		buf := new(bytes.Buffer)
//		io.Copy(buf, j)
		image := grab_image()

		mh.Set("Content-length", fmt.Sprintf("%d", len(image)))
		fm, err := m.CreatePart(mh)
		if err != nil {
			fmt.Println(err)
			return
		}
		fm.Write(image)
		time.Sleep(50 * time.Millisecond)
	}

	video_stop()

}

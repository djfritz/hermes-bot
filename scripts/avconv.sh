avconv -f video4linux2 -input_format mjpeg -framerate 15 -video_size 640x480 -i /dev/video1 -vcodec copy -f rawvideo udp://localhost:1234

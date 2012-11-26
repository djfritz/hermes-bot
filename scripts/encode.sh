avconv -f video4linux2 -input_format mjpeg -framerate 10 -video_size 640x480 -i /dev/video0 -vcodec copy -f rawvideo udp://localhost:1234

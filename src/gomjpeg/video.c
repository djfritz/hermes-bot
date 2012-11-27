#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <fcntl.h>
#include <errno.h>
#include <sys/ioctl.h>
#include <sys/types.h>
#include <sys/time.h>
#include <sys/mman.h>
#include <linux/videodev2.h>
#include <libv4l2.h>
#include "video.h"

struct v4l2_format              fmt;
struct v4l2_buffer              buf;
struct v4l2_requestbuffers      req;
enum v4l2_buf_type              type;
fd_set                          fds;
struct timeval                  tv;
int                             r, fd = -1;
unsigned int                    i, n_buffers;
char                            out_name[256];
FILE                            *fout;
struct buffer                   *buffers;

struct buffer {
        void   *start;
        size_t length;
};

static void xioctl(int fh, int request, void *arg) {
        int r;

        do {
                r = v4l2_ioctl(fh, request, arg);
        } while (r == -1 && ((errno == EINTR) || (errno == EAGAIN)));

        if (r == -1) {
                fprintf(stderr, "error %d, %s\n", errno, strerror(errno));
                exit(EXIT_FAILURE);
        }
}


int video_start(char *dev_name) {
        fd = v4l2_open(dev_name, O_RDWR | O_NONBLOCK, 0);
        if (fd < 0) {
		return 1;
        }

	memset(&fmt, 0, sizeof(fmt));
        fmt.type = V4L2_BUF_TYPE_VIDEO_CAPTURE;
        fmt.fmt.pix.width       = 640;
        fmt.fmt.pix.height      = 480;
        fmt.fmt.pix.pixelformat = V4L2_PIX_FMT_MJPEG;
        fmt.fmt.pix.field       = V4L2_FIELD_ANY;
        xioctl(fd, VIDIOC_S_FMT, &fmt);
        if (fmt.fmt.pix.pixelformat != V4L2_PIX_FMT_MJPEG) {
		v4l2_close(fd);
		return 2;
        }
        if ((fmt.fmt.pix.width != 640) || (fmt.fmt.pix.height != 480)) {
		v4l2_close(fd);
		return 3;
	}

        memset(&req, 0, sizeof(req));
        req.count = 2;
        req.type = V4L2_BUF_TYPE_VIDEO_CAPTURE;
        req.memory = V4L2_MEMORY_MMAP;
        xioctl(fd, VIDIOC_REQBUFS, &req);
        buffers = calloc(req.count, sizeof(*buffers));
        for (n_buffers = 0; n_buffers < req.count; ++n_buffers) {
                memset(&buf, 0, sizeof(buf));

                buf.type        = V4L2_BUF_TYPE_VIDEO_CAPTURE;
                buf.memory      = V4L2_MEMORY_MMAP;
                buf.index       = n_buffers;

                xioctl(fd, VIDIOC_QUERYBUF, &buf);

                buffers[n_buffers].length = buf.length;
                buffers[n_buffers].start = v4l2_mmap(NULL, buf.length,
                              PROT_READ | PROT_WRITE, MAP_SHARED,
                              fd, buf.m.offset);

                if (MAP_FAILED == buffers[n_buffers].start) {
			return 1;
                }
        }

        for (i = 0; i < n_buffers; ++i) {
                memset(&buf, 0, sizeof(buf));
                buf.type = V4L2_BUF_TYPE_VIDEO_CAPTURE;
                buf.memory = V4L2_MEMORY_MMAP;
                buf.index = i;
                xioctl(fd, VIDIOC_QBUF, &buf);
        }
        type = V4L2_BUF_TYPE_VIDEO_CAPTURE;

        xioctl(fd, VIDIOC_STREAMON, &type);
	return 0;
}

struct image grab_image(void) {
	struct image ret;
	ret.data = 0;
	do {
		FD_ZERO(&fds);
		FD_SET(fd, &fds);

		/* Timeout. */
		tv.tv_sec = 2;
		tv.tv_usec = 0;

		r = select(fd + 1, &fds, NULL, NULL, &tv);
	} while ((r == -1 && (errno = EINTR)));
	if (r == -1) {
		return ret;
	}

	memset(&buf, 0, sizeof(buf));
	buf.type = V4L2_BUF_TYPE_VIDEO_CAPTURE;
	buf.memory = V4L2_MEMORY_MMAP;
	xioctl(fd, VIDIOC_DQBUF, &buf);

	ret.data = buffers[buf.index].start;
	ret.length = buf.bytesused;

	xioctl(fd, VIDIOC_QBUF, &buf);

	return ret;
}

void video_stop() {
        type = V4L2_BUF_TYPE_VIDEO_CAPTURE;
        xioctl(fd, VIDIOC_STREAMOFF, &type);
        for (i = 0; i < n_buffers; ++i) {
                v4l2_munmap(buffers[i].start, buffers[i].length);
	}
        v4l2_close(fd);
}

#ifndef VIDEO
#define VIDEO

struct image {
	void *data;
	int length;
};

int video_start(char *dev_name);
void video_stop(void);
struct image grab_image(void);


#endif

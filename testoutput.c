#include <graphics.h>
#include <stdlib.h>
#include <stdio.h>

int main()
{
	initwindow(800, 800);

	line(400, 400, 394, 429);
	line(393, 418, 378, 443);
	line(380, 433, 357, 452);
	line(362, 443, 333, 453);
	line(342, 446, 312, 446);
	line(322, 442, 293, 431);
	line(304, 432, 281, 412);
	line(291, 416, 276, 390);
	line(284, 397, 278, 367);
	line(284, 377, 289, 347);
	line(290, 358, 305, 332);

	getch();
	closegraph();
	return 0;
}

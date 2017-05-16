#include <graphics.h>
#include <stdlib.h>
#include <stdio.h>

int main()
{
	int gdriver = DETECT, gmode, errorcode;
	initgraph(&gdriver, &gmode, NULL);

	line(400, 600, 350, 686);

	getch();
	closegraph();
	return 0;
}

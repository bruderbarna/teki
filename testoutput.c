#include <graphics.h>
#include <stdlib.h>
#include <stdio.h>

int main()
{
	initwindow(800, 800);

	for (int i = 0; i < 10; i++) {
	line(400, 400, 394, 429);
	}

	getch();
	closegraph();
	return 0;
}

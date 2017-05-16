all:
	-rm testoutput.c
	go build
	./teki test.teki testoutput.c
	cat test.teki
	@echo "\n\n\ncompiles to:"
	cat testoutput.c

run:
	g++ testoutput.c -lgraph -o testgraphicoutput
	./testgraphicoutput

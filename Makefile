sync: clean diosteama
	scp diosteama fary.pandacrew.net:
diosteama:
	CGO_ENABLED=0 go build
clean:
	rm diosteama

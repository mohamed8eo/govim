BINARY=govim
FILE = test.txt


build:
	go build -o $(BINARY) .
run:
	./$(BINARY) $(ARG)


test: build
	./$(BINARY) $(FILE)

build:
	go build -o main

run:
	go build -o main main.go
	./main ${file}
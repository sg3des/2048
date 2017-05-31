
name=2048

run: build
	./$(name) ${a}

build: 
	go build -o $(name)



.PHONY: build run
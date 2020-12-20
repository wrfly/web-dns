.PHONY: config lib route test

NAME := "web-dns"

test:
	go test -cover -v .

build:
	go build -o $(NAME) .

dev: build
	./$(NAME) --debug
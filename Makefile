.PHONY: config lib route test

NAME := "web-dns"

test:
	go test `glide nv`

build:
	go build -o $(NAME) .

dev: build
	./$(NAME) -d
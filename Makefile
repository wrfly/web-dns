.PHONY: config lib route test

NAME := "web-dns"

test:
	go test -cover -v `glide nv`

build:
	go build -o $(NAME) .

dev: build
	./$(NAME) --debug
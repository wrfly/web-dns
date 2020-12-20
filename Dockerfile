FROM golang:alpine AS build
RUN apk add git
COPY . /build
RUN cd /build && \
    go build -o web-dns .

FROM alpine
COPY --from=build /build/web-dns /usr/local/bin/
CMD [ "web-dns" ]

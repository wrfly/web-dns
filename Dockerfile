FROM golang:alpine as build
ENV PKG /go/src/github.com/wrfly/web-dns
COPY . ${PKG}
RUN glide i && \
    make test && \
    make build && \
    mv ${PKG}/web-dns /

FROM alpine
COPY --from=build /web-dns /usr/local/bin/
CMD [ "web-dns" ]

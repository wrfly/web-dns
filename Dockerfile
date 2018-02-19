FROM wrfly/glide AS build
ENV PKG /go/src/github.com/wrfly/web-dns
COPY . ${PKG}
RUN cd ${PKG} && \
    glide i && \
    make test && \
    make build && \
    mv ${PKG}/web-dns /

FROM alpine
COPY --from=build /web-dns /usr/local/bin/
CMD [ "web-dns" ]

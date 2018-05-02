FROM chapsuk/golang-baseimage:1.10.1
ADD . /go/src/github.com/chapsuk/miga
WORKDIR /go/src/github.com/chapsuk/miga
RUN make build

FROM alpine:3.7
COPY --from=0 /go/src/github.com/chapsuk/miga/bin/miga /miga
CMD ["/miga"]

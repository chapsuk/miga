FROM chapsuk/golang-baseimage:1.10.2
ADD . /go/src/github.com/chapsuk/miga
WORKDIR /go/src/github.com/chapsuk/miga
ARG VERSION
ENV VERSION=${VERSION}
RUN make build

FROM alpine:3.7
COPY --from=0 /go/src/github.com/chapsuk/miga/bin/miga /miga
ENTRYPOINT ["/miga"]

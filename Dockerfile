FROM chapsuk/golang-baseimage:1.10.2
ADD . /go/src/github.com/chapsuk/miga
WORKDIR /go/src/github.com/chapsuk/miga
RUN make build

FROM scratch
COPY --from=0 /go/src/github.com/chapsuk/miga/bin/miga /miga
ENTRYPOINT [ "/miga" ]

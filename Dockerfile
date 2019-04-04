FROM chapsuk/golang:1.12.3
ADD . /go/src/miga
WORKDIR /go/src/miga
RUN make build

FROM scratch
COPY --from=0 /go/src/miga/bin/miga /miga
ENTRYPOINT [ "/miga" ]

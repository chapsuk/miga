FROM golang:1.16.3-alpine3.13
RUN apk add make
ADD . /go/src/miga
WORKDIR /go/src/miga
RUN make build

FROM scratch
COPY --from=0 /go/src/miga/bin/miga /miga
ENTRYPOINT [ "/miga" ]

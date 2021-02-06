FROM golang:1.15.7-alpine3.13
RUN apk add make
ADD . /go/src/miga
WORKDIR /go/src/miga
RUN make build

FROM scratch
COPY --from=0 /go/src/miga/bin/miga /miga
ENTRYPOINT [ "/miga" ]

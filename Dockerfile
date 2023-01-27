FROM golang:1.19.5-alpine3.17
RUN apk add make
ADD . /go/src/miga
WORKDIR /go/src/miga
RUN make build

FROM scratch
COPY --from=0 /go/src/miga/bin/miga /miga
ENTRYPOINT [ "/miga" ]

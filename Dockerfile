FROM golang:1.21.4-alpine3.18 AS build
RUN apk add make tzdata
COPY . /go/src/miga
WORKDIR /go/src/miga
RUN make build

FROM scratch
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /go/src/miga/bin/miga /miga
ENTRYPOINT [ "/miga" ]

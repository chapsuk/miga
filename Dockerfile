FROM golang:1.23.0-alpine3.20 AS build
RUN apk add make tzdata
COPY . /go/src/miga
WORKDIR /go/src/miga
RUN make build

FROM scratch
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /go/src/miga/bin/miga /miga
ENTRYPOINT [ "/miga" ]

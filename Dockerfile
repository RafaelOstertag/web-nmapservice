FROM golang:alpine AS build

RUN apk add make git
COPY src /go/src/
COPY mocks /go/mocks/
COPY Makefile /go/
WORKDIR /go
RUN make nmapservice

FROM alpine:latest

RUN apk add nmap && rm -rf /var/cache/apk
COPY --from=build /go/bin/nmapservice /app/

WORKDIR /
ENTRYPOINT [ "/app/nmapservice" ]

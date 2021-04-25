FROM golang:alpine AS build

RUN apk add make git
COPY src /go/src/
COPY mocks /go/mocks/
WORKDIR /go/src
RUN make nmapservice

FROM alpine:latest

RUN apk add nmap && rm -rf /var/cache/apk
COPY --from=build /go/src/nmapservice /app/

WORKDIR /
ENTRYPOINT [ "/app/nmapservice" ]

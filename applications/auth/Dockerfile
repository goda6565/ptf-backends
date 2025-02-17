FROM golang:1.24.0-alpine3.21

RUN apk add --no-cache curl

WORKDIR /go/src/auth

COPY . .

RUN go mod download

RUN go build .

ENTRYPOINT ["./auth"]
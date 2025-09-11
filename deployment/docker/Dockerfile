FROM golang:1.23.2-alpine as builder

WORKDIR /home/builder

RUN apk add --no-cache make

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o bin/lgp-validator cmd/validator/main.go

FROM golang:1.22.7-alpine

RUN apk add --no-cache bash make git curl

RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.63.0
RUN go install github.com/vektra/mockery/v2@v2.46.0

COPY --from=builder /home/builder/bin/lgp-validator /usr/local/bin/
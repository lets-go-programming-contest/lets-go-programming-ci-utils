FROM golang:1.22.7-alpine

RUN apk add --no-cache bash make git task curl
RUN sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b /usr/bin v3.44.1 
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.63.0

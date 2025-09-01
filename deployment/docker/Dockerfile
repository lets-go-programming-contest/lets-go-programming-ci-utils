FROM golang:1.22.7-alpine

RUN apk add --no-cache bash make git
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.63.0
RUN golangci-lint --version

ENV TEST_DIR_COMMON=.ci_utils
ENV TEST_DIR_UTILS=.ci_common

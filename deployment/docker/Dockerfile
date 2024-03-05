FROM golang:1.21

WORKDIR /spbstu-2024-lgp-ci

RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
RUN golangci-lint --version

ENV TEST_DIR=/spbstu-2024-lgp-ci
ENV TEST_DIR_UTILS=/spbstu-2024-lgp-ci/utils
ENV TEST_DIR_COMMON=/spbstu-2024-lgp-ci/common
ENV TEST_MAKEFILE=/spbstu-2024-lgp-ci/utils/common.mk

COPY utils utils



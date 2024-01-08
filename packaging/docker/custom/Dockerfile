# syntax=docker/dockerfile:1

ARG BASE_IMAGE=alpine:3.18.3
ARG GO_IMAGE=golang:1.21.5-alpine3.18

ARG GO_SRC=go-builder
ARG JS_SRC=js-builder

FROM --platform=${JS_PLATFORM} ${JS_IMAGE} as js-builder

FROM ${GO_IMAGE} as go-builder
ARG GO_BUILD_TAGS="oss"
ARG WIRE_TAGS="oss"
ARG BINGO="true"

WORKDIR /tmp/grafana

COPY go.* ./
COPY .bingo .bingo

RUN go mod download

COPY Makefile build.go package.json ./
COPY pkg pkg

RUN make build WIRE_TAGS=${WIRE_TAGS}

EXPOSE 3000

ARG RUN_SH=./packaging/docker/run.sh
COPY ${RUN_SH} /run.sh

ENTRYPOINT [ "/run.sh" ]

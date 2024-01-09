ARG GO_IMAGE=golang:1.21.5-alpine3.18
ARG GO_SRC=go-builder
FROM ${GO_IMAGE} as builder

ARG GO_BUILD_TAGS="oss"
ARG WIRE_TAGS="oss"
ARG BINGO="true"

WORKDIR /build/dashboard-tsdata-bridge

COPY go.* ./
COPY .bingo .bingo
COPY Makefile main.go ./
COPY pkg pkg

# Install build dependencies
RUN if grep -i -q alpine /etc/issue; then \
      apk add --no-cache gcc g++ make git; \
    fi

RUN make gen

RUN go mod tidy \
    && go get -u -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-s -w' -o project

RUN make build

FROM scratch
COPY --from=builder /build/dashboard-tsdata-bridge /

ARG LOKI_URL
ENV LOKI_URL ${LOKI_URL}
ARG PROMETHEUS_URL
ENV PROMETHEUS_URL ${PROMETHEUS_URL}

EXPOSE 9090
CMD ["/project"]
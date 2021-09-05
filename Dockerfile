
ARG GO_VERSION=1.16.1
FROM golang:${GO_VERSION}-alpine AS builder
WORKDIR /build-dir

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build  \
    -o /currency_converter_job ./cmd/job

RUN CGO_ENABLED=0 go build \
    -o /currency_converter_web ./cmd/web

RUN CGO_ENABLED=0 go build \
    -o /cli ./cmd/cli

FROM alpine:latest
WORKDIR /app/

USER nobody

COPY --from=builder /currency_converter_web ./currency_converter_web
COPY --from=builder /currency_converter_job ./currency_converter_job
COPY --from=builder /cli ./build/cli

CMD  ["/app/currency_converter_web"]

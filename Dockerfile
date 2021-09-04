
ARG GO_VERSION=1.16.1
FROM golang:${GO_VERSION}-alpine AS builder
WORKDIR /build-dir
COPY . .

RUN CGO_ENABLED=0 go build  \
    -o /currency_convertor_jobs ./cmd/jobs

RUN CGO_ENABLED=0 go build \
    -o /currency_convertor_web ./cmd/web

RUN CGO_ENABLED=0 go build \
    -o /cli ./cmd/cli

FROM alpine:latest
WORKDIR /app/

USER nobody

COPY --from=builder /currency_convertor_web ./currency_convertor_web
COPY --from=builder /currency_convertor_jobs ./currency_convertor_jobs
COPY --from=builder /cli ./build/cli

CMD  ["/app/currency_convertor_web"]

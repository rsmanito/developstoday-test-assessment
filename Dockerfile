FROM golang:1.23.5-alpine AS builder
LABEL org.opencontainers.image.authors="rsmanito"

WORKDIR /build

COPY . .

RUN go mod download

RUN go build -o /build/sca-service ./cmd/main/main.go

FROM alpine:latest
WORKDIR /

COPY --from=builder /build/sca-service /sca-service

ENV PORT=3000

EXPOSE $PORT

ENTRYPOINT ["/sca-service"]

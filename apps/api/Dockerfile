# syntax=docker/dockerfile:1

FROM golang:1.26.1-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/pulsefeed-api ./cmd/feed
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/pulsefeed-worker ./cmd/worker

FROM alpine:3.22

WORKDIR /app

RUN apk add --no-cache ca-certificates ffmpeg tzdata

COPY --from=builder /out/pulsefeed-api ./pulsefeed-api
COPY --from=builder /out/pulsefeed-worker ./pulsefeed-worker
COPY configs ./configs

RUN mkdir -p uploads

EXPOSE 8080

CMD ["./pulsefeed-api"]

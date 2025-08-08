# build
FROM golang:1.24 AS builder

ENV NAME "usuf-bot"

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o ./bin/${NAME} ./cmd/main.go

# run
FROM debian:stable-slim

RUN apt-get update && \
    apt-get install -y \
    ffmpeg \
    python3 \
    python3-pip \
    && rm -rf /var/lib/apt/lists/*

# Устанавливаем yt-dlp через pip (с разрешением на установку в системный Python)
RUN pip3 install --no-cache-dir --break-system-packages yt-dlp

COPY --from=builder /app/bin/usuf-bot /usr/local/bin/app

ENTRYPOINT ["app"]

FROM golang:1.21.11 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o server ./cmd/server/main.go

FROM python:3.10-slim

RUN apt-get update && apt-get install -y \
    ffmpeg \
    curl \
    && rm -rf /var/lib/apt/lists/*

RUN pip install vosk

WORKDIR /app
COPY --from=builder /app/server /app/server
COPY web/ /app/web/
COPY transcribe.py /app/transcribe.py

EXPOSE 8080

CMD ["./server"]

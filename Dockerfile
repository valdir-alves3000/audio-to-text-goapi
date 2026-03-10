    FROM golang:1.21-alpine AS builder

    WORKDIR /app
    
    RUN apk add --no-cache git
    
    COPY go.mod go.sum ./
    RUN go mod download
    
    COPY . .
    
    RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
        go build -ldflags="-s -w" -o server ./cmd/server/main.go
    
    
    FROM python:3.10-slim
    
    WORKDIR /app
    
    RUN apt-get update && apt-get install -y \
        ffmpeg \
        curl \
        && rm -rf /var/lib/apt/lists/*
    
    RUN pip install --no-cache-dir --upgrade pip setuptools wheel
    
    RUN pip install --no-cache-dir \
        torch \
        --index-url https://download.pytorch.org/whl/cpu
    
    RUN pip install --no-cache-dir openai-whisper
    
    RUN python -c "import whisper; whisper.load_model('small')"
    
    COPY --from=builder /app/server /app/server
    COPY transcribe_worker.py /app/transcribe_worker.py
    COPY web/ /app/web/
    
    EXPOSE 8080
    
    CMD ["./server"]
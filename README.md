# Audio to Text GoAPI

Aplicação fullstack para transcrição de áudio e vídeo em tempo real,  
com backend em **Go + Python (OpenAI Whisper)** e frontend para upload e exibição progressiva do texto.

Permite upload de arquivos de áudio e vídeo (WAV, MP3, M4A, MP4, WEBM, AVI, MOV)  
e retorna a transcrição via **Server-Sent Events (SSE)**, exibindo o texto enquanto o áudio é processado.

---

## 🚀 Arquitetura

- Go → API HTTP + processamento de arquivos
- Python → Worker persistente com Whisper
- Whisper → Modelo carregado uma única vez
- SSE → Streaming de transcrição em tempo real

Fluxo:

```

Client
↓
Go API (Gin)
↓
Whisper Worker (Python persistente)
↓
Modelo Whisper (carregado 1x)
↓
Streaming de segmentos via SSE

````

---

## 📦 Requisitos

- Go 1.20+
- Python 3.8+
- pip
- ffmpeg (obrigatório para conversão de áudio)
- Linux/macOS ou WSL (Windows)

Instalar ffmpeg:

```bash
sudo apt install ffmpeg
````

---

## ⚙️ Configuração Local

### 1️⃣ Clone o repositório

```bash
git clone https://github.com/valdir-alves3000/audio-to-text-goapi.git
cd audio-to-text-goapi
```

---

### 2️⃣ Criar ambiente virtual Python

```bash
python3 -m venv .venv
source .venv/bin/activate
```

---

### 3️⃣ Instalar dependências Python

```bash
pip install openai-whisper
pip install torch
```

---

### 4️⃣ Executar servidor

```bash
go run cmd/server/main.go
```

Servidor:

```
http://localhost:8080
```

---

## 🐳 Execução com Docker

Pré-requisitos:

* Docker
* Docker Compose

Rodar:

```bash
docker compose up --build
```

---

## 🎯 Como usar

1. Acesse `http://localhost:8080`
2. Faça upload de um arquivo de áudio ou vídeo
3. Clique em **Transcrever**
4. Veja o texto sendo exibido progressivamente, letra por letra

---

## 🔥 Diferenciais Técnicos

* Modelo Whisper carregado apenas uma vez
* Worker Python persistente
* Streaming por segmento
* Cancelamento via context (cliente desconecta → processamento para)
* Conversão automática para WAV
* Suporte a múltiplos formatos
* Arquitetura desacoplada e pronta para escalar

---

## 🧠 Tecnologias

* Go (Gin)
* Python
* OpenAI Whisper
* ffmpeg
* Server-Sent Events (SSE)
* HTML / CSS / JavaScript

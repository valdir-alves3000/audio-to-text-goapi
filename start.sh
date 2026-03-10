#!/usr/bin/env bash

set -e 

VENV=".venv"

if [ ! -d "$VENV" ]; then
  echo "Criando virtualenv..."
  python3 -m venv "$VENV"
fi

source "$VENV/bin/activate"

echo "Atualizando pip..."
python -m pip install --upgrade pip

echo "Instalando dependencias..."

python -m pip install torch torchvision torchaudio \
  --index-url https://download.pytorch.org/whl/cpu

python -m pip install openai-whisper ffmpeg-python

echo "Verificando ffmpeg..."
if ! command -v ffmpeg &> /dev/null; then
  echo "ERRO: ffmpeg nao encontrado."
  echo "Instale com:"
  echo "  sudo apt install ffmpeg"
  exit 1
fi

echo "Iniciando servidor Go..."
go run cmd/server/main.go
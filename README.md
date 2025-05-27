# Audio to Text GoAPI

Aplicação fullstack para transcrição de áudio e vídeo, com backend em Go e Python (Vosk) e frontend para upload e exibição dos resultados.  
Permite o upload de arquivos de áudio e vídeo (WAV, MP3, M4A, MP4) e retorna a transcrição em texto.

---

## Requisitos

- [Go](https://golang.org/dl/) (versão 1.20+ recomendada)  
- [Python 3.8+](https://www.python.org/downloads/)  
- [pip](https://pip.pypa.io/en/stable/installation/) (gerenciador de pacotes Python)  
- Sistema Linux/macOS ou WSL para ativação do ambiente virtual

---

## Configuração e execução local

1. Clone o repositório:

```bash
git clone https://github.com/valdir-alves3000/audio-to-text-goapi.git
cd audio-to-text-goapi
````

2. Crie e ative o ambiente virtual Python:

```bash
python3 -m venv .venv
source .venv/bin/activate
```

3. Instale a dependência Python Vosk para reconhecimento de voz:

```bash
pip install vosk
```

4. Execute o servidor Go:

```bash
go run cmd/server/main.go
```

## Como usar

* Acesse `http://localhost:8080` no navegador
* Faça upload de arquivos de áudio ou vídeo suportados (WAV, MP3, M4A, MP4)
* Clique em **Transcrever** para obter o texto transcrito
* Veja o texto transcrito aparecer logo abaixo do player de mídia(em tempo real)

---

## Tecnologias

* Go (backend, servidor HTTP, manipulação de arquivos)
* Python + Vosk (transcrição de áudio)
* HTML/CSS/JavaScript (frontend)

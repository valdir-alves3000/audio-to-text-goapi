#!/bin/bash
source .venv/bin/activate
pip install vosk
go run cmd/server/main.go

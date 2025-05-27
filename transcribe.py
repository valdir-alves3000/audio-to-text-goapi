#!/usr/bin/env python3
import sys
import os
import json
import wave
import argparse
from vosk import Model, KaldiRecognizer

def download_model(lang):
    """Download a model if not already present."""
    from urllib.request import urlretrieve
    import zipfile
    
    model_mapping = {
        "en-us": {
            "url": "https://alphacephei.com/vosk/models/vosk-model-small-en-us-0.15.zip",
            "path": "models/vosk-model-small-en-us-0.15"
        },
        "pt": {
            "url": "https://alphacephei.com/vosk/models/vosk-model-small-pt-0.3.zip",
            "path": "models/vosk-model-small-pt-0.3"
        }
    }
    
    if lang not in model_mapping:
        print(f"Unsupported language: {lang}", file=sys.stderr)
        return None
    
    model_info = model_mapping[lang]
    model_path = model_info["path"]
    
    # Check if model already exists
    if os.path.exists(model_path):
        return model_path
    
    # Create models directory if it doesn't exist
    os.makedirs("models", exist_ok=True)
    
    # Download and extract the model
    zip_path = f"models/{os.path.basename(model_info['url'])}"
    print(f"Downloading model for {lang}...", file=sys.stderr)
    
    urlretrieve(model_info["url"], zip_path)
    
    print("Extracting model...", file=sys.stderr)
    with zipfile.ZipFile(zip_path, 'r') as zip_ref:
        zip_ref.extractall("models")
    
    # Clean up the zip file
    os.remove(zip_path)
    
    return model_path

def transcribe_audio(audio_file, lang="pt"):    
    """Transcribe audio file using Vosk."""
    model_path = download_model(lang)
    if not model_path:
        return ""
    
    model = Model(model_path)
    
    wf = wave.open(audio_file, "rb")
    if wf.getnchannels() != 1 or wf.getsampwidth() != 2 or wf.getcomptype() != "NONE":
        print("Audio file must be WAV format mono PCM.", file=sys.stderr)
        return ""
    
    rec = KaldiRecognizer(model, wf.getframerate())
    rec.SetWords(True)
    
    # Transcrição em tempo real
    while True:
        data = wf.readframes(4000)
        if len(data) == 0:
            break
        if rec.AcceptWaveform(data):
            part_result = json.loads(rec.Result())
            print(part_result.get('text', ''))  # Envia resultado parcial

    # Resultado final
    part_result = json.loads(rec.FinalResult())
    print(part_result.get('text', ''))  # Envia resultado final

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description='Transcribe audio file using Vosk')
    parser.add_argument('-i', '--input', required=True, help='Input audio file (WAV format)')
    parser.add_argument('-m', '--model', default="en-us", help='Language model to use (en-us, pt)')
    
    args = parser.parse_args()
    
    transcribe_audio(args.input, args.model)

#!/usr/bin/env python3
import sys
import json
import whisper

# Load model ONCE
MODEL = whisper.load_model("small")

print("READY", flush=True)

for line in sys.stdin:
    line = line.strip()
    if not line:
        continue

    try:
        payload = json.loads(line)

        audio_file = payload["file"]
        language = payload.get("lang", "pt")

        result = MODEL.transcribe(
            audio_file,
            language=language,
            fp16=False,
            verbose=False,
            beam_size=1,
            best_of=1
        )

        for segment in result["segments"]:
            print(segment["text"].strip(), flush=True)

        print("__END__", flush=True)

    except Exception as e:
        print(f"ERROR: {str(e)}", flush=True)
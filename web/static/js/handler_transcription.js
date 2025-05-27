import { VIEW } from './view.js';
import { getSelectedFile } from './file_handler.js';
import { showError } from './error.js';

let abortController = null;

export async function handleTranscription() {
    if (abortController) {
        abortController.abort();
    }

    const controller = new AbortController();
    abortController = controller;

    const file = getSelectedFile();

    if (!file) {
        showError('Por favor, selecione um arquivo primeiro.');
        return;
    }

    toggleLoading(true);
    clearTranscriptionResult();
    
    const formData = new FormData();
    formData.append('audio', file);
    formData.append('lang', VIEW.languageSelect.value);

    try {
        const response = await fetch('/api/transcribe', {
            method: 'POST',
            headers: {
                'Accept': 'text/event-stream'
            },
            body: formData,
            signal: controller.signal
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        if (!response.body) {
            throw new Error('ReadableStream not supported in this browser');
        }

        const reader = response.body.getReader();
        const decoder = new TextDecoder();
        let fullTranscript = '';

        while (true) {
            const { done, value } = await reader.read();
            if (done) break;

            const textChunk = decoder.decode(value, { stream: true });
            const lines = textChunk.split('\n');

            for (const line of lines) {
                if (line.startsWith('data: ')) {
                    const data = line.substring(6).trim();
                    if (data) {
                        fullTranscript += ` ${data}`;
                        updateTranscription(fullTranscript);
                    }
                } else if (line.includes('event: end')) {
                    toggleLoading(false);
                    return;
                }
            }
        }
    } catch (error) {
        showError('Falha na Transcrição do arquivo: ' + error.message);
        console.error('Transcription error:', error);
    } finally {
        toggleLoading(false);
    }
}

function toggleLoading(show) {
    VIEW.loading.style.display = show ? 'block' : 'none';
    VIEW.transcribeBtn.disabled = show;
}

function clearTranscriptionResult() {
    VIEW.transcriptionResult.style.display = 'none';
    VIEW.resultText.textContent = '';
}

function updateTranscription(text) {
    VIEW.resultText.textContent = text;
    VIEW.transcriptionResult.style.display = 'block';
    VIEW.transcriptionResult.scrollTop = VIEW.transcriptionResult.scrollHeight;
}


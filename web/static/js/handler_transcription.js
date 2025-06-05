import { VIEW } from './view.js?v=timestamp';
import { getSelectedFile } from './file_handler.js?v=timestamp';
import { showError } from './error.js?v=timestamp';

let abortController = null;

export async function handleTranscription() {
    if (abortController) {
        abortController.abort();
    }

    const controller = new AbortController();
    abortController = controller;

    const file = getSelectedFile();
    const formData = new FormData();
    formData.append('audio', file);
    formData.append('lang', VIEW.languageSelect.value);

    if (!file) {
        showError('Por favor, selecione um arquivo primeiro.');
        return;
    }

    toggleLoading(true);
    clearTranscriptionResult();

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
            let errorMsg = `HTTP error! status: ${response.status}`;
            try {
                const errorData = await response.json();
                if (errorData?.error) {
                    errorMsg = errorData.error; 
                }
            } catch (jsonErr) {
                console.error("Failed to parse error response:", jsonErr);
            }
        
            showError(errorMsg);
            return; 
        }
        

        if (!response.body) {
            showError('ReadableStream not supported in this browser');
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
                    break;
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
    setTimeout(() => {
        VIEW.resultText.scrollIntoView({ behavior: 'smooth', block: 'start' });
    }, 0);
}

window.addEventListener('beforeunload', () => {
    if (abortController) {
        abortController.abort();
    }
}); 
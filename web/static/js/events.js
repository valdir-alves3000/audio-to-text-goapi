import { VIEW } from './view.js';
import { processFile } from './file_handler.js';
import { handleTranscription } from './handler_transcription.js';

function preventDefaults(e) {
    e.preventDefault();
    e.stopPropagation();
}

function highlight() {
    VIEW.uploadArea.classList.add('drag-over');
}

function unhighlight() {
    VIEW.uploadArea.classList.remove('drag-over');
}

export function setupDragAndDrop() {
    ['dragenter', 'dragover', 'dragleave', 'drop'].forEach(event =>
        VIEW.uploadArea.addEventListener(event, preventDefaults, false)
    );

    ['dragenter', 'dragover'].forEach(event =>
        VIEW.uploadArea.addEventListener(event, highlight, false)
    );

    ['dragleave', 'drop'].forEach(event =>
        VIEW.uploadArea.addEventListener(event, unhighlight, false)
    );

    VIEW.uploadArea.addEventListener('drop', event => {
        const files = event.dataTransfer.files;
        if (files.length > 0) processFile(files[0]);
    }, false);
}

export function setupFileInput() {
    VIEW.fileInput.addEventListener('change', event => {
        if (event.target.files.length > 0) {
            processFile(event.target.files[0]);
        }
    });
}

export function setupTranscribeButton() {
    VIEW.transcribeBtn.addEventListener('click', handleTranscription);
}

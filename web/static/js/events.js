import { VIEW } from './view.js?v=timestamp';
import { processFile } from './file_handler.js?v=timestamp';
import { handleTranscription } from './handler_transcription.js?v=timestamp';

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

            VIEW.inputGroup.classList.add('active');
            VIEW.uploadText.classList.add('hidden');
            VIEW.fileInputLabel.classList.remove('bg-gradient-to-r', 'from-primary', 'to-secondary');
            VIEW.fileInputLabel.classList.add('bg-transparent', 'text-secondary', 'underline', 'shadow-none', 'hover:text-primary');
            VIEW.transcribeBtn.classList.remove('hidden');
            VIEW.transcribeBtn.classList.add('flex');
        }
    });
}

export function setupTranscribeButton() {
    VIEW.transcribeBtn.addEventListener('click', handleTranscription);
}

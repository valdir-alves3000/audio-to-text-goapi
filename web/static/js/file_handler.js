import { VIEW } from './view.js';
import { showError } from './error.js';

const SUPPORTED_FORMATS = ['.wav', '.mp3', '.m4a', '.mp4'];
let selectedFile = null;

export function getSelectedFile() {
    return selectedFile;
}

export function processFile(file) {
    const ext = '.' + file.name.split('.').pop().toLowerCase();

    if (!SUPPORTED_FORMATS.includes(ext)) {
        showError('Formato de arquivo não suportado. Use WAV, MP3, M4A ou MP4.');
        return;
    }

    selectedFile = file;

    resetUI();

    const fileURL = URL.createObjectURL(file);

    VIEW.mediaContainer.classList.remove("hidden")
    VIEW.uploadIcon.style.display = 'none';
    VIEW.inputGroup.classList.add("active")
    VIEW.fileInputLabel.textContent = "Alterar Arquivo";

    if (ext === '.mp4') {        
        VIEW.video.src = fileURL;
        VIEW.video.style.display = 'block';
    } else {
        VIEW.audio.src = fileURL;
        VIEW.audio.style.display = 'block';
    }
    
}

function resetUI() {
    VIEW.uploadIcon.style.display = 'block';    
    VIEW.video.style.display = 'none';
    VIEW.audio.style.display = 'none';    
    VIEW.inputGroup.classList.remove("active")
}
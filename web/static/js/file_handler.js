import { VIEW } from './view.js?v=timestamp';
import { showError } from './error.js?v=timestamp';

const SUPPORTED_FORMATS = ['.aac', '.flac', '.wav', '.mp3', '.m4a', '.opus', '.mp4', '.ogg', '.webm'];
const FORMATS_VIDEO = ['.mp4', '.webm']

let selectedFile = null;

export function getSelectedFile() {
    return selectedFile;
}

export function processFile(file) {
    const ext = '.' + file.name.split('.').pop().toLowerCase();

    if (!SUPPORTED_FORMATS.includes(ext)) {
        showError(`Formato de arquivo não suportado. Use WAV, MP3, M4A, MP4, WEBM, or OGG.`);
        return;
    }

    selectedFile = file;

    resetUI();

    const fileURL = URL.createObjectURL(file);

    VIEW.mediaContainer.classList.remove('hidden')
    VIEW.uploadIcon.style.display = 'none';
    VIEW.inputGroup.classList.remove('gap-4')
    VIEW.fileInputLabel.textContent = 'Alterar Arquivo';
    setTimeout(() => {
        VIEW.uploadArea.scrollIntoView({ behavior: 'smooth', block: 'start' });
    }, 0);

    if (FORMATS_VIDEO.includes(ext)) {
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
    VIEW.inputGroup.classList.remove('active')
}
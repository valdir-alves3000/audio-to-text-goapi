import { setupDragAndDrop, setupFileInput, setupTranscribeButton } from './events.js';
import { setupCopyButton } from './copy.js';

function app() {
    setupDragAndDrop();
    setupFileInput();
    setupTranscribeButton();
    setupCopyButton();
}

document.addEventListener('DOMContentLoaded', app);

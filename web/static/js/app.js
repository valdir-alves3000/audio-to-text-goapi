import { setupDragAndDrop, setupFileInput, setupTranscribeButton } from './events.js';
import { setupCopyButton } from './copy.js';
import { listenTranscript } from './listen_transcript.js';

function app() {
    setupDragAndDrop();
    setupFileInput();
    setupTranscribeButton();
    setupCopyButton();   
    listenTranscript();
}

document.addEventListener('DOMContentLoaded', app);

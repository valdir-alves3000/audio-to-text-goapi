import { VIEW } from './view.js';

export function showError(message) {
    VIEW.errorMessage.textContent = message;
    VIEW.errorModal.style.display = 'block';
    VIEW.errorClose.onclick = () => {
        VIEW.errorModal.style.display = 'none';
    };

    window.onclick = (event) => {
        if (event.target === VIEW.errorModal) {
            VIEW.errorModal.style.display = 'none';
        }
    };
}

export function hideError() {
    VIEW.errorModal.style.display = 'none';
}
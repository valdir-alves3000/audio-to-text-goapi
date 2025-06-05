import { showError } from "./error.js?v=timestamp";
import { VIEW } from "./view.js?v=timestamp";

let isPlaying = false;
let utterance = null;

export function listenTranscript() {
    VIEW.playBtn.addEventListener('click', togglePlay);
    VIEW.stopBtn.addEventListener('click', stop);
}

function play() {
    const text = VIEW.resultText.textContent.trim();
    const langSelect = VIEW.languageSelect.value;
    const lang = langSelect === 'pt' ? 'pt-BR' : 'en-US';

    if (!text) {
        showError('Nenhum texto disponível para leitura');
        return;
    }

    utterance = new SpeechSynthesisUtterance(text);
    utterance.lang = lang;

    utterance.onend = () => {
        stop();
    };

    speechSynthesis.cancel();
    speechSynthesis.speak(utterance);

    VIEW.playIcon.classList.remove('fa-volume-up');
    VIEW.playIcon.classList.add('fa-pause');
    isPlaying = true;
    toggleStopButton(true);
}


function pause() {
    speechSynthesis.pause();
    isPlaying = false;
    VIEW.playIcon.classList.add('fa-volume-up');
    VIEW.playIcon.classList.remove('fa-pause');
}

function stop() {
    isPlaying = false;
    speechSynthesis.cancel();
    toggleStopButton(false);    
    VIEW.playIcon.classList.add('fa-volume-up');
    VIEW.playIcon.classList.remove('fa-pause');
}

function togglePlay() {
    if (isPlaying) {
        pause()
        return
    }
    if (speechSynthesis.paused) {
        speechSynthesis.resume();
        VIEW.playIcon.classList.remove('fa-volume-up');
        VIEW.playIcon.classList.add('fa-pause');
        isPlaying = true;
        toggleStopButton(true);
        
        return;
    }
    play()
}

function toggleStopButton(visible) {
    VIEW.stopBtn.style.display = visible ? 'inline-block' : 'none';
}

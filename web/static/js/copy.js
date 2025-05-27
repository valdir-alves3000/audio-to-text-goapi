import { showError } from "./error.js";
import { VIEW } from "./view.js";

export function setupCopyButton(){
    VIEW.copyBtn.addEventListener('click',async()=>{
        try {
            const textToCopy = VIEW.resultText.textContent;

            if(!textToCopy.trim()){
                showError("No text to copy");
                return;
            }

            await navigator.clipboard.writeText(textToCopy)

            showCopyFeedback();
        } catch (err) {
            console.error('Failed to copy text',err);
            showError("Failed to copy text: " + err.message)
            
        }
    })
}

function showCopyFeedback(){
    const originalText = VIEW.copyBtn.innerHTML;
    VIEW.copyBtn.innerHTML = '<i class="fas fa-check"></i> Copiado';
    VIEW.copyBtn.classList.add('btn-success');
    
    setTimeout(() => {
        VIEW.copyBtn.innerHTML = originalText;
        VIEW.copyBtn.classList.remove('btn-success');
    }, 2000);
}


export function setupPlay() {
    // Variável para controlar o estado da síntese de voz
    let isPlaying = false;
    let utterance = null;

    VIEW.playBtn.addEventListener('click', async () => {
        const text = VIEW.resultText.textContent.trim();
        
        if (!text) {
            showError('Nenhum texto disponível para leitura');
            return;
        }

        const lang = VIEW.languageSelect.value;
        
        try {
            if (isPlaying) {
                // Pausar se já estiver tocando
                window.speechSynthesis.pause();
                VIEW.playBtn.textContent = '▶️ Play';
                isPlaying = false;
            } else {
                // Se estava pausado, retomar
                if (window.speechSynthesis.paused) {
                    window.speechSynthesis.resume();
                    VIEW.playBtn.textContent = '⏸️ Pause';
                    isPlaying = true;
                } 
                // Se não estava tocando, começar nova leitura
                else {
                    // Cancela qualquer fala em andamento
                    window.speechSynthesis.cancel();

                    utterance = new SpeechSynthesisUtterance(text);
                    
                    // Configuração do idioma e voz
                    utterance.lang = lang === 'pt' ? 'pt-BR' : 'en-US';
                    
                    // Tenta encontrar uma voz adequada
                    const voices = window.speechSynthesis.getVoices();
                    const preferredVoice = voices.find(voice => 
                        voice.lang.includes(lang === 'pt' ? 'pt' : 'en')
                    );
                    
                    if (preferredVoice) {
                        utterance.voice = preferredVoice;
                    }

                    utterance.onstart = () => {
                        isPlaying = true;
                        VIEW.playBtn.textContent = '⏸️ Pause';
                    };

                    utterance.onend = utterance.onerror = () => {
                        isPlaying = false;
                        VIEW.playBtn.textContent = '▶️ Play';
                    };

                    window.speechSynthesis.speak(utterance);
                }
            }
        } catch (error) {
            console.error('Erro na síntese de voz:', error);
            showError('Não foi possível ler o texto em voz alta');
            VIEW.playBtn.textContent = '▶️ Play';
            isPlaying = false;
        }
    });

    VIEW.stopBtn?.addEventListener('click', () => {
        window.speechSynthesis.cancel();
        isPlaying = false;
        VIEW.playBtn.textContent = '▶️ Play';
    });

    window.speechSynthesis.onvoiceschanged = () => {
        console.log('Vozes disponíveis:', window.speechSynthesis.getVoices());
    };
}
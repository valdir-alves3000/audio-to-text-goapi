import { showError } from "./error.js?v=timestamp";
import { VIEW } from "./view.js?v=timestamp";

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
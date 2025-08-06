function startListening() {
    const SpeechRecognition = window.SpeechRecognition || window.webkitSpeechRecognition;
    if (!SpeechRecognition) {
        alert("La reconnaissance vocale n'est pas supportée par ce navigateur.");
        return;
    }

    const recognition = new SpeechRecognition();
    recognition.lang = "fr-FR";
    recognition.interimResults = false;

    recognition.onresult = function(event) {
        const transcript = event.results[0][0].transcript.toLowerCase();
        document.getElementById("output").textContent = "Tu as dit : " + transcript;

        const spokenWords = transcript.split(/\s+/);

        const maskedWords = document.querySelectorAll('.masked');

        maskedWords.forEach(span => {
            const originalWord = span.dataset.word?.toLowerCase();
            if (spokenWords.includes(originalWord)) {
                span.classList.remove("masked");
                span.textContent = originalWord;
            }
        });
    };

    recognition.onerror = function(event) {
        console.error("Erreur de reconnaissance vocale :", event.error);
    };

    recognition.start();
}

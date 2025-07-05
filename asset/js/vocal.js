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
        const text = event.results[0][0].transcript;
        document.getElementById("output").textContent = "Tu as dit : " + text;
    };

    recognition.onerror = function(event) {
        console.error("Erreur :", event.error);
    };

    recognition.start();
}
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
function selectSong(songId) {
    return __awaiter(this, void 0, void 0, function* () {
        try {
            const response = yield fetch(`/api/lyrics/${songId}`);
            const data = yield response.json();
            const lyricsContainer = document.getElementById('lyrics-text');
            const sectionParoles = document.getElementById('lyrics-container');
            if (lyricsContainer && sectionParoles) {
                lyricsContainer.innerHTML = '';
                Object.entries(data.parole).forEach(([sectionName, lines]) => {
                    const sectionDiv = document.createElement('div');
                    sectionDiv.className = 'lyric-section';
                    // On ajoute un titre pour le couplet/refrain si tu veux
                    sectionDiv.innerHTML = `<h5>${sectionName}</h5>`;
                    lines.forEach(line => {
                        const p = document.createElement('p');
                        p.textContent = line;
                        sectionDiv.appendChild(p);
                    });
                    lyricsContainer.appendChild(sectionDiv);
                });
                sectionParoles.style.display = 'block';
            }
        }
        catch (error) {
            console.error("Erreur lors du chargement des paroles:", error);
        }
    });
}
window.initLyrics = function (songFileName, points, targetId) {
    return __awaiter(this, void 0, void 0, function* () {
        try {
            const response = yield fetch(`/api/lyrics/${encodeURIComponent(songFileName)}`);
            if (!response.ok)
                throw new Error("Erreur lors de la récupération du JSON");
            const data = yield response.json();
            console.log("Données reçues du serveur Go :", data);
            const container = document.getElementById(targetId);
            if (!container)
                return;
            container.innerHTML = '';
            // Logique de masquage
            Object.entries(data.parole).forEach(([section, lines]) => {
                const sectionDiv = document.createElement('div');
                sectionDiv.className = "lyric-section mb-3";
                const title = document.createElement('h5');
                title.textContent = section;
                title.style.visibility = "hidden";
                sectionDiv.appendChild(title);
                let isMasked = false;
                if (points >= 10 && points <= 50 && section.toLowerCase().includes("refrain2")) {
                    isMasked = true;
                }
                lines.forEach(line => {
                    const p = document.createElement('p');
                    if (isMasked) {
                        p.textContent = line.replace(/[a-zA-ZÀ-ÿ]+/g, '_');
                        p.classList.add('masked');
                    }
                    else {
                        p.textContent = line;
                    }
                    sectionDiv.appendChild(p);
                });
                container.appendChild(sectionDiv);
            });
        }
        catch (error) {
            console.error("Erreur lors du chargement des paroles:", error);
        }
    });
};
document.playAudio = function (filename) {
    const audio = new Audio(`https://asset.nolp-jeu.fr/musiques/${encodeURIComponent(filename)}`);
    audio.play().catch(error => {
        console.error("Erreur lors de la lecture de l'audio:", error);
    });
};
export {};

var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
class MusicGameClient {
    constructor() {
        this.currentAudio = null;
        this.audioPlayer = document.getElementById('audio-player');
        this.gameState = {
            currentSong: null,
            currentLyrics: '',
            currentLevel: 0,
            isPlaying: false,
            lyricsVisible: false,
            maskedLyrics: '',
            lyrics: {},
            currentPoints: 0
        };
        this.initializeEventListeners();
    }
    initializeEventListeners() {
        // Event listener pour la fin de lecture audio
        if (this.audioPlayer) {
            this.audioPlayer.addEventListener('ended', () => {
                this.gameState.isPlaying = false;
                this.updatePlayButton();
            });
            this.audioPlayer.addEventListener('play', () => {
                this.gameState.isPlaying = true;
                this.updatePlayButton();
            });
            this.audioPlayer.addEventListener('pause', () => {
                this.gameState.isPlaying = false;
                this.updatePlayButton();
            });
        }
    }
    // Rechercher de la musique
    searchMusic(title_1, artist_1) {
        return __awaiter(this, arguments, void 0, function* (title, artist, instrumental = false) {
            try {
                const params = new URLSearchParams({
                    title,
                    artist,
                    instrumental: instrumental.toString()
                });
                const response = yield fetch(`/api/search-music?${params}`);
                if (!response.ok) {
                    throw new Error(`Erreur de recherche: ${response.statusText}`);
                }
                return yield response.json();
            }
            catch (error) {
                console.error('Erreur lors de la recherche musicale:', error);
                return [];
            }
        });
    }
    // Prévisualiser une chanson
    previewSong(title, artist) {
        return __awaiter(this, void 0, void 0, function* () {
            console.log(`Aperçu de la chanson: ${title} par ${artist}`);
            // Afficher le lecteur
            const musicPlayer = document.getElementById('music-player');
            const currentSongInfo = document.getElementById('current-song-info');
            if (musicPlayer)
                musicPlayer.style.display = 'block';
            if (currentSongInfo)
                currentSongInfo.textContent = `Aperçu: ${title} par ${artist}`;
            // Arrêter l'audio précédent
            this.stopCurrentAudio();
            try {
                // Rechercher la musique
                const results = yield this.searchMusic(title, artist, false);
                if (results.length > 0) {
                    const result = results[0];
                    if (result.source === 'spotify' && result.preview_url) {
                        // Utiliser l'aperçu Spotify
                        this.audioPlayer.src = result.preview_url;
                    }
                    else if (result.source === 'youtube') {
                        // Pour YouTube, utiliser le proxy audio
                        this.audioPlayer.src = `/api/audio-proxy?url=${encodeURIComponent(result.source_url)}`;
                    }
                    else if (result.source === 'local') {
                        // Fichier local
                        this.audioPlayer.src = result.preview_url;
                    }
                    this.gameState.currentSong = result;
                }
                else {
                    // Utiliser un fichier de démonstration
                    this.audioPlayer.src = '/demo-audio.mp3';
                }
            }
            catch (error) {
                console.error('Erreur lors du chargement de l\'aperçu:', error);
                this.audioPlayer.src = '/demo-audio.mp3';
            }
            // Cacher les paroles pendant l'aperçu
            const lyricsContainer = document.getElementById('lyrics-container');
            if (lyricsContainer)
                lyricsContainer.style.display = 'none';
        });
    }
    // Sélectionner une chanson pour le jeu
    selectSong(level, songIndex, title, artist) {
        return __awaiter(this, void 0, void 0, function* () {
            console.log(`Chanson sélectionnée: ${title} par ${artist}, niveau: ${level}, index: ${songIndex}`);
            this.gameState.currentLevel = parseInt(level);
            // Afficher le lecteur
            const musicPlayer = document.getElementById('music-player');
            const currentSongInfo = document.getElementById('current-song-info');
            if (musicPlayer)
                musicPlayer.style.display = 'block';
            if (currentSongInfo) {
                currentSongInfo.textContent = `${title} par ${artist} (${level} points)`;
            }
            // Charger les paroles
            yield this.loadLyrics(level, songIndex);
            // Rechercher et charger la version instrumentale
            yield this.loadInstrumental(title, artist);
            // Afficher les paroles
            const lyricsContainer = document.getElementById('lyrics-container');
            if (lyricsContainer)
                lyricsContainer.style.display = 'block';
            this.gameState.lyricsVisible = true;
        });
    }
    // Charger les paroles depuis l'API
    loadLyrics(level, songIndex) {
        return __awaiter(this, void 0, void 0, function* () {
            try {
                const response = yield fetch(`/api/get-lyrics/${level}/${songIndex}`);
                if (response.ok) {
                    const lyricsData = yield response.json();
                    this.gameState.currentLyrics = lyricsData.parole || "Paroles non disponibles";
                    this.gameState.lyrics = this.splitLyricsBySections(this.gameState.currentLyrics);
                    this.displayMaskedLyrics(this.gameState.lyrics, this.gameState.currentPoints);
                }
                else {
                    this.gameState.currentLyrics = "Paroles non disponibles";
                    this.displayLyrics(this.gameState.currentLyrics);
                }
            }
            catch (error) {
                console.error('Erreur lors du chargement des paroles:', error);
                this.gameState.currentLyrics = "Erreur de chargement des paroles";
                this.displayLyrics(this.gameState.currentLyrics);
            }
        });
    }
    // Ajoute cette méthode dans ta classe :
    splitLyricsBySections(lyrics) {
        const sections = {};
        let currentSection = "default";
        const lines = lyrics.split('\n');
        for (const line of lines) {
            const trimmed = line.trim();
            if (trimmed.startsWith('[') && trimmed.endsWith(']')) {
                currentSection = trimmed.slice(1, -1);
                sections[currentSection] = [];
            }
            else if (trimmed) {
                if (!sections[currentSection])
                    sections[currentSection] = [];
                sections[currentSection].push(trimmed);
            }
        }
        return sections;
    }
    // Charger la version instrumentale
    loadInstrumental(title, artist) {
        return __awaiter(this, void 0, void 0, function* () {
            try {
                const results = yield this.searchMusic(title, artist, true);
                if (results.length > 0) {
                    const result = results[0];
                    if (result.source === 'local') {
                        this.audioPlayer.src = result.preview_url;
                    }
                    else if (result.source === 'youtube') {
                        this.audioPlayer.src = `/api/audio-proxy?url=${encodeURIComponent(result.source_url)}`;
                    }
                    else {
                        // Fallback sur l'original si pas d'instrumental trouvé
                        this.audioPlayer.src = result.preview_url || '/demo-instrumental.mp3';
                    }
                    this.gameState.currentSong = result;
                }
                else {
                    this.audioPlayer.src = '/demo-instrumental.mp3';
                }
            }
            catch (error) {
                console.error('Erreur lors du chargement de l\'instrumental:', error);
                this.audioPlayer.src = '/demo-instrumental.mp3';
            }
        });
    }
    displayMaskedLyrics(structuredLyrics, points) {
        const maskedLines = [];
        const isSectionVisible = (section) => {
            const lower = section.toLowerCase();
            if (points >= 40)
                return lower === "couplet1";
            if (points <= 30)
                return !lower.includes("refrain");
            return true;
        };
        for (const section in structuredLyrics) {
            const lines = structuredLyrics[section];
            for (const line of lines) {
                if (isSectionVisible(section)) {
                    maskedLines.push(line);
                }
                else {
                    const words = line.split(" ");
                    const masked = words.map(word => `<span class="masked-text">${'█'.repeat(word.length)}</span>`).join(" ");
                    maskedLines.push(masked);
                }
            }
            maskedLines.push("");
        }
        this.gameState.maskedLyrics = maskedLines.join("<br>");
        this.displayLyrics(this.gameState.maskedLyrics);
    }
    displayLyrics(lyrics) {
        const lyricsText = document.getElementById('lyrics-text');
        if (lyricsText) {
            lyricsText.innerHTML = lyrics;
        }
    }
    toggleLyrics() {
        const lyricsContainer = document.getElementById('lyrics-container');
        if (lyricsContainer) {
            if (lyricsContainer.style.display === 'none') {
                lyricsContainer.style.display = 'block';
                this.gameState.lyricsVisible = true;
            }
            else {
                lyricsContainer.style.display = 'none';
                this.gameState.lyricsVisible = false;
            }
        }
    }
    revealMoreLyrics() {
        if (!this.gameState.currentLyrics)
            return;
        const currentMaskPercentage = this.calculateCurrentMaskPercentage();
        const newMaskPercentage = Math.max(0, currentMaskPercentage - 0.1);
        this.displayMaskedLyricsWithPercentage(this.gameState.currentLyrics, newMaskPercentage);
    }
    revealAllLyrics() {
        this.displayLyrics(this.gameState.currentLyrics);
    }
    calculateCurrentMaskPercentage() {
        switch (this.gameState.currentLevel) {
            case 50: return 0.8;
            case 40: return 0.6;
            case 30: return 0.4;
            case 20: return 0.2;
            case 10: return 0.1;
            default: return 0.3;
        }
    }
    displayMaskedLyricsWithPercentage(lyrics, maskPercentage) {
        const words = lyrics.split(' ');
        const wordsToMask = Math.floor(words.length * maskPercentage);
        const indicesToMask = [];
        while (indicesToMask.length < wordsToMask) {
            const randomIndex = Math.floor(Math.random() * words.length);
            if (!indicesToMask.includes(randomIndex)) {
                indicesToMask.push(randomIndex);
            }
        }
        const maskedWords = words.map((word, index) => {
            if (indicesToMask.includes(index)) {
                const safeWord = word.replace(/"/g, '&quot;').replace(/</g, '&lt;').replace(/>/g, '&gt;');
                return `<span class="masked" data-word="${safeWord}">${'█'.repeat(word.length)}</span>`;
            }
            return word;
        });
        this.displayLyrics(maskedWords.join(' '));
    }
    stopCurrentAudio() {
        if (this.currentAudio) {
            this.currentAudio.pause();
            this.currentAudio.currentTime = 0;
        }
        if (this.audioPlayer) {
            this.audioPlayer.pause();
            this.audioPlayer.currentTime = 0;
        }
    }
    // Mettre à jour le bouton de lecture
    updatePlayButton() {
        // Ici vous pouvez ajouter la logique pour mettre à jour l'interface
        // selon l'état de lecture
    }
    // Obtenir l'état actuel du jeu
    getGameState() {
        return Object.assign({}, this.gameState);
    }
}
// Initialiser le client quand le DOM est chargé
document.addEventListener('DOMContentLoaded', () => {
    const musicClient = new MusicGameClient();
    // Rendre le client accessible globalement pour les fonctions inline du HTML
    window.musicClient = musicClient;
    // Redéfinir les fonctions globales pour utiliser le client TypeScript
    window.previewSong = (title, artist) => {
        musicClient.previewSong(title, artist);
    };
    window.selectSong = (level, songIndex, title, artist) => {
        musicClient.selectSong(level, songIndex, title, artist);
    };
    window.toggleLyrics = () => {
        musicClient.toggleLyrics();
    };
    // Nouvelles fonctions pour les indices
    window.revealMoreLyrics = () => {
        musicClient.revealMoreLyrics();
    };
    window.revealAllLyrics = () => {
        musicClient.revealAllLyrics();
    };
});
// Export pour utilisation en module
export { MusicGameClient };

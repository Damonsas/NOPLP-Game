// Types TypeScript pour la gestion de la musique
interface MusicSearchResult {
    title: string;
    artist: string;
    preview_url: string;
    source_url: string;
    source: 'youtube' | 'spotify' | 'local';
}

interface LyricsData {
    titre: string;
    artiste: string;
    parole: string;
}

interface GameState {
    currentSong: MusicSearchResult | null;
    currentLyrics: string;
    currentLevel: number;
    isPlaying: boolean;
    lyricsVisible: boolean;
    maskedLyrics: string;
    lyrics: { [section: string]: string[] };
    currentPoints: number;
}

class MusicGameClient {
    private audioPlayer: HTMLAudioElement;
    private gameState: GameState;
    private currentAudio: HTMLAudioElement | null = null;

    constructor() {
        this.audioPlayer = document.getElementById('audio-player') as HTMLAudioElement;
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

    private initializeEventListeners(): void {
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
    async searchMusic(title: string, artist: string, instrumental: boolean = false): Promise<MusicSearchResult[]> {
        try {
            const params = new URLSearchParams({
                title,
                artist,
                instrumental: instrumental.toString()
            });

            const response = await fetch(`/api/search-music?${params}`);
            if (!response.ok) {
                throw new Error(`Erreur de recherche: ${response.statusText}`);
            }

            return await response.json();
        } catch (error) {
            console.error('Erreur lors de la recherche musicale:', error);
            return [];
        }
    }

    // Prévisualiser une chanson
    async previewSong(title: string, artist: string): Promise<void> {
        console.log(`Aperçu de la chanson: ${title} par ${artist}`);
        
        // Afficher le lecteur
        const musicPlayer = document.getElementById('music-player');
        const currentSongInfo = document.getElementById('current-song-info');
        
        if (musicPlayer) musicPlayer.style.display = 'block';
        if (currentSongInfo) currentSongInfo.textContent = `Aperçu: ${title} par ${artist}`;

        // Arrêter l'audio précédent
        this.stopCurrentAudio();

        try {
            // Rechercher la musique
            const results = await this.searchMusic(title, artist, false);
            
            if (results.length > 0) {
                const result = results[0];
                
                if (result.source === 'spotify' && result.preview_url) {
                    // Utiliser l'aperçu Spotify
                    this.audioPlayer.src = result.preview_url;
                } else if (result.source === 'youtube') {
                    // Pour YouTube, utiliser le proxy audio
                    this.audioPlayer.src = `/api/audio-proxy?url=${encodeURIComponent(result.source_url)}`;
                } else if (result.source === 'local') {
                    // Fichier local
                    this.audioPlayer.src = result.preview_url;
                }
                
                this.gameState.currentSong = result;
            } else {
                // Utiliser un fichier de démonstration
                this.audioPlayer.src = '/demo-audio.mp3';
            }
        } catch (error) {
            console.error('Erreur lors du chargement de l\'aperçu:', error);
            this.audioPlayer.src = '/demo-audio.mp3';
        }

        // Cacher les paroles pendant l'aperçu
        const lyricsContainer = document.getElementById('lyrics-container');
        if (lyricsContainer) lyricsContainer.style.display = 'none';
    }

    // Sélectionner une chanson pour le jeu
    async selectSong(level: string, songIndex: number, title: string, artist: string): Promise<void> {
        console.log(`Chanson sélectionnée: ${title} par ${artist}, niveau: ${level}, index: ${songIndex}`);
        
        this.gameState.currentLevel = parseInt(level);
        
        // Afficher le lecteur
        const musicPlayer = document.getElementById('music-player');
        const currentSongInfo = document.getElementById('current-song-info');
        
        if (musicPlayer) musicPlayer.style.display = 'block';
        if (currentSongInfo) {
            currentSongInfo.textContent = `${title} par ${artist} (${level} points)`;
        }

        // Charger les paroles
        await this.loadLyrics(level, songIndex);

        // Rechercher et charger la version instrumentale
        await this.loadInstrumental(title, artist);

        // Afficher les paroles
        const lyricsContainer = document.getElementById('lyrics-container');
        if (lyricsContainer) lyricsContainer.style.display = 'block';
        
        this.gameState.lyricsVisible = true;
    }

    

    // Charger les paroles depuis l'API
    private async loadLyrics(level: string, songIndex: number): Promise<void> {
    try {
        const response = await fetch(`/api/get-lyrics/${level}/${songIndex}`);
        if (response.ok) {
            const lyricsData: LyricsData = await response.json();
            this.gameState.currentLyrics = lyricsData.parole || "Paroles non disponibles";
            this.gameState.lyrics = this.splitLyricsBySections(this.gameState.currentLyrics);
            this.displayMaskedLyrics(this.gameState.lyrics, this.gameState.currentPoints);
        } else {
            this.gameState.currentLyrics = "Paroles non disponibles";
            this.displayLyrics(this.gameState.currentLyrics);
        }
    } catch (error) {
        console.error('Erreur lors du chargement des paroles:', error);
        this.gameState.currentLyrics = "Erreur de chargement des paroles";
        this.displayLyrics(this.gameState.currentLyrics);
    }
}

// Ajoute cette méthode dans ta classe :
private splitLyricsBySections(lyrics: string): { [section: string]: string[] } {
    const sections: { [section: string]: string[] } = {};
    let currentSection = "default";
    const lines = lyrics.split('\n');
    for (const line of lines) {
        const trimmed = line.trim();
        if (trimmed.startsWith('[') && trimmed.endsWith(']')) {
            currentSection = trimmed.slice(1, -1);
            sections[currentSection] = [];
        } else if (trimmed) {
            if (!sections[currentSection]) sections[currentSection] = [];
            sections[currentSection].push(trimmed);
        }
    }
    return sections;
}

    // Charger la version instrumentale
    private async loadInstrumental(title: string, artist: string): Promise<void> {
        try {
            const results = await this.searchMusic(title, artist, true);
            
            if (results.length > 0) {
                const result = results[0];
                
                if (result.source === 'local') {
                    this.audioPlayer.src = result.preview_url;
                } else if (result.source === 'youtube') {
                    this.audioPlayer.src = `/api/audio-proxy?url=${encodeURIComponent(result.source_url)}`;
                } else {
                    // Fallback sur l'original si pas d'instrumental trouvé
                    this.audioPlayer.src = result.preview_url || '/demo-instrumental.mp3';
                }
                
                this.gameState.currentSong = result;
            } else {
                this.audioPlayer.src = '/demo-instrumental.mp3';
            }
        } catch (error) {
            console.error('Erreur lors du chargement de l\'instrumental:', error);
            this.audioPlayer.src = '/demo-instrumental.mp3';
        }
    }

    private displayMaskedLyrics(structuredLyrics: {[section: string]: string[]}, points: number): void {
    const maskedLines: string[] = [];

    const isSectionVisible = (section: string): boolean => {
        const lower = section.toLowerCase();
        if (points >= 40) return lower === "couplet1"; // 40/50 : juste couplet1 visible
        if (points <= 30) return !lower.includes("refrain"); // 30/20/10 : pas les refrains
        return true;
    };

    for (const section in structuredLyrics) {
        const lines = structuredLyrics[section];
        for (const line of lines) {
            if (isSectionVisible(section)) {
                maskedLines.push(line); // afficher normalement
            } else {
                const words = line.split(" ");
                const masked = words.map(word =>
                    `<span class="masked-text">${'█'.repeat(word.length)}</span>`
                ).join(" ");
                maskedLines.push(masked);
            }
        }

        maskedLines.push(""); // saut de ligne entre les sections
    }

    this.gameState.maskedLyrics = maskedLines.join("<br>");
    this.displayLyrics(this.gameState.maskedLyrics);
}


    private displayLyrics(lyrics: string): void {
        const lyricsText = document.getElementById('lyrics-text');
        if (lyricsText) {
            lyricsText.innerHTML = lyrics;
        }
    }

    toggleLyrics(): void {
        const lyricsContainer = document.getElementById('lyrics-container');
        if (lyricsContainer) {
            if (lyricsContainer.style.display === 'none') {
                lyricsContainer.style.display = 'block';
                this.gameState.lyricsVisible = true;
            } else {
                lyricsContainer.style.display = 'none';
                this.gameState.lyricsVisible = false;
            }
        }
    }

    revealMoreLyrics(): void {
        if (!this.gameState.currentLyrics) return;

        const currentMaskPercentage = this.calculateCurrentMaskPercentage();
        const newMaskPercentage = Math.max(0, currentMaskPercentage - 0.1);
        
        this.displayMaskedLyricsWithPercentage(this.gameState.currentLyrics, newMaskPercentage);
    }

    revealAllLyrics(): void {
        this.displayLyrics(this.gameState.currentLyrics);
    }

    private calculateCurrentMaskPercentage(): number {
        switch(this.gameState.currentLevel) {
            case 50: return 0.8;
            case 40: return 0.6;
            case 30: return 0.4;
            case 20: return 0.2;
            case 10: return 0.1;
            default: return 0.3;
        }
    }

    private displayMaskedLyricsWithPercentage(lyrics: string, maskPercentage: number): void {
        const words = lyrics.split(' ');
        const wordsToMask = Math.floor(words.length * maskPercentage);

        const indicesToMask: number[] = [];
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

    private stopCurrentAudio(): void {
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
    private updatePlayButton(): void {
        // Ici vous pouvez ajouter la logique pour mettre à jour l'interface
        // selon l'état de lecture
    }

    // Obtenir l'état actuel du jeu
    getGameState(): GameState {
        return { ...this.gameState };
    }
}



// Initialiser le client quand le DOM est chargé
document.addEventListener('DOMContentLoaded', () => {
    const musicClient = new MusicGameClient();
    
    // Rendre le client accessible globalement pour les fonctions inline du HTML
    (window as any).musicClient = musicClient;
    
    // Redéfinir les fonctions globales pour utiliser le client TypeScript
    (window as any).previewSong = (title: string, artist: string) => {
        musicClient.previewSong(title, artist);
    };
    
    (window as any).selectSong = (level: string, songIndex: number, title: string, artist: string) => {
        musicClient.selectSong(level, songIndex, title, artist);
    };
    
    (window as any).toggleLyrics = () => {
        musicClient.toggleLyrics();
    };
    
    // Nouvelles fonctions pour les indices
    (window as any).revealMoreLyrics = () => {
        musicClient.revealMoreLyrics();
    };
    
    (window as any).revealAllLyrics = () => {
        musicClient.revealAllLyrics();
    };
});

// Export pour utilisation en module
export { MusicGameClient, MusicSearchResult, LyricsData, GameState };

interface LyricsData {
    titre: string;
    artiste: string;
    parole: { [key: string]: string[] };
}

interface LyricsJSON {
    titre: string;
    artiste: string;
    parole: { [key: string]: string[] };
}

interface Song {
    titre: string;
    artiste: string;
    filename: string;
}

async function selectSong(songId: string) {
    try {
        const response = await fetch(`/api/lyrics/${songId}`);
        const data: LyricsData = await response.json();

        const lyricsContainer = document.getElementById('lyrics-text');
        const sectionParoles = document.getElementById('lyrics-container');

        if (lyricsContainer && sectionParoles) {
            lyricsContainer.innerHTML = ''; 

            Object.entries(data.parole).forEach(([sectionName, lines]) => {
                const sectionDiv = document.createElement('div');
                sectionDiv.className = 'lyric-section';
                
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
    } catch (error) {
        console.error("Erreur lors du chargement des paroles:", error);
    }
}

export {};

declare global {
    interface Window {
        initLyrics?: (songFileName: string, points: number | string, targetId: string) => Promise<void>;
    }
    interface Document {
        playAudio: (filename: string) => void;
    }
}

// Générateur de nombre aléatoire inclusif
function getRandomInt(min: number, max: number): number {
    return Math.floor(Math.random() * (max - min + 1)) + min;
}

function maskWordsInLine(line: string, wordsCountToMask: number): string {
    const tokens = line.split(/(\s+)/);
    
    const wordIndices: number[] = [];
    tokens.forEach((token, index) => {
        if (/[a-zA-ZÀ-ÿ]+/.test(token)) {
            wordIndices.push(index);
        }
    });

    if (wordIndices.length === 0) return line;

    const exactMaskCount = Math.min(wordsCountToMask, wordIndices.length);
    
    const maxStartIndex = wordIndices.length - exactMaskCount;
    const startIndex = getRandomInt(0, maxStartIndex);

    const targetIndicesToMask = wordIndices.slice(startIndex, startIndex + exactMaskCount);

    const maskedTokens = tokens.map((token, index) => {
        if (targetIndicesToMask.includes(index)) {
            return token.replace(/[a-zA-ZÀ-ÿ]/g, '_');
        }
        return token;
    });

    return maskedTokens.join('');
}

window.initLyrics = async function(songFileName: string, points: number | string, targetId: string) {
    try {
        const response = await fetch(`/api/lyrics/${encodeURIComponent(songFileName)}`);
        if (!response.ok) throw new Error("Erreur lors de la récupération du JSON des paroles");
        
        const data: LyricsJSON = await response.json();
        console.log("Données de paroles reçues :", data);
        
        const container = document.getElementById(targetId);
        if (!container) return;
        
        container.innerHTML = ''; 

        const sectionsKeys = Object.keys(data.parole);
        
        const refrains = sectionsKeys.filter(k => k.toLowerCase().includes('refrain'));
        const couplets = sectionsKeys.filter(k => k.toLowerCase().includes('couplet'));

        const firstRefrainKey = refrains[0] || null;
        const secondRefrainKey = refrains[1] || refrains[0] || null; // fallback au 1er si pas de 2ème
        const secondCoupletKey = couplets[1] || couplets[0] || null;

        let mode: 'points' | 'same-song' = 'points';
        let wordsToMask = 0;
        let targetSectionKey: string | null = null;

        const pointsStr = String(points);

        if (pointsStr === "same-song") {
            mode = 'same-song';
        } else {
            const pts = Number.parseInt(pointsStr, 10);
            if (pts === 50) {
                wordsToMask = getRandomInt(8, 10);
                targetSectionKey = secondRefrainKey;
            } else if (pts === 40) {
                wordsToMask = getRandomInt(5, 7);
                targetSectionKey = firstRefrainKey;
            } else if (pts === 30) {
                wordsToMask = getRandomInt(4, 5);
                targetSectionKey = secondCoupletKey;
            } else if (pts === 20) {
                wordsToMask = 3;
                targetSectionKey = secondCoupletKey;
            } else if (pts === 10) {
                wordsToMask = 2;
                targetSectionKey = secondCoupletKey;
            }
        }

        let sameSongCutLineIndex = -1;
        let sameSongCurrentLineCounter = 0;
        let sameSongIsCutting = false;

        if (mode === 'same-song') {
            sameSongCutLineIndex = getRandomInt(3, 5);
        }

        sectionsKeys.forEach((sectionName) => {
            const lines = data.parole[sectionName];
            const sectionDiv = document.createElement('div');
            sectionDiv.className = "lyric-section mb-3";

            const title = document.createElement('h5');
            title.textContent = sectionName;
            title.style.visibility = "hidden"; 
            sectionDiv.appendChild(title);

            const isTargetSection = (mode === 'points' && targetSectionKey && sectionName === targetSectionKey);
            
            let randomLineIndexToMask = -1;
            if (isTargetSection && lines.length > 0) {
                randomLineIndexToMask = getRandomInt(0, lines.length - 1);
            }

            lines.forEach((line, index) => {
                const p = document.createElement('p');

                if (mode === 'same-song') {
                    if (sameSongCurrentLineCounter >= sameSongCutLineIndex) {
                        sameSongIsCutting = true;
                    }

                    if (sameSongIsCutting) {
                        p.textContent = line.replace(/[a-zA-ZÀ-ÿ]/g, '_');
                        p.classList.add('masked');
                    } else {
                        p.textContent = line;
                    }
                    sameSongCurrentLineCounter++;

                } else if (isTargetSection && index === randomLineIndexToMask) {
                    // Logique classique par Points
                    p.textContent = maskWordsInLine(line, wordsToMask);
                    p.classList.add('masked');
                } else {
                    p.textContent = line;
                }

                sectionDiv.appendChild(p);
            });

            container.appendChild(sectionDiv);
        });

        const audioPlayer = document.getElementById(`audio-player-${points}`) as HTMLAudioElement | null;

        if (audioPlayer) {
            let audioFileName = songFileName.replace('.json', '.mp3');
            
            audioFileName = audioFileName
                .toLowerCase()
                .normalize("NFD")
                .replace(/[\u0300-\u036f]/g, "");

            console.log(`[Audio] Chargement du lecteur #audio-player-${points} avec :`, audioFileName);
            
            audioPlayer.src = `/api/musiques/${encodeURIComponent(audioFileName)}`;
            
            audioPlayer.load();
            audioPlayer.play().catch(error => {
                console.error("Erreur d'autoplay :", error);
            });
        } else {
            console.error(`Impossible de trouver le lecteur audio : #audio-player-${points}`);
        }

    } catch (error) {
        console.error("Erreur globale lors de l'initialisation du round :", error);
    }
}

document.playAudio = function(filename: string) {
    const audio = new Audio(`https://asset.nolp-jeu.fr/musiques/${encodeURIComponent(filename)}`);
    audio.play().catch(error => {
        console.error("Erreur lors de la lecture de l'audio:", error);
    });
}

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
    } catch (error) {
        console.error("Erreur lors du chargement des paroles:", error);
    }
}

export {};


declare global {
    interface Window {
        initLyrics?: (songFileName: string, points: number, targetId: string) => Promise<void>;
    }
     interface Document {
        playAudio: (filename: string) => void;
    }
}

window.initLyrics = async function(songFileName: string, points: number, targetId: string) {
    try {
        const response = await fetch(`/api/lyrics/${encodeURIComponent(songFileName)}`);
        if (!response.ok) throw new Error("Erreur lors de la récupération du JSON des paroles");
        
        const data: LyricsJSON = await response.json();
        console.log("Données de paroles reçues :", data);
        
        const container = document.getElementById(targetId);
        if (!container) return;
        
        container.innerHTML = ''; 

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
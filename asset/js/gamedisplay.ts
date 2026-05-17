
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
}

window.initLyrics = async function(songFileName: string, points: number, targetId: string) {
    try {
        const response = await fetch(`/api/lyrics/${encodeURIComponent(songFileName)}`);
        if (!response.ok) throw new Error("Erreur lors de la récupération du JSON");
        
        const data: LyricsJSON = await response.json();

        console.log("Données reçues du serveur Go :", data);
        
        const container = document.getElementById(targetId);
        if (!container) return;
        
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
                } else {
                    p.textContent = line;
                }
                sectionDiv.appendChild(p);
            });

            container.appendChild(sectionDiv);
        });

    } catch (error) {
        console.error("Erreur lors du chargement des paroles:", error);
    }
}
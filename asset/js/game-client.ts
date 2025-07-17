
interface Song {
    title: string;
    artist: string;
}

interface GameSession {
    currentSong: Song | null;
    lyricsVisible: boolean;
    lyricsContent: string | null;
}

const sessionId: string = 'ID_DE_LA_SESSION'; 

async function callApi<T>(url: string, body: T): Promise<GameSession> {
    const response = await fetch(url, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body)
    });
    if (!response.ok) {
        throw new Error(`Erreur API: ${response.statusText}`);
    }
    return response.json();
}

async function playSong(level: string, songIndex: number): Promise<void> {
    const updatedSession = await callApi(`/api/game-sessions/${sessionId}/start-song`, { level, songIndex });
    renderGame(updatedSession);
}

async function setLyricsVisibility(visible: boolean): Promise<void> {
    const updatedSession = await callApi(`/api/game-sessions/${sessionId}/lyrics-visibility`, { visible });
    renderGame(updatedSession);
}

function renderGame(session: GameSession): void {
    const songTitleElement = document.getElementById('song-title') as HTMLHeadingElement | null;
    const lyricsDiv = document.getElementById('lyrics') as HTMLDivElement | null;

    if (songTitleElement && session.currentSong) {
        songTitleElement.textContent = `${session.currentSong.title} - ${session.currentSong.artist}`;
    }

    if (lyricsDiv) {
        if (session.lyricsVisible && session.lyricsContent) {
            lyricsDiv.innerText = session.lyricsContent; // innerText est plus sûr pour afficher du texte brut
            lyricsDiv.style.display = 'block';
        } else {
            lyricsDiv.style.display = 'none';
        }
    }
}

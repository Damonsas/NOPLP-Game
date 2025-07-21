var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
const sessionId = 'ID_DE_LA_SESSION';
function callApi(url, body) {
    return __awaiter(this, void 0, void 0, function* () {
        const response = yield fetch(url, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(body)
        });
        if (!response.ok) {
            throw new Error(`Erreur API: ${response.statusText}`);
        }
        return response.json();
    });
}
function playSong(level, songIndex) {
    return __awaiter(this, void 0, void 0, function* () {
        const updatedSession = yield callApi(`/api/game-sessions/${sessionId}/start-song`, { level, songIndex });
        renderGame(updatedSession);
    });
}
function setLyricsVisibility(visible) {
    return __awaiter(this, void 0, void 0, function* () {
        const updatedSession = yield callApi(`/api/game-sessions/${sessionId}/lyrics-visibility`, { visible });
        renderGame(updatedSession);
    });
}
function renderGame(session) {
    const songTitleElement = document.getElementById('song-title');
    const lyricsDiv = document.getElementById('lyrics');
    if (songTitleElement && session.currentSong) {
        songTitleElement.textContent = `${session.currentSong.title} - ${session.currentSong.artist}`;
    }
    if (lyricsDiv) {
        if (session.lyricsVisible && session.lyricsContent) {
            lyricsDiv.innerText = session.lyricsContent; // innerText est plus sûr pour afficher du texte brut
            lyricsDiv.style.display = 'block';
        }
        else {
            lyricsDiv.style.display = 'none';
        }
    }
}

var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
function displaySong(titre, artiste, paroles) {
    return __awaiter(this, void 0, void 0, function* () {
        const songDisplay = document.getElementById('song-selection-section');
        const selectedSong = document.getElementsByClassName('btn-select');
        if (selectedSong) {
            if (songDisplay.style.display === 'none' || songDisplay.classList.contains('hidden')) {
                songDisplay.style.display = 'block';
                songDisplay.classList.remove('hidden');
            }
            else {
                songDisplay.style.display = 'none';
                songDisplay.classList.add('hidden');
            }
        }
    });
}

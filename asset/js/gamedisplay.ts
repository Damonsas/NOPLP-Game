async function displaySong (titre: string, artiste: string, paroles: string) {
    const songDisplay = document.getElementById('song-selection-section');
    const selectedSong = document.getElementsByClassName('btn-select');
    if (selectedSong){
        if (songDisplay.style.display === 'none' || songDisplay.classList.contains('hidden')) {
            songDisplay.style.display = 'block';
            songDisplay.classList.remove('hidden');  
        }
        else {
            songDisplay.style.display = 'none';
            songDisplay.classList.add('hidden');
        }
    }
}

function toggleElement(elementId) {
    const element = document.getElementById(elementId);
    if (element) {
        element.classList.toggle('visible');
    }
}

document.getElementById("startLyricsBtn").addEventListener("click", () => {
  document.getElementById("songSelect").style.display = "none";
  document.getElementById("lyricsSection").style.display = "block";
  startLyrics();
});


// Fonction pour toggle les chansons d'un niveau et fermer les autres
function toggleLevelSongs(levelId) {
    const element = document.getElementById(levelId);
    if (!element) return;
    
    const isCurrentlyVisible = element.classList.contains('visible');
    
    // Fermer tous les autres niveaux
    const allLevels = document.querySelectorAll('.songs-for-level');
    allLevels.forEach(function(level) {
        level.classList.remove('visible');
    });
    
    // Retirer la classe active de tous les boutons
    const allButtons = document.querySelectorAll('.point-button');
    allButtons.forEach(function(button) {
        button.classList.remove('active');
    });
    
    // Si l'élément n'était pas visible, l'afficher et marquer le bouton comme actif
    if (!isCurrentlyVisible) {
        element.classList.add('visible');
        // Trouver le bouton qui a été cliqué et le marquer comme actif
        const clickedButton = document.querySelector('[onclick*="' + levelId + '"]');
        if (clickedButton) {
            clickedButton.classList.add('active');
        }
    }
}
// Debug - vérifier que les fonctions sont chargées
console.log('Visibility functions loaded:', {
    toggleElement: typeof toggleElement,
    toggleLevelSongs: typeof toggleLevelSongs
});

        // Ajout des keyframes CSS
        const style = document.createElement('style');
        style.textContent = `
            @keyframes fadeInDown {
                from {
                    opacity: 0;
                    transform: translateY(-30px);
                }
                to {
                    opacity: 1;
                    transform: translateY(0);
                }
            }
            
            @keyframes fadeInUp {
                from {
                    opacity: 0;
                    transform: translateY(30px);
                }
                to {
                    opacity: 1;
                    transform: translateY(0);
                }
            }
        `;
        document.head.appendChild(style);
        window.toggleLevelSongs = toggleLevelSongs;
        window.toggleElement = toggleElement;
  
async function selectSong(level, songIndex, title, artist) {
    console.log("Chanson sélectionnée:", title, "par", artist, "niveau:", level, "index:", songIndex);
    
    currentLevel = level;
    
    // Afficher le lecteur
    document.getElementById('music-player').style.display = 'block';
    document.getElementById('current-song-info').textContent = title + " par " + artist + " (" + level + " points)";
    
    // Charger les paroles
    try {
        const response = await fetch('/api/get-lyrics/' + level + '/' + songIndex);
        if (response.ok) {
            const lyricsData = await response.json();
            currentLyrics = lyricsData.parole || lyricsData.lyrics || "Paroles non disponibles";
            displayMaskedLyrics(currentLyrics, parseInt(level));
        } else {
            currentLyrics = "Paroles non disponibles";
            document.getElementById('lyrics-text').textContent = currentLyrics;
        }
    } catch (error) {
        console.error("Erreur lors du chargement des paroles:", error);
        currentLyrics = "Erreur de chargement des paroles";
        document.getElementById('lyrics-text').textContent = currentLyrics;
    }
    
    // Charger l'audio (instrumental si possible)
    const audioPlayer = document.getElementById('audio-player');
    // Ici vous pourrez intégrer l'API de votre choix pour l'audio instrumental
    audioPlayer.src = "/demo-instrumental.mp3"; // Fichier de démonstration
    
    // Afficher les paroles
    document.getElementById('lyrics-container').style.display = 'block';
    
    // ✨ **ACTION : Rendre les boutons visibles**
    document.getElementById('action-buttons').style.visibility = 'visible';
}
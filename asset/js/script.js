
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
  

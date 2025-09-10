
function toggleElement(elementId) {
    const element = document.getElementById(elementId);
    if (element) {
        if (element.style.display === 'none' || element.classList.contains('hidden-element')) {
            element.style.display = 'block';
            element.classList.remove('hidden-element');
            element.classList.add('visible');
        } else {
            element.style.display = 'none';
            element.classList.add('hidden-element');
            element.classList.remove('visible');
        }
    }
}

function toggleLevelSongs(levelId) {
    const element = document.getElementById(levelId);
    if (!element) return;
    
    const isCurrentlyVisible = element.classList.contains('visible');
    
    const allLevels = document.querySelectorAll('.songs-for-level');
    allLevels.forEach(function(level) {
        level.classList.remove('visible');
        level.style.display = 'none';
        level.classList.add('hidden-element');
    });
    
    const allButtons = document.querySelectorAll('.point-button');
    allButtons.forEach(function(button) {
        button.classList.remove('active');
    });
    
    if (!isCurrentlyVisible) {
        element.classList.add('visible');
        element.style.display = 'block';
        element.classList.remove('hidden-element');
        const clickedButton = document.querySelector('[onclick*="' + levelId + '"]');
        if (clickedButton) {
            clickedButton.classList.add('active');
        }
    }
}

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

document.addEventListener('DOMContentLoaded', function() {
    const startBtn = document.querySelector('.startLyricsBtn');
    if (startBtn) {
        startBtn.addEventListener("click", () => {
            if (window.musicClient) {
                const gameState = window.musicClient.getGameState();
                if (gameState.currentSong && gameState.currentLyrics) {
                    console.log('Démarrage du jeu avec:', gameState);
                    alert('Jeu démarré avec la chanson sélectionnée !');
                } else {
                    alert(`Veuillez d'abord sélectionner une chanson`);
                }
            }
        });
    }
});

window.toggleElement = toggleElement;
window.toggleLevelSongs = toggleLevelSongs;

// console.log('Script.js (interface) loaded:', {
//     toggleElement: typeof toggleElement,
//     toggleLevelSongs: typeof toggleLevelSongs
// });
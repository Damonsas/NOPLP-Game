let preparedDuels = []; /// ça définie la grille

function createduel(prepduel) {
    const container = document.getElementsByClassName("Sectionduel")[0];

    if (prepduel.length === 0) {
        container.innerHTML = `
            <div class="alert alert-info">
                Aucune grille n'a été trouvée, veuillez en créer une via le bouton ci-dessous. 
            </div>
            ${getDuelHtml()}
        `;
        return;
    }

    let duelsHtml = getDuelHtml();
    duelsHtml += '<div class="duels-list">';
    
    prepduel.forEach((duel, index) => {
        duelsHtml += generateDuelCard(duel, index);
    });
    
    duelsHtml += '</div>';
    container.innerHTML = duelsHtml;
}

function getDuelHtml() {
    const container2 = document.getElementById("PrepGrille");
    
    if (container2) {
        container2.innerHTML = `
            <div class="button_prep_grille">
                <button onclick="openDuelCreation()">Préparer une grille</button>
                <button onclick="importDuelFromServer()">Importer depuis le serveur</button>
            </div>
        `;
    }
    
    return `
        <div class="button_prep_grille">
            <button onclick="openDuelCreation()">Préparer une grille</button>
            <button onclick="importDuelFromServer()">Importer depuis le serveur</button>
        </div>
    `;
}

function generateDuelCard(duel, index) {
    return `
        <div class="duel-card" data-index="${index}">
            <h3>${duel.name}</h3>
            <div class="duel-points">
                ${generatePointSection(50, duel.points[50])}
                ${generatePointSection(40, duel.points[40])}
                ${generatePointSection(30, duel.points[30])}
                ${generatePointSection(20, duel.points[20])}
                ${generatePointSection(10, duel.points[10])}
                <div class="point-section same-song">
                    <span class="points">Même chanson</span>
                    <span class="song">Chanson: ${duel.sameSong.title} - ${duel.sameSong.artist}</span>
                </div>
            </div>
            <div class="duel-actions">
                <button onclick="editDuel(${index})">Modifier</button>
                <button onclick="deleteDuel(${index})">Supprimer</button>
                <button onclick="startDuel(${index})">Commencer le duel</button>
                <button onclick="exportDuel(${index})">Exporter JSON</button>
            </div>
        </div>
    `;
}

function generatePointSection(points, pointData) {
    return `
        <div class="point-section">
            <span class="points">${points} pts</span>
            <span class="theme">Thème: ${pointData.theme}</span>
            <div class="songs-options">
                <div class="song-option">
                    <span class="song">Option 1: ${pointData.songs[0].title} - ${pointData.songs[0].artist}</span>
                </div>
                <div class="song-option">
                    <span class="song">Option 2: ${pointData.songs[1].title} - ${pointData.songs[1].artist}</span>
                </div>
            </div>
        </div>
    `;
}

function openDuelCreation() {
    const creationHtml = `
        <div class="duel-creation-form">
            <h3>Créer une nouvelle grille de duel</h3>
            <h4>Nom du Duel</h4>
            <input type="text" id="duelName" placeholder="Nom de duel" required>
            <form id="duelForm" onsubmit="saveDuel(event)">
                
                <div class="points-configuration">
                    ${generatePointConfigSection(50)}
                    ${generatePointConfigSection(40)}
                    ${generatePointConfigSection(30)}
                    ${generatePointConfigSection(20)}
                    ${generatePointConfigSection(10)}
                    
                    <div class="point-config same-song-config">
                        <h4>Même Chanson</h4>
                        <div class="song-inputs">
                            <input type="text" id="same-song-title" placeholder="Titre de la chanson" required>
                            <input type="text" id="same-song-artist" placeholder="Artiste" required>
                            <input type="url" id="same-song-audio" placeholder="Lien audio (optionnel)">
                            <input type="text" id="same-song-lyrics" placeholder="Fichier paroles JSON (optionnel)">
                        </div>
                        <small>Cette chanson sera interprétée par les deux joueurs sauf si clochette du premier. </small>
                    </div>
                </div>
                
                <div class="form-actions">
                    <button type="button" onclick="cancelDuelCreation()">Annuler</button>
                    <button type="submit">Créer le duel</button>
                </div>
            </form>
        </div>
    `;
    
    const container = document.getElementsByClassName("Sectionduel")[0];
    container.innerHTML = creationHtml;
}

function generatePointConfigSection(points) {
    return `
        <div class="point-config">
            <h4>${points} Points</h4>
            <input type="text" id="theme-${points}" placeholder="Thème" required>
            
            <div class="song-pair">
                <h5>Chanson 1</h5>
                <input type="text" id="song1-title-${points}" placeholder="Titre" required>
                <input type="text" id="song1-artist-${points}" placeholder="Artiste" required>
                <input type="url" id="song1-audio-${points}" placeholder="Lien audio (optionnel)">
                <input type="text" id="song1-lyrics-${points}" placeholder="Fichier paroles JSON (optionnel)">
                
                <h5>Chanson 2</h5>
                <input type="text" id="song2-title-${points}" placeholder="Titre" required>
                <input type="text" id="song2-artist-${points}" placeholder="Artiste" required>
                <input type="url" id="song2-audio-${points}" placeholder="Lien audio (optionnel)">
                <input type="text" id="song2-lyrics-${points}" placeholder="Fichier paroles JSON (optionnel)">
            </div>
        </div>
    `;
}

function saveDuel(event) {
    event.preventDefault();
    
    const formData = {
        name: document.getElementById('duelName').value,
        points: {},
        sameSong: {
            title: document.getElementById('same-song-title').value,
            artist: document.getElementById('same-song-artist').value,
            audioUrl: document.getElementById('same-song-audio').value || null,
            lyricsFile: document.getElementById('same-song-lyrics').value || null
        },
        createdAt: new Date().toISOString()
    };
    
    [50, 40, 30, 20, 10].forEach(points => {
        formData.points[points] = {
            theme: document.getElementById(`theme-${points}`).value,
            songs: [
                {
                    title: document.getElementById(`song1-title-${points}`).value,
                    artist: document.getElementById(`song1-artist-${points}`).value,
                    audioUrl: document.getElementById(`song1-audio-${points}`).value || null,
                    lyricsFile: document.getElementById(`song1-lyrics-${points}`).value || null
                },
                {
                    title: document.getElementById(`song2-title-${points}`).value,
                    artist: document.getElementById(`song2-artist-${points}`).value,
                    audioUrl: document.getElementById(`song2-audio-${points}`).value || null,
                    lyricsFile: document.getElementById(`song2-lyrics-${points}`).value || null
                }
            ]
        };
    });
    
    preparedDuels.push(formData);
    
    if (typeof(Storage) !== "undefined") {
        localStorage.setItem('preparedDuels', JSON.stringify(preparedDuels));
    }
    
    saveDuelToServer(formData);
    
    createduel(preparedDuels);
    
   showNotification('duel créé avec succes')

}

function showNotification(message, type = 'success') {
    const notif = document.createElement('div');
    notif.className = `notification ${type}`;
    notif.textContent = message;

    document.body.appendChild(notif);
    setTimeout(() => {
        notif.remove();
    }, 3000);
}
async function saveDuelToServer(duelData) {
    try {
        const response = await fetch('/api/duels', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(duelData)
        });
        
        if (!response.ok) {
            throw new Error('Erreur lors de la sauvegarde sur le serveur');
        }
        
        console.log('Duel sauvegardé sur le serveur avec succès');
    } catch (error) {
        console.error('Erreur lors de la sauvegarde:', error);
    }
}

async function importDuelFromServer() {
    try {
        const response = await fetch('/api/duels');
        
        if (!response.ok) {
            throw new Error('Erreur lors du chargement depuis le serveur');
        }
        
        const serverDuels = await response.json();
        
        preparedDuels = [...preparedDuels, ...serverDuels];
        
        if (typeof(Storage) !== "undefined") {
            localStorage.setItem('preparedDuels', JSON.stringify(preparedDuels));
        }
        
        createduel(preparedDuels);
        showNotification('duel importé avec succes')        
    } catch (error) {
        console.error('Erreur lors de l\'importation:', error);
        showNotification('Erreur lors de l\'importation des duels');
    }
}

// Fonction pour exporter un duel en JSON
function exportDuel(index) {
    const duel = preparedDuels[index];
    const dataStr = JSON.stringify(duel, null, 2);
    const dataUri = 'data:application/json;charset=utf-8,'+ encodeURIComponent(dataStr);
    
    const exportFileDefaultName = `duel_${duel.name.replace(/\s+/g, '_')}.json`;
    
    const linkElement = document.createElement('a');
    linkElement.setAttribute('href', dataUri);
    linkElement.setAttribute('download', exportFileDefaultName);
    linkElement.click();
}

// Fonction pour annuler la création
function cancelDuelCreation() {
    createduel(preparedDuels);
}

// Fonction pour supprimer un duel
function deleteDuel(index) {
    if (confirm('Êtes-vous sûr de vouloir supprimer ce duel ?')) {
        preparedDuels.splice(index, 1);
        
        // Mettre à jour le localStorage
        if (typeof(Storage) !== "undefined") {
            localStorage.setItem('preparedDuels', JSON.stringify(preparedDuels));
        }
        
        createduel(preparedDuels);
    }
}

// Fonction pour éditer un duel
function editDuel(index) {
    const duel = preparedDuels[index];
    
    // Ouvrir le formulaire avec les données pré-remplies
    openDuelCreation();
    
    // Remplir les champs avec les données existantes
    setTimeout(() => {
        document.getElementById('duelName').value = duel.name;
        
        // Remplir les données pour chaque niveau de points
        [50, 40, 30, 20, 10].forEach(points => {
            document.getElementById(`theme-${points}`).value = duel.points[points].theme;
            document.getElementById(`song1-title-${points}`).value = duel.points[points].songs[0].title;
            document.getElementById(`song1-artist-${points}`).value = duel.points[points].songs[0].artist;
            document.getElementById(`song1-audio-${points}`).value = duel.points[points].songs[0].audioUrl || '';
            document.getElementById(`song1-lyrics-${points}`).value = duel.points[points].songs[0].lyricsFile || '';
            
            document.getElementById(`song2-title-${points}`).value = duel.points[points].songs[1].title;
            document.getElementById(`song2-artist-${points}`).value = duel.points[points].songs[1].artist;
            document.getElementById(`song2-audio-${points}`).value = duel.points[points].songs[1].audioUrl || '';
            document.getElementById(`song2-lyrics-${points}`).value = duel.points[points].songs[1].lyricsFile || '';
        });
        
        // Remplir les données de la chanson commune
        document.getElementById('same-song-title').value = duel.sameSong.title;
        document.getElementById('same-song-artist').value = duel.sameSong.artist;
        document.getElementById('same-song-audio').value = duel.sameSong.audioUrl || '';
        document.getElementById('same-song-lyrics').value = duel.sameSong.lyricsFile || '';
        
        // Modifier le formulaire pour l'édition
        const form = document.getElementById('duelForm');
        form.onsubmit = function(event) {
            updateDuel(event, index);
        };
        
        const submitButton = form.querySelector('button[type="submit"]');
        submitButton.textContent = 'Mettre à jour';
    }, 100);
}

// Fonction pour mettre à jour un duel
function updateDuel(event, index) {
    event.preventDefault();
    
    const formData = {
        name: document.getElementById('duelName').value,
        points: {},
        sameSong: {
            title: document.getElementById('same-song-title').value,
            artist: document.getElementById('same-song-artist').value,
            audioUrl: document.getElementById('same-song-audio').value || null,
            lyricsFile: document.getElementById('same-song-lyrics').value || null
        },
        createdAt: preparedDuels[index].createdAt,
        updatedAt: new Date().toISOString()
    };
    
    // Construire les données pour chaque niveau de points
    [50, 40, 30, 20, 10].forEach(points => {
        formData.points[points] = {
            theme: document.getElementById(`theme-${points}`).value,
            songs: [
                {
                    title: document.getElementById(`song1-title-${points}`).value,
                    artist: document.getElementById(`song1-artist-${points}`).value,
                    audioUrl: document.getElementById(`song1-audio-${points}`).value || null,
                    lyricsFile: document.getElementById(`song1-lyrics-${points}`).value || null
                },
                {
                    title: document.getElementById(`song2-title-${points}`).value,
                    artist: document.getElementById(`song2-artist-${points}`).value,
                    audioUrl: document.getElementById(`song2-audio-${points}`).value || null,
                    lyricsFile: document.getElementById(`song2-lyrics-${points}`).value || null
                }
            ]
        };
    });
    
    // Mettre à jour le duel
    preparedDuels[index] = formData;
    
    // Sauvegarder dans le localStorage
    if (typeof(Storage) !== "undefined") {
        localStorage.setItem('preparedDuels', JSON.stringify(preparedDuels));
    }
    
    createduel(preparedDuels);
    showNotification('Duel mis à jour avec succès !');
}

function startDuel(index) {
    const duel = preparedDuels[index];
    
    console.log('Démarrage du duel:', duel);
    
    window.location.href = `/duel-game?duelId=${index}`;
}

function loadDuels() {
    if (typeof(Storage) !== "undefined") {
        const saved = localStorage.getItem('preparedDuel');
        if (saved) {
            preparedDuels = JSON.parse(saved);
        }
    }
    return preparedDuels;
}

document.addEventListener('DOMContentLoaded', function() {
    loadDuels();
    createduel(preparedDuels);
});
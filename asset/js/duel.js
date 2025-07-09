let preparedDuels = []; /// ça définie la grille
let tempDuelData = null; // Pour stocker temporairement les données en cours de création

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
                <button onclick="showImportJsonForm()">Importer fichier JSON</button>
            </div>
        `;
    }
    
    return `
        <div class="button_prep_grille">
            <button onclick="openDuelCreation()">Préparer une grille</button>
            <button onclick="importDuelFromServer()">Importer depuis le serveur</button>
            <button onclick="showImportJsonForm()">Importer fichier JSON</button>
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
                <div class="Final">
                    <span class="points">Final</span>
                    <span class="song">Chanson: ${duel.finalSong.title} - ${duel.finalSong.artist}</span>
                </div>    
            </div>
            <div class="duel-actions">
                <button onclick="editDuel(${index})">Modifier</button>
                <button onclick="deleteDuel(${index})">Supprimer</button>
                <button onclick="startDuel(${index})">Commencer le duel</button>
                <button onclick="exportDuelToServer(${index})">Exporter vers serveur</button>
                <button onclick="exportDuelToFile(${index})">Exporter JSON local</button>
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
                    ${pointData.songs[0].lyricsFile ? '<span class="lyrics-indicator">📝</span>' : ''}
                </div>
                <div class="song-option">
                    <span class="song">Option 2: ${pointData.songs[1].title} - ${pointData.songs[1].artist}</span>
                    ${pointData.songs[1].lyricsFile ? '<span class="lyrics-indicator">📝</span>' : ''}
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
                            <input type="text" id="same-song-lyrics" placeholder="Fichier paroles (optionnel)">
                            <button type="button" onclick="checkLyricsFile('same-song-lyrics')">Vérifier paroles</button>
                        </div>
                        <small>Cette chanson sera interprétée par les deux joueurs sauf si clochette du premier. </small>
                    </div>
                </div>
                
                <div class="form-actions">
                    <button type="button" onclick="saveTemporaryDuel()">Sauvegarder temporairement</button>
                    <button type="button" onclick="loadTemporaryDuel()">Charger sauvegarde temporaire</button>
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
                <input type="text" id="song1-lyrics-${points}" placeholder="Fichier paroles (optionnel)">
                <button type="button" onclick="checkLyricsFile('song1-lyrics-${points}')">Vérifier paroles</button>
                
                <h5>Chanson 2</h5>
                <input type="text" id="song2-title-${points}" placeholder="Titre" required>
                <input type="text" id="song2-artist-${points}" placeholder="Artiste" required>
                <input type="url" id="song2-audio-${points}" placeholder="Lien audio (optionnel)">
                <input type="text" id="song2-lyrics-${points}" placeholder="Fichier paroles (optionnel)">
                <button type="button" onclick="checkLyricsFile('song2-lyrics-${points}')">Vérifier paroles</button>
            </div>
        </div>
    `;
}

// Fonction pour vérifier si un fichier de paroles existe sur le serveur
async function checkLyricsFile(inputId) {
    const input = document.getElementById(inputId);
    const filename = input.value.trim();
    
    if (!filename) {
        showNotification('Veuillez entrer un nom de fichier', 'warning');
        return;
    }
    
    try {
        const response = await fetch(`/api/check-lyrics?filename=${encodeURIComponent(filename)}`);
        
        if (response.ok) {
            const data = await response.json();
            if (data.exists) {
                showNotification(`Fichier "${filename}" trouvé sur le serveur`, 'success');
                input.style.backgroundColor = '#d4edda';
                input.style.borderColor = '#c3e6cb';
                
                // Optionnel: charger automatiquement le contenu des paroles
                if (data.content) {
                    console.log('Contenu des paroles:', data.content);
                }
            } else {
                showNotification(`Fichier "${filename}" non trouvé sur le serveur`, 'error');
                input.style.backgroundColor = '#f8d7da';
                input.style.borderColor = '#f5c6cb';
            }
        } else {
            showNotification('Erreur lors de la vérification du fichier', 'error');
        }
    } catch (error) {
        console.error('Erreur lors de la vérification:', error);
        showNotification('Erreur de connexion au serveur', 'error');
    }
}

// Fonction pour sauvegarder temporairement le duel en cours de création
async function saveTemporaryDuel() {
    const formData = collectFormData();
    if (!formData) return;
    
    tempDuelData = formData;
    
    try {
        const response = await fetch('/api/temp-duel', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(formData)
        });
        
        if (response.ok) {
            showNotification('Duel sauvegardé temporairement', 'success');
        } else {
            showNotification('Erreur lors de la sauvegarde temporaire', 'error');
        }
    } catch (error) {
        console.error('Erreur lors de la sauvegarde temporaire:', error);
        showNotification('Erreur de connexion au serveur', 'error');
    }
}

// Fonction pour charger la sauvegarde temporaire
async function loadTemporaryDuel() {
    try {
        const response = await fetch('/api/temp-duel');
        
        if (response.ok) {
            const tempData = await response.json();
            fillFormWithData(tempData);
            showNotification('Sauvegarde temporaire chargée', 'success');
        } else {
            showNotification('Aucune sauvegarde temporaire trouvée', 'warning');
        }
    } catch (error) {
        console.error('Erreur lors du chargement temporaire:', error);
        showNotification('Erreur de connexion au serveur', 'error');
    }
}

// Fonction pour collecter les données du formulaire
function collectFormData() {
    const name = document.getElementById('duelName').value;
    if (!name) {
        showNotification('Veuillez entrer un nom de duel', 'warning');
        return null;
    }
    
    const formData = {
        name: name,
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
    
    return formData;
}

// Fonction pour remplir le formulaire avec des données
function fillFormWithData(data) {
    document.getElementById('duelName').value = data.name;
    
    [50, 40, 30, 20, 10].forEach(points => {
        document.getElementById(`theme-${points}`).value = data.points[points].theme;
        document.getElementById(`song1-title-${points}`).value = data.points[points].songs[0].title;
        document.getElementById(`song1-artist-${points}`).value = data.points[points].songs[0].artist;
        document.getElementById(`song1-audio-${points}`).value = data.points[points].songs[0].audioUrl || '';
        document.getElementById(`song1-lyrics-${points}`).value = data.points[points].songs[0].lyricsFile || '';
        
        document.getElementById(`song2-title-${points}`).value = data.points[points].songs[1].title;
        document.getElementById(`song2-artist-${points}`).value = data.points[points].songs[1].artist;
        document.getElementById(`song2-audio-${points}`).value = data.points[points].songs[1].audioUrl || '';
        document.getElementById(`song2-lyrics-${points}`).value = data.points[points].songs[1].lyricsFile || '';
    });
    
    document.getElementById('same-song-title').value = data.sameSong.title;
    document.getElementById('same-song-artist').value = data.sameSong.artist;
    document.getElementById('same-song-audio').value = data.sameSong.audioUrl || '';
    document.getElementById('same-song-lyrics').value = data.sameSong.lyricsFile || '';
}

function saveDuel(event) {
    event.preventDefault();
    
    const formData = collectFormData();
    if (!formData) return;
    
    preparedDuels.push(formData);
    
    // Sauvegarde locale (localStorage)
    if (typeof(Storage) !== "undefined") {
        localStorage.setItem('preparedDuels', JSON.stringify(preparedDuels));
    }
    
    // Sauvegarde sur le serveur
    saveDuelToServer(formData);
    
    createduel(preparedDuels);
    
    showNotification('Duel créé avec succès', 'success');
}

function showNotification(message, type = 'success') {
    const notif = document.createElement('div');
    notif.className = `notification ${type}`;
    notif.textContent = message;
    
    // Styles pour les notifications
    notif.style.cssText = `
        position: fixed;
        top: 20px;
        right: 20px;
        padding: 15px 20px;
        border-radius: 5px;
        color: white;
        font-weight: bold;
        z-index: 1000;
        max-width: 300px;
        word-wrap: break-word;
    `;
    
    switch(type) {
        case 'success':
            notif.style.backgroundColor = '#28a745';
            break;
        case 'error':
            notif.style.backgroundColor = '#dc3545';
            break;
        case 'warning':
            notif.style.backgroundColor = '#ffc107';
            notif.style.color = '#212529';
            break;
        default:
            notif.style.backgroundColor = '#17a2b8';
    }

    document.body.appendChild(notif);
    setTimeout(() => {
        notif.remove();
    }, 4000);
}

// Fonction pour exporter un duel vers le serveur
async function exportDuelToServer(index) {
    const duel = preparedDuels[index];
    
    try {
        const response = await fetch('/api/export-duel', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(duel)
        });
        
        if (response.ok) {
            showNotification(`Duel "${duel.name}" exporté vers le serveur`, 'success');
        } else {
            const error = await response.text();
            showNotification(`Erreur lors de l'export: ${error}`, 'error');
        }
    } catch (error) {
        console.error('Erreur lors de l\'export:', error);
        showNotification('Erreur de connexion au serveur', 'error');
    }
}

// Fonction pour exporter un duel en fichier JSON local
function exportDuelToFile(index) {
    const duel = preparedDuels[index];
    const dataStr = JSON.stringify(duel, null, 2);
    const dataUri = 'data:application/json;charset=utf-8,'+ encodeURIComponent(dataStr);
    
    const exportFileDefaultName = `duel_${duel.name.replace(/\s+/g, '_')}.json`;
    
    const linkElement = document.createElement('a');
    linkElement.setAttribute('href', dataUri);
    linkElement.setAttribute('download', exportFileDefaultName);
    linkElement.click();
}

// Fonction pour afficher le formulaire d'import JSON
function showImportJsonForm() {
    const importHtml = `
        <div class="import-form">
            <h3>Importer un duel depuis un fichier JSON</h3>
            <form id="importForm" onsubmit="importDuelFromFile(event)">
                <div class="file-input-container">
                    <label for="duelFile">Sélectionner un fichier JSON:</label>
                    <input type="file" id="duelFile" accept=".json" required>
                </div>
                <div class="form-actions">
                    <button type="button" onclick="cancelImport()">Annuler</button>
                    <button type="submit">Importer</button>
                </div>
            </form>
        </div>
    `;
    
    const container = document.getElementsByClassName("Sectionduel")[0];
    container.innerHTML = importHtml;
}

// Fonction pour importer un duel depuis un fichier local
function importDuelFromFile(event) {
    event.preventDefault();
    
    const fileInput = document.getElementById('duelFile');
    const file = fileInput.files[0];
    
    if (!file) {
        showNotification('Veuillez sélectionner un fichier', 'warning');
        return;
    }
    
    const reader = new FileReader();
    reader.onload = function(e) {
        try {
            const duelData = JSON.parse(e.target.result);
            
            // Validation basique
            if (!duelData.name || !duelData.points || !duelData.sameSong) {
                throw new Error('Format de fichier invalide');
            }
            
            // Ajouter le duel importé
            preparedDuels.push(duelData);
            
            // Sauvegarder localement
            if (typeof(Storage) !== "undefined") {
                localStorage.setItem('preparedDuels', JSON.stringify(preparedDuels));
            }
            
            createduel(preparedDuels);
            showNotification(`Duel "${duelData.name}" importé avec succès`, 'success');
            
        } catch (error) {
            console.error('Erreur lors de l\'import:', error);
            showNotification('Erreur: fichier JSON invalide', 'error');
        }
    };
    
    reader.readAsText(file);
}

function cancelImport() {
    createduel(preparedDuels);
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
        showNotification('Erreur lors de la sauvegarde sur le serveur', 'error');
    }
}

async function importDuelFromServer() {
    try {
        const response = await fetch('/api/duels');
        
        if (!response.ok) {
            throw new Error('Erreur lors du chargement depuis le serveur');
        }
        
        const serverDuels = await response.json();
        
        // Éviter les doublons
        const existingNames = preparedDuels.map(duel => duel.name);
        const newDuels = serverDuels.filter(duel => !existingNames.includes(duel.name));
        
        if (newDuels.length > 0) {
            preparedDuels = [...preparedDuels, ...newDuels];
            
            if (typeof(Storage) !== "undefined") {
                localStorage.setItem('preparedDuels', JSON.stringify(preparedDuels));
            }
            
            createduel(preparedDuels);
            showNotification(`${newDuels.length} duel(s) importé(s) avec succès`, 'success');
        } else {
            showNotification('Aucun nouveau duel à importer', 'info');
        }
        
    } catch (error) {
        console.error('Erreur lors de l\'importation:', error);
        showNotification('Erreur lors de l\'importation des duels', 'error');
    }
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
        showNotification('Duel supprimé avec succès', 'success');
    }
}

// Fonction pour éditer un duel
function editDuel(index) {
    const duel = preparedDuels[index];
    
    // Ouvrir le formulaire avec les données pré-remplies
    openDuelCreation();
    
    // Remplir les champs avec les données existantes
    setTimeout(() => {
        fillFormWithData(duel);
        
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
    
    const formData = collectFormData();
    if (!formData) return;
    
    // Conserver les données de création originales
    formData.createdAt = preparedDuels[index].createdAt;
    formData.updatedAt = new Date().toISOString();
    
    // Mettre à jour le duel
    preparedDuels[index] = formData;
    
    // Sauvegarder dans le localStorage
    if (typeof(Storage) !== "undefined") {
        localStorage.setItem('preparedDuels', JSON.stringify(preparedDuels));
    }
    
    createduel(preparedDuels);
    showNotification('Duel mis à jour avec succès !', 'success');
}

function startDuel(index) {
    const duel = preparedDuels[index];
    
    console.log('Démarrage du duel:', duel);
    
    window.location.href = `/duel-game?duelId=${index}`;
}

function loadDuels() {
    if (typeof(Storage) !== "undefined") {
        const saved = localStorage.getItem('preparedDuels');
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
import * as duelLogic from './duel.js';

const DUEL_POINTS_CATEGORIES = [50, 40, 30, 20, 10];

function renderDuelList() {
    const container = document.querySelector(".Sectionduel");
    if (!container) return;

    if (duelLogic.preparedDuels.length === 0) {
        container.innerHTML = `
            <div class="alert alert-info">
                Aucune grille n'a été trouvée, veuillez en créer une via le bouton ci-dessous. 
            </div>
            ${getMenuHtml()}
        `;
        return;
    }

    let duelsHtml = getMenuHtml();
    duelsHtml += '<div class="duels-list">';
    duelLogic.preparedDuels.forEach((duel, index) => {
        duelsHtml += generateDuelCard(duel, index);
    });
    duelsHtml += '</div>';
    container.innerHTML = duelsHtml;
}

function getMenuHtml() {
    return `
        <div class="button_prep_grille">
            <button id="create-duel-btn">Préparer une grille</button>
            <button id="import-server-btn">Importer depuis le serveur</button>
            <button id="import-json-btn">Importer fichier JSON</button>
        </div>
    `;
}

function generateDuelCard(duel, index) {
    const pointsSections = DUEL_POINTS_CATEGORIES
        .map(p => generatePointSection(p, duel.points[p]))
        .join('');

    return `
        <div class="duel-card" data-index="${index}">
            <h3>${duel.name}</h3>
            <div class="duel-points" style="display: none;">
                ${pointsSections}
                <div class="point-section same-song">
                    <span class="points">Même chanson</span>
                    <span class="song">Chanson: ${duel.sameSong.title} - ${duel.sameSong.artist}</span>
                </div>
            </div>
            <div class="duel-actions">
                <a href="/duel-game?duelId=${duel.id}" class="start-btn" style="text-decoration: none; padding: 10px 20px; background-color: #19206bff; color: white; border-radius: 5px;">Commencer le duel</a>
                <div class="duel-buttons">
                    <button class="edit-btn">Modifier</button>
                    <button class="delete-btn">Supprimer</button>
                    <button class="export-server-btn">Exporter vers serveur</button>
                    <button class="export-file-btn">Exporter JSON local</button>
                </div>
            </div>
        </div>
    `;
}

function generatePointSection(points, pointData) {
    const lyricsIndicator = (song) => song.lyricsFile ? '<span class="lyrics-indicator">📝</span>' : '';
    return `
        <div class="point-section">
            <span class="points">${points} pts</span>
            <span class="theme">Thème: ${pointData.theme}</span>
            <div class="songs-options">
                <div class="song-option">
                    <span class="song">Option 1: ${pointData.songs[0].title} - ${pointData.songs[0].artist}</span>
                    ${lyricsIndicator(pointData.songs[0])}
                </div>
                <div class="song-option">
                    <span class="song">Option 2: ${pointData.songs[1].title} - ${pointData.songs[1].artist}</span>
                    ${lyricsIndicator(pointData.songs[1])}
                </div>
            </div>
        </div>
    `;
}


function openDuelForm(duel = null, index = -1) {
    const isEditing = duel !== null;
    const pointConfigSections = DUEL_POINTS_CATEGORIES
        .map(p => generatePointConfigSection(p, duel ? duel.points[p] : null))
        .join('');

    const formHtml = `
        <div class="duel-creation-form">
            <h3>${isEditing ? 'Modifier la grille de duel' : 'Créer une nouvelle grille de duel'}</h3>
            <form id="duelForm" data-editing-index="${index}">
                <h4>Nom du Duel</h4>
                <input type="text" id="duelName" placeholder="Nom de duel" value="${duel ? duel.name : ''}" required>
                <div class="points-configuration">${pointConfigSections}</div>
                <div class="point-config same-song-config">
                    <h4>Même Chanson</h4>
                    <small>Cette chanson sera interprétée par les deux joueurs sauf si clochette du premier.</small>
                    <div class="song-inputs">
                        <input type="text" id="same-song-title" placeholder="Titre de la chanson" value="${duel?.sameSong.title || ''}" required>
                        <input type="text" id="same-song-artist" placeholder="Artiste" value="${duel?.sameSong.artist || ''}" required>
                        <input type="url" id="same-song-audio" placeholder="Lien audio (optionnel)" value="${duel?.sameSong.audioUrl || ''}">
                        <input type="text" class="lyrics-input" id="same-song-lyrics" placeholder="Fichier paroles (optionnel)" value="${duel?.sameSong.lyricsFile || ''}">
                        <button type="button" class="check-lyrics-btn">Vérifier paroles</button>
                    </div>
                </div>
                <div class="form-actions">
                    <button type="button" id="save-temp-btn">Sauvegarder temporairement</button>
                    <button type="button" id="load-temp-btn">Charger sauvegarde temporaire</button>
                    <button type="button" id="cancel-btn">Annuler</button>
                    <button type="submit">${isEditing ? 'Mettre à jour' : 'Créer le duel'}</button>
                </div>
            </form>
        </div>
    `;
    document.querySelector(".Sectionduel").innerHTML = formHtml;
}

function generatePointConfigSection(points, data) {
    return `
        <div class="point-config">
            <h4>${points} Points</h4>
            <input type="text" id="theme-${points}" placeholder="Thème" value="${data?.theme || ''}" required>
            <div class="song-pair">
                ${generateSongInputFields(points, 1, data?.songs[0])}
                ${generateSongInputFields(points, 2, data?.songs[1])}
            </div>
        </div>
    `;
}

function generateSongInputFields(points, num, songData) {
    return `
        <h5>Chanson ${num}</h5>
        <input type="text" id="song${num}-title-${points}" placeholder="Titre" value="${songData?.title || ''}" required>
        <input type="text" id="song${num}-artist-${points}" placeholder="Artiste" value="${songData?.artist || ''}" required>
        <input type="url" id="song${num}-audio-${points}" placeholder="Lien audio (optionnel)" value="${songData?.audioUrl || ''}">
        <input type="text" class="lyrics-input" id="song${num}-lyrics-${points}" placeholder="Fichier paroles (optionnel)" value="${songData?.lyricsFile || ''}">
        <button type="button" class="check-lyrics-btn">Vérifier paroles</button>
    `;
}

function showImportJsonForm() {
    const importHtml = `
        <div class="import-form">
            <h3>Importer un duel depuis un fichier JSON</h3>
            <form id="importForm">
                <div class="file-input-container">
                    <label for="duelFile">Sélectionner un fichier JSON:</label>
                    <input type="file" id="duelFile" accept=".json" required>
                </div>
                <div class="form-actions">
                    <button type="button" id="cancel-import-btn">Annuler</button>
                    <button type="submit">Importer</button>
                </div>
            </form>
        </div>
    `;
    document.querySelector(".Sectionduel").innerHTML = importHtml;
}


function collectFormData() {
    const name = document.getElementById('duelName').value;
    if (!name) {
        showNotification('Veuillez entrer un nom de duel', 'warning');
        return null;
    }
    const formData = {
        id: Date.now(),
        name,
        points: {},
        sameSong: {
            title: document.getElementById('same-song-title').value,
            artist: document.getElementById('same-song-artist').value,
            audioUrl: document.getElementById('same-song-audio').value || null,
            lyricsFile: document.getElementById('same-song-lyrics').value || null
        },
        createdAt: new Date().toISOString()
    };

    for (const points of DUEL_POINTS_CATEGORIES) {
        formData.points[String(points)] = {
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
    }
    return formData;
}

function fillFormWithData(data) {
    if (!data) return;
    openDuelForm(data, -1); // Simplement re-render le form avec les données
    document.getElementById('duelName').value = data.name;
    for (const points of DUEL_POINTS_CATEGORIES) {
        document.getElementById(`theme-${points}`).value = data.points[points].theme;
        for (let i = 0; i < 2; i++) {
            document.getElementById(`song${i+1}-title-${points}`).value = data.points[points].songs[i].title;
            document.getElementById(`song${i+1}-artist-${points}`).value = data.points[points].songs[i].artist;
            document.getElementById(`song${i+1}-audio-${points}`).value = data.points[points].songs[i].audioUrl || '';
            document.getElementById(`song${i+1}-lyrics-${points}`).value = data.points[points].songs[i].lyricsFile || '';
        }
    }
    document.getElementById('same-song-title').value = data.sameSong.title;
    document.getElementById('same-song-artist').value = data.sameSong.artist;
    document.getElementById('same-song-audio').value = data.sameSong.audioUrl || '';
    document.getElementById('same-song-lyrics').value = data.sameSong.lyricsFile || '';
}

// --- Notifications ---

function showNotification(message, type = 'success') {
    const notif = document.createElement('div');
    notif.className = `notification ${type}`;
    notif.textContent = message;
    notif.style.cssText = `
        position: fixed; top: 20px; right: 20px; padding: 15px 20px;
        border-radius: 5px; color: white; font-weight: bold; z-index: 1000;
        max-width: 300px; word-wrap: break-word;
    `;
    const colors = { success: '#28a745', error: '#dc3545', warning: '#ffc107', info: '#17a2b8' };
    notif.style.backgroundColor = colors[type] || colors.info;
    if (type === 'warning') notif.style.color = '#212529';
    document.body.appendChild(notif);
    setTimeout(() => notif.remove(), 4000);
}

//Gestionnaires d'événements

document.addEventListener('click', async (event) => {
    const target = event.target;
    const duelCard = target.closest('.duel-card');
    const index = duelCard ? parseInt(duelCard.dataset.index, 10) : -1;

    // Boutons dans la liste des duels
    if (target.matches('.delete-btn')) {
        if (confirm('Êtes-vous sûr de vouloir supprimer ce duel ?')) {
            duelLogic.deleteDuel(index);
            renderDuelList();
            showNotification('Duel supprimé avec succès', 'success');
        }
    } else if (target.matches('.edit-btn')) {
        openDuelForm(duelLogic.preparedDuels[index], index);
    } else if (target.matches('.export-file-btn')) {
        const duel = duelLogic.preparedDuels[index];
        const dataStr = JSON.stringify(duel, null, 2);
        const dataUri = 'data:application/json;charset=utf-8,' + encodeURIComponent(dataStr);
        const linkElement = document.createElement('a');
        linkElement.setAttribute('href', dataUri);
        linkElement.setAttribute('download', `duel_${duel.name.replace(/\s+/g, '_')}.json`);
        linkElement.click();
    } else if (target.matches('.export-server-btn')) {
        try {
            const message = await duelLogic.exportDuelToServer(index);
            showNotification(message, 'success');
        } catch (error) {
            showNotification(`Erreur lors de l'export: ${error.message}`, 'error');
        }
    }

    // Boutons du menu principal
    else if (target.matches('#create-duel-btn')) {
        openDuelForm();
    } else if (target.matches('#import-json-btn')) {
        showImportJsonForm();
    } else if (target.matches('#import-server-btn')) {
        try {
            const count = await duelLogic.importDuelsFromServer();
            if (count > 0) {
                renderDuelList();
                showNotification(`${count} duel(s) importé(s) avec succès`, 'success');
            } else {
                showNotification('Aucun nouveau duel à importer', 'info');
            }
        } catch (error) {
            showNotification(error.message, 'error');
        }
    }

    // Boutons du formulaire
    else if (target.matches('#cancel-btn') || target.matches('#cancel-import-btn')) {
        renderDuelList();
    } else if (target.matches('.check-lyrics-btn')) {
        const input = target.previousElementSibling;
        const filename = input.value.trim();
        try {
            const data = await duelLogic.checkLyricsFile(filename);
            if (data.exists) {
                showNotification(`Fichier "${filename}" trouvé`, 'success');
                input.style.backgroundColor = '#d4edda';
            } else {
                showNotification(`Fichier "${filename}" non trouvé`, 'error');
                input.style.backgroundColor = '#f8d7da';
            }
        } catch (error) {
            showNotification(error.message, 'error');
        }
    } else if (target.matches('#save-temp-btn')) {
        const formData = collectFormData();
        if (formData) {
            try {
                await duelLogic.saveTemporaryDuel(formData);
                showNotification('Duel sauvegardé temporairement', 'success');
            } catch (error) {
                showNotification(error.message, 'error');
            }
        }
    } else if (target.matches('#load-temp-btn')) {
        try {
            const tempData = await duelLogic.loadTemporaryDuel();
            fillFormWithData(tempData);
            showNotification('Sauvegarde temporaire chargée', 'success');
        } catch (error) {
            showNotification(error.message, 'warning');
        }
    }
});


document.addEventListener('submit', async (event) => { 
    event.preventDefault();
    const target = event.target;

    if (target.matches('#duelForm')) {
        const formData = collectFormData();
        if (!formData) return;

        const editingIndex = parseInt(target.dataset.editingIndex, 10);
        if (editingIndex > -1) {
            duelLogic.updateDuel(editingIndex, formData);
            showNotification('Duel mis à jour avec succès', 'success');
            renderDuelList();
        } else {
            try {
                await duelLogic.addDuel(formData); 
                showNotification('Duel créé avec succès', 'success');
                renderDuelList();
            } catch (error) {
                showNotification(`Erreur lors de la création : ${error.message}`, 'error');
            }
        }
    } else if (target.matches('#importForm')) {
        const fileInput = document.getElementById('duelFile');
        const file = fileInput.files[0];
        if (!file) {
            showNotification('Veuillez sélectionner un fichier', 'warning');
            return;
        }
        const reader = new FileReader();
        reader.onload = async function(e) { 
            try {
                const duelData = JSON.parse(e.target.result);
                if (!duelData.name || !duelData.points || !duelData.sameSong) {
                    throw new Error('Format de fichier invalide');
                }
                
                await duelLogic.addDuel(duelData); 
                
                renderDuelList();
                showNotification(`Duel "${duelData.name}" importé avec succès`, 'success');
            } catch (error) {
                showNotification(`Erreur: ${error.message}`, 'error');
            }
        };
        reader.readAsText(file);
    }
});
document.addEventListener('DOMContentLoaded', () => {
    duelLogic.loadDuelsFromStorage();
    renderDuelList();
});
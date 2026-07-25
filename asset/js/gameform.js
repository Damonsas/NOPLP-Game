var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
import { showNotification } from './gamenotification.js';
import { addDuel, loadDuelsFromStorage, preparedDuels } from './gamelogic.js';
const DUEL_POINTS_CATEGORIES = [50, 40, 30, 20, 10];
let editingDuelId = null;
function getLyricsListLocal() {
    return __awaiter(this, void 0, void 0, function* () {
        const indexPath = './data/serverdata/paroledata/index.json';
        try {
            console.log("Tentative de récupération des fichiers lyrics depuis", indexPath);
            const response = yield fetch(indexPath, { cache: 'no-store' });
            if (!response.ok) {
                console.warn(`index.json introuvable ou erreur (${response.status})`);
                return [];
            }
            const arr = yield response.json();
            if (!Array.isArray(arr)) {
                console.warn("Format inattendu pour index.json (pas un tableau)");
                return [];
            }
            const files = arr.map((item) => {
                if (!item)
                    return null;
                const raw = item.ligne || (item.artiste && item.titre ? `${item.artiste} - ${item.titre}` : null);
                if (!raw)
                    return null;
                return raw.endsWith('.json') ? raw : `${raw}.json`;
            }).filter(Boolean);
            return files;
        }
        catch (err) {
            console.error("Erreur lors de la récupération via index.json:", err);
            return [];
        }
    });
}
function isSoloMode() {
    return window.location.pathname.includes('solo');
}
function generateDuelCard(duel) {
    return `
    <div class="duel-card" id="duel-card-${duel.id}" data-duel-id="${duel.id}">
      <h3>${duel.name}</h3>
      <div class="duel-actions">
        <button type="button" class="play-duel-btn btn " data-duel-id="${duel.id}">Jouer</button>
        <button type="button" class="edit-duel-btn btn " data-duel-id="${duel.id}">
            Modifier <i class="fa-regular fa-pen-to-square"></i>
        </button>
        <button type="button" class="delete-duel-btn btn " data-duel-id="${duel.id}">
            Supprimer <i class="fa-regular fa-trash-can"></i>
        </button>
      </div>
    </div>
  `;
}
function handlePlayDuel(duelId) {
    return __awaiter(this, void 0, void 0, function* () {
        try {
            const id = Number.parseInt(duelId, 10);
            if (Number.isNaN(id))
                throw new Error('ID invalide');
            const duels = JSON.parse(localStorage.getItem('duels') || '[]');
            const duel = duels.find((d) => d.id === id);
            if (!duel)
                throw new Error('Duel non trouvé');
            const res = yield fetch('/api/duels', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(duel)
            });
            if (!res.ok)
                throw new Error(`Erreur serveur: ${res.status}`);
            const serverDuel = yield res.json();
            const url = isSoloMode() ? `/solo?id=${serverDuel.id}` : `/duel-game?id=${serverDuel.id}`;
            window.location.href = url;
        }
        catch (error) {
            showNotification(`Erreur: ${error}`, 'error');
        }
    });
}
function handleDeleteDuel(duelId) {
    return __awaiter(this, void 0, void 0, function* () {
        if (!confirm("Voulez-vous vraiment supprimer cette grille ?"))
            return;
        const id = Number.parseInt(duelId, 10);
        try {
            const res = yield fetch(`/api/duels/${id}`, { method: 'DELETE' });
            if (!res.ok) {
                console.warn(`L'API a répondu ${res.status}, suppression locale en cours...`);
            }
            const duels = JSON.parse(localStorage.getItem('duels') || '[]');
            const updatedDuels = duels.filter(d => d.id !== id);
            localStorage.setItem('duels', JSON.stringify(updatedDuels));
            const indexInMemory = preparedDuels.findIndex(d => d.id === id);
            if (indexInMemory !== -1) {
                preparedDuels.splice(indexInMemory, 1);
            }
            const cardEl = document.getElementById(`duel-card-${id}`);
            if (cardEl) {
                cardEl.remove();
            }
            else {
                renderDuelList();
            }
            showNotification("Grille supprimée avec succès !", "success");
            if (preparedDuels.length === 0) {
                renderDuelList();
            }
        }
        catch (error) {
            console.error("Erreur lors de la suppression:", error);
            showNotification("Erreur lors de la suppression.", "error");
        }
    });
}
function handleEditDuel(duelId) {
    return __awaiter(this, void 0, void 0, function* () {
        const id = Number.parseInt(duelId, 10);
        const duels = JSON.parse(localStorage.getItem('duels') || '[]');
        const duel = duels.find(d => d.id === id);
        if (!duel) {
            showNotification("Impossible de trouver la grille à modifier", "error");
            return;
        }
        editingDuelId = id;
        showCreateForm();
        setTimeout(() => {
            const nameInput = document.getElementById('duelName');
            const sameSongSelect = document.getElementById('sameSongFile');
            if (nameInput)
                nameInput.value = duel.name;
            if (sameSongSelect && duel.sameSong)
                sameSongSelect.value = duel.sameSong.lyricsFile || '';
            if (duel.points) {
                DUEL_POINTS_CATEGORIES.forEach(pts => {
                    var _a, _b;
                    const categoryData = duel.points[pts];
                    if (categoryData) {
                        const themeInput = document.querySelector(`input[name="theme-${pts}"]`);
                        const song1Select = document.querySelector(`select[name="song1-${pts}"]`);
                        const song2Select = document.querySelector(`select[name="song2-${pts}"]`);
                        if (themeInput && categoryData.theme)
                            themeInput.value = categoryData.theme;
                        if (song1Select && ((_a = categoryData.songs) === null || _a === void 0 ? void 0 : _a[0]))
                            song1Select.value = categoryData.songs[0].lyricsFile || '';
                        if (song2Select && ((_b = categoryData.songs) === null || _b === void 0 ? void 0 : _b[1]))
                            song2Select.value = categoryData.songs[1].lyricsFile || '';
                    }
                });
            }
            const formTitle = document.querySelector('#newDuelForm h3');
            if (formTitle)
                formTitle.textContent = "Modifier la grille";
        }, 100);
    });
}
function getMenuHtml() {
    return `<div class="button_prep_grille"><button id="create-duel-btn">Préparer une grille</button></div>`;
}
function renderDuelList() {
    const container = document.querySelector(".Sectionduel");
    if (!container)
        return;
    const existingPrepGrille = document.getElementById("PrepGrille");
    const prepGrilleContent = existingPrepGrille ? existingPrepGrille.outerHTML : '<div id="PrepGrille" style="display: none;"></div>';
    if (preparedDuels.length === 0) {
        container.innerHTML = `
      <div class="alert alert-info">
        Aucune grille n'a été trouvée, veuillez en créer une via le bouton ci-dessous.
      </div>
      ${getMenuHtml()}
      ${prepGrilleContent}
    `;
        return;
    }
    let duelsHtml = getMenuHtml();
    duelsHtml += '<div class="duels-list">';
    preparedDuels.forEach(duel => {
        duelsHtml += generateDuelCard(duel);
    });
    duelsHtml += '</div>';
    duelsHtml += prepGrilleContent;
    container.innerHTML = duelsHtml;
}
function generateSongSelectionHtml(points, lyricsFiles) {
    const soloMode = isSoloMode();
    const songOptions = lyricsFiles.map(file => `<option style="color: black" value="${file}">${file}</option>`).join('');
    if (soloMode) {
        return `
      <label>Chanson:</label>
      <select name="song1-${points}" required>
        <option style="color: black" value="">Sélectionner une chanson</option>
        ${songOptions}
      </select>
    `;
    }
    else {
        return `
      <label>Chanson 1:</label>
      <select name="song1-${points}" required>
        <option style="color: black" value="">Sélectionner une chanson</option>
        ${songOptions}
      </select>
      <label>Chanson 2:</label>
      <select name="song2-${points}" required>
        <option style="color: black" value="">Sélectionner une chanson</option>
        ${songOptions}
      </select>
    `;
    }
}
function attachUniqueSelectionHandlers(formOrContainer) {
    if (!formOrContainer)
        return;
    const songSelects = Array.from(formOrContainer.querySelectorAll('select[name^="sameSongFile"], select[name^="song1-"], select[name^="song2-"]'));
    function refreshDisabledOptions() {
        const selectedValues = new Set(songSelects
            .map(s => s.value)
            .filter(v => v && v.length > 0));
        songSelects.forEach(select => {
            const ownValue = select.value;
            Array.from(select.options).forEach(opt => {
                if (opt.value === ownValue) {
                    opt.disabled = false;
                    return;
                }
                opt.disabled = selectedValues.has(opt.value);
            });
        });
    }
    songSelects.forEach(select => {
        select.addEventListener('change', refreshDisabledOptions);
        select.addEventListener('input', refreshDisabledOptions);
    });
    refreshDisabledOptions();
}
function renderCreateDuelForm(lyricsFiles) {
    const container = document.getElementById("PrepGrille");
    if (!container)
        return;
    const soloMode = isSoloMode();
    const modeText = soloMode ? 'solo' : 'duel';
    const songOptions = lyricsFiles.map(file => `<option style="color: black" value="${file}">${file}</option>`).join('');
    let formHtml = `
    <div class="form-container">
      <h2 style="color: red;">Choisissez vos chansons</h2>
      <button id="back-to-list-btn" type="button">← Retour à la liste</button>
      <form id="newDuelForm">
        <h3>Créer une nouvelle grille de ${modeText}</h3>
        <label for="duelName">Nom de la grille:</label>
        <input type="text" id="duelName" name="duelName" required>

        <label for="sameSongFile">Sélectionner la chanson unique ("La Même Chanson") :</label>
        <select name="sameSongFile" id="sameSongFile" required>
          <option style="color: black" value="">Choisir la même chanson</option>
          ${songOptions}
        </select>
  `;
    DUEL_POINTS_CATEGORIES.forEach(points => {
        formHtml += `
      <div class="point-category">
        <h4>${points} Points</h4>
        <label>Thème:</label>
        <input type="text" name="theme-${points}" required>
        ${generateSongSelectionHtml(points, lyricsFiles)}
      </div>
    `;
    });
    formHtml += `
        <button type="submit">${editingDuelId ? 'Enregistrer les modifications' : 'Créer'}</button>
      </form>
    </div>
  `;
    container.innerHTML = formHtml;
    attachUniqueSelectionHandlers(document.getElementById('newDuelForm'));
}
function handleNewDuelFormSubmit(event) {
    return __awaiter(this, void 0, void 0, function* () {
        event.preventDefault();
        const form = event.target;
        const formData = new FormData(form);
        const duelData = { points: {} };
        const soloMode = isSoloMode();
        const duelPoints = duelData.points;
        const sameSongFileValue = formData.get('sameSongFile');
        for (const [key, value] of formData) {
            const parts = key.split('-');
            const fieldName = parts[0];
            const points = parts.length > 1 ? parts[1] : null;
            if (fieldName === 'duelName') {
                duelData.name = value;
            }
            else if (points) {
                const pointsKey = points;
                if (!duelPoints[pointsKey]) {
                    duelPoints[pointsKey] = {};
                }
                if (fieldName === 'theme') {
                    duelPoints[pointsKey].theme = value;
                }
                else if (fieldName.startsWith('song')) {
                    const songIndex = fieldName === 'song1' ? 0 : 1;
                    if (!duelPoints[pointsKey].songs) {
                        duelPoints[pointsKey].songs = soloMode ? [{}] : [{}, {}];
                    }
                    if (soloMode && songIndex === 0) {
                        duelPoints[pointsKey].songs[0] = { lyricsFile: value };
                    }
                    else if (!soloMode) {
                        duelPoints[pointsKey].songs[songIndex] = { lyricsFile: value };
                    }
                }
            }
        }
        const newDuel = {
            id: editingDuelId || Date.now(),
            name: duelData.name,
            points: duelData.points,
            sameSong: { title: 'N/A', artist: 'N/A', lyricsFile: sameSongFileValue },
            createdAt: new Date().toISOString()
        };
        try {
            yield addDuel(newDuel);
            showNotification(`Grille ${soloMode ? 'solo' : 'duel'} sauvegardée !`, 'success');
            editingDuelId = null;
            showDuelList();
        }
        catch (error) {
            if (error instanceof Error) {
                showNotification(`Erreur lors de la création : ${error.message}`, 'error');
            }
            else {
                showNotification('Une erreur inconnue est survenue.', 'error');
            }
        }
    });
}
function showCreateForm() {
    const formContainer = document.getElementById("PrepGrille");
    const listContent = document.querySelector('.duels-list');
    const alertContent = document.querySelector('.alert');
    const menuButton = document.querySelector('.button_prep_grille');
    if (formContainer && (!formContainer.innerHTML || formContainer.innerHTML.trim() === '')) {
        getLyricsListLocal().then(lyricsFiles => {
            if (lyricsFiles.length > 0) {
                renderCreateDuelForm(lyricsFiles);
            }
            showFormWithStyles(formContainer);
        }).catch(error => {
            console.error("Erreur lors de la récupération des fichiers:", error);
            showFormWithStyles(formContainer);
        });
    }
    else {
        showFormWithStyles(formContainer);
    }
    if (listContent)
        listContent.style.display = 'none';
    if (alertContent)
        alertContent.style.display = 'none';
    if (menuButton)
        menuButton.style.display = 'none';
}
function showFormWithStyles(formContainer) {
    if (formContainer) {
        formContainer.style.display = 'block';
        formContainer.style.visibility = 'visible';
        formContainer.style.opacity = '1';
        formContainer.style.height = 'auto';
        formContainer.style.position = 'relative';
        formContainer.style.zIndex = '1000';
    }
}
function showDuelList() {
    const formContainer = document.getElementById("PrepGrille");
    const alertContent = document.querySelector('.alert');
    const menuButton = document.querySelector('.button_prep_grille');
    if (formContainer) {
        formContainer.style.display = 'none';
    }
    if (alertContent) {
        alertContent.style.display = 'block';
    }
    if (menuButton) {
        menuButton.style.display = 'block';
    }
    editingDuelId = null;
    loadDuelsFromStorage();
    renderDuelList();
}
function handleImportFormSubmit(event) {
    return __awaiter(this, void 0, void 0, function* () {
        var _a;
        event.preventDefault();
        const fileInput = document.getElementById('duelFile');
        const file = (_a = fileInput.files) === null || _a === void 0 ? void 0 : _a[0];
        if (!file) {
            showNotification('Veuillez sélectionner un fichier', 'warning');
            return;
        }
        const reader = new FileReader();
        reader.onload = function (e) {
            return __awaiter(this, void 0, void 0, function* () {
                var _a;
                try {
                    const result = (_a = e.target) === null || _a === void 0 ? void 0 : _a.result;
                    if (typeof result !== 'string') {
                        throw new TypeError('Le contenu du fichier n\'est pas une chaîne de caractères.');
                    }
                    const duelData = JSON.parse(result);
                    yield addDuel(duelData);
                    renderDuelList();
                    showNotification(`Duel "${duelData.name}" importé avec succès`, 'success');
                }
                catch (error) {
                    if (error instanceof Error) {
                        showNotification(`Erreur: ${error.message}`, 'error');
                    }
                    else {
                        showNotification('Une erreur inconnue est survenue.', 'error');
                    }
                }
            });
        };
        reader.readAsText(file);
    });
}
document.addEventListener('submit', (event) => {
    const target = event.target;
    if (target.id === 'newDuelForm') {
        handleNewDuelFormSubmit(event);
    }
    else if (target.id === 'importForm') {
        handleImportFormSubmit(event);
    }
});
document.addEventListener('click', (event) => {
    const target = event.target;
    const playBtn = target.closest('.play-duel-btn');
    if (playBtn) {
        event.preventDefault();
        const id = playBtn.getAttribute('data-duel-id');
        if (id)
            handlePlayDuel(id);
        return;
    }
    const deleteBtn = target.closest('.delete-duel-btn');
    if (deleteBtn) {
        event.preventDefault();
        const id = deleteBtn.getAttribute('data-duel-id');
        if (id)
            handleDeleteDuel(id);
        return;
    }
    const editBtn = target.closest('.edit-duel-btn');
    if (editBtn) {
        event.preventDefault();
        const id = editBtn.getAttribute('data-duel-id');
        if (id)
            handleEditDuel(id);
        return;
    }
    if (target.id === 'create-duel-btn') {
        event.preventDefault();
        editingDuelId = null;
        showCreateForm();
    }
    else if (target.id === 'back-to-list-btn') {
        event.preventDefault();
        showDuelList();
    }
});
document.addEventListener('DOMContentLoaded', () => __awaiter(void 0, void 0, void 0, function* () {
    loadDuelsFromStorage();
    renderDuelList();
    setTimeout(() => __awaiter(void 0, void 0, void 0, function* () {
        let prepGrilleContainer = document.getElementById("PrepGrille");
        if (!prepGrilleContainer)
            return;
        try {
            const lyricsFiles = yield getLyricsListLocal();
            if (lyricsFiles.length > 0) {
                renderCreateDuelForm(lyricsFiles);
            }
        }
        catch (error) {
            console.error("Erreur d'initialisation des paroles:", error);
        }
    }), 200);
}));

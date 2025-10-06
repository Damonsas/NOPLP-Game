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
/**
 * Récupère la liste des fichiers de paroles locaux
 * @returns Promise avec la liste des noms de fichiers
 */
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
            // index.json est un tableau d'objets
            const arr = yield response.json();
            if (!Array.isArray(arr)) {
                console.warn("Format inattendu pour index.json (pas un tableau)");
                return [];
            }
            // On transforme chaque entrée en nom de fichier utilisable.
            // Priorité: si "ligne" existe on l'utilise, sinon on compose "Artiste - Titre".
            const files = arr.map((item) => {
                if (!item)
                    return null;
                // Nettoyage basique pour éviter slash etc.
                const raw = item.ligne || (item.artiste && item.titre ? `${item.artiste} - ${item.titre}` : null);
                if (!raw)
                    return null;
                // Ajout de l'extension .json si elle n'est pas déjà présente
                return raw.endsWith('.json') ? raw : `${raw}.json`;
            }).filter(Boolean);
            console.log("Fichiers lyrics trouvés via index.json:", files);
            return files;
        }
        catch (err) {
            console.error("Erreur lors de la récupération via index.json:", err);
            return [];
        }
    });
}
/**
 * Charge un fichier de paroles spécifique
 * @param filename Nom du fichier à charger
 * @returns Promise avec le contenu du fichier JSON
 */
function loadLyricsFile(filename) {
    return __awaiter(this, void 0, void 0, function* () {
        try {
            const response = yield fetch(`./data/paroledata/${filename}`);
            if (!response.ok) {
                throw new Error(`Impossible de charger ${filename}: ${response.status}`);
            }
            const data = yield response.json();
            console.log(`Fichier ${filename} chargé avec succès`);
            return data;
        }
        catch (error) {
            console.error(`Erreur lors du chargement de ${filename}:`, error);
            throw error;
        }
    });
}
/**
 * Détermine si on est en mode solo ou duel basé sur l'URL actuelle
 * @returns true si mode solo, false si mode duel
 */
function isSoloMode() {
    return window.location.pathname.includes('solo');
}
// === FONCTIONS DE RENDU (UI) ===
/**
 * Génère la carte HTML pour un duel donné.
 * @param duel Le duel à afficher.
 * @returns La chaîne HTML de la carte.
 */
function generateDuelCard(duel) {
    const themes = DUEL_POINTS_CATEGORIES.map(p => { var _a; return ((_a = duel.points[p.toString()]) === null || _a === void 0 ? void 0 : _a.theme); }).filter(t => t).join(', ');
    const currentMode = isSoloMode() ? 'solo' : 'duel';
    return `
    <div class="duel-card" data-duel-id="${duel.id}">
      <h3>${duel.name}</h3>
      <button onclick="window.location.href='/${currentMode}?id=${duel.id}'">Jouer</button>
    </div>
  `;
}
/**
 * Génère le bouton "Préparer une grille".
 * @returns La chaîne HTML du bouton.
 */
function getMenuHtml() {
    return `<div class="button_prep_grille"><button id="create-duel-btn">Préparer une grille</button></div>`;
}
/**
 * Affiche la liste des duels disponibles ou le message d'absence de duel.
 */
function renderDuelList() {
    const container = document.querySelector(".Sectionduel");
    if (!container)
        return;
    // Récupérer le contenu existant de PrepGrille pour le préserver
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
/**
 * Génère les champs de sélection de chansons selon le mode
 * @param points Points de la catégorie
 * @param lyricsFiles Liste des fichiers de paroles
 * @returns HTML des champs de sélection
 */
function generateSongSelectionHtml(points, lyricsFiles) {
    const soloMode = isSoloMode();
    const songOptions = lyricsFiles.map(file => `<option value="${file}">${file}</option>`).join('');
    if (soloMode) {
        return `
      <label>Chanson:</label>
      <select name="song1-${points}" required>
        <option value="">Sélectionner une musique</option>
        ${songOptions}

      </select>
    `;
    }
    else {
        return `
      <label>Chanson 1:</label>
      <select name="song1-${points}" required>
        <option value="">Sélectionner une musique</option>
        ${songOptions}
      </select>
      <label>Chanson 2:</label>
      <select name="song2-${points}" required>
        <option value="" >Sélectionner une musique</option>
        ${songOptions}

      </select>
    `;
    }
}
/**
 * Ajoute des listeners sur les selects de chansons du formulaire pour :
 *  - désactiver automatiquement les options déjà choisies dans les autres selects
 *  - autoriser la réactivation si une sélection est modifiée / supprimée
 *
 * Appelle : attachUniqueSelectionHandlers(document.getElementById('newDuelForm'));
 */
function attachUniqueSelectionHandlers(formOrContainer) {
    if (!formOrContainer)
        return;
    // Sélecteurs des selects de chansons (nom commençant par "song1-" ou "song2-")
    const songSelects = Array.from(formOrContainer.querySelectorAll('select[name^="song1-"], select[name^="song2-"]'));
    // Helper pour rafraîchir l'état disabled de toutes les options en fonction des valeurs sélectionnées
    function refreshDisabledOptions() {
        // valeurs actuellement sélectionnées (non vides)
        const selectedValues = songSelects
            .map(s => s.value)
            .filter(v => v && v.length > 0);
        // Pour chaque select, on active toutes les options puis on désactive celles qui sont sélectionnées ailleurs
        songSelects.forEach(select => {
            const ownValue = select.value;
            Array.from(select.options).forEach(opt => {
                // Toujours permettre la valeur courante du select (pour ne pas "bloquer" la sélection en cours)
                if (opt.value === ownValue) {
                    opt.disabled = false;
                    return;
                }
                // Désactive l'option si elle est choisie dans une autre select
                opt.disabled = selectedValues.includes(opt.value);
            });
        });
    }
    // Ajoute l'écoute "change" sur chaque select
    songSelects.forEach(select => {
        // Si l'utilisateur supprime la sélection (option vide), on gère aussi
        select.addEventListener('change', () => {
            refreshDisabledOptions();
        });
        // Add a keyboard-clear detection (optional) to re-enable options when cleared
        select.addEventListener('input', () => {
            refreshDisabledOptions();
        });
    });
    // Initial refresh (si certains selects ont déjà une valeur)
    refreshDisabledOptions();
}
/**
 * Génère et affiche le formulaire de création de duel avec les listes de musiques dynamiques.
 * @param lyricsFiles La liste des noms de fichiers de paroles.
 */
function renderCreateDuelForm(lyricsFiles) {
    const container = document.getElementById("PrepGrille");
    if (!container) {
        console.error("Container PrepGrille non trouvé");
        return;
    }
    const soloMode = isSoloMode();
    const modeText = soloMode ? 'solo' : 'duel';
    let formHtml = `
    <div class="form-container">
      <h2 style="color: red;">Choisissez vos musiques</h2>
      <button id="back-to-list-btn" type="button">← Retour à la liste</button>
      <form id="newDuelForm">
        <h3>Créer une nouvelle grille de ${modeText}</h3>
        <label for="duelName">Nom de la grille:</label>
        <input type="text" id="duelName" name="duelName" required >
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
        <button type="submit">Créer</button>
      </form>
    </div>
  `;
    container.innerHTML = formHtml;
    attachUniqueSelectionHandlers(document.getElementById('newDuelForm'));
    console.log("Formulaire généré avec styles");
    console.log("Contenu HTML du formulaire:", formHtml.substring(0, 300) + "...");
}
// === GESTIONNAIRES D'ÉVÉNEMENTS ===
/**
 * Gère la soumission du formulaire de création de duel.
 * @param event L'événement de soumission.
 */
function handleNewDuelFormSubmit(event) {
    return __awaiter(this, void 0, void 0, function* () {
        event.preventDefault();
        const form = event.target;
        const formData = new FormData(form);
        const duelData = {};
        const soloMode = isSoloMode();
        for (const [key, value] of formData) {
            const parts = key.split('-');
            const fieldName = parts[0];
            const points = parts.length > 1 ? parts[1] : null;
            if (fieldName === 'duelName') {
                duelData.name = value;
            }
            else if (points) {
                if (!duelData.points) {
                    duelData.points = {};
                }
                if (!duelData.points[points]) {
                    duelData.points[points] = {};
                }
                if (fieldName === 'theme') {
                    duelData.points[points].theme = value;
                }
                else if (fieldName.startsWith('song')) {
                    const songIndex = fieldName === 'song1' ? 0 : 1;
                    if (!duelData.points[points].songs) {
                        // En mode solo, on ne crée qu'une chanson
                        duelData.points[points].songs = soloMode ? [{}] : [{}, {}];
                    }
                    if (soloMode && songIndex === 0) {
                        // En mode solo, on ne stocke qu'une chanson
                        duelData.points[points].songs[0] = { lyricsFile: value };
                    }
                    else if (!soloMode) {
                        // En mode duel, on stocke les deux chansons
                        duelData.points[points].songs[songIndex] = { lyricsFile: value };
                    }
                }
            }
        }
        const allSongSelects = Array.from(form.querySelectorAll('select[name^="song1-"], select[name^="song2-"]'));
        const selectedValues = allSongSelects.map(s => s.value).filter(v => v && v.length > 0);
        const duplicates = selectedValues.reduce((acc, val) => {
            acc[val] = (acc[val] || 0) + 1;
            return acc;
        }, {});
        const dupKeys = Object.keys(duplicates).filter(k => duplicates[k] > 1);
        if (dupKeys.length > 0) {
            // Message lisible : afficher les 1 ou 2 premières chansons en doublon
            const example = dupKeys.slice(0, 3).join(', ');
            showNotification(`Erreur : la/les chanson(s) suivante(s) est/sont sélectionnée(s) plusieurs fois : ${example}`, 'error');
            return; // stop la soumission
        }
        const newDuel = {
            id: new Date().getTime().toString(),
            name: duelData.name,
            points: duelData.points,
            sameSong: { title: 'N/A', artist: 'N/A', lyricsFile: '' },
            createdAt: new Date().toISOString()
        };
        try {
            yield addDuel(newDuel);
            showNotification(`Grille ${soloMode ? 'solo' : 'duel'} créée et sauvegardée!`, 'success');
            showDuelList(); // Retour à la liste
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
/**
 * Affiche le formulaire de création et cache la liste
 */
function showCreateForm() {
    console.log("showCreateForm appelé");
    const formContainer = document.getElementById("PrepGrille");
    const listContent = document.querySelector('.duels-list');
    const alertContent = document.querySelector('.alert');
    const menuButton = document.querySelector('.button_prep_grille');
    console.log("Form container:", formContainer);
    console.log("Form container innerHTML:", formContainer === null || formContainer === void 0 ? void 0 : formContainer.innerHTML);
    // Si le formulaire est vide, le régénérer
    if (formContainer && (!formContainer.innerHTML || formContainer.innerHTML.trim() === '')) {
        console.log("Le formulaire est vide, régénération...");
        // Essayer de récupérer les fichiers locaux, sinon utiliser une liste par défaut
        getLyricsListLocal().then(lyricsFiles => {
            console.log("Fichiers récupérés avec succès:", lyricsFiles);
            if (lyricsFiles.length > 0) {
                renderCreateDuelForm(lyricsFiles);
            }
            else {
                renderCreateDuelFormWithError([
                    "Adele - Hello.json",
                    "Ed Sheeran - Shape of You.json",
                    "Billie Eilish - Bad Guy.json"
                ]);
            }
            showFormWithStyles(formContainer);
        }).catch(error => {
            console.error("Erreur lors de la récupération des fichiers:", error);
            console.log("Utilisation d'une liste par défaut pour le formulaire");
            // Créer le formulaire avec une liste par défaut
            renderCreateDuelFormWithError([
                "Artiste Exemple - Chanson 1.json",
                "Artiste Test - Chanson 2.json",
                "Demo Artist - Test Song.json"
            ]);
            showFormWithStyles(formContainer);
        });
    }
    else {
        // Le formulaire a déjà du contenu, juste l'afficher
        showFormWithStyles(formContainer);
    }
    // Cacher les autres éléments
    if (listContent)
        listContent.style.display = 'none';
    if (alertContent)
        alertContent.style.display = 'none';
    if (menuButton)
        menuButton.style.display = 'none';
}
/**
 * Applique les styles d'affichage au formulaire
 */
function showFormWithStyles(formContainer) {
    if (formContainer) {
        formContainer.style.display = 'block';
        formContainer.style.visibility = 'visible';
        formContainer.style.opacity = '1';
        formContainer.style.height = 'auto';
        formContainer.style.position = 'relative';
        formContainer.style.zIndex = '1000';
        console.log("Formulaire affiché avec styles");
        console.log("Contenu final:", formContainer.innerHTML.substring(0, 100) + "...");
    }
}
/**
 * Version de fallback pour créer le formulaire avec message d'erreur
 */
function renderCreateDuelFormWithError(fallbackFiles) {
    const container = document.getElementById("PrepGrille");
    if (!container)
        return;
    const soloMode = isSoloMode();
    const modeText = soloMode ? 'solo' : 'duel';
    let formHtml = `
        <div class="form-container">
            <h2 style="color: red;">FORMULAIRE DE TEST (MODE DÉBOGAGE)</h2>
            <div style="background: yellow; padding: 10px; margin: 10px 0; border: 1px solid orange;">
                ⚠️ Erreur : Impossible de charger la liste des musiques. Utilisation d'exemples.
            </div>
            <button id="back-to-list-btn" type="button" >← Retour à la liste</button>
            <form id="newDuelForm" >
                <h3>Créer une nouvelle grille de ${modeText}</h3>
                <label for="duelName">Nom de la grille:</label>
                <input type="text" id="duelName" name="duelName" required >
    `;
    DUEL_POINTS_CATEGORIES.forEach(points => {
        formHtml += `
            <div class="point-category" >
                <h4>${points} Points</h4>
                <label>Thème:</label>
                <input type="text" name="theme-${points}" required >
                ${generateSongSelectionHtmlFallback(points, fallbackFiles)}
            </div>
        `;
    });
    formHtml += `
            <button type="submit">Créer (Mode Test)</button>
        </form>
        </div>
    `;
    container.innerHTML = formHtml;
    console.log("Formulaire de fallback créé");
}
/**
 * Version fallback pour la génération des sélections de chansons
 */
function generateSongSelectionHtmlFallback(points, fallbackFiles) {
    const soloMode = isSoloMode();
    const songOptions = fallbackFiles.map(file => `<option value="${file}">${file} (exemple)</option>`).join('');
    if (soloMode) {
        return `
            <label>Chanson:</label>
            <select name="song1-${points}" required >
                <option value="">Sélectionner une musique</option>
                ${songOptions}
            </select>
        `;
    }
    else {
        return `
            <label>Chanson 1:</label>
            <select name="song1-${points}" required >
                <option value="">Sélectionner une musique</option>
                ${songOptions}
            </select>
            <label>Chanson 2:</label>
            <select name="song2-${points}" required >
                <option value="">Sélectionner une musique</option>
                ${songOptions}
            </select>
        `;
    }
}
/**
 * Affiche la liste des duels et cache le formulaire
 */
function showDuelList() {
    console.log("showDuelList appelé");
    const formContainer = document.getElementById("PrepGrille");
    const listContent = document.querySelector('.duels-list');
    const alertContent = document.querySelector('.alert');
    const menuButton = document.querySelector('.button_prep_grille');
    if (formContainer) {
        formContainer.style.display = 'none';
        console.log("Formulaire caché");
    }
    if (listContent) {
        listContent.style.display = 'block';
        console.log("Liste affichée");
    }
    if (alertContent) {
        alertContent.style.display = 'block';
        console.log("Alert affichée");
    }
    if (menuButton) {
        menuButton.style.display = 'block';
        console.log("Bouton menu affiché");
    }
    // Ne pas recharger la liste complètement pour éviter de perdre le PrepGrille
    loadDuelsFromStorage();
}
/**
 * Gère la soumission du formulaire d'importation de fichier.
 * @param event L'événement de soumission.
 */
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
                        throw new Error('Le contenu du fichier n\'est pas une chaîne de caractères.');
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
// === ÉCOUTEURS D'ÉVÉNEMENTS GLOBAUX ===
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
    console.log("Clic détecté sur:", target.id, target);
    if (target.id === 'create-duel-btn') {
        event.preventDefault();
        console.log("Bouton create-duel-btn cliqué");
        showCreateForm();
    }
    else if (target.id === 'back-to-list-btn') {
        event.preventDefault();
        console.log("Bouton back-to-list-btn cliqué");
        showDuelList();
    }
});
// === INITIALISATION ===
document.addEventListener('DOMContentLoaded', () => __awaiter(void 0, void 0, void 0, function* () {
    console.log("DOM chargé, initialisation...");
    loadDuelsFromStorage();
    // D'abord rendre la liste pour créer/recréer PrepGrille
    renderDuelList();
    // Attendre un peu que le DOM soit mis à jour
    setTimeout(() => __awaiter(void 0, void 0, void 0, function* () {
        // Vérifier que PrepGrille existe maintenant
        let prepGrilleContainer = document.getElementById("PrepGrille");
        console.log("PrepGrille après renderDuelList:", prepGrilleContainer);
        if (!prepGrilleContainer) {
            console.error("PrepGrille TOUJOURS non trouvé après renderDuelList !");
            return;
        }
        try {
            // Utiliser notre nouvelle fonction locale
            const lyricsFiles = yield getLyricsListLocal();
            console.log("Fichiers lyrics chargés avec succès:", lyricsFiles);
            if (lyricsFiles.length > 0) {
                renderCreateDuelForm(lyricsFiles);
            }
            else {
                throw new Error("Aucun fichier lyrics trouvé");
            }
        }
        catch (error) {
            console.error("Erreur lors du chargement des musiques:", error);
            console.log("Création du formulaire avec des exemples...");
            // Créer le formulaire avec des exemples même si getLyricsList échoue
            const fallbackFiles = [
                "Adele - Hello.json",
                "Ed Sheeran - Shape of You.json",
                "Billie Eilish - Bad Guy.json",
                "The Weeknd - Blinding Lights.json",
                "Dua Lipa - Levitating.json"
            ];
            renderCreateDuelFormWithError(fallbackFiles);
            showNotification('Impossible de charger la liste des musiques. Mode débogage activé.', 'warning');
        }
        // Vérifier que le contenu a bien été inséré
        const containerAfter = document.getElementById("PrepGrille");
        console.log("PrepGrille après generation:", containerAfter);
        console.log("Contenu du PrepGrille:", containerAfter === null || containerAfter === void 0 ? void 0 : containerAfter.innerHTML.substring(0, 200));
    }), 200);
}));

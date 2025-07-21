var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
// --- Variable de stockage ---
export let preparedDuels = [];
// --- Fonctions de gestion des données ---
/**
 * Charge les duels depuis le localStorage.
 */
export function loadDuelsFromStorage() {
    if (typeof (Storage) !== "undefined") {
        const saved = localStorage.getItem('preparedDuels');
        if (saved) {
            preparedDuels = JSON.parse(saved);
        }
    }
}
/**
 * Sauvegarde la liste des duels dans le localStorage.
 */
function saveDuelsToStorage() {
    if (typeof (Storage) !== "undefined") {
        localStorage.setItem('preparedDuels', JSON.stringify(preparedDuels));
    }
}
/**
 * Ajoute un nouveau duel à la liste et sauvegarde.
 * @param duelData Les données du duel à ajouter.
 */
export function addDuel(duelData) {
    preparedDuels.push(duelData);
    saveDuelsToStorage();
    saveDuelToServer(duelData).catch(error => console.error("Failed to save duel to server:", error));
}
/**
 * Met à jour un duel existant.
 * @param index L'index du duel à mettre à jour.
 * @param duelData Les nouvelles données du duel.
 */
export function updateDuel(index, duelData) {
    if (preparedDuels[index]) {
        // Conserve la date de création originale
        duelData.createdAt = preparedDuels[index].createdAt;
        duelData.updatedAt = new Date().toISOString();
        preparedDuels[index] = duelData;
        saveDuelsToStorage();
        // Potentiellement, ajouter une logique pour mettre à jour sur le serveur aussi
    }
}
/**
 * Supprime un duel de la liste.
 * @param index L'index du duel à supprimer.
 */
export function deleteDuel(index) {
    preparedDuels.splice(index, 1);
    saveDuelsToStorage();
}
// --- Fonctions d'interaction avec le serveur ---
/**
 * Vérifie si un fichier de paroles existe sur le serveur.
 * @param filename Le nom du fichier à vérifier.
 */
export function checkLyricsFile(filename) {
    return __awaiter(this, void 0, void 0, function* () {
        if (!filename) {
            throw new Error('Veuillez entrer un nom de fichier');
        }
        const response = yield fetch(`/api/check-lyrics?filename=${encodeURIComponent(filename)}`);
        if (!response.ok) {
            throw new Error('Erreur lors de la vérification du fichier');
        }
        return response.json();
    });
}
/**
 * Sauvegarde un duel sur le serveur.
 * @param duelData Les données du duel à sauvegarder.
 */
function saveDuelToServer(duelData) {
    return __awaiter(this, void 0, void 0, function* () {
        const response = yield fetch('/api/duels', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(duelData)
        });
        if (!response.ok) {
            throw new Error('Erreur lors de la sauvegarde sur le serveur');
        }
        console.log('Duel sauvegardé sur le serveur avec succès');
    });
}
/**
 * Exporte un duel spécifique vers le serveur.
 * @param index L'index du duel à exporter.
 */
export function exportDuelToServer(index) {
    return __awaiter(this, void 0, void 0, function* () {
        const duel = preparedDuels[index];
        const response = yield fetch('/api/export-duel', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(duel)
        });
        if (!response.ok) {
            const errorText = yield response.text();
            throw new Error(errorText || "Erreur lors de l'export");
        }
        return `Duel "${duel.name}" exporté vers le serveur`;
    });
}
/**
 * Importe les duels depuis le serveur, en évitant les doublons.
 */
export function importDuelsFromServer() {
    return __awaiter(this, void 0, void 0, function* () {
        const response = yield fetch('/api/duels');
        if (!response.ok) {
            throw new Error('Erreur lors du chargement depuis le serveur');
        }
        const serverDuels = yield response.json();
        const existingNames = new Set(preparedDuels.map(duel => duel.name));
        const newDuels = serverDuels.filter(duel => !existingNames.has(duel.name));
        if (newDuels.length > 0) {
            preparedDuels.push(...newDuels);
            saveDuelsToStorage();
        }
        return newDuels.length;
    });
}
/**
 * Sauvegarde temporaire d'un duel sur le serveur.
 * @param formData Les données du formulaire de duel.
 */
export function saveTemporaryDuel(formData) {
    return __awaiter(this, void 0, void 0, function* () {
        const response = yield fetch('/api/temp-duel', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(formData)
        });
        if (!response.ok) {
            throw new Error('Erreur lors de la sauvegarde temporaire');
        }
    });
}
/**
 * Charge la sauvegarde temporaire depuis le serveur.
 */
export function loadTemporaryDuel() {
    return __awaiter(this, void 0, void 0, function* () {
        const response = yield fetch('/api/temp-duel');
        if (response.status === 404) {
            throw new Error('Aucune sauvegarde temporaire trouvée');
        }
        if (!response.ok) {
            throw new Error('Erreur lors du chargement temporaire');
        }
        return response.json();
    });
}

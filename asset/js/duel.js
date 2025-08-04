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
 * Ajoute un nouveau duel.
 * 1. L'envoie au serveur qui lui assigne un ID définitif.
 * 2. Utilise la réponse du serveur pour l'ajouter à la liste locale.
 * 3. Sauvegarde la liste locale dans le localStorage.
 * @param duelData Les données du duel à ajouter.
 */
export function addDuel(duelData) {
    return __awaiter(this, void 0, void 0, function* () {
        try {
            const response = yield fetch('/api/duels', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify([duelData])
            });
            if (!response.ok) {
                const errorText = yield response.text();
                throw new Error(errorText || 'Erreur lors de la sauvegarde sur le serveur');
            }
            const savedDuels = yield response.json();
            // On vérifie que la réponse est correcte et on prend le premier élément
            if (savedDuels && savedDuels.length > 0) {
                preparedDuels.push(savedDuels[0]); // Ajouter le duel avec l'ID confirmé par le serveur
                saveDuelsToStorage();
            }
            else {
                throw new Error("La réponse du serveur était invalide après la création du duel.");
            }
        }
        catch (error) {
            console.error("Failed to save duel to server:", error);
            throw error;
        }
    });
}
/**
 * Met à jour un duel existant.
 * @param index L'index du duel à mettre à jour.
 * @param duelData Les nouvelles données du duel.
 */
export function updateDuel(index, duelData) {
    if (preparedDuels[index]) {
        duelData.createdAt = preparedDuels[index].createdAt;
        duelData.updatedAt = new Date().toISOString();
        preparedDuels[index] = duelData;
        saveDuelsToStorage();
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

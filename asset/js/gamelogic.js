var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
// gamelogic.ts
import { showNotification } from './gamenotification.js';
// === LOGIQUE GLOBALE ===
export let preparedDuels = [];
/**
 * Charge les duels depuis le stockage local.
 * Assure la persistance des données.
 */
export function loadDuelsFromStorage() {
    if (typeof localStorage !== 'undefined') {
        const duels = localStorage.getItem('duels');
        if (duels) {
            try {
                preparedDuels = JSON.parse(duels);
            }
            catch (error) {
                showNotification('Erreur lors du chargement des duels', 'error');
                console.error(error);
            }
        }
    }
}
/**
 * Ajoute un nouveau duel au stockage local.
 * @param duelData Les données du duel à ajouter.
 */
export function addDuel(duelData) {
    return __awaiter(this, void 0, void 0, function* () {
        // Simule une requête API (remplace par ton fetch Go)
        console.log("Ajout du duel", duelData);
        preparedDuels.push(duelData);
        if (typeof localStorage !== 'undefined') {
            localStorage.setItem('duels', JSON.stringify(preparedDuels));
        }
    });
}
/**
 * Récupère la liste des fichiers de paroles depuis le serveur Go.
 * @returns Une promesse qui résout avec un tableau de noms de fichiers.
 */
export function getLyricsList() {
    return __awaiter(this, void 0, void 0, function* () {
        const response = yield fetch('/api/lyrics-list');
        if (!response.ok) {
            throw new Error('Impossible de charger la liste des paroles.');
        }
        return response.json();
    });
}

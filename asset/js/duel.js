"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
var __generator = (this && this.__generator) || function (thisArg, body) {
    var _ = { label: 0, sent: function() { if (t[0] & 1) throw t[1]; return t[1]; }, trys: [], ops: [] }, f, y, t, g = Object.create((typeof Iterator === "function" ? Iterator : Object).prototype);
    return g.next = verb(0), g["throw"] = verb(1), g["return"] = verb(2), typeof Symbol === "function" && (g[Symbol.iterator] = function() { return this; }), g;
    function verb(n) { return function (v) { return step([n, v]); }; }
    function step(op) {
        if (f) throw new TypeError("Generator is already executing.");
        while (g && (g = 0, op[0] && (_ = 0)), _) try {
            if (f = 1, y && (t = op[0] & 2 ? y["return"] : op[0] ? y["throw"] || ((t = y["return"]) && t.call(y), 0) : y.next) && !(t = t.call(y, op[1])).done) return t;
            if (y = 0, t) op = [op[0] & 2, t.value];
            switch (op[0]) {
                case 0: case 1: t = op; break;
                case 4: _.label++; return { value: op[1], done: false };
                case 5: _.label++; y = op[1]; op = [0]; continue;
                case 7: op = _.ops.pop(); _.trys.pop(); continue;
                default:
                    if (!(t = _.trys, t = t.length > 0 && t[t.length - 1]) && (op[0] === 6 || op[0] === 2)) { _ = 0; continue; }
                    if (op[0] === 3 && (!t || (op[1] > t[0] && op[1] < t[3]))) { _.label = op[1]; break; }
                    if (op[0] === 6 && _.label < t[1]) { _.label = t[1]; t = op; break; }
                    if (t && _.label < t[2]) { _.label = t[2]; _.ops.push(op); break; }
                    if (t[2]) _.ops.pop();
                    _.trys.pop(); continue;
            }
            op = body.call(thisArg, _);
        } catch (e) { op = [6, e]; y = 0; } finally { f = t = 0; }
        if (op[0] & 5) throw op[1]; return { value: op[0] ? op[1] : void 0, done: true };
    }
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.preparedDuels = void 0;
exports.loadDuelsFromStorage = loadDuelsFromStorage;
exports.addDuel = addDuel;
exports.updateDuel = updateDuel;
exports.deleteDuel = deleteDuel;
exports.checkLyricsFile = checkLyricsFile;
exports.exportDuelToServer = exportDuelToServer;
exports.importDuelsFromServer = importDuelsFromServer;
exports.saveTemporaryDuel = saveTemporaryDuel;
exports.loadTemporaryDuel = loadTemporaryDuel;
// --- Variable de stockage ---
exports.preparedDuels = [];
// --- Fonctions de gestion des données ---
/**
 * Charge les duels depuis le localStorage.
 */
function loadDuelsFromStorage() {
    if (typeof (Storage) !== "undefined") {
        var saved = localStorage.getItem('preparedDuels');
        if (saved) {
            exports.preparedDuels = JSON.parse(saved);
        }
    }
}
/**
 * Sauvegarde la liste des duels dans le localStorage.
 */
function saveDuelsToStorage() {
    if (typeof (Storage) !== "undefined") {
        localStorage.setItem('preparedDuels', JSON.stringify(exports.preparedDuels));
    }
}
/**
 * Ajoute un nouveau duel à la liste et sauvegarde.
 * @param duelData Les données du duel à ajouter.
 */
function addDuel(duelData) {
    exports.preparedDuels.push(duelData);
    saveDuelsToStorage();
    saveDuelToServer(duelData).catch(function (error) { return console.error("Failed to save duel to server:", error); });
}
/**
 * Met à jour un duel existant.
 * @param index L'index du duel à mettre à jour.
 * @param duelData Les nouvelles données du duel.
 */
function updateDuel(index, duelData) {
    if (exports.preparedDuels[index]) {
        // Conserve la date de création originale
        duelData.createdAt = exports.preparedDuels[index].createdAt;
        duelData.updatedAt = new Date().toISOString();
        exports.preparedDuels[index] = duelData;
        saveDuelsToStorage();
        // Potentiellement, ajouter une logique pour mettre à jour sur le serveur aussi
    }
}
/**
 * Supprime un duel de la liste.
 * @param index L'index du duel à supprimer.
 */
function deleteDuel(index) {
    exports.preparedDuels.splice(index, 1);
    saveDuelsToStorage();
}
// --- Fonctions d'interaction avec le serveur ---
/**
 * Vérifie si un fichier de paroles existe sur le serveur.
 * @param filename Le nom du fichier à vérifier.
 */
function checkLyricsFile(filename) {
    return __awaiter(this, void 0, void 0, function () {
        var response;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    if (!filename) {
                        throw new Error('Veuillez entrer un nom de fichier');
                    }
                    return [4 /*yield*/, fetch("/api/check-lyrics?filename=".concat(encodeURIComponent(filename)))];
                case 1:
                    response = _a.sent();
                    if (!response.ok) {
                        throw new Error('Erreur lors de la vérification du fichier');
                    }
                    return [2 /*return*/, response.json()];
            }
        });
    });
}
/**
 * Sauvegarde un duel sur le serveur.
 * @param duelData Les données du duel à sauvegarder.
 */
function saveDuelToServer(duelData) {
    return __awaiter(this, void 0, void 0, function () {
        var response;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0: return [4 /*yield*/, fetch('/api/duels', {
                        method: 'POST',
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify(duelData)
                    })];
                case 1:
                    response = _a.sent();
                    if (!response.ok) {
                        throw new Error('Erreur lors de la sauvegarde sur le serveur');
                    }
                    console.log('Duel sauvegardé sur le serveur avec succès');
                    return [2 /*return*/];
            }
        });
    });
}
/**
 * Exporte un duel spécifique vers le serveur.
 * @param index L'index du duel à exporter.
 */
function exportDuelToServer(index) {
    return __awaiter(this, void 0, void 0, function () {
        var duel, response, errorText;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    duel = exports.preparedDuels[index];
                    return [4 /*yield*/, fetch('/api/export-duel', {
                            method: 'POST',
                            headers: { 'Content-Type': 'application/json' },
                            body: JSON.stringify(duel)
                        })];
                case 1:
                    response = _a.sent();
                    if (!!response.ok) return [3 /*break*/, 3];
                    return [4 /*yield*/, response.text()];
                case 2:
                    errorText = _a.sent();
                    throw new Error(errorText || "Erreur lors de l'export");
                case 3: return [2 /*return*/, "Duel \"".concat(duel.name, "\" export\u00E9 vers le serveur")];
            }
        });
    });
}
/**
 * Importe les duels depuis le serveur, en évitant les doublons.
 */
function importDuelsFromServer() {
    return __awaiter(this, void 0, void 0, function () {
        var response, serverDuels, existingNames, newDuels;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0: return [4 /*yield*/, fetch('/api/duels')];
                case 1:
                    response = _a.sent();
                    if (!response.ok) {
                        throw new Error('Erreur lors du chargement depuis le serveur');
                    }
                    return [4 /*yield*/, response.json()];
                case 2:
                    serverDuels = _a.sent();
                    existingNames = new Set(exports.preparedDuels.map(function (duel) { return duel.name; }));
                    newDuels = serverDuels.filter(function (duel) { return !existingNames.has(duel.name); });
                    if (newDuels.length > 0) {
                        exports.preparedDuels.push.apply(exports.preparedDuels, newDuels);
                        saveDuelsToStorage();
                    }
                    return [2 /*return*/, newDuels.length];
            }
        });
    });
}
/**
 * Sauvegarde temporaire d'un duel sur le serveur.
 * @param formData Les données du formulaire de duel.
 */
function saveTemporaryDuel(formData) {
    return __awaiter(this, void 0, void 0, function () {
        var response;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0: return [4 /*yield*/, fetch('/api/temp-duel', {
                        method: 'POST',
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify(formData)
                    })];
                case 1:
                    response = _a.sent();
                    if (!response.ok) {
                        throw new Error('Erreur lors de la sauvegarde temporaire');
                    }
                    return [2 /*return*/];
            }
        });
    });
}
/**
 * Charge la sauvegarde temporaire depuis le serveur.
 */
function loadTemporaryDuel() {
    return __awaiter(this, void 0, void 0, function () {
        var response;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0: return [4 /*yield*/, fetch('/api/temp-duel')];
                case 1:
                    response = _a.sent();
                    if (response.status === 404) {
                        throw new Error('Aucune sauvegarde temporaire trouvée');
                    }
                    if (!response.ok) {
                        throw new Error('Erreur lors du chargement temporaire');
                    }
                    return [2 /*return*/, response.json()];
            }
        });
    });
}

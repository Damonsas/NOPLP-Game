
interface Song {
  title: string;
  artist: string;
  audioUrl?: string | null;
  lyricsFile?: string | null;
}

interface PointCategory {
  theme: string;
  songs: [Song, Song];
}

interface Duel {
  id: string;
  name: string;
  points: {
    '10': PointCategory;
    '20': PointCategory;
    '30': PointCategory;
    '40': PointCategory;
    '50': PointCategory;
  };
  sameSong: Song;
  createdAt: string;
  updatedAt?: string;
}

interface GameSession {
  id: string;
  duelId: number;
  currentLevel: string;
  selectedSongs: { [key: string]: number }; 
  team1Score: number;
  team2Score: number;
  startedAt: string; 
  status: 'playing' | 'paused' | 'finished';

  currentSong?: Song;
  lyricsContent?: string;
  lyricsVisible: boolean;
}


export let preparedDuels: Duel[] = [];



export function loadDuelsFromStorage(): void {
    if (typeof(Storage) !== "undefined") {
        const saved = localStorage.getItem('preparedDuels');
        if (saved) {
            preparedDuels = JSON.parse(saved);
        }
    }
}


function saveDuelsToStorage(): void {
    if (typeof(Storage) !== "undefined") {
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
export async function addDuel(duelData: Duel): Promise<void> {
    try {
        const response = await fetch('/api/duels', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
           
            body: JSON.stringify([duelData])
        });

        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(errorText || 'Erreur lors de la sauvegarde sur le serveur');
        }

        const savedDuels: Duel[] = await response.json();

        if (savedDuels && savedDuels.length > 0) {
            preparedDuels.push(savedDuels[0]); // Ajouter le duel avec l'ID confirmé par le serveur
            saveDuelsToStorage();
        } else {
            throw new Error("La réponse du serveur était invalide après la création du duel.");
        }

    } catch (error) {
        console.error("Failed to save duel to server:", error);
        throw error;
    }
}
/**
 * Met à jour un duel existant.
 * @param index L'index du duel à mettre à jour.
 * @param duelData Les nouvelles données du duel.
 */
export function updateDuel(index: number, duelData: Duel): void {
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
export function deleteDuel(index: number): void {
    preparedDuels.splice(index, 1);
    saveDuelsToStorage();
}


/**
 * Vérifie si un fichier de paroles existe sur le serveur.
 * @param filename Le nom du fichier à vérifier.
 */
export async function checkLyricsFile(filename: string): Promise<{ exists: boolean; content?: string }> {
    if (!filename) {
        throw new Error('Veuillez entrer un nom de fichier');
    }
    const response = await fetch(`/api/check-lyrics?filename=${encodeURIComponent(filename)}`);
    if (!response.ok) {
        throw new Error('Erreur lors de la vérification du fichier');
    }
    return response.json();
}

/**
 * Sauvegarde un duel sur le serveur.
 * @param duelData Les données du duel à sauvegarder.
 */
async function saveDuelToServer(duelData: Duel): Promise<void> {
    const response = await fetch('/api/duels', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(duelData)
    });
    if (!response.ok) {
        throw new Error('Erreur lors de la sauvegarde sur le serveur');
    }
    console.log('Duel sauvegardé sur le serveur avec succès');
}

/**
 * Exporte un duel spécifique vers le serveur.
 * @param index L'index du duel à exporter.
 */
export async function exportDuelToServer(index: number): Promise<string> {
    const duel = preparedDuels[index];
    const response = await fetch('/api/export-duel', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(duel)
    });
    if (!response.ok) {
        const errorText = await response.text();
        throw new Error(errorText || "Erreur lors de l'export");
    }
    return `Duel "${duel.name}" exporté vers le serveur`;
}

/**
 * Importe les duels depuis le serveur, en évitant les doublons.
 */
export async function importDuelsFromServer(): Promise<number> {
    const response = await fetch('/api/duels');
    if (!response.ok) {
        throw new Error('Erreur lors du chargement depuis le serveur');
    }
    const serverDuels: Duel[] = await response.json();

    const existingNames = new Set(preparedDuels.map(duel => duel.name));
    const newDuels = serverDuels.filter(duel => !existingNames.has(duel.name));

    if (newDuels.length > 0) {
        preparedDuels.push(...newDuels);
        saveDuelsToStorage();
    }
    return newDuels.length;
}

/**
 * Sauvegarde temporaire d'un duel sur le serveur.
 * @param formData Les données du formulaire de duel.
 */
export async function saveTemporaryDuel(formData: Duel): Promise<void> {
    const response = await fetch('/api/temp-duel', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(formData)
    });
    if (!response.ok) {
        throw new Error('Erreur lors de la sauvegarde temporaire');
    }
}

/**
 * Charge la sauvegarde temporaire depuis le serveur.
 */
export async function loadTemporaryDuel(): Promise<Duel> {
    const response = await fetch('/api/temp-duel');
    if (response.status === 404) {
        throw new Error('Aucune sauvegarde temporaire trouvée');
    }
    if (!response.ok) {
        throw new Error('Erreur lors du chargement temporaire');
    }
    return response.json();
}
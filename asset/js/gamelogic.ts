// gamelogic.ts
import { showNotification } from './gamenotification.js';

// === INTERFACES ET TYPES ===
export interface GameFormData {
    [key: string]: FormDataEntryValue;
}

export interface Song {
  title: string;
  artist: string;
  audioUrl?: string | null;
  lyricsFile?: string | null;
}

export interface PointCategory {
  theme: string;
  songs: [Song, Song];
}

export interface Duel {
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
  mode?: 'solo' | 'duel';
}

export interface GameSession {
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

// === LOGIQUE GLOBALE ===
export let preparedDuels: Duel[] = [];

/**
 * Charge les duels depuis le stockage local.
 * Assure la persistance des données.
 */
export function loadDuelsFromStorage(): void {
  if (typeof localStorage !== 'undefined') {
    const duels = localStorage.getItem('duels');
    if (duels) {
      try {
        preparedDuels = JSON.parse(duels);
      } catch (error) {
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
export async function addDuel(duelData: Duel): Promise<void> {
  // Simule une requête API (remplace par ton fetch Go)
  console.log("Ajout du duel", duelData);
  preparedDuels.push(duelData);
  if (typeof localStorage !== 'undefined') {
    localStorage.setItem('duels', JSON.stringify(preparedDuels));
  }
}

/**
 * Récupère la liste des fichiers de paroles depuis le serveur Go.
 * @returns Une promesse qui résout avec un tableau de noms de fichiers.
 */
export async function getLyricsList(): Promise<string[]> {
    const response = await fetch('/api/lyrics-list');
    if (!response.ok) {
        throw new Error('Impossible de charger la liste des paroles.');
    }
    return response.json();
}
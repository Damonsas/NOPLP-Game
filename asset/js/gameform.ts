import { showNotification } from './gamenotification.js';
import { addDuel, loadDuelsFromStorage, preparedDuels, Duel, getLyricsList } from './gamelogic.js';

const DUEL_POINTS_CATEGORIES = [50, 40, 30, 20, 10];


/**
 * Récupère la liste des fichiers de paroles locaux
 * @returns Promise avec la liste des noms de fichiers
 */
async function getLyricsListLocal(): Promise<string[]> {
  const indexPath = './data/serverdata/paroledata/index.json';
  try {
    console.log("Tentative de récupération des fichiers lyrics depuis", indexPath);

    const response = await fetch(indexPath, { cache: 'no-store' });
    if (!response.ok) {
      console.warn(`index.json introuvable ou erreur (${response.status})`);
      return [];
    }

    const arr = await response.json();
    if (!Array.isArray(arr)) {
      console.warn("Format inattendu pour index.json (pas un tableau)");
      return [];
    }

   
    const files = arr.map((item: any) => {
      if (!item) return null;
      const raw = item.ligne || (item.artiste && item.titre ? `${item.artiste} - ${item.titre}` : null);
      if (!raw) return null;
      return raw.endsWith('.json') ? raw : `${raw}.json`;
    }).filter(Boolean) as string[];

    console.log("Fichiers lyrics trouvés via index.json:", files);
    return files;
  } catch (err) {
    console.error("Erreur lors de la récupération via index.json:", err);
    return [];
  }
}


/**
 * Charge un fichier de paroles spécifique
 * @param filename Nom du fichier à charger
 * @returns Promise avec le contenu du fichier JSON
 */
async function loadLyricsFile(filename: string): Promise<any> {
    try {
        const response = await fetch(`./data/paroledata/${filename}`);
        if (!response.ok) {
            throw new Error(`Impossible de charger ${filename}: ${response.status}`);
        }
        const data = await response.json();
        console.log(`Fichier ${filename} chargé avec succès`);
        return data;
    } catch (error) {
        console.error(`Erreur lors du chargement de ${filename}:`, error);
        throw error;
    }
}

/**
 * Détermine si on est en mode solo ou duel basé sur l'URL actuelle
 * @returns true si mode solo, false si mode duel
 */
function isSoloMode(): boolean {
  return window.location.pathname.includes('solo');
}


/**
 * Génère la carte HTML pour un duel donné.
 * @param duel Le duel à afficher.
 * @returns La chaîne HTML de la carte.
 */
function generateDuelCard(duel: Duel): string {
  const themes = DUEL_POINTS_CATEGORIES.map(p => (duel.points[p.toString()]?.theme)).filter(t => t).join(', ');
  const currentMode = isSoloMode();
  
  
  const playUrl = currentMode 
    ? `/solo?id=${duel.id}` 
    : `/duel-game?id=${duel.id}`;  

  const supprform = `/duel-delete?id=${duel.id}`;
  const modifform = `/duel-edit?id=${duel.id}`;

  return `
    <div class="duel-card" data-duel-id="${duel.id}">
      <h3>${duel.name}</h3>
      <button class="play-duel-btn" data-duel-id="${duel.id}">Jouer</button>
      <button onclick="window.location.href='${supprform}'"> Supprimer <i class="fa-regular fa-trash-can"></i> </button>
      <button onclick="window.location.href='${modifform}'"> Modifier </button>
    </div>
  `;

}

async function handlePlayDuel(duelId: string) {
  try {
    const id = parseInt(duelId);
    if (isNaN(id)) throw new Error('ID invalide');

    const duels = JSON.parse(localStorage.getItem('duels') || '[]');
    const duel = duels.find((d: any) => d.id == id);
    if (!duel) throw new Error('Duel non trouvé');

    const res = await fetch('/api/duels', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(duel)
    });

    if (!res.ok) throw new Error(`Erreur serveur: ${res.status}`);
    
    const serverDuel = await res.json();
    const url = isSoloMode() ? `/solo?id=${serverDuel.id}` : `/duel-game?id=${serverDuel.id}`;
    
    window.location.href = url;
  } catch (error) {
    showNotification(`Erreur: ${error}`, 'error');
  }
}

/**
 * Génère le bouton "Préparer une grille".
 * @returns La chaîne HTML du bouton.
 */
function getMenuHtml(): string {
  return `<div class="button_prep_grille"><button id="create-duel-btn">Préparer une grille</button></div>`;
}

/**
 * Affiche la liste des duels disponibles ou le message d'absence de duel.
 */
function renderDuelList(): void {
  const container = document.querySelector(".Sectionduel") as HTMLElement;
  if (!container) return;

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
function generateSongSelectionHtml(points: number, lyricsFiles: string[]): string {
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
  } else {
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

function attachUniqueSelectionHandlers(formOrContainer: HTMLElement | null): void {
  if (!formOrContainer) return;

  const songSelects = Array.from(formOrContainer.querySelectorAll('select[name^="song1-"], select[name^="song2-"]')) as HTMLSelectElement[];

  function refreshDisabledOptions() {
    const selectedValues = songSelects
      .map(s => s.value)
      .filter(v => v && v.length > 0);

    songSelects.forEach(select => {
      const ownValue = select.value;
      Array.from(select.options).forEach(opt => {
        if (opt.value === ownValue) {
          opt.disabled = false;
          return;
        }
        opt.disabled = selectedValues.includes(opt.value);
      });
    });
  }

  songSelects.forEach(select => {
    select.addEventListener('change', () => {
      refreshDisabledOptions();
    });

    select.addEventListener('input', () => {
      refreshDisabledOptions();
    });
  });

  refreshDisabledOptions();
}


/**
 * Génère et affiche le formulaire de création de duel avec les listes de musiques dynamiques.
 * @param lyricsFiles La liste des noms de fichiers de paroles.
 */
function renderCreateDuelForm(lyricsFiles: string[]): void {
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
}

// === GESTIONNAIRES D'ÉVÉNEMENTS ===

/**
 * Gère la soumission du formulaire de création de duel.
 * @param event L'événement de soumission.
 */
async function handleNewDuelFormSubmit(event: Event): Promise<void> {
    event.preventDefault();
    const form = event.target as HTMLFormElement;
    const formData = new FormData(form);

    const duelData: Partial<Duel> = {};
    const soloMode = isSoloMode();

    for (const [key, value] of formData as any) {
        const parts = key.split('-');
        const fieldName = parts[0];
        const points = parts.length > 1 ? parts[1] : null;

        if (fieldName === 'duelName') {
            duelData.name = value as string;
        } else if (points) {
            if (!duelData.points) {
                duelData.points = {} as any;
            }
            if (!duelData.points[points]) {
                duelData.points[points] = {} as any;
            }
            if (fieldName === 'theme') {
                (duelData.points[points] as any).theme = value as string;
            } else if (fieldName.startsWith('song')) {
                const songIndex = fieldName === 'song1' ? 0 : 1;
                if (!(duelData.points[points] as any).songs) {
                    (duelData.points[points] as any).songs = soloMode ? [{}] : [{}, {}];
                }
                if (soloMode && songIndex === 0) {
                    (duelData.points[points] as any).songs[0] = { lyricsFile: value as string } as any;
                } else if (!soloMode) {
                    (duelData.points[points] as any).songs[songIndex] = { lyricsFile: value as string } as any;
                }
            }
        }
    }

    const allSongSelects = Array.from((form as HTMLFormElement).querySelectorAll('select[name^="song1-"], select[name^="song2-"]')) as HTMLSelectElement[];
    const selectedValues = allSongSelects.map(s => s.value).filter(v => v && v.length > 0);

    const duplicates = selectedValues.reduce((acc: Record<string, number>, val) => {
    acc[val] = (acc[val] || 0) + 1;
    return acc;
     }, {});

    const dupKeys = Object.keys(duplicates).filter(k => duplicates[k] > 1);
    if (dupKeys.length > 0) {
    const example = dupKeys.slice(0, 3).join(', ');
    showNotification(`Erreur : la/les chanson(s) suivante(s) est/sont sélectionnée(s) plusieurs fois : ${example}`, 'error');
    return; // stop la soumission
    }

    const newDuel: Duel = {
      id: Date.now(),
      name: duelData.name as string,
      points: duelData.points as any,
      sameSong: { title: 'N/A', artist: 'N/A', lyricsFile: '' },
      createdAt: new Date().toISOString()
    };

    try {
        await addDuel(newDuel);
        showNotification(`Grille ${soloMode ? 'solo' : 'duel'} créée et sauvegardée!`, 'success');
        showDuelList(); // Retour à la liste
    } catch (error: unknown) {
        if (error instanceof Error) {
            showNotification(`Erreur lors de la création : ${error.message}`, 'error');
        } else {
            showNotification('Une erreur inconnue est survenue.', 'error');
        }
    }
}

/**
 * Affiche le formulaire de création et cache la liste
 */
function showCreateForm(): void {
    
    const formContainer = document.getElementById("PrepGrille");
    const listContent = document.querySelector('.duels-list') as HTMLElement;
    const alertContent = document.querySelector('.alert') as HTMLElement;
    const menuButton = document.querySelector('.button_prep_grille') as HTMLElement;
    
    if (formContainer && (!formContainer.innerHTML || formContainer.innerHTML.trim() === '')) {        
        getLyricsListLocal().then(lyricsFiles => {
            if (lyricsFiles.length > 0) {
                renderCreateDuelForm(lyricsFiles);
            } else {
                ([
                    "Adele - Hello.json",
                    "Ed Sheeran - Shape of You.json",
                    "Billie Eilish - Bad Guy.json"
                ]);
            }
            showFormWithStyles(formContainer);
        }).catch(error => {
            console.error("Erreur lors de la récupération des fichiers:", error);
            console.log("Utilisation d'une liste par défaut pour le formulaire");
            
            ([
                "Artiste Exemple - Chanson 1.json",
                "Artiste Test - Chanson 2.json",
                "Demo Artist - Test Song.json"
            ]);
            showFormWithStyles(formContainer);
        });
    } else {
        showFormWithStyles(formContainer);
    }

    if (listContent) listContent.style.display = 'none';
    if (alertContent) alertContent.style.display = 'none';
    if (menuButton) menuButton.style.display = 'none';
}

/**
 * Applique les styles d'affichage au formulaire
 */
function showFormWithStyles(formContainer: HTMLElement | null): void {
    if (formContainer) {
        formContainer.style.display = 'block';
        formContainer.style.visibility = 'visible';
        formContainer.style.opacity = '1';
        formContainer.style.height = 'auto';
        formContainer.style.position = 'relative';
        formContainer.style.zIndex = '1000';
    }
}


/**
 * Affiche la liste des duels et cache le formulaire
 */
function showDuelList(): void {
    console.log("showDuelList appelé");
    
    const formContainer = document.getElementById("PrepGrille");
    const listContent = document.querySelector('.duels-list') as HTMLElement;
    const alertContent = document.querySelector('.alert') as HTMLElement;
    const menuButton = document.querySelector('.button_prep_grille') as HTMLElement;

    if (formContainer) {
        formContainer.style.display = 'none';
    }

    if (listContent) {
        listContent.style.display = 'block';
    }

    if (alertContent) {
        alertContent.style.display = 'block';
    }

    if (menuButton) {
        menuButton.style.display = 'block';
    }

    // Ne pas recharger la liste complètement pour éviter de perdre le PrepGrille
    loadDuelsFromStorage();
}

/**
 * Gère la soumission du formulaire d'importation de fichier.
 * @param event L'événement de soumission.
 */
async function handleImportFormSubmit(event: Event): Promise<void> {
    event.preventDefault();
    const fileInput = document.getElementById('duelFile') as HTMLInputElement;
    const file = fileInput.files?.[0];

    if (!file) {
        showNotification('Veuillez sélectionner un fichier', 'warning');
        return;
    }

    const reader = new FileReader();
    reader.onload = async function(e: ProgressEvent<FileReader>) {
        try {
            const result = e.target?.result;
            if (typeof result !== 'string') {
                throw new Error('Le contenu du fichier n\'est pas une chaîne de caractères.');
            }
            const duelData = JSON.parse(result) as Duel;

            await addDuel(duelData);
            renderDuelList();
            showNotification(`Duel "${duelData.name}" importé avec succès`, 'success');
        } catch (error: unknown) {
            if (error instanceof Error) {
                showNotification(`Erreur: ${error.message}`, 'error');
            } else {
                showNotification('Une erreur inconnue est survenue.', 'error');
            }
        }
    };
    reader.readAsText(file);
}

// === ÉCOUTEURS D'ÉVÉNEMENTS GLOBAUX ===
document.addEventListener('submit', (event) => {
    const target = event.target as HTMLElement;
    if (target.id === 'newDuelForm') {
        handleNewDuelFormSubmit(event);
    } else if (target.id === 'importForm') {
        handleImportFormSubmit(event);
    }
});


// a enlenver si tout fonctionnel
document.addEventListener('click', (event) => {
    const target = event.target as HTMLElement;

    if (target.classList.contains('play-duel-btn')) {
        event.preventDefault();
        const id = target.getAttribute('data-duel-id');
        if (id) {
            console.log('🎯 Clic sur Jouer détecté, duelId:', id);
            handlePlayDuel(id);
        } else {
            console.error('❌ Pas de duelId trouvé sur le bouton');
        }
        return;
    }
    
    if (target.id === 'create-duel-btn') {
        event.preventDefault();
        console.log("Bouton create-duel-btn cliqué");
        showCreateForm();
    } else if (target.id === 'back-to-list-btn') {
        event.preventDefault();
        console.log("Bouton back-to-list-btn cliqué");
        showDuelList();
    }
});

document.addEventListener('DOMContentLoaded', async () => {
    
    loadDuelsFromStorage();
    
    renderDuelList();
    
    setTimeout(async () => {
        let prepGrilleContainer = document.getElementById("PrepGrille");
        console.log("PrepGrille après renderDuelList:", prepGrilleContainer);
        
        if (!prepGrilleContainer) {
            console.error("PrepGrille TOUJOURS non trouvé après renderDuelList !");
            return;
        }

        try {
            const lyricsFiles = await getLyricsListLocal();
            
            if (lyricsFiles.length > 0) {
                renderCreateDuelForm(lyricsFiles);
            } else {
                throw new Error("Aucun fichier lyrics trouvé");
            }
            
        } catch (error) {
            
            // Créer le formulaire avec des exemples même si getLyricsList échoue
            const fallbackFiles = [
                "Adele - Hello.json",
                "Ed Sheeran - Shape of You.json", 
                "Billie Eilish - Bad Guy.json",
                "The Weeknd - Blinding Lights.json",
                "Dua Lipa - Levitating.json"
            ];
        }
        
        // Vérifier que le contenu a bien été inséré
        const containerAfter = document.getElementById("PrepGrille");
        
    }, 200);
});

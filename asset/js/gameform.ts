import { showNotification } from './gamenotification.js';
import { addDuel, loadDuelsFromStorage, preparedDuels, Duel, getLyricsList } from './gamelogic.js';

// === CONSTANTES ===
const DUEL_POINTS_CATEGORIES = [50, 40, 30, 20, 10];
const sameSong = "memechanson";

// === FONCTIONS UTILITAIRES POUR LYRICS ===

/**
 * Récupère la liste des fichiers de paroles locaux
 * @returns Promise avec la liste des noms de fichiers
 */
async function getLyricsListLocal(): Promise<string[]> {
    try {
        console.log("Tentative de récupération des fichiers lyrics locaux...");
        
        // Essayer de lire le contenu du dossier data/serverdata/paroledata/
        // Cette approche dépend de votre serveur - voici plusieurs méthodes :
        
        // Méthode 1 : Si vous avez un endpoint API qui liste les fichiers
        const response = await fetch('/api/lyrics-list');
        if (response.ok) {
            const files = await response.json();
            console.log("Fichiers lyrics trouvés via API:", files);
            return files;
        }
    } catch (error) {
        console.error("Erreur lors de la récupération via API:", error);
    }
    
    try {
        // Méthode 2 : Si vous avez un fichier index.json qui liste tous les fichiers
        const response = await fetch('/data/serverdata/paroledata/index.json');
        if (response.ok) {
            const data = await response.json();
            console.log("Fichiers lyrics trouvés via index.json:", data.files);
            return data.files || [];
        }
    } catch (error) {
        console.error("Erreur lors de la récupération via index.json:", error);
    }
    
    try {
        // Méthode 3 : Liste hardcodée temporaire (à remplacer par vos vrais fichiers)
        console.log("Utilisation de la liste hardcodée temporaire");
        return [
            "Adele - Hello.json",
            "Ed Sheeran - Shape of You.json",
            "Billie Eilish - Bad Guy.json",
            "The Weeknd - Blinding Lights.json",
            "Dua Lipa - Levitating.json",
            "Post Malone - Circles.json",
            "Taylor Swift - Anti-Hero.json",
            "Harry Styles - As It Was.json",
            "Olivia Rodrigo - drivers license.json",
            "BTS - Dynamite.json"
        ];
    } catch (error) {
        console.error("Erreur même avec la liste hardcodée:", error);
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
        const response = await fetch(`/data/serverdata/paroledata/${filename}`);
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

// === FONCTIONS DE RENDU (UI) ===

/**
 * Génère la carte HTML pour un duel donné.
 * @param duel Le duel à afficher.
 * @returns La chaîne HTML de la carte.
 */
function generateDuelCard(duel: Duel): string {
  const themes = DUEL_POINTS_CATEGORIES.map(p => (duel.points[p.toString()]?.theme)).filter(t => t).join(', ');
  const currentMode = isSoloMode() ? 'solo' : 'duel';
  return `
    <div class="duel-card" data-duel-id="${duel.id}">
      <h3>${duel.name}</h3>
      <p>Thèmes : ${themes}</p>
      <button onclick="window.location.href='/${currentMode}?id=${duel.id}'">Jouer</button>
    </div>
  `;
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

function generateSameSongHtml(lyricsFiles: string[]): string {
  const options = ['<option value="">Aucune</option>']
    .concat(lyricsFiles.map(f => `<option value="${f}">${f}</option>`))
    .join('');
  return `
    <div class="same-song-block" style="margin-bottom: 10px;">
      <label for="samesong">Même chanson (optionnel) :</label>
      <select id="samesong" name="samesong" style="display:block; margin:5px 0; padding:5px; width:320px;">
        ${options}
      </select>
    </div>
  `;
}

function generateSongSelectionHtml(points: number, lyricsFiles: string[]): string {
  const soloMode = isSoloMode();
  const songOptions = lyricsFiles.map(file => `<option value="${file}">${file}</option>`).join('');

  if (soloMode) {
    return `
      <label>Chanson:</label>
      <select name="song1-${points}" required style="display: block; margin: 5px 0; padding: 5px; width: 200px;">
        <option value="">Sélectionner une musique</option>
        ${songOptions}
      </select>
    `;
  } else {
    return `
      <label>Chanson 1:</label>
      <select name="song1-${points}" required style="display: block; margin: 5px 0; padding: 5px; width: 200px;">
        <option value="">Sélectionner une musique</option>
        ${songOptions}
      </select>
      <label>Chanson 2:</label>
      <select name="song2-${points}" required style="display: block; margin: 5px 0; padding: 5px; width: 200px;">
        <option value="">Sélectionner une musique</option>
        ${songOptions}
      </select>
    `;
  }
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
      <h2>FORMULAIRE DE TEST</h2>
      <button id="back-to-list-btn" type="button" style="margin-bottom: 20px; background: blue; color: white; padding: 10px;">← Retour à la liste</button>
      <form id="newDuelForm">
        <h3>Créer une nouvelle grille de ${modeText}</h3>
        <label for="duelName">Nom de la grille:</label>
        <input type="text" id="duelName" name="duelName" required style="display: block; margin: 10px 0; padding: 5px;">
  `;

  DUEL_POINTS_CATEGORIES.forEach(points => {
    formHtml += `
      <div class="point-category">
        <h4>${points} Points</h4>
        <label>Thème:</label>
        <input type="text" name="theme-${points}" required style="display: block; margin: 5px 0; padding: 5px;">
        ${generateSongSelectionHtml(points, lyricsFiles)}
      </div>
    `;
  });

  formHtml += `
        <button type="submit" style="background: green; color: white; padding: 10px 20px; margin: 10px 0;">Créer</button>
      </form>
    </div>
  `;
  
  container.innerHTML = formHtml;
  console.log("Formulaire généré avec styles de debug");
  console.log("Contenu HTML du formulaire:", formHtml.substring(0, 300) + "...");
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

  // initialise la structure points
  duelData.points = {} as any;

  // lire chaque champ proprement
  formData.forEach((value, key) => {
    if (!value) return;
    // clé "samesong" est globale
    if (key === 'samesong') {
      const v = value.toString();
      if (v) {
        duelData.sameSong = {
          title: v.replace(/\.[^.]*$/, ''), // retirer extension pour titre
          artist: 'Inconnu',
          lyricsFile: v
        } as any;
      }
      return;
    }

    // clé standard : theme-50, song1-30, song2-30, duelName, ...
    const parts = key.split('-');
    const fieldName = parts[0];
    const points = parts.length > 1 ? parts[1] : null;

    if (fieldName === 'duelName') {
      duelData.name = value.toString();
      return;
    }

    if (!points) return; 

    if (!duelData.points) duelData.points = {} as any;
    if (!duelData.points[points]) duelData.points[points] = {} as any;

    if (fieldName === 'theme') {
      (duelData.points[points] as any).theme = value.toString();
    } else if (fieldName === 'song1' || fieldName === 'song2') {
      if (!(duelData.points[points] as any).songs) {
        (duelData.points[points] as any).songs = soloMode ? [{}] : [{}, {}];
      }
      const index = fieldName === 'song1' ? 0 : 1;
      if (soloMode) {
        // en solo, on ne considère que song1 (index 0)
        (duelData.points[points] as any).songs[0] = { lyricsFile: value.toString() } as any;
      } else {
        (duelData.points[points] as any).songs[index] = { lyricsFile: value.toString() } as any;
      }
    }
  })

    const newDuel: Duel = {
      id: new Date().getTime().toString(),
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
    console.log("showCreateForm appelé");
    
    const formContainer = document.getElementById("PrepGrille");
    const listContent = document.querySelector('.duels-list') as HTMLElement;
    const alertContent = document.querySelector('.alert') as HTMLElement;
    const menuButton = document.querySelector('.button_prep_grille') as HTMLElement;

    console.log("Form container:", formContainer);
    console.log("Form container innerHTML:", formContainer?.innerHTML);
    
    // Si le formulaire est vide, le régénérer
    if (formContainer && (!formContainer.innerHTML || formContainer.innerHTML.trim() === '')) {
        console.log("Le formulaire est vide, régénération...");
        
        // Essayer de récupérer les fichiers locaux, sinon utiliser une liste par défaut
        getLyricsListLocal().then(lyricsFiles => {
            console.log("Fichiers récupérés avec succès:", lyricsFiles);
            if (lyricsFiles.length > 0) {
                renderCreateDuelForm(lyricsFiles);
            } else {
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
    } else {
        // Le formulaire a déjà du contenu, juste l'afficher
        showFormWithStyles(formContainer);
    }

    // Cacher les autres éléments
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
        console.log("Formulaire affiché avec styles");
        console.log("Contenu final:", formContainer.innerHTML.substring(0, 100) + "...");
    }
}

/**
 * Version de fallback pour créer le formulaire avec message d'erreur
 */
function renderCreateDuelFormWithError(fallbackFiles: string[]): void {
    const container = document.getElementById("PrepGrille");
    if (!container) return;

    const soloMode = isSoloMode();
    const modeText = soloMode ? 'solo' : 'duel';

    let formHtml = `
        <div class="form-container">
            <h2>FORMULAIRE DE TEST (MODE DÉBOGAGE)</h2>
            <div style="background: yellow; padding: 10px; margin: 10px 0; border: 1px solid orange;">
                ⚠️ Erreur : Impossible de charger la liste des musiques. Utilisation d'exemples.
            </div>
            <button id="back-to-list-btn" type="button" style="margin-bottom: 20px; background: blue; color: white; padding: 10px;">← Retour à la liste</button>
            <form id="newDuelForm" style="border: 1px solid green; padding: 15px;">
                <h3>Créer une nouvelle grille de ${modeText}</h3>
                <label for="duelName">Nom de la grille:</label>
                <input type="text" id="duelName" name="duelName" required style="display: block; margin: 10px 0; padding: 5px;">
    `;

    DUEL_POINTS_CATEGORIES.forEach(points => {
        formHtml += `
            <div class="point-category" style="border: 1px solid blue; margin: 10px 0; padding: 10px;">
                <h4>${points} Points</h4>
                <label>Thème:</label>
                <input type="text" name="theme-${points}" required style="display: block; margin: 5px 0; padding: 5px;">
                ${generateSongSelectionHtmlFallback(points, fallbackFiles)}
            </div>
        `;
    });

    formHtml += `
            <button type="submit" style="background: green; color: white; padding: 10px 20px; margin: 10px 0;">Créer (Mode Test)</button>
        </form>
        </div>
    `;
    
    container.innerHTML = formHtml;
    console.log("Formulaire de fallback créé");
}

/**
 * Version fallback pour la génération des sélections de chansons
 */
function generateSongSelectionHtmlFallback(points: number, fallbackFiles: string[]): string {
  const soloMode = isSoloMode();
  const songOptions = fallbackFiles.map(file => `<option value="${file}">${file} (exemple)</option>`).join('');

  if (soloMode) {
    return `
      <label>Chanson:</label>
      <select name="song1-${points}" required style="display: block; margin: 5px 0; padding: 5px; width: 200px;">
        <option value="">Sélectionner une musique</option>
        ${songOptions}
      </select>
    `;
  } else {
    return `
      <label>Chanson 1:</label>
      <select name="song1-${points}" required style="display: block; margin: 5px 0; padding: 5px; width: 200px;">
        <option value="">Sélectionner une musique</option>
        ${songOptions}
      </select>
      <label>Chanson 2:</label>
      <select name="song2-${points}" required style="display: block; margin: 5px 0; padding: 5px; width: 200px;">
        <option value="">Sélectionner une musique</option>
        ${songOptions}
      </select>
    `;
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

document.addEventListener('click', (event) => {
    const target = event.target as HTMLElement;
    console.log("Clic détecté sur:", target.id, target);
    
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

// === INITIALISATION ===
document.addEventListener('DOMContentLoaded', async () => {
    console.log("DOM chargé, initialisation...");
    
    loadDuelsFromStorage();
    
    // D'abord rendre la liste pour créer/recréer PrepGrille
    renderDuelList();
    
    // Attendre un peu que le DOM soit mis à jour
    setTimeout(async () => {
        // Vérifier que PrepGrille existe maintenant
        let prepGrilleContainer = document.getElementById("PrepGrille");
        console.log("PrepGrille après renderDuelList:", prepGrilleContainer);
        
        if (!prepGrilleContainer) {
            console.error("PrepGrille TOUJOURS non trouvé après renderDuelList !");
            return;
        }

        try {
            // Utiliser notre nouvelle fonction locale
            const lyricsFiles = await getLyricsListLocal();
            console.log("Fichiers lyrics chargés avec succès:", lyricsFiles);
            
            if (lyricsFiles.length > 0) {
                renderCreateDuelForm(lyricsFiles);
            } else {
                throw new Error("Aucun fichier lyrics trouvé");
            }
            
        } catch (error) {
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
        console.log("Contenu du PrepGrille:", containerAfter?.innerHTML.substring(0, 200));
        
    }, 200);
});
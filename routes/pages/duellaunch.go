package game

import (
	"encoding/json"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type LyricsData struct {
	Titre   string `json:"titre"`
	Artiste string `json:"artiste"`
	Parole  string `json:"parole"`
}

// DisplayDuel affiche le contenu d'un duel préparé
func DisplayDuel(w http.ResponseWriter, r *http.Request) {
	fmt.Println(">>> Requête reçue pour /duel-game, traitement par DisplayDuel...")

	duelID := r.URL.Query().Get("duelId")
	if duelID == "" {
		http.Error(w, "ID de duel manquant dans les paramètres de la requête", http.StatusBadRequest)
		return
	}

	// Convertir l'ID en entier
	id, err := strconv.Atoi(duelID)
	if err != nil {
		http.Error(w, "ID de duel invalide", http.StatusBadRequest)
		return
	}

	// Trouver le duel correspondant
	var selectedDuel *Duel
	for i := range duels {
		if duels[i].ID == id {
			selectedDuel = &duels[i]
			break
		}
	}

	if selectedDuel == nil {
		http.Error(w, "Duel non trouvé", http.StatusNotFound)
		return
	}

	// Vérifier l'existence des paroles pour "La Même Chanson"
	sameSongLyricsExists := false
	if selectedDuel.SameSong.LyricsFile != nil && *selectedDuel.SameSong.LyricsFile != "" {
		filePath := filepath.Join(paroleDataPath, *selectedDuel.SameSong.LyricsFile)
		if _, err := os.Stat(filePath); err == nil {
			sameSongLyricsExists = true
		}
	}

	templateData := struct {
		Duel                 *Duel
		LevelsOrder          []string
		LyricsExists         map[string]map[int]bool
		SameSongLyricsExists bool
	}{
		Duel:                 selectedDuel,
		LevelsOrder:          []string{"50", "40", "30", "20", "10"},
		LyricsExists:         make(map[string]map[int]bool),
		SameSongLyricsExists: sameSongLyricsExists,
	}

	for _, level := range templateData.LevelsOrder {
		templateData.LyricsExists[level] = make(map[int]bool)
		if pointLevel, exists := selectedDuel.Points[level]; exists {
			for i, song := range pointLevel.Songs {
				if song.LyricsFile != nil && *song.LyricsFile != "" {
					filePath := filepath.Join(paroleDataPath, *song.LyricsFile)
					if _, err := os.Stat(filePath); err == nil {
						templateData.LyricsExists[level][i] = true
					}
				}
			}
		}
	}

	if r.Method == http.MethodPost {
		action := r.FormValue("action")
		switch action {
		case "start_session":
			http.Redirect(w, r, fmt.Sprintf("/duel-game?duelId=%d", id), http.StatusSeeOther)
			return
		case "export":
			http.Redirect(w, r, fmt.Sprintf("/api/export-duel-server/%d", id), http.StatusSeeOther)
			return
		}
	}

	tmpl := `
		<!DOCTYPE html>
		<html lang="fr">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Duel: {{.Duel.Name}}</title>
			<link rel="stylesheet" href="/asset/scss/style.css">
		</head>
		<body>
    <div class="duelContainer">
        <div class="duel-header">
            <h1>{{.Duel.Name}}</h1>
            <div class="metadata">
                <p>Créé le: {{.Duel.CreatedAt.Format "02/01/2006 à 15:04"}}</p>
                {{with .Duel.UpdatedAt}}
                    <p>Mis à jour le: {{.Format "02/01/2006 à 15:04"}}</p>
                {{end}}
            </div>
        </div>

        <section class="duel-form">
            <!-- Section "La Même Chanson" tout en haut -->
            <div class="same-song-section">
                <button class="same-song-button" onclick="toggleElement('same-song-details')">
                    🎵 La Même Chanson
                </button>
                <div id="same-song-details" class="hidden-element song-card" style="max-width: 500px; margin: 0 auto;">
                    <div class="song-info">
                        <div class="song-title">{{.Duel.SameSong.Title}}</div>
                        <div class="song-artist">par {{.Duel.SameSong.Artist}}</div>
                    </div>
                    {{if .Duel.SameSong.AudioURL}}
                    <div>
                        <strong>Audio:</strong> <a href="{{.Duel.SameSong.AudioURL}}" target="_blank">Écouter</a>
                    </div>
                    {{end}}
                    <div class="lyrics-status">
                        {{if .SameSongLyricsExists}}
                        <span class="lyrics-available">✓ Paroles disponibles</span>
                        {{else}}
                        <span class="lyrics-missing">✗ Paroles non disponibles</span>
                        {{end}}
                    </div>
                </div>
            </div>

            <div class="points-grid">
                {{range .LevelsOrder}}
                {{$level := .}}
                {{$pointLevel := index $.Duel.Points $level}}
                <button class="point-button" onclick="toggleLevelSongs('level-{{$level}}')">
                    <div>{{$level}} Points</div>
                    <div style="font-size: 14px; margin-top: 5px;">{{$pointLevel.Theme}}</div>
                </button>
                {{end}}
            </div>

            {{range .LevelsOrder}}
            {{$level := .}}
            {{$pointLevel := index $.Duel.Points $level}}
            <div id="level-{{$level}}" class="hidden-element songs-for-level">
                <div class="level-section">
                    <div class="level-header">
                        <h3>{{$level}} Points - {{$pointLevel.Theme}}</h3>
                    </div>
                    <div class="level-songs-container">
                        {{range $index, $song := $pointLevel.Songs}}
                        <div class="song-card" onclick="previewSong('{{$song.Title}}', '{{$song.Artist}}')">
                            <div class="song-info">
                                <div class="song-title">{{$song.Title}}</div>
                                <div class="song-artist">par {{$song.Artist}}</div>
                            </div>
                            {{if $song.AudioURL}}
                            <div>
                                <strong>Audio:</strong> <a href="{{$song.AudioURL}}" target="_blank">Écouter</a>
                            </div>
                            {{end}}
                            <div class="lyrics-status">
                                {{if index (index $.LyricsExists $level) $index}}
                                <span class="lyrics-available">✓ Paroles disponibles</span>
                                {{else}}
                                <span class="lyrics-missing">✗ Paroles non disponibles</span>
                                {{end}}
                            </div>
                            <div class="song-actions">
                                <button class="btn-preview" onclick="event.stopPropagation(); previewSong('{{$song.Title}}', '{{$song.Artist}}')">
                                    🎵 Aperçu
                                </button>
                                <button class="btn-select" onclick="event.stopPropagation(); selectSong('{{$level}}', {{$index}}, '{{$song.Title}}', '{{$song.Artist}}')">
                                    Sélectionner
                                </button>
                            </div>
                        </div>
                        {{end}}
                    </div>
                </div>
            </div>
            {{end}}
        </section>

        <section class="songSelect duel-container">
			<div id="music-player" class="music-player" style="display: none;">
				<h4 id="current-song-info">Aucune chanson sélectionnée</h4>
				<div class="audio-controls">
					<audio id="audio-player" controls style="width: 100%;">
						Votre navigateur ne supporte pas l'élément audio.
					</audio>
				</div>
				// <button onclick="toggleLyrics()" class="btn btn-secondary">Afficher/Masquer les paroles</button>
			</div>

			<div id="lyrics-container" class="lyrics-container" style="display: none;">
				<h4>Paroles</h4>
				<div id="lyrics-text" class="lyrics-text"></div>
			</div>
			<div class="actions">
            <form method="POST" style="display: inline;">
                <input type="hidden" name="action" value="start_session">
                <button type="submit" class="btn btn-success">Démarrer une partie</button>
            </form>
            
            // <form method="POST" style="display: inline;">
            //     <input type="hidden" name="action" value="export">
            //     <button type="submit" class="btn btn-primary">Exporter ce duel</button>
            // </form>
            
            <a href="/duel" class="btn btn-secondary">Retour aux duels</a>
        </div>
		</section>

        
    </div>

    <script>
        let currentAudio = null;
        let currentLyrics = "";
        let currentLevel = "";
        
        // Fonction pour basculer l'affichage des éléments
        function toggleElement(elementId) {
            const element = document.getElementById(elementId);
            if (element.style.display === 'none' || element.classList.contains('hidden-element')) {
                element.style.display = 'block';
                element.classList.remove('hidden-element');
            } else {
                element.style.display = 'none';
                element.classList.add('hidden-element');
            }
        }

        function toggleLevelSongs(levelId) {
            toggleElement(levelId);
        }

        // Fonction pour obtenir l'URL d'aperçu de YouTube
        async function getYouTubePreview(title, artist) {
            try {
                // Utilisation de l'API YouTube Search (nécessite une clé API)
                const query = encodeURIComponent(title + " " + artist + " instrumental");
                const searchUrl = "https://www.googleapis.com/youtube/v3/search?part=snippet&q=" + query + "&type=video&key=YOUR_YOUTUBE_API_KEY";
                
                // Pour le moment, utilisons une URL de démonstration
                // En production, vous devrez remplacer par votre clé API YouTube
                console.log("Recherche YouTube pour:", title, artist);
                return null; // Retourner null pour utiliser un fichier audio local ou une autre source
            } catch (error) {
                console.error("Erreur lors de la recherche YouTube:", error);
                return null;
            }
        }

        // Fonction pour prévisualiser une chanson
        async function previewSong(title, artist) {
            console.log("Aperçu de la chanson:", title, "par", artist);
            
            // Afficher le lecteur
            document.getElementById('music-player').style.display = 'block';
            document.getElementById('current-song-info').textContent = "Aperçu: " + title + " par " + artist;
            
            // Arrêter l'audio précédent s'il y en a un
            if (currentAudio) {
                currentAudio.pause();
                currentAudio.currentTime = 0;
            }
            
            // Pour l'instant, utilisons un fichier audio de démonstration
            // En production, vous pourrez intégrer l'API YouTube ou Spotify
            const audioPlayer = document.getElementById('audio-player');
            audioPlayer.src = "/demo-audio.mp3"; // Fichier de démonstration
            
            // Cacher les paroles pendant l'aperçu
            document.getElementById('lyrics-container').style.display = 'none';
        }

        // Fonction pour sélectionner une chanson et charger les paroles
        async function selectSong(level, songIndex, title, artist) {
            console.log("Chanson sélectionnée:", title, "par", artist, "niveau:", level, "index:", songIndex);
            
            currentLevel = level;
            
            // Afficher le lecteur
            document.getElementById('music-player').style.display = 'block';
            document.getElementById('current-song-info').textContent = title + " par " + artist + " (" + level + " points)";
            
            // Charger les paroles
            try {
                const response = await fetch('/api/get-lyrics/' + level + '/' + songIndex);
                if (response.ok) {
                    const lyricsData = await response.json();
                    currentLyrics = lyricsData.parole || lyricsData.lyrics || "Paroles non disponibles";
                    displayMaskedLyrics(currentLyrics, parseInt(level));
                } else {
                    currentLyrics = "Paroles non disponibles";
                    document.getElementById('lyrics-text').textContent = currentLyrics;
                }
            } catch (error) {
                console.error("Erreur lors du chargement des paroles:", error);
                currentLyrics = "Erreur de chargement des paroles";
                document.getElementById('lyrics-text').textContent = currentLyrics;
            }
            
            // Charger l'audio (instrumental si possible)
            const audioPlayer = document.getElementById('audio-player');
            // Ici vous pourrez intégrer l'API de votre choix pour l'audio instrumental
            audioPlayer.src = "/demo-instrumental.mp3"; // Fichier de démonstration
            
            // Afficher les paroles
            document.getElementById('lyrics-container').style.display = 'block';
        }

        // Fonction pour masquer les paroles selon le niveau de points
        function displayMaskedLyrics(lyrics, points) {
            if (!lyrics) return;
            
            // Calculer le pourcentage de masquage selon les points
            let maskPercentage;
            switch(points) {
                case 50: maskPercentage = 0.8; break;  // 80% masqué pour 50 points
                case 40: maskPercentage = 0.6; break;  // 60% masqué pour 40 points
                case 30: maskPercentage = 0.4; break;  // 40% masqué pour 30 points
                case 20: maskPercentage = 0.2; break;  // 20% masqué pour 20 points
                case 10: maskPercentage = 0.1; break;  // 10% masqué pour 10 points
                default: maskPercentage = 0.3; break;
            }
            
            const words = lyrics.split(' ');
            const wordsToMask = Math.floor(words.length * maskPercentage);
            
            // Créer un array d'indices à masquer de manière aléatoire
            const indicesToMask = [];
            while (indicesToMask.length < wordsToMask) {
                const randomIndex = Math.floor(Math.random() * words.length);
                if (!indicesToMask.includes(randomIndex)) {
                    indicesToMask.push(randomIndex);
                }
            }
            
            // Appliquer le masquage
            const maskedWords = words.map((word, index) => {
                if (indicesToMask.includes(index)) {
					return '<span class="masked" data-word="' + word + '">█'.repeat(word.length) + '</span>';
                }
                return word;
            });
            
            document.getElementById('lyrics-text').innerHTML = maskedWords.join(' ');
        }

        // Fonction pour basculer l'affichage des paroles
        function toggleLyrics() {
            const lyricsContainer = document.getElementById('lyrics-container');
            if (lyricsContainer.style.display === 'none') {
                lyricsContainer.style.display = 'block';
            } else {
                lyricsContainer.style.display = 'none';
            }
        }
    </script>
		</body>
		</html>`

	t, err := template.New("duel").Parse(tmpl)
	if err != nil {
		http.Error(w, "Erreur lors du parsing du template", http.StatusInternalServerError)
		return
	}

	if err := t.Execute(w, templateData); err != nil {
		fmt.Println("Erreur lors de l'exécution du template :", err)
		http.Error(w, "Erreur lors de l'exécution du template", http.StatusInternalServerError)
		return
	}

}

// GetLyrics - Handler pour récupérer les paroles d'une chanson
func GetLyrics(w http.ResponseWriter, r *http.Request) {
	// Extraire les paramètres de l'URL
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 {
		http.Error(w, "Paramètres manquants", http.StatusBadRequest)
		return
	}

	level := parts[3]
	songIndexStr := parts[4]

	songIndex, err := strconv.Atoi(songIndexStr)
	if err != nil {
		http.Error(w, "Index de chanson invalide", http.StatusBadRequest)
		return
	}

	// Trouver le duel actuel (vous devrez adapter selon votre logique)
	// Pour cet exemple, je prends le premier duel disponible
	if len(duels) == 0 {
		http.Error(w, "Aucun duel disponible", http.StatusNotFound)
		return
	}

	duel := &duels[0] // Adaptez selon votre logique de session

	pointLevel, ok := duel.Points[level]
	if !ok {
		http.Error(w, "Niveau de points invalide", http.StatusBadRequest)
		return
	}

	if songIndex < 0 || songIndex >= len(pointLevel.Songs) {
		http.Error(w, "Index de chanson invalide", http.StatusBadRequest)
		return
	}

	song := pointLevel.Songs[songIndex]

	// Charger les paroles depuis le fichier JSON
	var lyricsData LyricsData
	if song.LyricsFile != nil && *song.LyricsFile != "" {
		filePath := filepath.Join(paroleDataPath, *song.LyricsFile)
		if !filepath.IsAbs(filePath) {
			content, err := os.ReadFile(filePath)
			if err == nil {
				if err := json.Unmarshal(content, &lyricsData); err == nil {
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(lyricsData)
					return
				}
			}
		}
	}

	// Si pas de paroles trouvées
	lyricsData = LyricsData{
		Titre:   song.Title,
		Artiste: song.Artist,
		Parole:  "Paroles non disponibles",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lyricsData)
}

func MaskLyrics(lyrics string, points int) string {
	if lyrics == "" {
		return lyrics
	}

	sections := splitLyricsBySections(lyrics)

	var targetSection string
	switch {
	case points >= 40:
		targetSection = "Couplet 2"
	case points >= 10 && points <= 30:
		targetSection = "Refrain"
	default:
		targetSection = "" // Pas de section ciblée
	}

	// Si la section existe, appliquer le masquage uniquement sur elle
	if content, ok := sections[targetSection]; ok {
		sections[targetSection] = maskSectionContent(content, points)
	}

	// Reconstruire les paroles avec sections masquées
	var rebuiltLyrics strings.Builder
	for section, content := range sections {
		rebuiltLyrics.WriteString("[" + section + "]\n")
		rebuiltLyrics.WriteString(content + "\n\n")
	}

	return strings.TrimSpace(rebuiltLyrics.String())
}

// Découpe les paroles par section
func splitLyricsBySections(lyrics string) map[string]string {
	lines := strings.Split(lyrics, "\n")
	currentSection := "Intro"
	sections := make(map[string]string)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentSection = strings.Trim(line, "[]")
		} else {
			sections[currentSection] += line + "\n"
		}
	}

	return sections
}

// Applique le masquage à une section de texte
func maskSectionContent(text string, points int) string {
	var maskPercentage float64
	switch points {
	case 50:
		maskPercentage = 0.8
	case 40:
		maskPercentage = 0.6
	case 30:
		maskPercentage = 0.4
	case 20:
		maskPercentage = 0.2
	case 10:
		maskPercentage = 0.1
	default:
		maskPercentage = 0.3
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return text
	}

	wordsToMask := int(float64(len(words)) * maskPercentage)
	if wordsToMask == 0 && maskPercentage > 0 {
		wordsToMask = 1
	}

	rand.Seed(time.Now().UnixNano())
	indicesToMask := make(map[int]bool)
	for len(indicesToMask) < wordsToMask && len(indicesToMask) < len(words) {
		randomIndex := rand.Intn(len(words))
		indicesToMask[randomIndex] = true
	}

	for i, word := range words {
		if indicesToMask[i] {
			words[i] = strings.Repeat("_", len(word))
		}
	}

	return strings.Join(words, " ")
}

// StartSong démarre une chanson avec ses paroles
func StartSong(sessionID string, level string, songIndex int) (*GameSession, error) {
	session, ok := gameSessions[sessionID]
	if !ok {
		return nil, fmt.Errorf("session non trouvée : %s", sessionID)
	}

	var duel *Duel
	for i := range duels {
		if duels[i].ID == session.DuelID {
			duel = &duels[i]
			break
		}
	}
	if duel == nil {
		return nil, fmt.Errorf("duel non trouvé pour cette session")
	}

	pointLevel, ok := duel.Points[level]
	if !ok {
		return nil, fmt.Errorf("niveau de points invalide : %s", level)
	}
	if songIndex < 0 || songIndex >= len(pointLevel.Songs) {
		return nil, fmt.Errorf("index de chanson invalide : %d", songIndex)
	}
	song := pointLevel.Songs[songIndex]
	session.CurrentSong = &song

	session.LyricsContent = "Paroles non disponibles."
	if song.LyricsFile != nil && *song.LyricsFile != "" {
		filePath := filepath.Join(paroleDataPath, *song.LyricsFile)
		if !filepath.IsAbs(filePath) {
			content, err := os.ReadFile(filePath)
			if err == nil {
				session.LyricsContent = string(content)
			}
		}
	}

	session.LyricsVisible = true

	return session, nil
}

func SetLyricsVisibility(sessionID string, visible bool) (*GameSession, error) {
	session, ok := gameSessions[sessionID]
	if !ok {
		return nil, fmt.Errorf("session non trouvée")
	}

	session.LyricsVisible = visible
	return session, nil
}

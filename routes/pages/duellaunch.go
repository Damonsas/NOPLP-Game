package game

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type LyricsData struct {
	Titre   string              `json:"titre"`
	Artiste string              `json:"artiste"`
	Parole  map[string][]string `json:"parole"`
}

type LyricsStructure struct {
	Titre   string              `json:"titre"`
	Artiste string              `json:"artiste"`
	Parole  map[string][]string `json:"parole"`
}

func DisplayDuel(w http.ResponseWriter, r *http.Request) {
	fmt.Println(">>> Requête reçue pour /duel-game, traitement par DisplayDuel...")

	duelID := r.URL.Query().Get("duelId")
	if duelID == "" {
		http.Error(w, "ID de duel manquant dans les paramètres de la requête", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(duelID)
	if err != nil {
		http.Error(w, "ID de duel invalide", http.StatusBadRequest)
		return
	}

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
			<link rel="apple-touch-icon" sizes="180x180" href="../asset/ressource/img/favicon/apple-touch-icon.png">
    		<link rel="icon" type="image/png" sizes="32x32" href="../asset/ressource/img/favicon/favicon-32x32.png">
    		<link rel="icon" type="image/png" sizes="16x16" href="../asset/ressource/img/favicon/favicon-16x16.png">
			<script defer src="/asset/js/script.js"></script>
			<script src="/asset/js/vocal.js"></script>
    		<script type="module" src="/asset/js/duel.js"></script>
    		<script type="module" src="/asset/js/ui.js"></script>
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
            <div class="same-song-section">
                <button class="same-song-button" onclick="toggleElement('same-song-details')" >
                    🎵C'est La Même Chanson
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
                <button class="point-button fade-in-left-normal" onclick="toggleLevelSongs('level-{{$level}}')">
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
				
			</div>

			<div id="lyrics-container" class="lyrics-container" style="display: none;">
				<h4>Paroles</h4>
				<div id="lyrics-text" class="lyrics-text"></div>
			</div>
			<div class="actions" id="action-buttons">
				<button class="startLyricsBtn" class="start-lyrics-button">Demarrer</button>

				<a class="startLyricsBtn" href="/duel" class="btn btn-secondary">Retour aux duels</a>
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

        async function previewSong(title, artist) {
            console.log("Aperçu de la chanson:", title, "par", artist);
            
            document.getElementById('music-player').style.display = 'block';
            document.getElementById('current-song-info').textContent = "Aperçu: " + title + " par " + artist;
            
            if (currentAudio) {
                currentAudio.pause();
                currentAudio.currentTime = 0;
            }
            
            
            const audioPlayer = document.getElementById('audio-player');
            audioPlayer.src = "/demo-audio.mp3"; // Fichier de démonstration
            
            
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

func convertStructuredLyricsToText(paroleMap map[string][]string) string {
	var result strings.Builder

	sectionOrder := []string{"couplet1", "refrain1", "couplet2", "refrain2", "outro"}

	for _, sectionName := range sectionOrder {
		if lines, exists := paroleMap[sectionName]; exists {
			result.WriteString("[" + strings.Title(sectionName) + "]\n")
			for _, line := range lines {
				result.WriteString(line + "\n")
			}
			result.WriteString("\n")
		}
	}

	for sectionName, lines := range paroleMap {
		found := false
		for _, orderedSection := range sectionOrder {
			if orderedSection == sectionName {
				found = true
				break
			}
		}
		if !found {
			result.WriteString("[" + strings.Title(sectionName) + "]\n")
			for _, line := range lines {
				result.WriteString(line + "\n")
			}
			result.WriteString("\n")
		}
	}

	return strings.TrimSpace(result.String())
}

func GetLyricsFromAPI(trackID string) (string, error) {
	endpoint := fmt.Sprintf("track.lyrics.get?track_id=%s&apikey=%s", trackID, musixmatchAPIKey)
	resp, err := http.Get(musixmatchBaseURL + endpoint)
	if err != nil {
		return "", fmt.Errorf("erreur HTTP: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result struct {
		Message struct {
			Body struct {
				Lyrics struct {
					LyricsBody string `json:"lyrics_body"`
				} `json:"lyrics"`
			} `json:"body"`
		} `json:"message"`
	}
	json.Unmarshal(body, &result)
	if result.Message.Body.Lyrics.LyricsBody == "" {
		return "", errors.New("paroles non trouvées")
	}
	return result.Message.Body.Lyrics.LyricsBody, nil
}

func GetLyricsdata(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	level := vars["level"]
	songIndexStr := vars["songIndex"]
	songIndex, err := strconv.Atoi(songIndexStr)
	if err != nil {
		http.Error(w, "Index de chanson invalide", http.StatusBadRequest)
		return
	}

	duelIDStr := r.URL.Query().Get("duelId")
	var selectedDuel *Duel
	if duelIDStr != "" {
		duelID, err := strconv.Atoi(duelIDStr)
		if err == nil {
			for i := range duels {
				if duels[i].ID == duelID {
					selectedDuel = &duels[i]
					break
				}
			}
		}
	}
	if selectedDuel == nil && len(duels) > 0 {
		selectedDuel = &duels[0]
	}
	if selectedDuel == nil {
		http.Error(w, "Aucun duel disponible", http.StatusNotFound)
		return
	}

	pointLevel, ok := selectedDuel.Points[level]
	if !ok {
		http.Error(w, "Niveau de points invalide", http.StatusBadRequest)
		return
	}
	if songIndex < 0 || songIndex >= len(pointLevel.Songs) {
		http.Error(w, "Index de chanson invalide", http.StatusBadRequest)
		return
	}
	song := pointLevel.Songs[songIndex]

	fmt.Printf("DEBUG - Recherche paroles pour: %s - %s\n", song.Title, song.Artist)
	fmt.Printf("DEBUG - LyricsFile: %v\n", song.LyricsFile)

	if song.LyricsFile != nil && *song.LyricsFile != "" {
		filePath := filepath.Join(paroleDataPath, *song.LyricsFile)
		fmt.Printf("DEBUG - Chemin du fichier: %s\n", filePath)

		if _, err := os.Stat(filePath); err == nil {
			content, err := os.ReadFile(filePath)
			if err != nil {
				fmt.Printf("DEBUG - Erreur lecture fichier: %v\n", err)
			} else {
				fmt.Printf("DEBUG - Contenu fichier lu, taille: %d bytes\n", len(content))

				var structuredLyrics LyricsStructure
				if err := json.Unmarshal(content, &structuredLyrics); err == nil {
					lyricsText := convertStructuredLyricsToText(structuredLyrics.Parole)
					lyricsData := map[string]interface{}{
						"titre":   structuredLyrics.Titre,
						"artiste": structuredLyrics.Artiste,
						"parole":  lyricsText,
					}
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(lyricsData)
					return
				}

				var lyricsData map[string]interface{}
				if err := json.Unmarshal(content, &lyricsData); err != nil {
					fmt.Printf("DEBUG - Erreur parsing JSON: %v\n", err)
					lyricsData = map[string]interface{}{
						"titre":   song.Title,
						"artiste": song.Artist,
						"parole":  string(content),
					}
				} else {
					if parole, exists := lyricsData["parole"]; exists {
						if paroleStr, ok := parole.(string); ok && paroleStr != "" {
						} else {
							lyricsData["parole"] = "Format de paroles non reconnu"
						}
					} else {
						lyricsData["parole"] = "Clé 'parole' manquante dans le JSON"
					}
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(lyricsData)
				return
			}
		} else {
			fmt.Printf("DEBUG - Fichier n'existe pas: %s\n", filePath)
		}
	}

	trackID, err := SearchTrack(song.Title, song.Artist)
	if err != nil {
		fmt.Printf("DEBUG - Erreur recherche track: %v\n", err)
		lyricsData := map[string]interface{}{
			"titre":   song.Title,
			"artiste": song.Artist,
			"parole":  "Paroles non disponibles",
			"erreur":  err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(lyricsData)
		return
	}

	lyrics, err := GetLyricsFromAPI(trackID)
	if err != nil {
		fmt.Printf("DEBUG - Erreur API externe: %v\n", err)
		lyricsData := map[string]interface{}{
			"titre":   song.Title,
			"artiste": song.Artist,
			"parole":  "Paroles non disponibles",
			"erreur":  err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(lyricsData)
		return
	}

	if song.LyricsFile != nil && *song.LyricsFile != "" {
		filePath := filepath.Join(paroleDataPath, *song.LyricsFile)
		os.WriteFile(filePath, []byte(lyrics), 0644)
		fmt.Printf("DEBUG - Paroles sauvegardées dans %s\n", filePath)
	}

	lyricsData := map[string]interface{}{
		"titre":   song.Title,
		"artiste": song.Artist,
		"parole":  lyrics,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lyricsData)

	lyrics, err = GetLyricsFromAPI(trackID)
	if err != nil {
		fmt.Printf("DEBUG - Erreur API externe: %v\n", err)
		lyricsData := map[string]interface{}{
			"titre":   song.Title,
			"artiste": song.Artist,
			"parole":  "Paroles non disponibles",
			"erreur":  err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(lyricsData)
		return
	}

	if song.LyricsFile != nil && *song.LyricsFile != "" {
		filePath := filepath.Join(paroleDataPath, *song.LyricsFile)
		os.WriteFile(filePath, []byte(lyrics), 0644)
		fmt.Printf("DEBUG - Paroles sauvegardées dans %s\n", filePath)
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
		targetSection = ""
	}

	if content, ok := sections[targetSection]; ok {
		lines := strings.Split(strings.TrimSpace(content), "\n")
		sections[targetSection] = MaskedSectionContent(targetSection, lines, points)
	}

	var rebuiltLyrics strings.Builder
	for section, content := range sections {
		rebuiltLyrics.WriteString("[" + section + "]\n")
		rebuiltLyrics.WriteString(content + "\n\n")
	}

	return strings.TrimSpace(rebuiltLyrics.String())
}

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

func MaskedSectionContent(sectionName string, lines []string, points int) string {
	showSection := false

	switch points {
	case 50, 40:
		showSection = strings.ToLower(sectionName) == "couplet1"
	case 30, 20, 10:
		showSection = !strings.Contains(strings.ToLower(sectionName), "refrain")
	default:
		showSection = true
	}

	if showSection {
		return strings.Join(lines, "<br>")
	}

	var maskedLines []string
	for _, line := range lines {
		words := strings.Fields(line)
		maskedWords := make([]string, len(words))
		for i, word := range words {
			maskedWords[i] = strings.Repeat("█", len(word))
		}
		maskedLines = append(maskedLines, strings.Join(maskedWords, " "))
	}

	return strings.Join(maskedLines, "<br>")
}

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

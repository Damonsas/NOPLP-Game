package game

type LyricsStructure struct {
	Titre   string              `json:"titre"`
	Artiste string              `json:"artiste"`
	Parole  map[string][]string `json:"parole"`
}

type LyricsData struct {
	Titre   string              `json:"titre"`
	Artiste string              `json:"artiste"`
	Parole  map[string][]string `json:"parole"`
}

type LyricsFileData struct {
	Titre   string              `json:"titre"`
	Artiste string              `json:"artiste"`
	Parole  map[string][]string `json:"parole"`
}

type LyricsCheckResponse struct {
	Exists  bool   `json:"exists"`
	Content string `json:"content,omitempty"`
}

// func GetLyricsData(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("========================================")
// 	fmt.Println(">>> API REACHED: GetLyricsData appelée !")
// 	fmt.Println("URL complète:", r.URL.String())
// 	fmt.Println("========================================")

// 	song, err := extractSongFromRequest(r)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	fmt.Printf("DEBUG - Recherche paroles pour: %s - %s\n", song.Artiste, song.Titre)
// 	fmt.Printf("DEBUG - LyricsFile: %v\n", song.LyricsFile)

// 	if lyrics, found := tryLoadFromFile(song); found {
// 		respondWithLyrics(w, song, lyrics)
// 		return
// 	}

// 	if lyrics, err := tryLoadFromAPI(song); err == nil {
// 		if song.LyricsFile != nil && *song.LyricsFile != "" {
// 			savelyrics(song, lyrics)
// 		}
// 		respondWithLyrics(w, song, lyrics)
// 		return
// 	}

// 	respondWithError(w, song, "Paroles non disponibles")
// }

// func extractSongFromRequest(r *http.Request) (*Song, error) {
// 	vars := mux.Vars(r)
// 	level := vars["level"]
// 	songIndexStr := vars["songIndex"]

// 	songIndex, err := strconv.Atoi(songIndexStr)
// 	if err != nil {
// 		return nil, errors.New("Index de chanson invalide")
// 	}

// 	selectedDuel := selectDuel(r.URL.Query().Get("duelId"))
// 	if selectedDuel == nil {
// 		return nil, errors.New("Aucun duel disponible")
// 	}

// 	pointLevel, ok := selectedDuel.Points[level]
// 	if !ok {
// 		return nil, errors.New("Niveau de points invalide")
// 	}

// 	if songIndex < 0 || songIndex >= len(pointLevel.Songs) {
// 		return nil, errors.New("Index de chanson invalide")
// 	}

// 	return &pointLevel.Songs[songIndex], nil
// }

// func selectDuel(duelIDStr string) *Duel {
// 	if duelIDStr != "" {
// 		duelID, err := strconv.Atoi(duelIDStr)
// 		if err == nil {
// 			for i := range duels {
// 				if duels[i].ID == duelID {
// 					return &duels[i]
// 				}
// 			}
// 		}
// 	}
// 	if len(duels) > 0 {
// 		return &duels[0]
// 	}
// 	return nil
// }

// func tryLoadFromFile(song *Song) (string, bool) {
// 	if song.LyricsFile == nil || *song.LyricsFile == "" {
// 		fmt.Printf("DEBUG - Pas de fichier de paroles configuré\n")
// 		return "", false
// 	}

// 	path := filepath.Join(paroleDataPath, *song.LyricsFile)
// 	absPath, _ := filepath.Abs(path)
// 	fmt.Printf("DEBUG - Chemin du fichier: %s\n", absPath)

// 	content, err := os.ReadFile(absPath)
// 	if err != nil {
// 		fmt.Printf("DEBUG - Fichier n'existe pas: %s (erreur: %v)\n", absPath, err)
// 		return "", false
// 	}

// 	fmt.Printf("DEBUG - Contenu fichier lu, taille: %d bytes\n", len(content))
// 	fmt.Printf("DEBUG - Contenu brut: %s\n", string(content)[:min(200, len(content))])

// 	var structuredLyrics LyricsStructure
// 	if err := json.Unmarshal(content, &structuredLyrics); err == nil {
// 		fmt.Printf("DEBUG - Structure imbriquée détectée\n")
// 		return convertStructuredLyricsToText(structuredLyrics.Parole), true
// 	}

// 	var lyricsData map[string]interface{}
// 	if err := json.Unmarshal(content, &lyricsData); err == nil {
// 		if parole, exists := lyricsData["parole"]; exists {
// 			if paroleStr, ok := parole.(string); ok && paroleStr != "" {
// 				fmt.Printf("DEBUG - Paroles trouvées dans structure simple\n")
// 				return paroleStr, true
// 			}
// 		}
// 	}

// 	return string(content), true
// }

// func tryLoadFromAPI(song *Song) (string, error) {
// 	fmt.Printf("DEBUG - Tentative de récupération via API externe...\n")
// 	trackID, err := handlers.SearchTrack(song.Titre, song.Artiste)
// 	if err != nil {
// 		fmt.Printf("DEBUG - Erreur recherche track: %v\n", err)
// 		return "", err
// 	}

// 	lyrics, err := GetLyricsFromAPI(trackID)
// 	if err != nil {
// 		fmt.Printf("DEBUG - Erreur API externe: %v\n", err)
// 		return "", err
// 	}

// 	return lyrics, nil
// }

// func savelyrics(song *Song, lyrics string) {
// 	if song.LyricsFile != nil && *song.LyricsFile != "" {
// 		filePath := filepath.Join(paroleDataPath, *song.LyricsFile)
// 		os.WriteFile(filePath, []byte(lyrics), 0644)
// 		fmt.Printf("DEBUG - Paroles sauvegardées dans %s\n", filePath)
// 	}
// }

// func respondWithLyrics(w http.ResponseWriter, song *Song, lyrics string) {
// 	lyricsData := map[string]interface{}{
// 		"titre":   song.Titre,
// 		"artiste": song.Artiste,
// 		"parole":  lyrics,
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(lyricsData)
// }

// func respondWithError(w http.ResponseWriter, song *Song, message string) {
// 	lyricsData := map[string]interface{}{
// 		"titre":   song.Titre,
// 		"artiste": song.Artiste,
// 		"parole":  message,
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusInternalServerError)
// 	json.NewEncoder(w).Encode(lyricsData)
// }

// func min(a, b int) int {
// 	if a < b {
// 		return a
// 	}
// 	return b
// }

// func GetLyricsByTitleAndArtist(Titre, Artiste string) (string, error) {
// 	// 1. On liste tous les fichiers du dossier de paroles
// 	files, err := os.ReadDir(paroleDataPath)
// 	if err != nil {
// 		return "", fmt.Errorf("erreur lecture dossier: %v", err)
// 	}

// 	// 2. On prépare les termes de recherche (minuscules et sans espaces inutiles)
// 	artisteRecherche := strings.ToLower(strings.TrimSpace(Artiste))
// 	titreRecherche := strings.ToLower(strings.TrimSpace(Titre))

// 	// On extrait le nom de famille (le dernier mot de l'artiste)
// 	parts := strings.Fields(artisteRecherche)
// 	nomDeFamille := ""
// 	if len(parts) > 0 {
// 		nomDeFamille = parts[len(parts)-1]
// 	}

// 	for _, f := range files {
// 		if f.IsDir() {
// 			continue
// 		}
// 		fileName := strings.ToLower(f.Name())

// 		// 3. LA LOGIQUE DE CORRESPONDANCE
// 		// On cherche si le fichier contient le TITRE
// 		if strings.Contains(fileName, titreRecherche) {
// 			// ET s'il contient soit l'artiste complet, soit juste le nom de famille
// 			if strings.Contains(fileName, artisteRecherche) || (nomDeFamille != "" && strings.Contains(fileName, nomDeFamille)) {
// 				// FICHIER TROUVÉ !
// 				content, err := os.ReadFile(filepath.Join(paroleDataPath, f.Name()))
// 				if err != nil {
// 					return "", err
// 				}
// 				return string(content), nil
// 			}
// 		}
// 	}

// 	return "", fmt.Errorf("paroles non trouvées pour %s - %s", Artiste, Titre)
// }

// // normalizeForSearch normalise une chaîne pour la comparaison
// // (minuscules, sans accents, sans caractères spéciaux)
// func normalizeForSearch(s string) string {
// 	s = strings.ToLower(s)
// 	s = removeAccents(s)
// 	// Remplacer les caractères spéciaux par des espaces
// 	s = regexp.MustCompile(`[^a-z0-9\s]+`).ReplaceAllString(s, " ")
// 	// Nettoyer les espaces multiples
// 	s = regexp.MustCompile(`\s+`).ReplaceAllString(s, " ")
// 	return strings.TrimSpace(s)
// }

// func GetLyricsBySongHandler(w http.ResponseWriter, r *http.Request) {
// 	Titre := r.URL.Query().Get("titre")
// 	Artiste := r.URL.Query().Get("artiste")

// 	if Titre == "" || Artiste == "" {
// 		http.Error(w, "Titre ou artiste manquant", http.StatusBadRequest)
// 		return
// 	}

// 	content, err := GetLyricsByTitleAndArtist(Titre, Artiste)
// 	if err != nil {
// 		http.Error(w, "Paroles non trouvées", http.StatusNotFound)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write([]byte(content))
// }

// func convertStructuredLyricsToText(paroleMap map[string][]string) string {
// 	var result strings.Builder

// 	sectionOrder := []string{"couplet1", "refrain1", "couplet2", "refrain2", "outro"}

// 	for _, sectionName := range sectionOrder {
// 		if lines, exists := paroleMap[sectionName]; exists {
// 			result.WriteString("[" + strings.Title(sectionName) + "]\n")
// 			for _, line := range lines {
// 				result.WriteString(line + "\n")
// 			}
// 			result.WriteString("\n")
// 		}
// 	}

// 	for sectionName, lines := range paroleMap {
// 		found := false
// 		for _, orderedSection := range sectionOrder {
// 			if orderedSection == sectionName {
// 				found = true
// 				break
// 			}
// 		}
// 		if !found {
// 			result.WriteString("[" + strings.Title(sectionName) + "]\n")
// 			for _, line := range lines {
// 				result.WriteString(line + "\n")
// 			}
// 			result.WriteString("\n")
// 		}
// 	}

// 	return strings.TrimSpace(result.String())
// }

// func GetLyricsFromAPI(trackID string) (string, error) {
// 	endpoint := fmt.Sprintf("track.lyrics.get?track_id=%s&apikey=%s", trackID, "YOUR_API_KEY")
// 	resp, err := http.Get(endpoint)
// 	if err != nil {
// 		return "", fmt.Errorf("erreur HTTP: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	body, _ := io.ReadAll(resp.Body)
// 	var result struct {
// 		Message struct {
// 			Body struct {
// 				Lyrics struct {
// 					LyricsBody string `json:"lyrics_body"`
// 				} `json:"lyrics"`
// 			} `json:"body"`
// 		} `json:"message"`
// 	}
// 	json.Unmarshal(body, &result)
// 	if result.Message.Body.Lyrics.LyricsBody == "" {
// 		return "", errors.New("paroles non trouvées")
// 	}
// 	return result.Message.Body.Lyrics.LyricsBody, nil
// }

// func CheckLyricsFile(w http.ResponseWriter, r *http.Request) {
// 	filename := r.URL.Query().Get("filename")
// 	if filename == "" {
// 		http.Error(w, "Nom de fichier manquant", http.StatusBadRequest)
// 		return
// 	}

// 	filename = normalizeName(filename)

// 	if !strings.HasSuffix(filename, ".json") {
// 		filename += ".json"
// 	}

// 	filePath := filepath.Join(paroleDataPath, filename)

// 	response := LyricsCheckResponse{
// 		Exists: false,
// 	}

// 	if fileContent, err := os.ReadFile(filePath); err == nil {
// 		response.Exists = true
// 		response.Content = string(fileContent)
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(response)
// }

// func normalizeName(input string) string {
// 	input = strings.ToLower(input)
// 	input = removeAccents(input)
// 	input = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(input, " ")
// 	input = strings.TrimSpace(input)
// 	input = strings.ReplaceAll(input, " ", "_")
// 	return input
// }

// func removeAccents(s string) string {
// 	var sb strings.Builder
// 	for _, r := range s {
// 		if r > unicode.MaxASCII {
// 			r = unicode.To(unicode.LowerCase, r)
// 			switch r {
// 			case 'à', 'â', 'ä':
// 				r = 'a'
// 			case 'ç':
// 				r = 'c'
// 			case 'é', 'è', 'ê', 'ë':
// 				r = 'e'
// 			case 'î', 'ï':
// 				r = 'i'
// 			case 'ô', 'ö':
// 				r = 'o'
// 			case 'ù', 'û', 'ü':
// 				r = 'u'
// 			default:
// 				r = '-'
// 			}
// 		}
// 		sb.WriteRune(r)
// 	}
// 	return sb.String()
// }

// func GetLyricsFilesList(w http.ResponseWriter, r *http.Request) {
// 	files, err := os.ReadDir(paroleDataPath)
// 	if err != nil {
// 		http.Error(w, "Erreur lors de la lecture du dossier paroles", http.StatusInternalServerError)
// 		return
// 	}

// 	var lyricsFiles []string
// 	for _, file := range files {
// 		if !file.IsDir() && strings.HasSuffix(file.Name(), ".txt") {
// 			lyricsFiles = append(lyricsFiles, file.Name())
// 		}
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(lyricsFiles)
// }

// func HandleLyricsVisibility(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	sessionID := vars["id"]

// 	var request struct {
// 		Visible bool `json:"visible"`
// 	}
// 	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
// 		http.Error(w, "Requête invalide", http.StatusBadRequest)
// 		return
// 	}

// 	session, err := SetLyricsVisibility(sessionID, request.Visible)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusNotFound)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(session)
// }

// func SetLyricsVisibility(sessionID string, visible bool) (*GameSession, error) {
// 	session, ok := gameSessions[sessionID]
// 	if !ok {
// 		return nil, fmt.Errorf("session non trouvée")
// 	}

// 	session.LyricsVisible = visible
// 	return session, nil
// }

// // gesttion des paroles

// func MaskLyrics(lyrics string, points int) string {
// 	if lyrics == "" {
// 		return lyrics
// 	}

// 	sections := splitLyricsBySections(lyrics)

// 	var targetSection string
// 	switch {
// 	case points >= 40:
// 		targetSection = "Couplet 2"
// 	case points >= 10 && points <= 30:
// 		targetSection = "Refrain"
// 	default:
// 		targetSection = ""
// 	}

// 	if content, ok := sections[targetSection]; ok {
// 		lines := strings.Split(strings.TrimSpace(content), "\n")
// 		sections[targetSection] = MaskedSectionContent(targetSection, lines, points)
// 	}

// 	var rebuiltLyrics strings.Builder
// 	for section, content := range sections {
// 		rebuiltLyrics.WriteString("[" + section + "]\n")
// 		rebuiltLyrics.WriteString(content + "\n\n")
// 	}

// 	return strings.TrimSpace(rebuiltLyrics.String())
// }

// func splitLyricsBySections(lyrics string) map[string]string {
// 	lines := strings.Split(lyrics, "\n")
// 	currentSection := "Intro"
// 	sections := make(map[string]string)

// 	for _, line := range lines {
// 		line = strings.TrimSpace(line)
// 		if line == "" {
// 			continue
// 		}
// 		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
// 			currentSection = strings.Trim(line, "[]")
// 		} else {
// 			sections[currentSection] += line + "\n"
// 		}
// 	}

// 	return sections
// }

// func MaskedSectionContent(sectionName string, lines []string, points int) string {
// 	showSection := false

// 	switch points {
// 	case 50, 40:
// 		showSection = strings.ToLower(sectionName) == "couplet1"
// 	case 30, 20, 10:
// 		showSection = !strings.Contains(strings.ToLower(sectionName), "refrain")
// 	default:
// 		showSection = true
// 	}

// 	if showSection {
// 		return strings.Join(lines, "<br>")
// 	}

// 	var maskedLines []string
// 	for _, line := range lines {
// 		words := strings.Fields(line)
// 		maskedWords := make([]string, len(words))
// 		for i, word := range words {
// 			maskedWords[i] = strings.Repeat("█", len(word))
// 		}
// 		maskedLines = append(maskedLines, strings.Join(maskedWords, " "))
// 	}

// 	return strings.Join(maskedLines, "<br>")
// }

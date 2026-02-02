package game

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
)

// généralité : import/ export/ création

func loadSongMetadataFromLyricsFile(lyricsFileName string) (Titre string, Artiste string, err error) {
	if lyricsFileName == "" {
		return "", "", fmt.Errorf("nom de fichier de paroles vide")
	}

	path := filepath.Join(paroleDataPath, lyricsFileName)
	absPath, _ := filepath.Abs(path)

	if _, statErr := os.Stat(absPath); os.IsNotExist(statErr) {
		return parseLyricsFilename(lyricsFileName)
	}

	// Lire le fichier JSON
	fileContent, err := os.ReadFile(absPath)
	if err != nil {
		return parseLyricsFilename(lyricsFileName)
	}

	// Parser le JSON
	var lyricsData LyricsFileData
	if err := json.Unmarshal(fileContent, &lyricsData); err != nil {
		return parseLyricsFilename(lyricsFileName)
	}

	return lyricsData.Titre, lyricsData.Artiste, nil
}

func parseLyricsFilename(filename string) (Titre string, Artiste string, err error) {
	nameWithoutExt := filename
	if len(filename) > 5 && filename[len(filename)-5:] == ".json" {
		nameWithoutExt = filename[:len(filename)-5]
	}

	// Chercher " - " pour séparer artiste et titre
	separatorIndex := -1
	for i := 0; i < len(nameWithoutExt)-2; i++ {
		if nameWithoutExt[i:i+3] == " - " {
			separatorIndex = i
			break
		}
	}

	if separatorIndex > 0 {
		Artiste = nameWithoutExt[:separatorIndex]
		Titre = nameWithoutExt[separatorIndex+3:]
		return Titre, Artiste, nil
	}

	// Si pas de séparateur trouvé
	return nameWithoutExt, "Artiste inconnu", nil
}

// Fonction pour compléter les métadonnées d'un duel
func enrichDuelWithMetadata(duel *Duel) error {
	// Enrichir les chansons dans chaque niveau de points
	for levelKey, pointLevel := range duel.Points {
		for i := range pointLevel.Songs {
			song := &pointLevel.Songs[i]

			// Si on a un lyricsFile mais pas de title/artist
			if song.LyricsFile != nil && *song.LyricsFile != "" {
				if song.Titre == "" || song.Artiste == "" {
					title, artist, err := loadSongMetadataFromLyricsFile(*song.LyricsFile)
					if err == nil {
						if song.Titre == "" {
							song.Titre = title
						}
						if song.Artiste == "" {
							song.Artiste = artist
						}
					}
				}
			}
		}
		duel.Points[levelKey] = pointLevel
	}

	// Enrichir sameSong si nécessaire
	if duel.SameSong.LyricsFile != nil && *duel.SameSong.LyricsFile != "" {
		if duel.SameSong.Titre == "" || duel.SameSong.Artiste == "" ||
			duel.SameSong.Titre == "N/A" || duel.SameSong.Artiste == "N/A" {
			title, artist, err := loadSongMetadataFromLyricsFile(*duel.SameSong.LyricsFile)
			if err == nil {
				if duel.SameSong.Titre == "" || duel.SameSong.Titre == "N/A" {
					duel.SameSong.Titre = title
				}
				if duel.SameSong.Artiste == "" || duel.SameSong.Artiste == "N/A" {
					duel.SameSong.Artiste = artist
				}
			}
		}
	}

	return nil
}

// MODIFIER la fonction CreateDuel pour enrichir les métadonnées
func CreateDuel(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Erreur lors de la lecture du corps de la requête", http.StatusBadRequest)
		return
	}
	var singleDuel Duel
	if err := json.Unmarshal(body, &singleDuel); err == nil {

		// AJOUTER ICI : Enrichir avec les métadonnées
		if err := enrichDuelWithMetadata(&singleDuel); err != nil {
			fmt.Println("Avertissement : erreur lors de l'enrichissement des métadonnées:", err)
		}

		if err := validateDuelForClient(&singleDuel); err != nil {
			fmt.Println("3 Validation échouée:", err)
			http.Error(w, fmt.Sprintf("Données de duel invalides: %v", err), http.StatusBadRequest)
			return
		}

		updateNextDuelID()
		if (singleDuel.ID != 0 && isIDTaken(singleDuel.ID)) || singleDuel.ID == 0 {
			singleDuel.ID = nextDuelID
			nextDuelID++
		}

		duels = append(duels, singleDuel)
		updateNextDuelID()

		if err := saveDuelsToServer(); err != nil {
			http.Error(w, "Erreur lors de la sauvegarde du fichier", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(singleDuel)
		return
	}

	// Gestion du tableau de duels (code existant)
	var duelsToCreate []Duel
	if err := json.Unmarshal(body, &duelsToCreate); err != nil {
		http.Error(w, "Erreur lors du décodage JSON : un duel ou un tableau de duels est attendu.", http.StatusBadRequest)
		return
	}

	if len(duelsToCreate) != 1 {
		http.Error(w, "La création ne peut concerner qu'un seul duel à la fois.", http.StatusBadRequest)
		return
	}

	newDuel := duelsToCreate[0]

	// Enrichir avec les métadonnées
	if err := enrichDuelWithMetadata(&newDuel); err != nil {
		fmt.Println("Avertissement : erreur lors de l'enrichissement des métadonnées:", err)
	}

	if err := validateDuel(&newDuel); err != nil {
		http.Error(w, fmt.Sprintf("Données de duel invalides: %v", err), http.StatusBadRequest)
		return
	}

	updateNextDuelID()
	if (newDuel.ID != 0 && isIDTaken(newDuel.ID)) || newDuel.ID == 0 {
		newDuel.ID = nextDuelID
		nextDuelID++
	}

	duels = append(duels, newDuel)
	updateNextDuelID()

	if err := saveDuelsToServer(); err != nil {
		http.Error(w, "Erreur lors de la sauvegarde du fichier", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode([]Duel{newDuel})

	fmt.Println("6 Duel créé avec succès avec l'ID :", newDuel.ID)
}

// MODIFIER aussi LoadDuelFromJSON pour enrichir les duels importés
func LoadDuelFromJSON(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Erreur lors de l'analyse du formulaire multipart", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("duelFile")
	if err != nil {
		http.Error(w, "Erreur lors de la lecture du fichier 'duelFile'", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileContent, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Erreur lors de la lecture du contenu du fichier", http.StatusInternalServerError)
		return
	}

	var loadedDuel Duel
	if err := json.Unmarshal(fileContent, &loadedDuel); err != nil {
		http.Error(w, "Erreur lors du décodage JSON du fichier", http.StatusBadRequest)
		return
	}

	// Enrichir avec les métadonnées
	if err := enrichDuelWithMetadata(&loadedDuel); err != nil {
		fmt.Println("Avertissement : erreur lors de l'enrichissement des métadonnées:", err)
	}

	if err := validateDuel(&loadedDuel); err != nil {
		http.Error(w, fmt.Sprintf("Fichier JSON de duel invalide: %v", err), http.StatusBadRequest)
		return
	}

	loadedDuel.ID = nextDuelID
	nextDuelID++

	duels = append(duels, loadedDuel)

	if err := saveDuelsToServer(); err != nil {
		http.Error(w, "Erreur lors de la sauvegarde", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(loadedDuel)
}

func createDirectories() {
	dirs := []string{duelSaveDataPath, prepDuelDataPath, paroleDataPath}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("Erreur lors de la création du dossier %s: %v\n", dir, err)
		}
	}
}

func updateNextDuelID() {
	maxID := 0
	for _, duel := range duels {
		if duel.ID > maxID {
			maxID = duel.ID
		}
	}
	nextDuelID = maxID + 1
}

func isIDTaken(id int) bool {
	for _, duel := range duels {
		if duel.ID == id {
			return true
		}
	}
	return false
}

func GetDuels(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(duels); err != nil {
		http.Error(w, "Erreur lors de l'encodage JSON des duels", http.StatusInternalServerError)
	}
}

func GetDuelByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	duelID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID du duel invalide", http.StatusBadRequest)
		return
	}

	for _, duel := range duels {
		if duel.ID == duelID {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(duel)
			return
		}
	}

	http.Error(w, "Duel non trouvé", http.StatusNotFound)
}

func prepareDuelForExport(duel *Duel) error {
	updateNextDuelID()

	needsSave := false

	if duel.ID == 0 {
		duel.ID = nextDuelID
		nextDuelID++
		needsSave = true
	}

	if needsSave {
		return saveDuelsToServer()
	}

	return nil
}

func saveDuelsToServer() error {
	filePath := filepath.Join(duelSaveDataPath, "duels.json")

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(duels)
}

func loadDuelsFromServer() error {
	filePath := filepath.Join(duelSaveDataPath, "duels.json")

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	var loadedDuels []Duel
	if err := json.Unmarshal(fileContent, &loadedDuels); err != nil {
		return err
	}

	duels = loadedDuels

	for _, duel := range duels {
		if duel.ID >= nextDuelID {
			nextDuelID = duel.ID + 1
		}
	}

	return nil
}

// MaJ et validation

func DownloadDuel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	duelID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	var selected *Duel
	for i := range duels {
		if duels[i].ID == duelID {
			selected = &duels[i]
			break
		}
	}

	if selected == nil {
		http.Error(w, "Duel introuvable", http.StatusNotFound)
		return
	}

	if err := prepareDuelForExport(selected); err != nil {
		http.Error(w, "Erreur lors de la préparation du duel", http.StatusInternalServerError)
		return
	}

	fileName := fmt.Sprintf("%d_id.json", selected.ID)

	exportDuel := *selected

	jsonData, err := json.MarshalIndent(exportDuel, "", "  ")
	if err != nil {
		http.Error(w, "Erreur d'encodage JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	w.Header().Set("Content-Length", strconv.Itoa(len(jsonData)))

	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func UpdateDuel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	duelID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID du duel invalide", http.StatusBadRequest)
		return
	}

	duelIndex := -1
	for i, duel := range duels {
		if duel.ID == duelID {
			duelIndex = i
			break
		}
	}

	if duelIndex == -1 {
		http.Error(w, "Duel non trouvé pour la mise à jour", http.StatusNotFound)
		return
	}

	var updatedDuel Duel
	if err := json.NewDecoder(r.Body).Decode(&updatedDuel); err != nil {
		http.Error(w, "Erreur lors du décodage JSON pour la mise à jour", http.StatusBadRequest)
		return
	}

	if err := validateDuel(&updatedDuel); err != nil {
		http.Error(w, fmt.Sprintf("Données de mise à jour invalides: %v", err), http.StatusBadRequest)
		return
	}

	updatedDuel.ID = duelID

	duels[duelIndex] = updatedDuel

	if err := saveDuelsToServer(); err != nil {
		http.Error(w, "Erreur lors de la sauvegarde", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedDuel)
}

func DeleteDuel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	duelID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID du duel invalide pour la suppression", http.StatusBadRequest)
		return
	}

	for i, duel := range duels {
		if duel.ID == duelID {
			duels = append(duels[:i], duels[i+1:]...)

			if err := saveDuelsToServer(); err != nil {
				http.Error(w, "Erreur lors de la sauvegarde après suppression", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "Duel non trouvé pour la suppression", http.StatusNotFound)
}

func validateDuelForClient(duel *Duel) error {
	if duel.Name == "" {
		return fmt.Errorf("le nom du duel est requis")
	}

	requiredLevels := []string{"50", "40", "30", "20", "10"}
	if len(duel.Points) != len(requiredLevels) {
		return fmt.Errorf("le nombre de niveaux de points est incorrect. Requis: %v", requiredLevels)
	}

	for _, level := range requiredLevels {
		pointLevel, exists := duel.Points[level]
		if !exists {
			return fmt.Errorf("le niveau %s points est manquant", level)
		}

		if pointLevel.Theme == "" {
			return fmt.Errorf("le thème pour %s points est requis", level)
		}

		if len(pointLevel.Songs) < 1 || len(pointLevel.Songs) > 2 {
			return fmt.Errorf("1 ou 2 chansons sont requises pour le niveau %s points", level)
		}

		for _, song := range pointLevel.Songs {
			if song.LyricsFile != nil && *song.LyricsFile == "" {
				continue
			}
		}
	}

	return nil
}

func validateDuel(duel *Duel) error {
	if duel.Name == "" {
		return fmt.Errorf("le nom du duel est requis")
	}

	requiredLevels := []string{"50", "40", "30", "20", "10"}
	if len(duel.Points) != len(requiredLevels) {
		return fmt.Errorf("le nombre de niveaux de points est incorrect. Requis: %v", requiredLevels)
	}

	for _, level := range requiredLevels {
		pointLevel, exists := duel.Points[level]
		if !exists {
			return fmt.Errorf("le niveau %s points est manquant", level)
		}

		if pointLevel.Theme == "" {
			return fmt.Errorf("le thème pour %s points est requis", level)
		}

		if len(pointLevel.Songs) != 2 {
			return fmt.Errorf("exactement 2 chansons sont requises pour le niveau %s points", level)
		}

		for i, song := range pointLevel.Songs {
			if song.Titre == "" {
				return fmt.Errorf("le titre de la chanson %d pour %s points est requis", i+1, level)
			}
			if song.Artiste == "" {
				return fmt.Errorf("l'artiste de la chanson %d pour %s points est requis", i+1, level)
			}
		}
	}

	if duel.SameSong.Titre == "" {
		return fmt.Errorf("le titre de 'La Même Chanson' est requis")
	}
	if duel.SameSong.Artiste == "" {
		return fmt.Errorf("l'artiste de 'La Même Chanson' est requis")
	}

	return nil
}

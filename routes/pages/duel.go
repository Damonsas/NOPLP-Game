package game

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type Song struct {
	Title      string  `json:"title"`
	Artist     string  `json:"artist"`
	AudioURL   *string `json:"audioUrl,omitempty"`
	LyricsFile *string `json:"lyricsFile,omitempty"`
}

type PointLevel struct {
	Theme string `json:"theme"`
	Songs []Song `json:"songs"`
}

type Duel struct {
	ID        int                   `json:"id,omitempty"`
	Name      string                `json:"name"`
	Points    map[string]PointLevel `json:"points"`
	SameSong  Song                  `json:"sameSong"`
	CreatedAt time.Time             `json:"createdAt"`
	UpdatedAt *time.Time            `json:"updatedAt,omitempty"`
}

type GameSession struct {
	ID            string         `json:"id"`
	DuelID        int            `json:"duelId"`
	CurrentLevel  string         `json:"currentLevel"`
	SelectedSongs map[string]int `json:"selectedSongs"`
	Team1Score    int            `json:"team1Score"`
	Team2Score    int            `json:"team2Score"`
	StartedAt     time.Time      `json:"startedAt"`
	Status        string         `json:"status"` // "playing", "paused", "finished"

	CurrentSong   *Song  `json:"currentSong,omitempty"`   // La chanson actuellement jouée
	LyricsContent string `json:"lyricsContent,omitempty"` // Le texte des paroles
	LyricsVisible bool   `json:"lyricsVisible"`
}

type LyricsCheckResponse struct {
	Exists  bool   `json:"exists"`
	Content string `json:"content,omitempty"`
}

var duels []Duel
var gameSessions map[string]*GameSession
var nextDuelID int = 1

const (
	duelSaveDataPath = "data/serverdata/duelsavedata"
	prepDuelDataPath = "data/serverdata/prepdueldata"
	paroleDataPath   = "data/serverdata/paroledata"
)

func init() {
	duels = make([]Duel, 0)
	gameSessions = make(map[string]*GameSession)

	createDirectories()

	loadDuelsFromServer()
}

func createDirectories() {
	dirs := []string{duelSaveDataPath, prepDuelDataPath, paroleDataPath}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("Erreur lors de la création du dossier %s: %v\n", dir, err)
		}
	}
}

func DuelMaestroChallenger(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "routes/pages/duel.html")
}

func DuelGamePage(w http.ResponseWriter, r *http.Request) {
	duelID := r.URL.Query().Get("duelId")
	if duelID == "" {
		http.Error(w, "ID de duel manquant dans les paramètres de la requête", http.StatusBadRequest)
		return
	}
}

func GetDuels(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(duels); err != nil {
		http.Error(w, "Erreur lors de l'encodage JSON des duels", http.StatusInternalServerError)
	}
}

func CreateDuel(w http.ResponseWriter, r *http.Request) {
	var newDuel Duel
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Erreur lors de la lecture du corps de la requête", http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &newDuel); err != nil {
		http.Error(w, "Erreur lors du décodage JSON", http.StatusBadRequest)
		return
	}

	if err := validateDuel(&newDuel); err != nil {
		http.Error(w, fmt.Sprintf("Données de duel invalides: %v", err), http.StatusBadRequest)
		return
	}

	newDuel.ID = nextDuelID
	nextDuelID++
	newDuel.CreatedAt = time.Now()

	duels = append(duels, newDuel)

	if err := saveDuelsToServer(); err != nil {
		http.Error(w, "Erreur lors de la sauvegarde du fichier", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newDuel)
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
	updatedDuel.CreatedAt = duels[duelIndex].CreatedAt
	now := time.Now()
	updatedDuel.UpdatedAt = &now

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

	if err := validateDuel(&loadedDuel); err != nil {
		http.Error(w, fmt.Sprintf("Fichier JSON de duel invalide: %v", err), http.StatusBadRequest)
		return
	}

	loadedDuel.ID = nextDuelID
	nextDuelID++
	loadedDuel.CreatedAt = time.Now()

	duels = append(duels, loadedDuel)

	if err := saveDuelsToServer(); err != nil {
		http.Error(w, "Erreur lors de la sauvegarde", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(loadedDuel)
}

// prepareDuelForExport prépare un duel pour l'export en s'assurant qu'il a un ID et des métadonnées valides
func prepareDuelForExport(duel *Duel) error {
	updateNextDuelID()

	needsSave := false

	// S'assurer que le duel a un ID valide
	if duel.ID == 0 {
		duel.ID = nextDuelID
		nextDuelID++
		needsSave = true
	}

	// S'assurer que createdAt est défini
	if duel.CreatedAt.IsZero() {
		duel.CreatedAt = time.Now()
		needsSave = true
	}

	// Sauvegarder les modifications si nécessaire
	if needsSave {
		return saveDuelsToServer()
	}

	return nil
}

// ImportDuelFromServer charge un duel depuis le système de fichiers du serveur
func ImportDuelFromServer(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("filename")
	if filename == "" {
		http.Error(w, "Nom de fichier manquant", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join(duelSaveDataPath, filename)
	if !strings.HasSuffix(filePath, ".json") {
		filePath += ".json"
	}

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		http.Error(w, "Fichier non trouvé", http.StatusNotFound)
		return
	}

	// S'assurer que nextDuelID est à jour avant l'import
	updateNextDuelID()

	// Essayer d'abord de charger comme un tableau (format duels.json complet)
	var loadedDuels []Duel
	if err := json.Unmarshal(fileContent, &loadedDuels); err == nil && len(loadedDuels) > 0 {
		// Format tableau complet - importer tous les duels
		var importedDuels []Duel

		for _, loadedDuel := range loadedDuels {
			if err := validateDuel(&loadedDuel); err != nil {
				http.Error(w, fmt.Sprintf("Fichier JSON de duel invalide: %v", err), http.StatusBadRequest)
				return
			}

			// Générer un nouvel ID unique même si le duel en avait déjà un
			loadedDuel.ID = nextDuelID
			nextDuelID++
			loadedDuel.CreatedAt = time.Now()
			// Réinitialiser UpdatedAt pour un nouvel import
			loadedDuel.UpdatedAt = nil

			duels = append(duels, loadedDuel)
			importedDuels = append(importedDuels, loadedDuel)
		}

		// Sauvegarder immédiatement pour persister les nouveaux IDs
		if err := saveDuelsToServer(); err != nil {
			http.Error(w, "Erreur lors de la sauvegarde après import", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		response := map[string]interface{}{
			"message": fmt.Sprintf("%d duel(s) importé(s) avec succès", len(importedDuels)),
			"duels":   importedDuels,
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Sinon, essayer de charger comme un duel unique
	var loadedDuel Duel
	if err := json.Unmarshal(fileContent, &loadedDuel); err != nil {
		http.Error(w, "Erreur lors du décodage JSON du fichier", http.StatusBadRequest)
		return
	}

	if err := validateDuel(&loadedDuel); err != nil {
		http.Error(w, fmt.Sprintf("Fichier JSON de duel invalide: %v", err), http.StatusBadRequest)
		return
	}

	// Générer un nouvel ID unique même si le duel en avait déjà un
	loadedDuel.ID = nextDuelID
	nextDuelID++
	loadedDuel.CreatedAt = time.Now()
	// Réinitialiser UpdatedAt pour un nouvel import
	loadedDuel.UpdatedAt = nil

	duels = append(duels, loadedDuel)

	// Sauvegarder immédiatement pour persister le nouvel ID
	if err := saveDuelsToServer(); err != nil {
		http.Error(w, "Erreur lors de la sauvegarde après import", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(loadedDuel)
}

func ExportDuelToServer(w http.ResponseWriter, r *http.Request) {
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

	// Préparer le duel pour l'export
	if err := prepareDuelForExport(selected); err != nil {
		http.Error(w, "Erreur lors de la préparation du duel", http.StatusInternalServerError)
		return
	}

	// Nettoyer le nom pour le nom de fichier
	cleanName := sanitizeFileName(selected.Name)
	fileName := fmt.Sprintf("%s_id_%d.json", cleanName, selected.ID)
	filePath := filepath.Join(duelSaveDataPath, fileName)

	// Créer un duel individuel (format compatible avec votre structure)
	exportDuel := *selected

	file, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Erreur création du fichier", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(exportDuel); err != nil {
		http.Error(w, "Erreur d'encodage JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"message":  "Duel sauvegardé avec succès sur le serveur",
		"filename": fileName,
		"path":     filePath,
		"duelId":   selected.ID,
		"duelName": selected.Name,
	}
	json.NewEncoder(w).Encode(response)
}

// DownloadDuel permet au client de télécharger un duel au format JSON
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

	// Préparer le duel pour l'export
	if err := prepareDuelForExport(selected); err != nil {
		http.Error(w, "Erreur lors de la préparation du duel", http.StatusInternalServerError)
		return
	}

	// Nettoyer le nom pour le nom de fichier
	cleanName := sanitizeFileName(selected.Name)
	fileName := fmt.Sprintf("%s_id_%d.json", cleanName, selected.ID)

	// Créer un duel individuel (format compatible avec votre structure)
	exportDuel := *selected

	// Encoder en JSON
	jsonData, err := json.MarshalIndent(exportDuel, "", "  ")
	if err != nil {
		http.Error(w, "Erreur d'encodage JSON", http.StatusInternalServerError)
		return
	}

	// Définir les headers pour le téléchargement
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	w.Header().Set("Content-Length", strconv.Itoa(len(jsonData)))

	// Envoyer le fichier
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// sanitizeFileName nettoie un nom pour qu'il soit valide comme nom de fichier
func sanitizeFileName(name string) string {
	// Remplacer les caractères problématiques par des underscores
	replacements := []string{" ", "/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	cleanName := name
	for _, char := range replacements {
		cleanName = strings.ReplaceAll(cleanName, char, "_")
	}

	// Supprimer les underscores multiples consécutifs
	for strings.Contains(cleanName, "__") {
		cleanName = strings.ReplaceAll(cleanName, "__", "_")
	}

	// Supprimer les underscores en début/fin
	cleanName = strings.Trim(cleanName, "_")

	// S'assurer qu'il n'est pas vide
	if cleanName == "" {
		cleanName = "duel"
	}

	return cleanName
}

// updateNextDuelID met à jour nextDuelID pour éviter les conflits d'IDs
func updateNextDuelID() {
	maxID := 0
	for _, duel := range duels {
		if duel.ID > maxID {
			maxID = duel.ID
		}
	}
	nextDuelID = maxID + 1
}

// validateDuel vérifie que la structure d'un duel est complète et correcte.
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
			if song.Title == "" {
				return fmt.Errorf("le titre de la chanson %d pour %s points est requis", i+1, level)
			}
			if song.Artist == "" {
				return fmt.Errorf("l'artiste de la chanson %d pour %s points est requis", i+1, level)
			}
		}
	}

	if duel.SameSong.Title == "" {
		return fmt.Errorf("le titre de 'La Même Chanson' est requis")
	}
	if duel.SameSong.Artist == "" {
		return fmt.Errorf("l'artiste de 'La Même Chanson' est requis")
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

// SaveTemporaryDuel sauvegarde un duel temporaire
func SaveTemporaryDuel(w http.ResponseWriter, r *http.Request) {
	var tempDuel Duel
	if err := json.NewDecoder(r.Body).Decode(&tempDuel); err != nil {
		http.Error(w, "Erreur lors du décodage JSON", http.StatusBadRequest)
		return
	}

	filename := "temp_duel.json"
	filePath := filepath.Join(prepDuelDataPath, filename)

	file, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Erreur lors de la création du fichier temporaire", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(tempDuel); err != nil {
		http.Error(w, "Erreur lors de l'encodage JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"message": "Duel temporaire sauvegardé avec succès",
	}
	json.NewEncoder(w).Encode(response)
}

// LoadTemporaryDuel charge un duel temporaire
func LoadTemporaryDuel(w http.ResponseWriter, r *http.Request) {
	filename := "temp_duel.json"
	filePath := filepath.Join(prepDuelDataPath, filename)

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		http.Error(w, "Aucun duel temporaire trouvé", http.StatusNotFound)
		return
	}

	var tempDuel Duel
	if err := json.Unmarshal(fileContent, &tempDuel); err != nil {
		http.Error(w, "Erreur lors du décodage JSON du fichier temporaire", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tempDuel)
}

// CheckLyricsFile vérifie si un fichier de paroles existe
func CheckLyricsFile(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("filename")
	if filename == "" {
		http.Error(w, "Nom de fichier manquant", http.StatusBadRequest)
		return
	}

	// Ajouter l'extension .txt si elle n'est pas présente
	if !strings.HasSuffix(filename, ".txt") {
		filename += ".txt"
	}

	filePath := filepath.Join(paroleDataPath, filename)

	response := LyricsCheckResponse{
		Exists: false,
	}

	if fileContent, err := os.ReadFile(filePath); err == nil {
		response.Exists = true
		response.Content = string(fileContent)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetLyricsFilesList(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir(paroleDataPath)
	if err != nil {
		http.Error(w, "Erreur lors de la lecture du dossier paroles", http.StatusInternalServerError)
		return
	}

	var lyricsFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".txt") {
			lyricsFiles = append(lyricsFiles, file.Name())
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lyricsFiles)
}

func GetServerDuelsList(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir(duelSaveDataPath)
	if err != nil {
		http.Error(w, "Erreur lors de la lecture du dossier duels", http.StatusInternalServerError)
		return
	}

	var duelFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			duelFiles = append(duelFiles, file.Name())
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(duelFiles)
}

func StartGameSession(w http.ResponseWriter, r *http.Request) {
	var request struct {
		DuelID int `json:"duelId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Erreur lors du décodage JSON de la requête", http.StatusBadRequest)
		return
	}

	var selectedDuel *Duel
	for i := range duels {
		if duels[i].ID == request.DuelID {
			selectedDuel = &duels[i]
			break
		}
	}

	if selectedDuel == nil {
		http.Error(w, "Duel non trouvé pour démarrer une session", http.StatusNotFound)
		return
	}

	sessionID := fmt.Sprintf("session_%d_%d", request.DuelID, time.Now().UnixNano())
	session := &GameSession{
		ID:            sessionID,
		DuelID:        request.DuelID,
		CurrentLevel:  "50",
		SelectedSongs: make(map[string]int),
		Team1Score:    0,
		Team2Score:    0,
		StartedAt:     time.Now(),
		Status:        "playing",
	}

	gameSessions[sessionID] = session

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(session)
}

func GetGameSession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["id"]

	session, exists := gameSessions[sessionID]
	if !exists {
		http.Error(w, "Session de jeu non trouvée", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}

func SelectSongForLevel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["id"]

	session, exists := gameSessions[sessionID]
	if !exists {
		http.Error(w, "Session non trouvée", http.StatusNotFound)
		return
	}

	var request struct {
		Level     string `json:"level"`
		SongIndex int    `json:"songIndex"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Erreur lors du décodage JSON", http.StatusBadRequest)
		return
	}

	validLevels := []string{"50", "40", "30", "20", "10"}
	isValidLevel := false
	for _, level := range validLevels {
		if level == request.Level {
			isValidLevel = true
			break
		}
	}
	if !isValidLevel {
		http.Error(w, "Niveau invalide", http.StatusBadRequest)
		return
	}

	if request.SongIndex < 0 || request.SongIndex > 1 {
		http.Error(w, "Index de chanson invalide (doit être 0 ou 1)", http.StatusBadRequest)
		return
	}

	session.SelectedSongs[request.Level] = request.SongIndex

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}

func UpdateGameScore(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["id"]

	session, exists := gameSessions[sessionID]
	if !exists {
		http.Error(w, "Session non trouvée", http.StatusNotFound)
		return
	}

	var request struct {
		Team   int `json:"team"` // 1 ou 2
		Points int `json:"points"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Erreur lors du décodage JSON", http.StatusBadRequest)
		return
	}

	if request.Team != 1 && request.Team != 2 {
		http.Error(w, "L'équipe doit être 1 ou 2", http.StatusBadRequest)
		return
	}

	if request.Team == 1 {
		session.Team1Score += request.Points
	} else {
		session.Team2Score += request.Points
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}

func FinishGameSession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["id"]

	session, exists := gameSessions[sessionID]
	if !exists {
		http.Error(w, "Session non trouvée", http.StatusNotFound)
		return
	}

	session.Status = "finished"

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}

func loadDuelsFromServer() error {
	filePath := filepath.Join(duelSaveDataPath, "duels.json")

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		// Si le fichier n'existe pas, ce n'est pas une erreur
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

func HandleStartSong(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["id"]

	var request struct {
		Level     string `json:"level"`
		SongIndex int    `json:"songIndex"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Requête invalide", http.StatusBadRequest)
		return
	}

	// Appel de la logique depuis duellaunch.go
	session, err := StartSong(sessionID, request.Level, request.SongIndex)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}

// HandleLyricsVisibility est appelé pour montrer/cacher les paroles.
func HandleLyricsVisibility(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["id"]

	var request struct {
		Visible bool `json:"visible"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Requête invalide", http.StatusBadRequest)
		return
	}

	// Appel de la logique depuis duellaunch.go
	session, err := SetLyricsVisibility(sessionID, request.Visible)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}

func SetupDuelRoutes(r *mux.Router) {
	r.HandleFunc("/duel", DuelMaestroChallenger).Methods("GET")
	r.HandleFunc("/duel-game", DisplayDuel).Methods("GET", "POST")

	r.HandleFunc("/duel-display", DisplayDuel).Methods("GET", "POST")

	r.HandleFunc("/api/duels", GetDuels).Methods("GET")
	r.HandleFunc("/api/duels", CreateDuel).Methods("POST")
	r.HandleFunc("/api/duels/{id:[0-9]+}", GetDuelByID).Methods("GET")
	r.HandleFunc("/api/duels/{id:[0-9]+}", UpdateDuel).Methods("PUT")
	r.HandleFunc("/api/duels/{id:[0-9]+}", DeleteDuel).Methods("DELETE")
	r.HandleFunc("/api/upload-duel", LoadDuelFromJSON).Methods("POST")

	r.HandleFunc("/api/import-duel-server", ImportDuelFromServer).Methods("GET")
	r.HandleFunc("/api/export-duel-server/{id:[0-9]+}", ExportDuelToServer).Methods("POST")
	r.HandleFunc("/api/download-duel/{id:[0-9]+}", DownloadDuel).Methods("GET")
	r.HandleFunc("/api/server-duels-list", GetServerDuelsList).Methods("GET")

	r.HandleFunc("/api/temp-duel", SaveTemporaryDuel).Methods("POST")
	r.HandleFunc("/api/temp-duel", LoadTemporaryDuel).Methods("GET")

	r.HandleFunc("/api/check-lyrics", CheckLyricsFile).Methods("GET")
	r.HandleFunc("/api/lyrics-list", GetLyricsFilesList).Methods("GET")

	r.HandleFunc("/api/game-sessions", StartGameSession).Methods("POST")
	r.HandleFunc("/api/game-sessions/{id}", GetGameSession).Methods("GET")
	r.HandleFunc("/api/game-sessions/{id}/select-song", SelectSongForLevel).Methods("POST")
	r.HandleFunc("/api/game-sessions/{id}/update-score", UpdateGameScore).Methods("POST")
	r.HandleFunc("/api/game-sessions/{id}/finish", FinishGameSession).Methods("POST")

	r.HandleFunc("/api/game-sessions/{id}/start-song", HandleStartSong).Methods("POST")
	r.HandleFunc("/api/game-sessions/{id}/lyrics-visibility", HandleLyricsVisibility).Methods("POST")
}

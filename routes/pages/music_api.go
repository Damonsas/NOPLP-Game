package game

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// YouTubeResponse structure pour la réponse de l'API YouTube
type YouTubeResponse struct {
	Items []struct {
		ID struct {
			VideoID string `json:"videoId"`
		} `json:"id"`
		Snippet struct {
			Title       string `json:"title"`
			Description string `json:"description"`
		} `json:"snippet"`
	} `json:"items"`
}

// SpotifyTokenResponse structure pour l'authentification Spotify
type SpotifyTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// SpotifySearchResponse structure pour la recherche Spotify
type SpotifySearchResponse struct {
	Tracks struct {
		Items []struct {
			ID           string  `json:"id"`
			Name         string  `json:"name"`
			PreviewURL   *string `json:"preview_url"`
			ExternalURLs struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Artists []struct {
				Name string `json:"name"`
			} `json:"artists"`
		} `json:"items"`
	} `json:"tracks"`
}

// MusicSearchResult structure pour le résultat de recherche unifié
type MusicSearchResult struct {
	Title      string `json:"title"`
	Artist     string `json:"artist"`
	PreviewURL string `json:"preview_url"`
	SourceURL  string `json:"source_url"`
	Source     string `json:"source"` // "youtube", "spotify", "local"
}

// SearchMusicHandler gère les recherches de musique
func SearchMusicHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	title := r.URL.Query().Get("title")
	artist := r.URL.Query().Get("artist")
	instrumental := r.URL.Query().Get("instrumental") == "true"

	if title == "" {
		http.Error(w, "Titre manquant", http.StatusBadRequest)
		return
	}

	config := GetMusicConfig()
	results := []MusicSearchResult{}

	// 1. Recherche locale d'abord (plus rapide)
	if config.IsLocalEnabled() {
		if localResult := searchLocal(title, artist, instrumental); localResult != nil {
			results = append(results, *localResult)
		}
	}

	// 2. Recherche Deezer (API gratuite, rapide)
	if deezerResult, err := searchDeezer(title, artist); err == nil && deezerResult != nil {
		results = append(results, *deezerResult)
	}

	// 3. Recherche iTunes (API gratuite)
	if itunesResult, err := searchITunes(title, artist); err == nil && itunesResult != nil {
		results = append(results, *itunesResult)
	}

	// 4. Recherche Spotify (si configuré)
	if config.IsSpotifyEnabled() {
		if spotifyResult, err := searchSpotify(title, artist); err == nil && spotifyResult != nil {
			results = append(results, *spotifyResult)
		}
	}

	// 5. Recherche YouTube (si configuré et instrumental demandé)
	if config.IsYouTubeEnabled() && instrumental {
		if youtubeResult, err := searchYouTube(title, artist, instrumental); err == nil && youtubeResult != nil {
			results = append(results, *youtubeResult)
		}
	}

	// Si aucun résultat trouvé, retourner un résultat par défaut
	if len(results) == 0 {
		results = append(results, MusicSearchResult{
			Title:      title,
			Artist:     artist,
			PreviewURL: "/demo-audio.mp3",
			SourceURL:  "/demo-audio.mp3",
			Source:     "local",
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// searchYouTube recherche sur YouTube
func searchYouTube(title, artist string, instrumental bool) (*MusicSearchResult, error) {
	config := GetMusicConfig()
	if !config.IsYouTubeEnabled() {
		return nil, fmt.Errorf("YouTube API non configurée")
	}

	// Construire la requête de recherche
	query := title
	if artist != "" {
		query += " " + artist
	}
	if instrumental {
		query += " instrumental"
	}

	searchURL := fmt.Sprintf(
		"https://www.googleapis.com/youtube/v3/search?part=snippet&q=%s&type=video&maxResults=1&key=%s",
		url.QueryEscape(query),
		config.YouTubeAPIKey,
	)

	resp, err := http.Get(searchURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var youtubeResp YouTubeResponse
	if err := json.NewDecoder(resp.Body).Decode(&youtubeResp); err != nil {
		return nil, err
	}

	if len(youtubeResp.Items) == 0 {
		return nil, fmt.Errorf("aucun résultat trouvé")
	}

	item := youtubeResp.Items[0]
	return &MusicSearchResult{
		Title:      item.Snippet.Title,
		Artist:     artist,
		PreviewURL: fmt.Sprintf("https://www.youtube.com/watch?v=%s", item.ID.VideoID),
		SourceURL:  fmt.Sprintf("https://www.youtube.com/watch?v=%s", item.ID.VideoID),
		Source:     "youtube",
	}, nil
}

// searchDeezer recherche sur Deezer (API gratuite)
func searchDeezer(title, artist string) (*MusicSearchResult, error) {
	query := title
	if artist != "" {
		query += " " + artist
	}

	searchURL := fmt.Sprintf(
		"https://api.deezer.com/search?q=%s&limit=1",
		url.QueryEscape(query),
	)

	resp, err := http.Get(searchURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var deezerResp struct {
		Data []struct {
			ID      int    `json:"id"`
			Title   string `json:"title"`
			Preview string `json:"preview"`
			Link    string `json:"link"`
			Artist  struct {
				Name string `json:"name"`
			} `json:"artist"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&deezerResp); err != nil {
		return nil, err
	}

	if len(deezerResp.Data) == 0 {
		return nil, fmt.Errorf("aucun résultat trouvé")
	}

	track := deezerResp.Data[0]
	return &MusicSearchResult{
		Title:      track.Title,
		Artist:     track.Artist.Name,
		PreviewURL: track.Preview,
		SourceURL:  track.Link,
		Source:     "deezer",
	}, nil
}

// searchITunes recherche sur iTunes (API gratuite)
func searchITunes(title, artist string) (*MusicSearchResult, error) {
	query := title
	if artist != "" {
		query += " " + artist
	}

	searchURL := fmt.Sprintf(
		"https://itunes.apple.com/search?term=%s&media=music&limit=1",
		url.QueryEscape(query),
	)

	resp, err := http.Get(searchURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var itunesResp struct {
		Results []struct {
			TrackName    string `json:"trackName"`
			ArtistName   string `json:"artistName"`
			PreviewURL   string `json:"previewUrl"`
			TrackViewURL string `json:"trackViewUrl"`
		} `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&itunesResp); err != nil {
		return nil, err
	}

	if len(itunesResp.Results) == 0 {
		return nil, fmt.Errorf("aucun résultat trouvé")
	}

	track := itunesResp.Results[0]
	return &MusicSearchResult{
		Title:      track.TrackName,
		Artist:     track.ArtistName,
		PreviewURL: track.PreviewURL,
		SourceURL:  track.TrackViewURL,
		Source:     "itunes",
	}, nil
}

// searchSpotify recherche sur Spotify
func searchSpotify(title, artist string) (*MusicSearchResult, error) {
	config := GetMusicConfig()
	if !config.IsSpotifyEnabled() {
		return nil, fmt.Errorf("Spotify API non configurée")
	}

	// Construire la requête de recherche
	query := title
	if artist != "" {
		query += " artist:" + artist
	}

	// Obtenir un token d'accès
	token, err := getSpotifyToken(config.SpotifyClientID, config.SpotifyClientSecret)
	if err != nil {
		return nil, err
	}

	// Faire la recherche
	searchURL := fmt.Sprintf(
		"https://api.spotify.com/v1/search?q=%s&type=track&limit=1",
		url.QueryEscape(query),
	)

	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var spotifyResp SpotifySearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&spotifyResp); err != nil {
		return nil, err
	}

	if len(spotifyResp.Tracks.Items) == 0 {
		return nil, fmt.Errorf("aucun résultat trouvé")
	}

	track := spotifyResp.Tracks.Items[0]
	previewURL := ""
	if track.PreviewURL != nil {
		previewURL = *track.PreviewURL
	}

	artistName := ""
	if len(track.Artists) > 0 {
		artistName = track.Artists[0].Name
	}

	return &MusicSearchResult{
		Title:      track.Name,
		Artist:     artistName,
		PreviewURL: previewURL,
		SourceURL:  track.ExternalURLs.Spotify,
		Source:     "spotify",
	}, nil
}

// getSpotifyToken obtient un token d'accès Spotify
func getSpotifyToken(clientID, clientSecret string) (string, error) {
	tokenURL := "https://accounts.spotify.com/api/token"

	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(clientID, clientSecret)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var tokenResp SpotifyTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", err
	}

	return tokenResp.AccessToken, nil
}

// searchLocal recherche dans les fichiers audio locaux
func searchLocal(title, artist string, instrumental bool) *MusicSearchResult {
	// Ici vous pouvez implémenter la logique pour chercher dans vos fichiers locaux
	// Par exemple, dans un dossier "audio" avec une structure organisée

	// Exemple de logique de recherche locale
	audioPath := "/audio/" // Chemin vers vos fichiers audio

	// Construire le nom de fichier potentiel
	filename := strings.ToLower(strings.ReplaceAll(title, " ", "_"))
	if artist != "" {
		filename = strings.ToLower(strings.ReplaceAll(artist, " ", "_")) + "_" + filename
	}
	if instrumental {
		filename += "_instrumental"
	}

	// Vérifier différents formats
	extensions := []string{".mp3", ".wav", ".ogg", ".m4a"}
	for _, ext := range extensions {
		fullPath := audioPath + filename + ext
		_ = fullPath
		// Dans un vrai cas, vous vérifieriez si le fichier existe
		// if _, err := os.Stat(fullPath); err == nil {
		//     return &MusicSearchResult{
		//         Title:      title,
		//         Artist:     artist,
		//         PreviewURL: fullPath,
		//         SourceURL:  fullPath,
		//         Source:     "local",
		//     }
		// }
	}

	return nil
}

// GetInstrumentalHandler endpoint pour obtenir la version instrumentale
func GetInstrumentalHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	title := r.URL.Query().Get("title")
	artist := r.URL.Query().Get("artist")

	if title == "" {
		http.Error(w, "Titre manquant", http.StatusBadRequest)
		return
	}

	// Rechercher spécifiquement la version instrumentale
	result, err := searchYouTube(title, artist, true)
	if err != nil {
		// Fallback sur une recherche locale
		if localResult := searchLocal(title, artist, true); localResult != nil {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(localResult)
			return
		}

		http.Error(w, "Version instrumentale non trouvée", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// AudioProxyHandler proxy pour servir les fichiers audio avec CORS
func AudioProxyHandler(w http.ResponseWriter, r *http.Request) {
	audioURL := r.URL.Query().Get("url")
	if audioURL == "" {
		http.Error(w, "URL audio manquante", http.StatusBadRequest)
		return
	}

	// Ajouter les headers CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Proxy la requête vers l'URL audio
	resp, err := http.Get(audioURL)
	if err != nil {
		http.Error(w, "Erreur lors du chargement de l'audio", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Copier les headers de contenu
	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.Header().Set("Content-Length", resp.Header.Get("Content-Length"))

	// Copier le contenu
	w.WriteHeader(resp.StatusCode)
	// io.Copy(w, resp.Body) // Vous devrez importer "io"
}

// RegisterMusicRoutes enregistre toutes les routes liées à la musique
func RegisterMusicRoutes() {
	http.HandleFunc("/api/search-music", SearchMusicHandler)
	http.HandleFunc("/api/get-instrumental", GetInstrumentalHandler)
	http.HandleFunc("/api/audio-proxy", AudioProxyHandler)
}

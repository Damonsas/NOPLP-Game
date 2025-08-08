package game

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

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

type SpotifyTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

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

type MusicSearchResult struct {
	Title      string `json:"title"`
	Artist     string `json:"artist"`
	PreviewURL string `json:"preview_url"`
	SourceURL  string `json:"source_url"`
	Source     string `json:"source"`
}

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

	if config.IsLocalEnabled() {
		if localResult := searchLocal(title, artist, instrumental); localResult != nil {
			results = append(results, *localResult)
		}
	}

	if deezerResult, err := searchDeezer(title, artist); err == nil && deezerResult != nil {
		results = append(results, *deezerResult)
	}

	if itunesResult, err := searchITunes(title, artist); err == nil && itunesResult != nil {
		results = append(results, *itunesResult)
	}

	if config.IsSpotifyEnabled() {
		if spotifyResult, err := searchSpotify(title, artist); err == nil && spotifyResult != nil {
			results = append(results, *spotifyResult)
		}
	}

	if config.IsYouTubeEnabled() && instrumental {
		if youtubeResult, err := searchYouTube(title, artist, instrumental); err == nil && youtubeResult != nil {
			results = append(results, *youtubeResult)
		}
	}

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

func AudioProxyHandler(w http.ResponseWriter, r *http.Request) {
	audioURL := r.URL.Query().Get("url")
	if audioURL == "" {
		http.Error(w, "URL audio manquante", http.StatusBadRequest)
		return
	}

	// Autoriser les requêtes cross-origin (CORS si front distant)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Télécharger le flux depuis l'URL (streaming)
	resp, err := http.Get(audioURL)
	if err != nil {
		http.Error(w, "Erreur lors du chargement de l'audio", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Copier les headers (type MIME, taille...)
	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.Header().Set("Content-Length", resp.Header.Get("Content-Length"))

	// Écrire le contenu audio directement dans la réponse
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		http.Error(w, "Erreur lors de la transmission", http.StatusInternalServerError)
	}
}

func searchYouTube(title, artist string, instrumental bool) (*MusicSearchResult, error) {
	config := GetMusicConfig()
	if !config.IsYouTubeEnabled() {
		return nil, fmt.Errorf("YouTube API non configurée")
	}

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

func searchSpotify(title, artist string) (*MusicSearchResult, error) {
	config := GetMusicConfig()
	if !config.IsSpotifyEnabled() {
		return nil, fmt.Errorf("spotify API non configurée")
	}

	query := title
	if artist != "" {
		query += " artist:" + artist
	}

	token, err := getSpotifyToken(config.SpotifyClientID, config.SpotifyClientSecret)
	if err != nil {
		return nil, err
	}

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

func searchLocal(title, artist string, instrumental bool) *MusicSearchResult {

	audioPath := "/audio/"

	filename := strings.ToLower(strings.ReplaceAll(title, " ", "_"))
	if artist != "" {
		filename = strings.ToLower(strings.ReplaceAll(artist, " ", "_")) + "_" + filename
	}
	if instrumental {
		filename += "_instrumental"
	}

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

	result, err := searchYouTube(title, artist, true)
	if err != nil {
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

func RegisterMusicRoutes() {
	http.HandleFunc("/api/search-music", SearchMusicHandler)
	http.HandleFunc("/api/get-instrumental", GetInstrumentalHandler)
	http.HandleFunc("/api/audio-proxy", AudioProxyHandler)
}

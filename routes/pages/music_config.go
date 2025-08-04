package game

import (
	"os"
)

// MusicConfig contient la configuration pour les API musicales
type MusicConfig struct {
	YouTubeAPIKey       string
	SpotifyClientID     string
	SpotifyClientSecret string
	LocalAudioPath      string
	EnableYouTube       bool
	EnableSpotify       bool
	EnableLocal         bool
}

// Singleton pour la configuration
var musicConfig *MusicConfig

// GetMusicConfig retourne la configuration musicale
func GetMusicConfig() *MusicConfig {
	if musicConfig == nil {
		musicConfig = loadMusicConfig()
	}
	return musicConfig
}

// loadMusicConfig charge la configuration depuis les variables d'environnement
func loadMusicConfig() *MusicConfig {
	config := &MusicConfig{
		YouTubeAPIKey:       os.Getenv("YOUTUBE_API_KEY"),
		SpotifyClientID:     os.Getenv("SPOTIFY_CLIENT_ID"),
		SpotifyClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
		LocalAudioPath:      os.Getenv("LOCAL_AUDIO_PATH"),
		EnableYouTube:       os.Getenv("ENABLE_YOUTUBE") != "false",
		EnableSpotify:       os.Getenv("ENABLE_SPOTIFY") != "false",
		EnableLocal:         os.Getenv("ENABLE_LOCAL") != "false",
	}

	// Valeurs par défaut
	if config.LocalAudioPath == "" {
		config.LocalAudioPath = "./data/audio"
	}

	return config
}

// IsYouTubeEnabled vérifie si YouTube est configuré et activé
func (c *MusicConfig) IsYouTubeEnabled() bool {
	return c.EnableYouTube && c.YouTubeAPIKey != ""
}

// IsSpotifyEnabled vérifie si Spotify est configuré et activé
func (c *MusicConfig) IsSpotifyEnabled() bool {
	return c.EnableSpotify && c.SpotifyClientID != "" && c.SpotifyClientSecret != ""
}

// IsLocalEnabled vérifie si la recherche locale est activée
func (c *MusicConfig) IsLocalEnabled() bool {
	return c.EnableLocal && c.LocalAudioPath != ""
}

// Alternative API configuration pour utiliser des APIs gratuites ou alternatives

// FreeMusicAPIs contient les configurations pour les APIs gratuites
type FreeMusicAPIs struct {
	// Deezer (gratuit pour les aperçus)
	DeezerAPIURL string

	// Last.fm (gratuit)
	LastFMAPIKey string
	LastFMAPIURL string

	// MusicBrainz (gratuit)
	MusicBrainzAPIURL string

	// iTunes/Apple Music (gratuit pour la recherche)
	iTunesAPIURL string

	// FreeSound (pour les effets sonores)
	FreeSoundAPIKey string
	FreeSoundAPIURL string
}

// GetFreeMusicAPIs retourne la configuration des APIs gratuites
func GetFreeMusicAPIs() *FreeMusicAPIs {
	return &FreeMusicAPIs{
		DeezerAPIURL:      "https://api.deezer.com",
		LastFMAPIKey:      os.Getenv("LASTFM_API_KEY"),
		LastFMAPIURL:      "https://ws.audioscrobbler.com/2.0/",
		MusicBrainzAPIURL: "https://musicbrainz.org/ws/2",
		iTunesAPIURL:      "https://itunes.apple.com/search",
		FreeSoundAPIKey:   os.Getenv("FREESOUND_API_KEY"),
		FreeSoundAPIURL:   "https://freesound.org/apiv2",
	}
}

// InitializeMusicAPIs initialise les routes et la configuration
func InitializeMusicAPIs() {
	config := GetMusicConfig()

	// Log de la configuration
	if config.IsYouTubeEnabled() {
		println("✓ YouTube API activée")
	} else {
		println("⚠ YouTube API non configurée (définissez YOUTUBE_API_KEY)")
	}

	if config.IsSpotifyEnabled() {
		println("✓ Spotify API activée")
	} else {
		println("⚠ Spotify API non configurée (définissez SPOTIFY_CLIENT_ID et SPOTIFY_CLIENT_SECRET)")
	}

	if config.IsLocalEnabled() {
		println("✓ Recherche locale activée:", config.LocalAudioPath)
	} else {
		println("⚠ Recherche locale non configurée")
	}

	// Enregistrer les routes
	RegisterMusicRoutes()
}

// Instructions pour configurer les APIs
const ConfigurationInstructions = `
# Configuration des API Musicales

## YouTube Data API v3 (Recommandé)
1. Allez sur https://console.developers.google.com/
2. Créez un nouveau projet ou sélectionnez un projet existant
3. Activez l'API "YouTube Data API v3"
4. Créez une clé API
5. Définissez la variable d'environnement: YOUTUBE_API_KEY=votre_cle_api

## Spotify Web API (Pour les aperçus audio)
1. Allez sur https://developer.spotify.com/dashboard/
2. Créez une nouvelle application
3. Obtenez votre Client ID et Client Secret
4. Définissez les variables d'environnement:
   SPOTIFY_CLIENT_ID=votre_client_id
   SPOTIFY_CLIENT_SECRET=votre_client_secret

## APIs Gratuites Alternatives

### Deezer API (Gratuit)
- Aucune clé API requise pour la recherche basique
- Aperçus audio de 30 secondes disponibles
- URL: https://api.deezer.com

### iTunes Search API (Gratuit)
- Aucune clé API requise
- Aperçus audio de 30 secondes
- URL: https://itunes.apple.com/search

### Last.fm API (Gratuit)
1. Créez un compte sur https://www.last.fm/api/account/create
2. Obtenez votre clé API
3. Définissez: LASTFM_API_KEY=votre_cle_api

### MusicBrainz API (Gratuit)
- Aucune clé API requise
- Métadonnées musicales complètes
- URL: https://musicbrainz.org/ws/2

## Configuration des fichiers locaux
1. Créez un dossier pour vos fichiers audio: ./data/audio
2. Organisez vos fichiers par structure:
   - ./data/audio/instrumentals/
   - ./data/audio/previews/
3. Définissez: LOCAL_AUDIO_PATH=./data/audio

## Variables d'environnement optionnelles
- ENABLE_YOUTUBE=true/false (défaut: true)
- ENABLE_SPOTIFY=true/false (défaut: true)
- ENABLE_LOCAL=true/false (défaut: true)
`

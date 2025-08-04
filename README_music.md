# 🎵 Système Musical pour le Jeu de Duel

Ce système permet d'intégrer la lecture de musique et l'affichage de paroles masquées dans votre jeu de duel musical.

## 🚀 Installation et Configuration

### 1. Fichiers ajoutés/modifiés

- `duellaunch.go` - Version complétée avec le lecteur musical et les paroles
- `music_api.go` - API pour la gestion de la musique
- `music_config.go` - Configuration des API musicales
- `music-client.ts` - Client TypeScript pour une meilleure UX
- Structure de dossiers recommandée :
  ```
  ./data/
  ├── audio/
  │   ├── instrumentals/
  │   └── previews/
  └── serverdata/
      └── paroledata/
          └── *.json (fichiers de paroles)
  ```

### 2. Configuration des API

#### Option 1 : APIs Gratuites (Recommandé pour débuter)

**Deezer API** (Gratuit, aucune clé requise)
- Aperçus audio de 30 secondes
- Métadonnées complètes
- Aucune configuration nécessaire

**iTunes Search API** (Gratuit, aucune clé requise)
- Aperçus audio de 30 secondes
- Large catalogue
- Aucune configuration nécessaire

#### Option 2 : APIs Premium

**YouTube Data API v3**
```bash
export YOUTUBE_API_KEY="votre_cle_youtube"
```

**Spotify Web API**
```bash
export SPOTIFY_CLIENT_ID="votre_client_id"
export SPOTIFY_CLIENT_SECRET="votre_client_secret"
```

### 3. Variables d'environnement

```bash
# APIs (optionnelles)
export YOUTUBE_API_KEY="votre_cle_youtube"
export SPOTIFY_CLIENT_ID="votre_client_id"
export SPOTIFY_CLIENT_SECRET="votre_client_secret"
export LASTFM_API_KEY="votre_cle_lastfm"

# Configuration locale
export LOCAL_AUDIO_PATH="./data/audio"

# Activation/désactivation des services
export ENABLE_YOUTUBE="true"
export ENABLE_SPOTIFY="true"
export ENABLE_LOCAL="true"
```

## 🎮 Utilisation

### 1. Structure des paroles (JSON)

Vos fichiers de paroles doivent suivre cette structure :
```json
{
    "titre": "Balance ton quoi",
    "artiste": "Angèle",
    "parole": "Ils parlent tous comme des animaux..."
}
```

### 2. Fonctionnalités disponibles

#### Aperçu des chansons
- Cliquez sur une carte de chanson pour un aperçu audio
- Recherche automatique sur Deezer/iTunes/Spotify
- Lecture d'un extrait de 30 secondes

#### Sélection pour le jeu
- Bouton "Sélectionner" charge la version instrumentale
- Affichage des paroles masquées selon le niveau de difficulté
- Contrôles de lecture intégrés

#### Masquage des paroles par niveau
- **50 points** : 80% des mots masqués (très difficile)
- **40 points** : 60% des mots masqués
- **30 points** : 40% des mots masqués
- **20 points** : 20% des mots masqués
- **10 points** : 10% des mots masqués (facile)

### 3. Routes API disponibles

```
GET /api/search-music?title=...&artist=...&instrumental=true
GET /api/get-lyrics/{level}/{songIndex}
GET /api/get-instrumental?title=...&artist=...
GET /api/audio-proxy?url=...
```

## 🛠 Intégration dans votre code Go

### 1. Initialisation

```go
package main

import (
    "your-project/game"
)

func main() {
    // Initialiser les APIs musicales
    game.InitializeMusicAPIs()
    
    // Vos autres initialisations...
    
    // Démarrer le serveur
    http.ListenAndServe(":8080", nil)
}
```

### 2. Enregistrement des routes

```go
func init() {
    // Routes existantes...
    http.HandleFunc("/duel-game", game.DisplayDuel)
    http.HandleFunc("/api/get-lyrics/", game.GetLyrics)
    
    // Nouvelles routes musicales
    game.RegisterMusicRoutes()
}
```

## 📁 Organisation des fichiers audio

### Structure recommandée
```
./data/audio/
├── instrumentals/
│   ├── angele_balance_ton_quoi_instrumental.mp3
│   └── ...
├── previews/
│   ├── angele_balance_ton_quoi_preview.mp3
│   └── ...
└── demos/
    ├── demo-audio.mp3
    └── demo-instrumental.mp3
```

### Nommage des fichiers
- Format : `{artiste}_{titre}_{type}.{extension}`
- Exemples :
  - `angele_balance_ton_quoi_instrumental.mp3`
  - `stromae_alors_on_danse_preview.mp3`

## 🎨 Personnalisation CSS

Les classes CSS disponibles pour la personnalisation :

```css
.music-player { /* Lecteur principal */ }
.lyrics-container { /* Container des paroles */ }
.lyrics-text { /* Texte des paroles */ }
.masked-text { /* Mots masqués */ }
.btn-select { /* Bouton sélectionner */ }
.btn-preview { /* Bouton aperçu */ }
.audio-controls { /* Contrôles audio */ }
```

## 🔧 Fonctions JavaScript disponibles

```javascript
// Client TypeScript (si utilisé)
const musicClient = new MusicGameClient();

// Fonctions globales
previewSong(title, artist);
selectSong(level, songIndex, title, artist);
toggleLyrics();
revealMoreLyrics(); // Révéler plus de paroles (indices)
revealAllLyrics(); // Révéler toutes les paroles
```

## 🐛 Dépannage

### Problèmes courants

1. **Pas d'audio** : Vérifiez que les fichiers de démonstration existent
2. **Paroles non chargées** : Vérifiez le format JSON et les chemins
3. **API ne fonctionne pas** : Vérifiez les clés API et les quotas

### Logs de débogage

Activez les logs pour voir quelle API est utilisée :
```bash
export DEBUG_MUSIC_API="true"
```

### Mode développement

Pour tester sans APIs externes :
```bash
export ENABLE_YOUTUBE="false"
export ENABLE_SPOTIFY="false"
export ENABLE_LOCAL="true"
```

## 📈 Améliorations possibles

1. **Cache audio** : Mise en cache des recherches fréquentes
2. **Playlists** : Génération automatique de playlists
3. **Synchronisation** : Synchronisation paroles/audio
4. **Effets** : Effets audio en temps réel
5. **Statistiques** : Tracking des chansons les plus jouées

## 🤝 Contribution

Pour contribuer au développement :
1. Ajoutez de nouveaux fournisseurs d'API dans `music_api.go`
2. Améliorez l'algorithme de masquage des paroles
3. Ajoutez des fonctionnalités au client TypeScript

## 📄 Licence

Ce code est fourni à des fins éducatives et de développement. Respectez les conditions d'utilisation des APIs tierces utilisées.
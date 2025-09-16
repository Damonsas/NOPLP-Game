package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const musixmatchAPIKey = "num api"
const musixmatchBaseURL = "lyrics.ovh"

type LyricsResponse struct {
	Message struct {
		Body struct {
			Lyrics struct {
				LyricsBody string `json:"lyrics_body"`
			} `json:"lyrics"`
		} `json:"body"`
	} `json:"message"`
}

func GetLyrics(trackID string) (string, error) {

	endpoint := fmt.Sprintf("track.lyrics.get?track_id=%s&apikey=%s", trackID, musixmatchAPIKey)
	resp, err := http.Get(musixmatchBaseURL + endpoint)
	if err != nil {
		return "", fmt.Errorf("erreur HTTP: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("erreur lecture réponse: %v", err)
	}

	var result LyricsResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("erreur parsing JSON: %v", err)
	}

	if result.Message.Body.Lyrics.LyricsBody == "" {
		return "", errors.New("paroles non trouvées")
	}

	return result.Message.Body.Lyrics.LyricsBody, nil
}

func SearchTrack(title, artist string) (string, error) {
	query := url.QueryEscape(fmt.Sprintf("%s %s", title, artist))
	endpoint := fmt.Sprintf("track.search?q_track=%s&q_artist=%s&apikey=%s", query, artist, musixmatchAPIKey)
	resp, err := http.Get(musixmatchBaseURL + endpoint)
	if err != nil {
		return "", fmt.Errorf("erreur HTTP: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("erreur lecture réponse: %v", err)
	}

	var searchResult struct {
		Message struct {
			Body struct {
				TrackList []struct {
					Track struct {
						TrackID int `json:"track_id"`
					} `json:"track"`
				} `json:"track_list"`
			} `json:"body"`
		} `json:"message"`
	}
	if err := json.Unmarshal(body, &searchResult); err != nil {
		return "", fmt.Errorf("erreur parsing JSON: %v", err)
	}

	if len(searchResult.Message.Body.TrackList) == 0 {
		return "", errors.New("piste non trouvée")
	}

	return fmt.Sprintf("%d", searchResult.Message.Body.TrackList[0].Track.TrackID), nil
}

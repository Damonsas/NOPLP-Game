package game

import (
	"strings"
)

type MemeChansonData struct {
	Titre   string `json:"titre"`
	Artiste string `json:"artiste"`
	Parole  string `json:"parole"`
}

func Duelmemechanson(lines []string) string {
	lyrics := strings.Join(lines, " ")
	if lyrics == "" {
		return lyrics
	}

	maskPercentage := 0.8
	words := strings.Fields(lyrics)
	wordsToMask := int(float64(len(words)) * maskPercentage)
	maskedWords := make([]string, len(words))

	for i, word := range words {
		if i < wordsToMask {
			maskedWords[i] = strings.Repeat("█", len(word))
		} else {
			maskedWords[i] = word
		}
	}

	return strings.Join(maskedWords, " ")
}

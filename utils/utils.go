package utils

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

type Album struct {
	Title    string
	Hashtags []string
	Songs    []Song
}

type Song struct {
	Title     string
	SpotifyID string `json:"spotify_id"`
	Link      string
	Emoji     string
	Hashtags  []string
	Lyrics    [][]string
}

func GetHashtagsAsString(hashtags []string) (hashtagString string) {
	return "#" + strings.Join(hashtags, " #")
}

func GetAlbum(name string) (album Album) {
	rawData, err := os.ReadFile("albums/" + name + ".json")

	if err != nil {
		log.Panic(err)
	}

	if err := json.Unmarshal(rawData, &album); err != nil {
		log.Panic(err)
	}

	return album
}

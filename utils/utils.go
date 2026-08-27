package utils

import (
	"encoding/json/v2"
	"log"
	"os"
	"strings"
)

type Album struct {
	Title    string   `json:"title"`
	Hashtags []string `json:"hastags"`
	Songs    []Song   `json:"songs"`
}

type Song struct {
	Title     string     `json:"title"`
	SpotifyID string     `json:"spotify_id"`
	Link      string     `json:"link"`
	Emoji     string     `json:"emoji"`
	Hashtags  []string   `json:"hashtags"`
	Lyrics    [][]string `json:"lyrics"`
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

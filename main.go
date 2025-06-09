package main

import (
	"encoding/json"
	"flag"
	"fmt"
	_ "image/png"
	"math/rand"
	"os"
	"strings"
)

const WIDTH int = 1080
const MARGIN int = 78
const SPOTIFY_TRACK_LINK = "https://open.spotify.com/track/%s"

type Album struct {
	Title string
	Songs []Song
}

type Song struct {
	Title     string
	SpotifyID string `json:"spotify_id"`
	Lyrics    string
	Emoji     string
	Link      string
}

var albums = []string{
	"days_of_thunder",
	"endless_summer",
	"heroes",
	"horror_show",
	"kids",
	"monsters",
	"nocturnal",
	"songs",
}

func main() {
	shouldGenerateAllImages := flag.Bool("generate-all-images", false, "Generate all images")
	flag.Parse()

	if *shouldGenerateAllImages {
		generateAllImages()
		os.Exit(0)
	}

	randomAlbumName := albums[rand.Intn(len(albums))]
	album := getAlbum(randomAlbumName)
	song := album.Songs[rand.Intn(len(album.Songs))]

	lyricParts := strings.Split(song.Lyrics, "|")
	lyrics := lyricParts[rand.Intn(len(lyricParts))]

	var reply, link string

	fmt.Printf("Album title: %s\n", album.Title)
	fmt.Printf("Song title: %s\n", song.Title)

	if song.SpotifyID != "" {
		fmt.Printf(SPOTIFY_TRACK_LINK+"\n", song.SpotifyID)
		link = fmt.Sprintf(SPOTIFY_TRACK_LINK, song.SpotifyID)
	}

	if song.Link != "" {
		fmt.Printf("Link: %s\n", song.Link)
		link = song.Link
	}

	if song.Emoji != "" {
		fmt.Printf("Emoji: %s\n", song.Emoji)
		reply = song.Emoji
	}

	fmt.Printf("Lyrics: \n%s\n", lyrics)

	//generateImage(randomAlbumName, album, song, lyrics, "preview.png")

	if os.Getenv("BOTSKY_HANDLE") != "" && os.Getenv("BOTSKY_APPKEY") != "" {
		postToBluesky(lyrics, reply, link)
	}
}

func getAlbum(name string) (album Album) {
	rawData, err := os.ReadFile("albums/" + name + ".json")

	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(rawData, &album); err != nil {
		panic(err)
	}

	return album
}

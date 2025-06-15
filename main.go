package main

import (
	"encoding/json"
	"flag"
	"fmt"
	_ "image/png"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

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
	generateAllImagesFlag := flag.Bool("generate-all-images", false, "Generate all images (for testing only)")
	fakeMidnightFlag := flag.Bool("fake-midnight", false, "Fake that it's midnight (for testing only)")
	flag.Parse()

	if *generateAllImagesFlag {
		generateAllImages()
		os.Exit(0)
	}

	randomAlbumName := albums[rand.Intn(len(albums))]
	album := getAlbum(randomAlbumName)
	song := album.Songs[rand.Intn(len(album.Songs))]

	//Special post at midnight
	if (time.Now().Hour() == 0 && time.Now().Minute() == 0) || *fakeMidnightFlag {
		randomAlbumName = "songs"
		song = Song{"Midnight", "", "We are one beating heart", "💓", ""}
		album = Album{"Midnight", []Song{song}}
	}

	lyricParts := strings.Split(song.Lyrics, "|")
	lyrics := lyricParts[rand.Intn(len(lyricParts))]

	var reply, link string

	log.Printf("Album title: %s\n", album.Title)
	log.Printf("Song title: %s\n", song.Title)

	if song.SpotifyID != "" {
		log.Printf("Spotify link: "+SPOTIFY_TRACK_LINK+"\n", song.SpotifyID)
		link = fmt.Sprintf(SPOTIFY_TRACK_LINK, song.SpotifyID)
	}

	if song.Link != "" {
		log.Printf("Link: %s\n", song.Link)
		link = song.Link
	}

	if song.Emoji != "" {
		log.Printf("Emoji: %s\n", song.Emoji)
		reply = song.Emoji
	}

	log.Printf("Lyrics: \n%s\n", lyrics)

	if os.Getenv("BOTSKY_HANDLE") != "" && os.Getenv("BOTSKY_APPKEY") != "" {
		postToBluesky(lyrics, reply, link)
	}

	if os.Getenv("THREADS_ACCESS_TOKEN") != "" {
		postToThreads(lyrics, reply, link)
	}
}

func getAlbum(name string) (album Album) {
	rawData, err := os.ReadFile("albums/" + name + ".json")

	if err != nil {
		log.Panic(err)
	}

	if err := json.Unmarshal(rawData, &album); err != nil {
		log.Panic(err)
	}

	return album
}

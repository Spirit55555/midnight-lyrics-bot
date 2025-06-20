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
	Title    string
	Songs    []Song
	Hashtags []string
}

type Song struct {
	Title     string
	SpotifyID string `json:"spotify_id"`
	Lyrics    string
	Emoji     string
	Link      string
	Hashtags  []string
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

var defaultHashtags = []string{
	"themidnight",
	"synthwave",
}

func main() {
	blueskyFlag := flag.Bool("bluesky", false, "If it should post to Bluesky")
	threadsFlag := flag.Bool("threads", false, "If it should post to Threads")
	instagramFlag := flag.Bool("instagram", false, "If it should post to Instagram")
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
		song = Song{"Midnight", "", "We are one beating heart", "💓", "", []string{}}
		album = Album{"Midnight", []Song{song}, []string{}}
	}

	lyricParts := strings.Split(song.Lyrics, "|")
	lyrics := lyricParts[rand.Intn(len(lyricParts))]

	var reply, link string
	var hashtags []string

	log.Printf("Album title: %s\n", album.Title)

	if len(album.Hashtags) > 0 {
		log.Printf("Album hashtags: %s\n", getHashtagsAsString(album.Hashtags))
		hashtags = append(hashtags, album.Hashtags...)
	}

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

	if len(song.Hashtags) > 0 {
		log.Printf("Song hashtags: %s\n", getHashtagsAsString(song.Hashtags))
		hashtags = append(hashtags, song.Hashtags...)
	}

	log.Printf("Lyrics: \n%s\n", lyrics)

	//Add the default hashtags
	hashtags = append(hashtags, defaultHashtags...)

	if *blueskyFlag && os.Getenv("BOTSKY_HANDLE") != "" && os.Getenv("BOTSKY_APPKEY") != "" {
		postToBluesky(lyrics, reply, link, hashtags)
	}

	if *threadsFlag && os.Getenv("THREADS_ACCESS_TOKEN") != "" {
		postToThreads(lyrics, reply, link)
	}

	if *instagramFlag && os.Getenv("INSTAGRAM_ACCESS_TOKEN") != "" && os.Getenv("INSTAGRAM_IMAGES_URL") != "" {
		imagePath := generateImage(randomAlbumName, album, song, lyrics)
		postToInstagram(lyrics, reply, imagePath, hashtags)
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

func getHashtagsAsString(hashtags []string) (hashtagString string) {
	return "#" + strings.Join(hashtags, " #")
}

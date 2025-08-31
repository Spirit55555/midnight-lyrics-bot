package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

const SPOTIFY_TRACK_LINK = "https://open.spotify.com/track/%s"

const ALT_TEXT_SONG = "Lyrics from a The Midnight song. \nSong title: %s. \nLyrics: %s"
const ALT_TEXT_ALBUM_SONG = "Lyrics from a The Midnight song. \nAlbum title: %s. \nSong title: %s. \nLyrics: %s"

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

var albums = []string{
	"days_of_thunder",
	"endless_summer",
	"heroes",
	"horror_show",
	"kids",
	"monsters",
	"nocturnal",
	"syndicate",
	"night_drive",
	"silence",
	"land_locked_heart",
	"vehlinggo"
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

	args := flag.Args()

	if *generateAllImagesFlag {
		generateAllImages()
		os.Exit(0)
	}

	albumName := albums[rand.Intn(len(albums))]

	if (len(args) == 1 || len(args) == 2) && args[0] != "" && slices.Contains(albums, args[0]) {
		albumName = args[0]
	}

	album := getAlbum(albumName)
	song := album.Songs[rand.Intn(len(album.Songs))]

	if len(args) == 2 && args[1] != "" {
		songId, _ := strconv.Atoi(args[1])
		if len(album.Songs) > (songId - 1) {
			song = album.Songs[songId-1]
		}
	}

	//Special post at midnight
	if (time.Now().Hour() == 0 && time.Now().Minute() == 0) || *fakeMidnightFlag {
		albumName = "songs"
		song = Song{"Midnight", "", "", "💓", []string{}, [][]string{{"We are one beating heart"}}}
		album = Album{"Midnight", []string{}, []Song{song}}
	}

	lyrics := strings.Join(song.Lyrics[rand.Intn(len(song.Lyrics))], "\n")

	var reply, link string
	var hashtags []string
	var altText string

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

	//Generate alt text for images
	if albumName == "songs" {
		altText = fmt.Sprintf(ALT_TEXT_SONG, song.Title, lyrics)
	} else {
		altText = fmt.Sprintf(ALT_TEXT_ALBUM_SONG, album.Title, song.Title, lyrics)
	}

	log.Printf("Alt text: %s\n", altText)

	if *blueskyFlag && os.Getenv("BOTSKY_HANDLE") != "" && os.Getenv("BOTSKY_APPKEY") != "" {
		postToBluesky(lyrics, reply, link, hashtags)
	}

	if *threadsFlag && os.Getenv("THREADS_ACCESS_TOKEN") != "" {
		postToThreads(lyrics, reply, link)
	}

	if *instagramFlag && os.Getenv("INSTAGRAM_ACCESS_TOKEN") != "" && os.Getenv("INSTAGRAM_IMAGES_URL") != "" {
		imagePath := generateImage(albumName, album, song, lyrics)
		postToInstagram(reply, imagePath, altText, hashtags)
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

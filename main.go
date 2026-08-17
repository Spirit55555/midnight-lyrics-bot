package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/Spirit55555/midnight-lyrics-bot/bluesky"
	"github.com/Spirit55555/midnight-lyrics-bot/image"
	"github.com/Spirit55555/midnight-lyrics-bot/meta/instagram"
	"github.com/Spirit55555/midnight-lyrics-bot/meta/threads"
	"github.com/Spirit55555/midnight-lyrics-bot/utils"
)

const SPOTIFY_TRACK_LINK = "https://open.spotify.com/track/%s"

const ALT_TEXT_SONG = "Lyrics from a The Midnight song. \nSong title: %s. \nLyrics: %s"
const ALT_TEXT_ALBUM_SONG = "Lyrics from a The Midnight song. \nAlbum title: %s. \nSong title: %s. \nLyrics: %s"

var albums = []string{
	"cold_pizza",
	"days_of_thunder",
	"endless_summer",
	"heroes",
	"horror_show",
	"kids",
	"land_locked_heart",
	"monsters",
	"neon_oddysey",
	"night_drive",
	"nocturnal",
	"silence",
	"songs",
	"syndicate",
	"syndicate_deluxe",
	"tonight_tonight",
	"vehlinggo",
}

var defaultHashtags = []string{
	"themidnight",
	"synthwave",
}

func main() {
	blueskyFlag := flag.Bool("bluesky", false, "If it should post to Bluesky")
	threadsFlag := flag.Bool("threads", false, "If it should post to Threads")
	instagramFlag := flag.Bool("instagram", false, "If it should post to Instagram")
	refreshTokensFlag := flag.Bool("refresh-tokens", false, "Refresh and print new tokens for Threads and Instagram")
	generateImageFlag := flag.Bool("generate-image", false, "Generate the lyrics image, useful when specifying the album and/or song (for testing only)")
	generateAllImagesFlag := flag.Bool("generate-all-images", false, "Generate all images (for testing only)")
	fakeMidnightFlag := flag.Bool("fake-midnight", false, "Fake that it's midnight (for testing only)")
	flag.Parse()

	args := flag.Args()

	if *generateAllImagesFlag {
		image.GenerateAll(albums)
		os.Exit(0)
	}

	if *refreshTokensFlag {
		if os.Getenv("THREADS_ACCESS_TOKEN") != "" {
			threadsResponse := threads.RefreshToken()

			if threadsResponse.AccessToken != "" {
				log.Printf("Threads new token: %s", threadsResponse.AccessToken)
			} else {
				log.Printf("Threads error message: %s", threadsResponse.Error.Message)
			}
		}

		if os.Getenv("INSTAGRAM_ACCESS_TOKEN") != "" {
			instagramResponse := instagram.RefreshToken()

			if instagramResponse.AccessToken != "" {
				log.Printf("Instagram new token: %s", instagramResponse.AccessToken)
			} else {
				log.Printf("Instagram error message: %s", instagramResponse.Error.Message)
			}
		}

		os.Exit(0)
	}

	albumName := albums[rand.Intn(len(albums))]

	if (len(args) == 1 || len(args) == 2) && args[0] != "" && slices.Contains(albums, args[0]) {
		albumName = args[0]
	}

	album := utils.GetAlbum(albumName)
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
		song = utils.Song{Title: "Midnight", SpotifyID: "", Link: "", Emoji: "💓", Hashtags: []string{}, Lyrics: [][]string{{"We are one beating heart"}}}
		album = utils.Album{Title: "Midnight", Hashtags: []string{}, Songs: []utils.Song{song}}
	}

	lyrics := strings.Join(song.Lyrics[rand.Intn(len(song.Lyrics))], "\n")

	var reply, link string
	var hashtags []string
	var altText string

	log.Printf("Album title: %s\n", album.Title)

	if len(album.Hashtags) > 0 {
		log.Printf("Album hashtags: %s\n", utils.GetHashtagsAsString(album.Hashtags))
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
		log.Printf("Song hashtags: %s\n", utils.GetHashtagsAsString(song.Hashtags))
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

	if *generateImageFlag {
		imagePath := image.Generate(albumName, album, song, lyrics)
		log.Printf("Image generated: %s\n", imagePath)
	}

	if *blueskyFlag && os.Getenv("BOTSKY_HANDLE") != "" && os.Getenv("BOTSKY_APPKEY") != "" {
		bluesky.Post(lyrics, reply, link, hashtags)
	}

	if *threadsFlag && os.Getenv("THREADS_ACCESS_TOKEN") != "" {
		threads.Post(lyrics, reply, link)
	}

	if *instagramFlag && os.Getenv("INSTAGRAM_ACCESS_TOKEN") != "" && os.Getenv("INSTAGRAM_IMAGES_URL") != "" {
		imagePath := image.Generate(albumName, album, song, lyrics)
		instagram.Post(reply, imagePath, altText, hashtags)
	}
}

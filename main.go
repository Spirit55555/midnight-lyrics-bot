package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"image"
	_ "image/png"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"git.sr.ht/~sbinet/gg"
)

const WIDTH int = 1080
const MARGIN int = 78

type Album struct {
	Title string
	Songs []Song
}

type Song struct {
	Title     string
	SpotifyID string `json:"spotify_id"`
	Lyrics    string
	Emoji     string
	Reply     string
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

	fmt.Printf("Album title: %s\n", album.Title)
	fmt.Printf("Song title: %s\n", song.Title)
	fmt.Printf("Spotify link: https://open.spotify.com/track/%s\n", song.SpotifyID)
	fmt.Printf("Lyrics: \n%s\n", lyrics)

	generateImage(randomAlbumName, album, song, lyrics, "preview.png")
}

func generateImage(albumName string, album Album, song Song, lyrics string, finalImagePath string) {
	albumFile, _ := os.Open("images/" + albumName + ".png")
	defer albumFile.Close()

	albumBackground, _, _ := image.Decode(albumFile)

	title := song.Title

	if albumName != "songs" {
		title = album.Title + "\n" + title
	}

	dc := gg.NewContextForImage(albumBackground)
	dc.SetHexColor("F4F4F4")

	dc.LoadFontFace("fonts/Optiker-K.ttf", 62)
	dc.DrawStringWrapped(lyrics, float64(WIDTH/2), float64(WIDTH/2), 0.5, 0.5, float64(WIDTH-(MARGIN*2)), 1.2, gg.AlignCenter)

	dc.LoadFontFace("fonts/Optiker-K.ttf", 40)
	dc.DrawStringWrapped(title, float64(MARGIN), float64(WIDTH-170), 0, 0, 800, 1.5, gg.AlignLeft)

	dc.SavePNG(finalImagePath)
}

func generateAllImages() {
	for _, albumName := range albums {
		album := getAlbum(albumName)

		fmt.Printf("Album: %s\n", album.Title)

		for _, song := range album.Songs {
			songFolder := filepath.Join("generated_images", albumName, strings.ToLower(strings.ReplaceAll(song.Title, " ", "_")))
			os.MkdirAll(songFolder, os.ModePerm)

			lyricParts := strings.Split(song.Lyrics, "|")

			fmt.Printf("Song: %s\n", song.Title)
			fmt.Printf("Total lyrics: %d\n", len(lyricParts))

			for i, lyrics := range lyricParts {
				fmt.Printf("(%d/%d)\n", i+1, len(lyricParts))

				hash := sha1.Sum([]byte(lyrics))

				generateImage(albumName, album, song, lyrics, filepath.Join(songFolder, hex.EncodeToString(hash[:])+".png"))
			}
		}
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

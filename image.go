package main

import (
	"crypto/sha1"
	"encoding/hex"
	"image"
	"log"
	"os"
	"path/filepath"
	"strings"

	"git.sr.ht/~sbinet/gg"
)

const WIDTH int = 1090
const HEIGHT int = 1350
const MARGIN int = 78

const FONT = "fonts/Optiker-K.ttf"
const TEXT_COLOR = "F4F4F4"
const SHADOW_COLOR = "000000"
const SHADOW_OFFSET = 2

func generateImage(albumName string, album Album, song Song, lyrics string) (imagePath string) {
	albumFile, _ := os.Open("images/" + albumName + ".png")
	defer albumFile.Close()

	albumBackground, _, _ := image.Decode(albumFile)

	title := song.Title

	if albumName != "songs" {
		title = album.Title + "\n" + title
	}

	dc := gg.NewContextForImage(albumBackground)

	//Lyrics
	dc.LoadFontFace(FONT, 62)

	dc.SetHexColor(SHADOW_COLOR)
	dc.DrawStringWrapped(lyrics, float64(WIDTH/2)+SHADOW_OFFSET, float64(HEIGHT/2)+SHADOW_OFFSET, 0.5, 0.5, float64(WIDTH-(MARGIN*2)), 1.2, gg.AlignCenter)

	dc.SetHexColor(TEXT_COLOR)
	dc.DrawStringWrapped(lyrics, float64(WIDTH/2), float64(HEIGHT/2), 0.5, 0.5, float64(WIDTH-(MARGIN*2)), 1.2, gg.AlignCenter)

	//Album/song title
	dc.LoadFontFace(FONT, 40)

	dc.SetHexColor(SHADOW_COLOR)
	dc.DrawStringWrapped(title, float64(MARGIN)+SHADOW_OFFSET, float64(HEIGHT-175)+SHADOW_OFFSET, 0, 0, 700, 1.5, gg.AlignLeft)

	dc.SetHexColor(TEXT_COLOR)
	dc.DrawStringWrapped(title, float64(MARGIN), float64(HEIGHT-175), 0, 0, 700, 1.5, gg.AlignLeft)

	hash := sha1.Sum([]byte(lyrics))
	imageName := hex.EncodeToString(hash[:]) + ".jpg"
	songFolder := filepath.Join(albumName, strings.ToLower(strings.ReplaceAll(song.Title, " ", "_")))
	os.MkdirAll(filepath.Join("generated_images", songFolder), os.ModePerm)

	dc.SaveJPG(filepath.Join("generated_images", songFolder, imageName), 100)

	return filepath.Join(songFolder, imageName)
}

func generateAllImages() {
	for _, albumName := range albums {
		album := getAlbum(albumName)

		log.Printf("Album: %s\n", album.Title)

		for _, song := range album.Songs {
			lyricParts := song.Lyrics

			log.Printf("Song: %s\n", song.Title)
			log.Printf("Total lyrics: %d\n", len(lyricParts))

			for i, lyrics := range lyricParts {
				log.Printf("(%d/%d)\n", i+1, len(lyricParts))

				generateImage(albumName, album, song, strings.Join(lyrics, "\n"))
			}
		}
	}
}

package main

import (
	"errors"
	"os"
	"testing"
)

func TestIfCoverImagesExist(t *testing.T) {
	for _, albumName := range albums {
		if _, err := os.Stat("images/" + albumName + ".png"); errors.Is(err, os.ErrNotExist) {
			t.Errorf("Cover image not found for album: %s", albumName)
		}
	}
}

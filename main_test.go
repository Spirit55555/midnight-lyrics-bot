package main

import (
	"errors"
	"os"
	"slices"
	"testing"
)

func TestIfCoverImagesExist(t *testing.T) {
	for _, albumName := range albums {
		if _, err := os.Stat("images/" + albumName + ".png"); errors.Is(err, os.ErrNotExist) {
			t.Errorf("Cover image not found for album: %s", albumName)
		}
	}
}

func TestAlbumsIsSorted(t *testing.T) {
	if !slices.IsSorted(albums) {
		sortedALbums := slices.Clone(albums)
		slices.Sort(sortedALbums)

		t.Errorf("Albums are not sorted in ascending order.\nFound:    %v\nExpected: %v", albums, sortedALbums)
	}
}

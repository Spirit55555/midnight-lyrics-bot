package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"maps"
	"net/http"
	"net/url"
	"os"
)

type InstagramResponse struct {
	Id    string
	Error struct {
		Message string
	}
}

func postToInstagram(caption, imagePath, altText string, hashtags []string) {
	imageURL := os.Getenv("INSTAGRAM_IMAGES_URL") + "/" + imagePath

	//Add #hashtags to caption
	if len(hashtags) > 0 {
		hashtagString := getHashtagsAsString(hashtags)

		if caption != "" {
			caption = caption + " " + hashtagString
		} else {
			caption = hashtagString
		}
	}

	//Create post
	container := makeRequestToInstagram("media", url.Values{"image_url": {imageURL}, "caption": {caption}, "alt_text": {altText}})

	if container.Id != "" {
		log.Printf("Instagram container ID: %s", container.Id)
	} else {
		log.Printf("Instagram container error: %s", container.Error.Message)
	}

	//Publish post
	publish := makeRequestToInstagram("media_publish", url.Values{"creation_id": {container.Id}})

	if publish.Id != "" {
		log.Printf("Instagram publish ID: %s", publish.Id)
	} else {
		log.Printf("Instagram publish error: %s", container.Error.Message)
	}
}

func makeRequestToInstagram(endpoint string, params url.Values) (response InstagramResponse) {
	accessToken := os.Getenv("INSTAGRAM_ACCESS_TOKEN")
	maps.Copy(params, url.Values{"access_token": {accessToken}})

	resp, err := http.PostForm(fmt.Sprintf("https://graph.instagram.com/v23.0/me/%s", endpoint), params)

	if err != nil {
		log.Panic(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("Instagram returned \"%s\" for request to %s", resp.Status, resp.Request.URL)

		return response
	}

	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Panic(err)
	}

	if err := json.Unmarshal(respBody, &response); err != nil {
		log.Panic(err)
	}

	return response
}

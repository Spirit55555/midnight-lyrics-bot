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
	Id string
}

func postToInstagram(post, caption, imagePath string) {
	imageURL := os.Getenv("INSTAGRAM_IMAGES_URL") + "/" + imagePath

	//Create post
	container := makeRequestToInstagram("media", url.Values{"image_url": {imageURL}, "caption": {caption}, "alt_text": {post}})
	log.Printf("Instagram container ID: %s", container.Id)

	//Publish post
	publish := makeRequestToInstagram("media_publish", url.Values{"creation_id": {container.Id}})
	log.Printf("Instagram publish ID: %s", publish.Id)
}

func makeRequestToInstagram(endpoint string, params url.Values) (response InstagramResponse) {
	accessToken := os.Getenv("INSTAGRAM_ACCESS_TOKEN")
	maps.Copy(params, url.Values{"access_token": {accessToken}})

	resp, err := http.PostForm(fmt.Sprintf("https://graph.instagram.com/v23.0/me/%s", endpoint), params)

	if err != nil {
		log.Panic(err)
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Panic(err)
	}

	if err := json.Unmarshal(respBody, &response); err != nil {
		log.Panic(err)
	}

	return response
}

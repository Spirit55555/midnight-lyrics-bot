package instagram

import (
	"encoding/json"
	"fmt"
	"log"
	"maps"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/Spirit55555/midnight-lyrics-bot/utils"
)

type InstagramResponse struct {
	Id    string
	Error struct {
		Message string
	}

	// Used by refresh_access_token endpoint
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func Post(caption, imagePath, altText string, hashtags []string) {
	imageURL := os.Getenv("INSTAGRAM_IMAGES_URL") + "/" + imagePath

	//Add #hashtags to caption
	if len(hashtags) > 0 {
		hashtagString := utils.GetHashtagsAsString(hashtags)

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

	//Sleep for 5 secs, container should be ready by then
	time.Sleep(5 * time.Second)

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
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Panic(err)
	}

	return response
}

func RefreshToken() (response InstagramResponse) {
	accessToken := os.Getenv("INSTAGRAM_ACCESS_TOKEN")

	url, _ := url.Parse(fmt.Sprintf("https://graph.instagram.com/v23.0/%s", "refresh_access_token"))
	query := url.Query()

	query.Add("grant_type", "ig_refresh_token")
	query.Add("access_token", accessToken)

	url.RawQuery = query.Encode()

	resp, err := http.Get(url.String())

	if err != nil {
		log.Panic(err)
	}

	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Panic(err)
	}

	return response
}

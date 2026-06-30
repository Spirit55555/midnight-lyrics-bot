package instagram

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/Spirit55555/midnight-lyrics-bot/meta"
	"github.com/Spirit55555/midnight-lyrics-bot/utils"
)

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
	container := makeRequest("media", url.Values{"image_url": {imageURL}, "caption": {caption}, "alt_text": {altText}})

	if container.Id != "" {
		log.Printf("Instagram container ID: %s", container.Id)
	} else {
		log.Printf("Instagram container error: %s", container.Error.Message)
	}

	//Sleep for 5 secs, container should be ready by then
	time.Sleep(5 * time.Second)

	//Publish post
	publish := makeRequest("media_publish", url.Values{"creation_id": {container.Id}})

	if publish.Id != "" {
		log.Printf("Instagram publish ID: %s", publish.Id)
	} else {
		log.Printf("Instagram publish error: %s", container.Error.Message)
	}
}

func makeRequest(endpoint string, params url.Values) (response meta.Response) {
	accessToken := os.Getenv("INSTAGRAM_ACCESS_TOKEN")
	endpoint = fmt.Sprintf("me/%s", endpoint)

	response = meta.MakePOSTRequest(meta.INSTAGRAM, accessToken, endpoint, params)

	return response
}

func RefreshToken() (response meta.RefreshTokenResponse) {
	accessToken := os.Getenv("INSTAGRAM_ACCESS_TOKEN")

	response = meta.RefreshToken(meta.INSTAGRAM, accessToken)

	return response
}

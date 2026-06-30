package threads

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/Spirit55555/midnight-lyrics-bot/meta"
)

func Post(post, reply, link string) {
	//Add link to reply
	if reply != "" && link != "" {
		reply = reply + " " + link
	} else if reply == "" {
		reply = link
	}

	//Create post
	collection := makeRequest("threads", url.Values{"media_type": {"TEXT"}, "text": {post}})

	if collection.Id != "" {
		log.Printf("Threads collection ID: %s", collection.Id)
	} else {
		log.Printf("Threads collection error: %s", collection.Error.Message)
	}

	//Sleep for 5 secs, collection should be ready by then
	time.Sleep(5 * time.Second)

	//Publish post
	publish := makeRequest("threads_publish", url.Values{"creation_id": {collection.Id}})

	if publish.Id != "" {
		log.Printf("Threads publish ID: %s", publish.Id)
	} else {
		log.Printf("Threads public error: %s", publish.Error.Message)
	}

	//Create reply
	replyCollection := makeRequest("threads", url.Values{"media_type": {"TEXT"}, "text": {reply}, "reply_to_id": {publish.Id}})

	if replyCollection.Id != "" {
		log.Printf("Threads reply collection ID: %s", replyCollection.Id)
	} else {
		log.Printf("Threads reply collection error: %s", replyCollection.Error.Message)
	}

	//Sleep for 5 secs, reply should be ready by then
	time.Sleep(5 * time.Second)

	//Post reply
	replyPublish := makeRequest("threads_publish", url.Values{"creation_id": {replyCollection.Id}})

	if replyPublish.Id != "" {
		log.Printf("Threads reply publish ID: %s", replyPublish.Id)
	} else {
		log.Printf("Threads reply publish error: %s", replyPublish.Error.Message)
	}
}

func makeRequest(endpoint string, params url.Values) (response meta.Response) {
	accessToken := os.Getenv("THREADS_ACCESS_TOKEN")
	endpoint = fmt.Sprintf("me/%s", endpoint)

	response = meta.MakePOSTRequest(meta.THREADS, accessToken, endpoint, params)

	return response
}

func RefreshToken() (response meta.RefreshTokenResponse) {
	accessToken := os.Getenv("THREADS_ACCESS_TOKEN")

	response = meta.RefreshToken(meta.THREADS, accessToken)

	return response
}

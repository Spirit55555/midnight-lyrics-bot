package threads

import (
	"encoding/json"
	"fmt"
	"log"
	"maps"
	"net/http"
	"net/url"
	"os"
	"time"
)

type ThreadsResponse struct {
	Id    string
	Error struct {
		Message string
	}

	// Used by refresh_access_token endpoint
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func Post(post, reply, link string) {
	//Add link to reply
	if reply != "" && link != "" {
		reply = reply + " " + link
	} else if reply == "" {
		reply = link
	}

	//Create post
	collection := makeRequestToThreads("threads", url.Values{"media_type": {"TEXT"}, "text": {post}})

	if collection.Id != "" {
		log.Printf("Threads collection ID: %s", collection.Id)
	} else {
		log.Printf("Threads collection error: %s", collection.Error.Message)
	}

	//Sleep for 5 secs, collection should be ready by then
	time.Sleep(5 * time.Second)

	//Publish post
	publish := makeRequestToThreads("threads_publish", url.Values{"creation_id": {collection.Id}})

	if publish.Id != "" {
		log.Printf("Threads publish ID: %s", publish.Id)
	} else {
		log.Printf("Threads public error: %s", publish.Error.Message)
	}

	//Create reply
	replyCollection := makeRequestToThreads("threads", url.Values{"media_type": {"TEXT"}, "text": {reply}, "reply_to_id": {publish.Id}})

	if replyCollection.Id != "" {
		log.Printf("Threads reply collection ID: %s", replyCollection.Id)
	} else {
		log.Printf("Threads reply collection error: %s", replyCollection.Error.Message)
	}

	//Sleep for 5 secs, reply should be ready by then
	time.Sleep(5 * time.Second)

	//Post reply
	replyPublish := makeRequestToThreads("threads_publish", url.Values{"creation_id": {replyCollection.Id}})

	if replyPublish.Id != "" {
		log.Printf("Threads reply publish ID: %s", replyPublish.Id)
	} else {
		log.Printf("Threads reply publish error: %s", replyPublish.Error.Message)
	}
}

func makeRequestToThreads(endpoint string, params url.Values) (response ThreadsResponse) {
	accessToken := os.Getenv("THREADS_ACCESS_TOKEN")
	maps.Copy(params, url.Values{"access_token": {accessToken}})

	resp, err := http.PostForm(fmt.Sprintf("https://graph.threads.net/v1.0/me/%s", endpoint), params)

	if err != nil {
		log.Panic(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("Threads returned \"%s\" for request to %s", resp.Status, resp.Request.URL)
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Panic(err)
	}

	return response
}

func RefreshToken() (response ThreadsResponse) {
	accessToken := os.Getenv("THREADS_ACCESS_TOKEN")

	url, _ := url.Parse(fmt.Sprintf("https://graph.threads.net/v1.0/%s", "refresh_access_token"))
	query := url.Query()

	query.Add("grant_type", "th_refresh_token")
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

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

type ThreadsResponse struct {
	Id string
}

func postToThreads(post, reply, link string) {
	//Add link to reply
	if reply != "" && link != "" {
		reply = reply + " " + link
	} else if reply == "" {
		reply = link
	}

	//Create post
	collection := makeRequestToThreads("threads", url.Values{"media_type": {"TEXT"}, "text": {post}})
	log.Printf("Threads collection ID: %s", collection.Id)

	//Publish post
	publish := makeRequestToThreads("threads_publish", url.Values{"creation_id": {collection.Id}})
	log.Printf("Threads publish ID: %s", publish.Id)

	//Create reply
	replyCollection := makeRequestToThreads("threads", url.Values{"media_type": {"TEXT"}, "text": {reply}, "reply_to_id": {publish.Id}})
	log.Printf("Threads reply collection ID: %s", replyCollection.Id)

	//Post reply
	replyPublish := makeRequestToThreads("threads_publish", url.Values{"creation_id": {replyCollection.Id}})
	log.Printf("Threads reply publish ID: %s", replyPublish.Id)
}

func makeRequestToThreads(endpoint string, params url.Values) (response ThreadsResponse) {
	accessToken := os.Getenv("THREADS_ACCESS_TOKEN")
	maps.Copy(params, url.Values{"access_token": {accessToken}})

	resp, err := http.PostForm(fmt.Sprintf("https://graph.threads.net/v1.0/me/%s", endpoint), params)

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

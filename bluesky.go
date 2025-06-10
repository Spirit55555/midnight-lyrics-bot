package main

import (
	"context"
	"log"

	"github.com/davhofer/botsky/pkg/botsky"
)

func postToBluesky(post, reply, link string) {
	handle, appkey, err := botsky.GetEnvCredentials()

	if err != nil {
		log.Panic(err)
	}

	ctx := context.Background()

	// Set up a client
	client, err := botsky.NewClient(ctx, handle, appkey)

	if err != nil {
		log.Panic(err)
	}

	err = client.Authenticate(ctx)

	if err != nil {
		log.Panic(err)
	}

	log.Println("Bluesky authentication successful")

	pb := botsky.NewPostBuilder(post).AddLanguage("en-US")

	if link != "" {
		pb = pb.AddEmbedLink(link)
	}

	cid, uri, err := client.Post(ctx, pb)
	if err != nil {
		log.Println("Error:", err)
	} else {
		log.Println("Posted:", cid, uri)
	}

	if reply != "" {
		rpb := botsky.NewPostBuilder(reply).AddLanguage("en-US").ReplyTo(uri)

		rcid, ruri, rerr := client.Post(ctx, rpb)
		if rerr != nil {
			log.Println("Reply error:", rerr)
		} else {
			log.Println("Reply posted:", rcid, ruri)
		}
	}
}

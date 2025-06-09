package main

import (
	"context"
	"fmt"

	"github.com/davhofer/botsky/pkg/botsky"
)

func postToBluesky(post, reply, link string) {
	handle, appkey, err := botsky.GetEnvCredentials()

	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	// Set up a client
	client, err := botsky.NewClient(ctx, handle, appkey)

	if err != nil {
		panic(err)
	}

	err = client.Authenticate(ctx)

	if err != nil {
		panic(err)
	}

	fmt.Println("Bluesky authentication successful")

	pb := botsky.NewPostBuilder(post)

	if link != "" {
		pb = pb.AddEmbedLink(link)
	}

	cid, uri, err := client.Post(ctx, pb)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Posted:", cid, uri)
	}

	if reply != "" {
		rpb := botsky.NewPostBuilder(reply).ReplyTo(uri)

		rcid, ruri, rerr := client.Post(ctx, rpb)
		if rerr != nil {
			fmt.Println("Reply error:", rerr)
		} else {
			fmt.Println("Reply posted:", rcid, ruri)
		}
	}
}

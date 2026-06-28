package meta

import (
	"encoding/json"
	"fmt"
	"log"
	"maps"
	"net/http"
	"net/url"
)

type Service int

const (
	INSTAGRAM Service = iota
	THREADS
)

const (
	INSTAGRAM_ENDPOINT = "https://graph.instagram.com/v23.0/%s"
	THREADS_ENDPOINT   = "https://graph.threads.net/v1.0/%s"
)

type Response struct {
	Id    string
	Error struct {
		Message string
	}

	// Used by refresh_access_token endpoint
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func MakePOSTRequest(service Service, accessToken string, endpoint string, params url.Values) (response Response) {
	var serviceName string

	switch service {
	case INSTAGRAM:
		endpoint = fmt.Sprintf(INSTAGRAM_ENDPOINT, endpoint)
		serviceName = "Instagram"
	case THREADS:
		endpoint = fmt.Sprintf(THREADS_ENDPOINT, endpoint)
		serviceName = "Threads"
	}

	maps.Copy(params, url.Values{"access_token": {accessToken}})

	resp, err := http.PostForm(endpoint, params)

	if err != nil {
		log.Panic(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("%s returned \"%s\" for request to %s", serviceName, resp.Status, resp.Request.URL)
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Panic(err)
	}

	return response
}

func RefreshToken(service Service, accessToken string) (response Response) {
	var endpoint string
	var grantType string

	switch service {
	case INSTAGRAM:
		endpoint = INSTAGRAM_ENDPOINT
		grantType = "ig_refresh_token"
	case THREADS:
		endpoint = THREADS_ENDPOINT
		grantType = "th_refresh_token"
	}

	url, _ := url.Parse(fmt.Sprintf(endpoint, "refresh_access_token"))
	query := url.Query()

	query.Add("grant_type", grantType)
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

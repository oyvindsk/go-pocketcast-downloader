package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Scope    string `json:"scope"`
}

type loginRes struct {
	Token string `json:"token"`
	UUID  string `json:"uuid"`
}

func login(email, pass string) (string, error) {

	url := "https://api.pocketcasts.com/user/login"

	req := loginReq{
		Email:    email,
		Password: pass,
		Scope:    "webplayer",
	}

	// log.Printf("login: req: %+v", req)

	reqJson, err := json.Marshal(&req)
	if err != nil {
		return "", fmt.Errorf("login: %s", err)
	}

	res, err := http.Post(url, "application/json", bytes.NewReader(reqJson))
	if err != nil {
		return "", fmt.Errorf("login: %s", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("login: http status: %s", res.Status)
	}

	dec := json.NewDecoder(res.Body)

	var resJson loginRes

	err = dec.Decode(&resJson)
	if err != nil {
		return "", fmt.Errorf("login: %s", err)
	}

	// log.Printf("http res:\n%#v\n", resJson)

	return resJson.Token, nil
}

type resultStarred struct {
	Total    int `json:"total"`
	Episodes []struct {
		UUID          string    `json:"uuid"`
		URL           string    `json:"url"`
		Published     time.Time `json:"published"`
		Duration      int       `json:"duration"`
		FileType      string    `json:"fileType"`
		Title         string    `json:"title"`
		Size          string    `json:"size"`
		PlayingStatus int       `json:"playingStatus"`
		PlayedUpTo    int       `json:"playedUpTo"`
		Starred       bool      `json:"starred"`
		PodcastUUID   string    `json:"podcastUuid"`
		PodcastTitle  string    `json:"podcastTitle"`
		EpisodeType   string    `json:"episodeType"`
		EpisodeSeason int       `json:"episodeSeason"`
		EpisodeNumber int       `json:"episodeNumber"`
		IsDeleted     bool      `json:"isDeleted"`
	} `json:"episodes"`
}

func getStarred(authToken string) (*resultStarred, error) {

	url := "https://api.pocketcasts.com/user/starred"

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, fmt.Errorf("getStarred: %s", err)
	}

	req.Header.Add("Authorization", "Bearer "+authToken)
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("getStarred: %s", err)
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("getStarred: http status: %s", res.Status)
	}

	defer res.Body.Close()

	dec := json.NewDecoder(res.Body)

	var starred resultStarred

	err = dec.Decode(&starred)
	if err != nil {
		return nil, fmt.Errorf("getStarred: %s", err)
	}

	return &starred, nil

}

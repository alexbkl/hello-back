package oauth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type GithubUser struct {
	Name  string
	Email string
	Photo string
}

func GetGithubUser(token string) (*GithubUser, error) {
	req, err := http.NewRequest("GET", github_api_url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	client := http.Client{
		Timeout: time.Second * 30,
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, errors.New("could not retrieve user")
	}

	var resBody bytes.Buffer
	_, err = io.Copy(&resBody, res.Body)
	if err != nil {
		return nil, err
	}

	var GithubUserRes map[string]interface{}

	if err := json.Unmarshal(resBody.Bytes(), &GithubUserRes); err != nil {
		return nil, err
	}

	userBody := &GithubUser{
		Name:  GithubUserRes["login"].(string),
		Email: GithubUserRes["email"].(string),
		Photo: GithubUserRes["avatar_url"].(string),
	}

	return userBody, nil
}

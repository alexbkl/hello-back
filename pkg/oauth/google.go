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

type GoogleUser struct {
	Sub           string `json:"sub"`
	Name          string `json:"name" binding:"required"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Email         string `json:"email" binding:"required"`
	EmailVerified bool   `json:"email_verified"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

func GetGoogleUser(token string) (*GoogleUser, error) {
	req, err := http.NewRequest("GET", google_api_url, nil)
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

	var GoogleUserRes map[string]interface{}

	if err := json.Unmarshal(resBody.Bytes(), &GoogleUserRes); err != nil {
		return nil, err
	}

	userBody := &GoogleUser{
		Sub:           GoogleUserRes["sub"].(string),
		Name:          GoogleUserRes["name"].(string),
		Email:         GoogleUserRes["email"].(string),
		GivenName:     GoogleUserRes["given_name"].(string),
		FamilyName:    GoogleUserRes["family_name"].(string),
		EmailVerified: GoogleUserRes["email_verified"].(bool),
		Picture:       GoogleUserRes["picture"].(string),
		Locale:        GoogleUserRes["locale"].(string),
	}

	return userBody, nil
}

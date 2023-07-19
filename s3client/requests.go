package s3client

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

func SendRequest(method string, url string, headers map[string]string, body string) (string, error) {
	// Create a new request using http
	req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return "", err
	}

	// Add headers to the request
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Send the request via a client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	// Defer the closing of the body
	defer resp.Body.Close()

	// Read the content into a byte array
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Convert the byte array to a string
	bodyString := string(bodyBytes)

	return bodyString, nil
}

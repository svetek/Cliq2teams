package msteams

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
)

func DownloadUrlToBase64(url string) (string, error) {
	// Send GET request
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	// Check if response status is OK
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected response status: %s", response.Status)
	}

	// Read the response body
	fileBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(fileBytes), nil
}

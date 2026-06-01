package api

import (
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"os"
)

func FetchVouchersJSON[T any](url string, target *T) error {
	apiKey := os.Getenv("API_KEY")
	apiSecret := os.Getenv("API_SECRET")

	// Create a new GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	// Add Basic Auth header
	req.SetBasicAuth(apiKey, apiSecret)

	// Send the request
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Read response body using io.ReadAll
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// fmt.Println("Status:", resp.Status)
	// fmt.Println("Response:", string(body))
	return json.Unmarshal(body, target)
}

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"twitter/pkg/routes"
)

const (
	TOPIC = "services"
)

func TestGetUserData(t *testing.T) {

	resp, err := postTweet("http://localhost:8080", "userTest", "messageTest")
	if err != nil {
		t.Fatalf("postTweet returned error: %v", err)
	}
	defer resp.Body.Close()

	// status code is 200 OK
	if resp.StatusCode != http.StatusOK {
		t.Errorf("unexpected status code: got %d, want %d", resp.StatusCode, http.StatusOK)
	}

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf(" reading response")
	}

	var actual routes.Response
	if err := json.Unmarshal(responseData, (&actual)); err != nil {
		t.Errorf("decoding json")
	}

	expected := "Nuevo twitter userTest messageTest"

	if actual.Message != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", actual, expected)
	} else {
		t.Log("Test ok")
	}
}

func postTweet(baseURL, userName string, message string) (*http.Response, error) {

	url := fmt.Sprintf("%s/%s/%s", baseURL, userName, message)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	client := &http.Client{}
	return client.Do(req)
}

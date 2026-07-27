package apexapi

import (
	// "os"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

var client = &http.Client{Timeout: 30 * time.Second}
var mapRotationURL = "https://api.mozambiquehe.re/maprotation"

// SetClient overrides the HTTP client used for Apex API requests.
func SetClient(httpClient *http.Client) {
	if httpClient == nil {
		client = &http.Client{Timeout: 30 * time.Second}
		return
	}

	client = httpClient
}

// SetBaseURL overrides the map rotation endpoint used for Apex API requests.
// An empty url resets it to the default mozambiquehe.re endpoint. Intended
// for pointing at a mock server in tests/local runs.
func SetBaseURL(url string) {
	if url == "" {
		mapRotationURL = "https://api.mozambiquehe.re/maprotation"
		return
	}

	mapRotationURL = url
}

// Map is a structure containing information about a scheduled map
type Map struct {
	Start int64  `json:"start"`
	End   int64  `json:"end"`
	Map   string `json:"map"`
	Code  string `json:"code"`
	Asset string `json:"asset"`
}

// Maps is a structure containing information about current and next Map structures
type Maps struct {
	Current Map `json:"current"`
	Next    Map `json:"next"`
}

// Modes is a structure containing information about game modes
type Modes struct {
	Ranked Maps `json:"ranked"`
	Ltm    Maps `json:"ltm"`
	Pub    Maps `json:"battle_royale"`
}

// Update method to get information from mozambiquehe.re api
func (v *Modes) Update(apiKey string) error {
	u, err := url.Parse(mapRotationURL)
	if err != nil {
		return err
	}
	q := u.Query()
	q.Set("auth", apiKey)
	q.Set("version", "2")
	u.RawQuery = q.Encode()

	resp, err := client.Get(u.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &v)
	if err != nil {
		return fmt.Errorf("%s. body: %s", err, body)
	}
	return nil
}

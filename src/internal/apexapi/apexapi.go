package apexapi

import (
	// "os"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

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

// Update method to get information from mozambiquehe.re api
func (m *Maps) Update(apiKey string) error {
	// apiKey := os.Getenv("APEXLEGENDS_STATUS_API_KEY")
	url := "https://api.mozambiquehe.re/maprotation?auth=" + apiKey

	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &m)
	if err != nil {
		return fmt.Errorf("%s. body: %s", err, body)
	}
	return nil
}

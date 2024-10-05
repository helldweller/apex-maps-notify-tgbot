package apexapi_v2

import (
	// "os"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
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

type V2 struct {
	Ranked Maps `json:"ranked"`
	Ltm    Maps `json:"ltm"`
	Pub    Maps `json:"battle_royale"`
}

// Update method to get information from mozambiquehe.re api
func (v *V2) Update(apiKey string) error {
	u, err := url.Parse("https://api.mozambiquehe.re/maprotation")
	if err != nil {
		return err
	}
	q := u.Query()
	q.Set("auth", apiKey)
	q.Set("version", "2")
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return err
	}

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

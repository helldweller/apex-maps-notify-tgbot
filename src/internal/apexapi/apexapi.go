package apexapi

import (
    // "os"
    "fmt"
    "errors"
    "net/http"
    "io/ioutil"
    "encoding/json"    
)

type Map struct {
    Start       int64   `json:'start'`
    End         int64   `json:'end'`
    Map         string  `json:'map'`
    Code        string  `json:'code'`
    Asset       string  `json:'asset'`
}

type Maps struct {
    Current     Map     `json:'current'`
    Next        Map     `json:'next'`
}

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
        return errors.New(fmt.Sprintf("%s. body: %s", err, body))
    }
    return nil
}

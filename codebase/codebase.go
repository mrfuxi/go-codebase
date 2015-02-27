package codebase

import (
    "encoding/xml"
    "log"
    "net/http"
    "net/url"
    "path"

    "github.com/google/go-querystring/query"
)

type CodeBaseAPI struct {
    Project      string
    userNameAuth string
    apiKey       string
    apiScheme    string
    apiHost      string
}

type baseQueryOptions struct {
    Query string `url:"query,omitempty"`
    Page  uint   `url:"page,omitempty"` // Pages are 1 based
}

// Build new Codebase client to access specific project.
// Project can be changes later on
func NewCodeBaseClient(username, apiKey, project string) *CodeBaseAPI {
    client := &CodeBaseAPI{
        apiScheme:    "http",
        apiHost:      "api3.codebasehq.com",
        userNameAuth: username,
        apiKey:       apiKey,
        Project:      project,
    }

    return client
}

func (c *CodeBaseAPI) fetchFromCodebase(endpoint string, objects interface{}, queryOpts interface{}) (err error) {
    client := &http.Client{}

    query_value, _ := query.Values(queryOpts)
    api_url := url.URL{
        Scheme:   c.apiScheme,
        Host:     c.apiHost,
        Path:     path.Join(c.Project, endpoint),
        RawQuery: query_value.Encode(),
    }

    req, _ := http.NewRequest("GET", api_url.String(), nil)
    req.SetBasicAuth(c.userNameAuth, c.apiKey)
    req.Header.Set("Accept", "application/xml")

    resp, e := client.Do(req)
    if e != nil {
        log.Println("Err: ", e.Error())
        return e
    }

    decoder := xml.NewDecoder(resp.Body)
    err = decoder.Decode(&objects)
    return
}

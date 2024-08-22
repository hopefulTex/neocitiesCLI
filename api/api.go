package api

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Config struct {
	IsDefault   bool   `json:"is_default"`
	Domain      string `json:"domain"`
	IsSubdomain bool   `json:"is_subdomain"`
	APIKey      string `json:"api_key"`
}

type Connection struct {
	client *http.Client
	// domain  string
	api_key string
	headers map[string]string
}

const (
	GET  string = "GET"
	POST string = "POST"
)

const baseURL = "https://neocities.org/api/"

func NewConnection(cfg Config) *Connection {

	conn := &Connection{
		client: &http.Client{},
		// domain:  cfg.Domain,
		api_key: cfg.APIKey,
		headers: map[string]string{
			"Authorization": "Bearer " + cfg.APIKey,
		},
	}
	return conn
}

func (c *Connection) SetAPIKey(key string) {
	c.api_key = key
}

func (c *Connection) Request(method, endpoint string, params []string, body io.Reader) (Response, error) {
	var response Response

	uri := baseURL + endpoint

	if len(params) > 0 {
		_, err := url.ParseRequestURI(params[0])
		if err == nil {
			uri = params[0]
		} else {
			uri = uri + "?"
			for _, param := range params {
				uri = uri + param + "&"
			}
			uri = strings.TrimSuffix(uri, "&")
		}
	}

	request, err := http.NewRequest(string(method), uri, body)
	if err != nil {
		return response, err
	}
	for k, v := range c.headers {
		if v != "" {
			request.Header.Set(k, v)
		}
	}

	reply, err := c.client.Do(request)
	if err != nil {
		return response, err
	}
	defer reply.Body.Close()

	err = json.NewDecoder(reply.Body).Decode(&response)
	if err != nil {
		return response, err
	}
	return response, nil
}

type Response struct {
	Response string     `json:"response"`
	Info     SiteInfo   `json:"info"`
	Files    []ListItem `json:"files"`
	APIKey   string     `json:"api_key"`
}

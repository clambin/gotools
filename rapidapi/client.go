package rapidapi

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
)

// Client represents a RapidAPI client
type Client struct {
	Client   *http.Client
	HostName string
	APIKey   string
}

// New RapidAPI client
// Currently don't need this as we embed the struct anonymously in the API Client
// func New(hostName, apiKey string) *Client {
//	return newWithHTTPClient(hostName, apiKey, &http.Client{})
//}

// NewWithHTTPClient creates a new RapidAPI client with a specified http.Client
// Used to stub server calls during unit tests
func NewWithHTTPClient(client *http.Client, hostName, apiKey string) *Client {
	return &Client{Client: client, HostName: hostName, APIKey: apiKey}
}

// Call an endpoint on the API
func (client *Client) Call(endpoint string) (body []byte, err error) {
	return client.CallWithContext(context.Background(), endpoint)
}

// CallWithContext calls an endpoint on the API with a provided context
func (client *Client) CallWithContext(ctx context.Context, endpoint string) (body []byte, err error) {
	url := "https://" + client.HostName + endpoint
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	req.Header.Add("x-rapidapi-key", client.APIKey)
	req.Header.Add("x-rapidapi-host", client.HostName)

	var resp *http.Response
	resp, err = client.Client.Do(req)

	if err == nil {
		if resp.StatusCode == 200 {
			body, err = ioutil.ReadAll(resp.Body)
		} else {
			err = errors.New(resp.Status)
		}
		_ = resp.Body.Close()
	}
	return
}

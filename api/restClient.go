package api

import (
	"io"
	"net/http"

	"github.com/streame-gg/go-discord-wrapper/types/common"
)

type RestClient struct {
	BaseURL string
	token   string
	Version common.APIVersion

	httpClient *http.Client
}

type RestClientOption func(*RestClient)

func WithBaseURL(baseURL string) RestClientOption {
	return func(c *RestClient) {
		c.BaseURL = baseURL
	}
}

func WithApiVersion(version common.APIVersion) RestClientOption {
	return func(c *RestClient) {
		c.Version = version
	}
}

func WithHttpClient(client *http.Client) RestClientOption {
	return func(c *RestClient) {
		c.httpClient = client
	}
}

func NewRestClient(token string, options ...RestClientOption) *RestClient {
	c := &RestClient{
		BaseURL:    "https://discord.com/api",
		token:      token,
		Version:    common.APIVersion10,
		httpClient: http.DefaultClient,
	}

	for _, option := range options {
		option(c)
	}

	return c
}

func (c *RestClient) buildURL() string {
	return c.BaseURL + "/v" + c.Version.ToString()
}

func (c *RestClient) generateRequest(method, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, c.buildURL()+path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bot "+c.token)
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func (c *RestClient) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return nil, nil
}

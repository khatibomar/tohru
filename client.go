package tohru

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	BaseAPI = "https://anslayer.com"
)

type errorRes struct {
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

type service struct {
	client *TohruClient
}

type TohruClient struct {
	cfg     *Config
	client  *http.Client
	header  http.Header
	service service

	AnimeService   *AnimeService
	EpisodeService *EpisodeService
}

func NewTohruClient(cfg *Config) *TohruClient {
	client := http.Client{}

	header := http.Header{}
	header.Set("Content-Type", "application/json")
	header.Set("Accept", "*/*")
	header.Set("Client-Id", cfg.clientID)
	header.Set("Client-Secret", cfg.clientSecret)

	tohru := &TohruClient{
		client: &client,
		header: header,
		cfg:    cfg,
	}

	tohru.service.client = tohru

	tohru.AnimeService = (*AnimeService)(&tohru.service)
	tohru.EpisodeService = (*EpisodeService)(&tohru.service)

	return tohru
}

func (c *TohruClient) request(ctx context.Context, method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	if body != nil {
		c.header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	req.Header = c.header

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != 200 {
		_, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		var errRes errorRes
		err = json.NewDecoder(resp.Body).Decode(&errRes)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		return nil, fmt.Errorf("%s : %s", errRes.Title, errRes.Detail)
	}
	return resp, nil
}

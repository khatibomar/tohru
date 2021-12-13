package angoslayer

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
	client *AngoClient
}

type ResponseType interface {
	GetResult() interface{}
}

type AngoClient struct {
	cfg     *Config
	client  *http.Client
	header  http.Header
	service service

	Anime *AnimeService
}

func NewAngoClient(cfg *Config) *AngoClient {
	client := http.Client{}

	header := http.Header{}
	header.Set("Content-Type", "application/json")
	header.Set("Accept", "*/*")
	header.Set("Client-Id", cfg.clientID)
	header.Set("Client-Secret", cfg.clientSecret)

	ango := &AngoClient{
		client: &client,
		header: header,
	}

	ango.service.client = ango

	ango.Anime = (*AnimeService)(&ango.service)

	return ango
}

func (c *AngoClient) Request(ctx context.Context, method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header = c.header

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != 200 {
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

func (c *AngoClient) RequestAndDecode(ctx context.Context, method, url string, body io.Reader, rt ResponseType) error {
	resp, err := c.Request(ctx, method, url, body)
	if err != nil {
		return err
	}

	err = json.NewDecoder(resp.Body).Decode(rt)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	return err
}

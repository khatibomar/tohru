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

type ResponseType interface {
	GetResult() string
}

type AngoClient struct {
	client *http.Client
	header http.Header
}

func NewAngoClient() *AngoClient {
	client := http.Client{}

	header := http.Header{}
	header.Set("Content-Type", "application/json")

	ango := &AngoClient{
		client: &client,
		header: header,
	}

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
		return nil, fmt.Errorf("non-200 status code -> %s", err)
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

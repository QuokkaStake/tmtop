package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

type Client struct {
	Logger zerolog.Logger
	Host   string
}

func NewClient(logger zerolog.Logger, invoker, host string) *Client {
	return &Client{
		Logger: logger.With().
			Str("component", "http").
			Str("invoker", invoker).
			Logger(),
		Host: host,
	}
}

func (c *Client) GetInternal(relativeURL string) (io.ReadCloser, error) {
	client := &http.Client{Timeout: 300 * time.Second}
	start := time.Now()

	fullURL := fmt.Sprintf("%s%s", c.Host, relativeURL)

	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "tmtop")

	c.Logger.Debug().Str("url", fullURL).Msg("Doing a query...")

	res, err := client.Do(req)
	if err != nil {
		c.Logger.Warn().Str("url", fullURL).Err(err).Msg("Query failed")
		return nil, err
	}

	c.Logger.Debug().Str("url", fullURL).Dur("duration", time.Since(start)).Msg("Query is finished")

	return res.Body, nil
}

func (c *Client) Get(relativeURL string, target interface{}) error {
	body, err := c.GetInternal(relativeURL)
	if err != nil {
		return err
	}

	if err := json.NewDecoder(body).Decode(target); err != nil {
		return err
	}

	return body.Close()
}

func (c *Client) GetPlain(relativeURL string) ([]byte, error) {
	body, err := c.GetInternal(relativeURL)

	if err != nil {
		return nil, err
	}

	bytes, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}

	if err := body.Close(); err != nil {
		return nil, err
	}

	return bytes, nil
}

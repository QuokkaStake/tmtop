package http

import (
	"encoding/json"
	"fmt"
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

func (c *Client) Get(relativeURL string, target interface{}) error {
	client := &http.Client{Timeout: 300 * time.Second}
	start := time.Now()

	fullURL := fmt.Sprintf("%s%s", c.Host, relativeURL)

	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", "tmtop")

	c.Logger.Debug().Str("url", fullURL).Msg("Doing a query...")

	res, err := client.Do(req)
	if err != nil {
		c.Logger.Warn().Str("url", fullURL).Err(err).Msg("Query failed")
		return err
	}
	defer res.Body.Close()

	c.Logger.Debug().Str("url", fullURL).Dur("duration", time.Since(start)).Msg("Query is finished")

	return json.NewDecoder(res.Body).Decode(target)
}

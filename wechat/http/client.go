package http

import (
	"io"
	"net/http"
)

type Client struct {
	Request *http.Request
	Client  *http.Client
}

func NewClient(method string, urlStr string, body io.Reader) (client *Client, err error) {
	req, err := http.NewRequest(method, urlStr, body)

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.86 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Connection", "keep-alive")

	return &Client{
		Request: req,
		Client:  &http.Client{},
	}, err
}

func (c *Client) SetHeader(key string, value string) *Client {
	c.Request.Header.Set(key, value)
	return c
}

func (c *Client) SetUserAgent(ua string) *Client {
	if ua == "" {
		c.Request.Header.Del("User-Agent")
	} else {
		c.Request.Header.Set("User-Agent", ua)
	}
	return c
}

func (c *Client) SetCookies(cookies []*http.Cookie) *Client {
	for _, v := range cookies {
		c.Request.AddCookie(v)
	}
	return c
}

func (c *Client) Do() (*http.Response, error) {
	resp, err := c.Client.Do(c.Request)
	return resp, err
}

func init() {
}

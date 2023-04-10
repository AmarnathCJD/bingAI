package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	sessionName string
	initFlag    bool
	defaultConv *bingConv
	convs       map[string]*bingConv
	cookies     []*http.Cookie
	wss         *soc
}

type Config struct {
	AccessToken      string
	Cookies          string
	SessionName      string
	BaseURL          string
	DisableCacheTODO bool
}

func NewClient(config *Config) *Client {
	if config.SessionName == "" {
		config.SessionName = genUUID4()
	}
	if config.BaseURL == "" {
		config.BaseURL = "wss://sydney.bing.com/sydney/ChatHub"
	}
	client := &Client{
		sessionName: config.SessionName,
		convs:       map[string]*bingConv{},
		wss: &soc{
			url: config.BaseURL,
		},
	}

	if config.AccessToken != "" {
		client.cookies = append(client.cookies, accessTokenToCookie(config.AccessToken))
	} else if config.Cookies != "" {
		client.cookies = resolveCookies(config.Cookies)
	}
	return client
}

var AUTH_HEADERS = map[string]string{
	"authority":                 "edgeservices.bing.com",
	"accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
	"accept-language":           "en-US,en;q=0.9",
	"cache-control":             "max-age=0",
	"upgrade-insecure-requests": "1",
	"user-agent":                "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36 Edg/110.0.1587.69",
	"x-edge-shopping-flag":      "1",
	"x-forwarded-for":           fakeIp,
}

func (c *Client) initConv(ctx context.Context) (*bingConv, error) {
	req, _ := http.NewRequestWithContext(ctx, "GET", "https://edgeservices.bing.com/edgesvc/turing/conversation/create", nil)
	for k, v := range AUTH_HEADERS {
		req.Header.Set(k, v)
	}
	for _, cookie := range c.cookies {
		req.AddCookie(cookie)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	RETRY_COUNT := 0
	for resp.StatusCode != 200 && RETRY_COUNT < 3 {
		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		RETRY_COUNT++
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to init auth after %d retries, status code: %d", RETRY_COUNT, resp.StatusCode)
	}
	var bingResp bingConv
	if err := json.NewDecoder(resp.Body).Decode(&bingResp); err != nil {
		return nil, err
	}
	return &bingResp, nil
}

func (c *Client) Start(ctx context.Context) error {
	if err := c.wss.Connect(ctx); err != nil {
		return err
	}
	if err := c.wss.initialHandshake(); err != nil {
		return err
	}
	if conv, err := c.initConv(ctx); err != nil {
		return err
	} else {
		c.convs[conv.ConversationID] = conv
		c.defaultConv = conv
	}
	c.initFlag = true
	return nil
}

func (c *Client) ExportAuthToken() string {
	return cookieToAccessToken(c.cookies[0])
}

func (c *Client) NewConversation(ctx context.Context) (*bingConv, error) {
	if !c.initFlag {
		return nil, fmt.Errorf("client not initialized")
	}
	conv, err := c.initConv(ctx)
	if err != nil {
		return nil, err
	}
	c.convs[conv.ConversationID] = conv
	return conv, nil
}

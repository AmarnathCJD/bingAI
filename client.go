package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	mrand "math/rand"
	"net/http"
	"strings"
)

type Client struct {
	sessionName string
	initFlag    bool
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

func (c *Client) initAuth(ctx context.Context) error {
	req, _ := http.NewRequestWithContext(ctx, "GET", "https://edgeservices.bing.com/edgesvc/turing/conversation/create", nil)
	for k, v := range AUTH_HEADERS {
		req.Header.Set(k, v)
	}
	for _, cookie := range c.cookies {
		req.AddCookie(cookie)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		req, _ := http.NewRequestWithContext(ctx, "GET", "https://edge.churchless.tech/edgesvc/turing/conversation/create", nil)
		for k, v := range AUTH_HEADERS {
			req.Header.Set(k, v)
		}
		for _, cookie := range c.cookies {
			req.AddCookie(cookie)
		}
		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			return fmt.Errorf("failed to create conversation: %s", resp.Status)
		}
	}
	var bingResp bingConv
	if err := json.NewDecoder(resp.Body).Decode(&bingResp); err != nil {
		return err
	}
	c.convs["default"] = &bingResp
	return nil
}

func (c *Client) Start(ctx context.Context) error {
	if err := c.wss.Connect(ctx); err != nil {
		return err
	}
	if err := c.wss.initialHandshake(); err != nil {
		return err
	}
	if err := c.initAuth(ctx); err != nil {
		return err
	}
	c.initFlag = true
	return nil
}

func (c *Client) ExportAuthToken() string {
	return cookieToAccessToken(c.cookies[0])
}

// ------------- Helper functions -------------

func genRandomIP() string {
	return fmt.Sprintf("%d.%d.%d.%d", mrand.Intn(255), mrand.Intn(255), mrand.Intn(255), mrand.Intn(255))
}

func resolveCookies(c string) []*http.Cookie {
	cookies := []*http.Cookie{}
	for _, cookie := range strings.Split(c, ";") {
		cookie = strings.TrimSpace(cookie)
		split := strings.Split(cookie, "=")
		cookies = append(cookies, &http.Cookie{
			Name:  split[0],
			Value: split[1],
		})
	}
	return cookies
}

func genRandHex(n int) string {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}

func genUUID4() string {
	return fmt.Sprintf("%s-%s-4%s-%s-%s", genRandHex(8), genRandHex(4), genRandHex(3), genRandHex(4), genRandHex(12))
}

func cookieToAccessToken(cookie *http.Cookie) string {
	cookieValue := cookie.Value
	urlsafe_b64encode := func(s string) string {
		return strings.TrimRight(base64.URLEncoding.EncodeToString([]byte(s)), "=")
	}
	return urlsafe_b64encode(cookieValue)
}

func accessTokenToCookie(token string) *http.Cookie {
	urlsafe_b64decode := func(s string) string {
		missing_padding := len(s) % 4
		if missing_padding != 0 {
			s += strings.Repeat("=", 4-missing_padding)
		}
		b, _ := base64.URLEncoding.DecodeString(s)
		return string(b)
	}
	return &http.Cookie{
		Name:  "_U",
		Value: urlsafe_b64decode(token),
	}
}

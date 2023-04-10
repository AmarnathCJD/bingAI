package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	mrand "math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type cfg struct {
	apiKey      int32
	apiHash     string
	botToken    string
	accessToken string
	cookies     string
}

func getInt32Env(key string) int32 {
	if v, err := strconv.Atoi(os.Getenv(key)); err == nil {
		return int32(v)
	}
	return 0
}

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
	for _, cookie := range cookies {
		if cookie.Name == "_U" {
			return []*http.Cookie{cookie}
		}
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

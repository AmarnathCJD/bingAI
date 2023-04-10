package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"nhooyr.io/websocket"
)

var fakeIp = genRandomIP() // fake ip for the websocket connection

// SOC_HEADERS are the headers that are sent with the websocket connection
var SOC_HEADERS = map[string]string{
	"accept":                 "application/json",
	"accept-language":        "en-US,en;q=0.9",
	"content-type":           "application/json",
	"x-ms-client-request-id": genUUID4(),
	"x-ms-useragent":         "azsdk-js-api-client-factory/1.0.0-beta.1 core-rest-pipeline/1.10.0 OS/Win32",
	"Referer":                "https://www.bing.com/search?q=Bing+AI&showconv=1&FORM=hpcodx",
	"Referrer-Policy":        "origin-when-cross-origin",
	"x-forwarded-for":        fakeIp,
}

// soc is a wrapper around websocket.Conn to make it easier to use
type soc struct {
	sync.Mutex
	ctx        context.Context
	url        string
	connection *websocket.Conn
}

func (w *soc) IsActive() bool {
	return w.connection != nil
}

func (w *soc) Connect(ctx context.Context) error {
	var headers = map[string][]string{}
	for k, v := range SOC_HEADERS {
		headers[k] = []string{v}
	}

	conn, _, err := websocket.Dial(ctx, w.url, &websocket.DialOptions{
		HTTPClient: &http.Client{},
		HTTPHeader: map[string][]string{},
	})
	if err != nil {
		return err
	}
	w.connection = conn
	w.ctx = ctx
	return nil
}

func (w *soc) Close() error {
	if w.connection != nil {
		return w.connection.Close(websocket.StatusNormalClosure, "")
	}
	return nil
}

func (w *soc) Send(data []byte) error {
	if w.connection == nil {
		return fmt.Errorf("connection is not active")
	}
	// add "\x1e" to the end of the message
	data = append(data, "\x1e"...)
	return w.connection.Write(w.ctx, websocket.MessageText, data)
}

func (w *soc) Receive() ([]byte, error) {
	if w.connection == nil {
		return nil, fmt.Errorf("connection is not active")
	}
	_, data, err := w.connection.Read(w.ctx)
	if err != nil {
		return nil, err
	}
	// remove "\x1e" from the end of the message
	final_data := data[:len(data)-1]
	if strings.Contains(string(final_data), "\x1e") {
		final_data = []byte(strings.Split(string(final_data), "\x1e")[0])
	}
	return []byte(final_data), nil
}

func (w *soc) initialHandshake() error {
	if w.connection == nil {
		return fmt.Errorf("connection is not active")
	}
	err := w.Send([]byte(`{"protocol":"json","version":1}`))
	if err != nil {
		return err
	}
	_, err = w.Receive()
	return err
}

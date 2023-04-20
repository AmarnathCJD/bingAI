package main

import (
	"context"
	"encoding/json"
	"fmt"
)

// bingConv is a struct to hold the conversation data
type bingConv struct {
	ClientID              string `json:"clientId,omitempty"`
	ConversationID        string `json:"conversationId,omitempty"`
	ConversationSignature string `json:"conversationSignature,omitempty"`
	invocationID          int
}

func (b *bingConv) String() string {
	return fmt.Sprintf("clientId: %s, conversationId: %s, conversationSignature: %s", b.ClientID, b.ConversationID, b.ConversationSignature)
}

func (b *bingConv) MakePayload(prompt string) []byte {
	payload := map[string]interface{}{
		"arguments": []interface{}{
			map[string]interface{}{
				"source": "cib",
				"optionsSets": []string{
					"nlu_direct_response_filter",
					"deepleo",
					"disable_emoji_spoken_text",
					"responsible_ai_policy_235",
					"enablemm",
					"h3relaxedimg",
					"dtappid",
					"cricinfo",
					"cricinfov2",
					"dv3sugg",
				},
				"sliceIds": []string{
					"222dtappid",
					"225cricinfo",
					"224locals0",
				},
				"traceId":          genRandHex(32),
				"isStartOfSession": b.invocationID == 0,
				"message": map[string]interface{}{
					"author":      "user",
					"inputMethod": "Keyboard",
					"text":        prompt,
					"messageType": "Chat",
				},
				"conversationSignature": b.ConversationSignature,
				"participant": map[string]interface{}{
					"id": b.ClientID,
				},
				"conversationId": b.ConversationID,
			},
		},
		"invocationId": fmt.Sprint(b.invocationID),
		"target":       "chat",
		"type":         4,
	}
	json_payload, _ := json.Marshal(payload)
	return json_payload
}

// bingResp is a struct to hold the response data
type bingResp struct {
	Type         int    `json:"type"`
	InvocationID string `json:"invocationId"`
	Item         struct {
		Messages []struct {
			Text          string `json:"text"`
			Author        string `json:"author"`
			AdaptiveCards []struct {
				Type    string `json:"type"`
				Version string `json:"version"`
				Body    []struct {
					Type string `json:"type"`
					Text string `json:"text"`
					Wrap bool   `json:"wrap"`
				} `json:"body"`
			} `json:"adaptiveCards,omitempty"`
		} `json:"messages"`
	} `json:"item"`
}

type ChatResp struct {
	InvocationID   string
	Message        string
	ConversationID string
}

func (b *bingResp) String() string {
	return fmt.Sprintf("type: %d, invocationId: %s, messages: %v", b.Type, b.InvocationID, b.Item.Messages)
}

func (c *Client) Ask(ctx context.Context, prompt string, conversationId ...string) (*ChatResp, error) {
	if !c.initFlag {
		return nil, fmt.Errorf("client not initialized, call Start() first")
	}
	wss := NewSoc("wss://sydney.bing.com/sydney/ChatHub")
	if err := wss.Connect(ctx); err != nil {
		return nil, err
	}
	defer wss.Close()
	if err := wss.Send([]byte(`{"protocol":"json","version":1}`)); err != nil {
		return nil, err
	}
	var conv *bingConv
	var convoId string
	if len(conversationId) == 0 {
		if c.defaultConv.invocationID+1 > 19 {
			c.defaultConv, _ = c.initConv(ctx)
		}
		conv = c.defaultConv
		convoId = c.defaultConv.ConversationID
	} else {
		if c.convs[conversationId[0]].invocationID+1 > 19 {
			newConv, err := c.initConv(ctx)
			if err != nil {
				return nil, err
			}
			c.convs[conversationId[0]] = newConv
		}
		conv = c.convs[conversationId[0]]
		convoId = conversationId[0]
	}
	if conv == nil {
		return nil, fmt.Errorf("conversation not found")
	}
	payload := conv.MakePayload(prompt)
	c.convs[convoId].invocationID++ // increment invocation id
	if err := wss.Send(payload); err != nil {
		return nil, err
	}
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			resp, err := wss.Receive()
			if err != nil {
				return nil, err
			}
			var bingResp bingResp
			if err := json.Unmarshal(resp, &bingResp); err != nil {
				return nil, err
			}
			if bingResp.Type == 2 {
				if len(bingResp.Item.Messages) == 0 {
					return nil, fmt.Errorf("no messages in response")
				}
				return &ChatResp{InvocationID: bingResp.InvocationID, Message: bingResp.Item.Messages[len(bingResp.Item.Messages)-1].Text, ConversationID: convoId}, nil
			}
		}
	}
}

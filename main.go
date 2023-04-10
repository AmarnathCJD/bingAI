package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/amarnathcjd/gogram/telegram"
	dotenv "github.com/joho/godotenv"
)

var botID int64 = 0
var cli *Client

func main() {
	dotenv.Load()

	cfg := &cfg{
		getInt32Env("API_KEY"),
		os.Getenv("API_HASH"),
		os.Getenv("BOT_TOKEN"),
		os.Getenv("ACCESS_TOKEN"),
		os.Getenv("COOKIES"),
	}

	for _, v := range []string{"API_KEY", "API_HASH", "BOT_TOKEN"} {
		if os.Getenv(v) == "" {
			log.Fatalf("Missing %s", v)
		}
	}

	if os.Getenv("ACCESS_TOKEN") == "" && os.Getenv("COOKIES") == "" {
		log.Fatal("Missing ACCESS_TOKEN or COOKIES")
	}

	tg, _ := telegram.NewClient(
		telegram.ClientConfig{
			AppID:   cfg.apiKey,
			AppHash: cfg.apiHash,
		},
	)
	if err := tg.Connect(); err != nil {
		log.Fatal(err)
	}
	if err := tg.LoginBot(cfg.botToken); err != nil {
		log.Fatal(err)
	}

	bot, err := tg.GetMe()
	if err != nil {
		log.Fatal(err)
	}
	botID = bot.ID

	client := NewClient(&Config{
		AccessToken: cfg.accessToken,
		Cookies:     cfg.cookies,
	})

	if err := client.Start(context.Background()); err != nil {
		log.Fatal(err)
	}
	cli = client

	log.Println("BingAI is ready to answer your questions!")

	tg.AddMessageHandler(telegram.OnNewMessage, bingAIHandler)
	tg.AddMessageHandler("/start", startHandler)
	tg.Idle()
}

func triggerDetector(m *telegram.NewMessage) (string, bool) {
	if strings.HasPrefix(m.Text(), "/bingai") || strings.HasPrefix(m.Text(), "/ai") {
		return m.Args(), true
	}
	if m.IsPrivate() && !strings.HasPrefix(m.Text(), "/") {
		return m.Text(), true
	} else if m.IsGroup() && m.IsReply() {
		r, err := m.GetReplyMessage()
		if err != nil {
			return "", false
		}
		if r.SenderID() == botID {
			return m.Text(), true
		}
	}
	return "", false
}

func bingAIHandler(m *telegram.NewMessage) error {
	if query, ok := triggerDetector(m); ok {
		m.SendAction("typing")
		if resp, err := cli.Ask(context.Background(), query); err == nil {
			if _, err := m.Reply(resp.Message); err != nil {
				log.Println(err)
			}
		} else {
			log.Println(err)
		}
	}
	return nil
}

var startMessage = `Hi, I'm BingAI, a bot that answers your questions using Bing Chatbot.`

func startHandler(m *telegram.NewMessage) error {
	if !m.IsPrivate() {
		m.Reply(startMessage)
	}
	b := telegram.Button{}
	m.Reply(startMessage, telegram.SendOptions{
		ReplyMarkup: b.Keyboard(
			b.Row(
				b.URL("Contact Developer", "https://t.me/roseloverx"),
			),
			b.Row(
				b.URL("Bing Chatbot", "https://www.bing.com/chat"),
				b.URL("Source Code", "https://github.com/amarnathcjd/bingai"),
			),
		),
	})
	return nil
}

# BingAI

BingAI is a simple AI Embeded with Telegram bot. It can answer your questions and do some simple tasks.
based on Bing Sydney AI Lab's [Bing](https://www.bing.com/chat) chatbot

## Usage
```bash
git clone github.com/amarnathcjd/bingai
cd bingai
go run main.go
```

## Requirements
- Microsoft account with preview access to Bing AI
- Telegram bot token, get it from [@BotFather](https://telegram.me/BotFather)
- Telegram AppID and AppHash, get it from [my.telegram.org](https://my.telegram.org)
- Go 1.8 or above

## Get Cookies
- Go to [Bing](https://www.bing.com/chat) on Latest Edge browser (or set Edge UA in Chrome or Firefox)
- Login with your Microsoft account with preview access to Bing AI
- Open Developer Tools (F12) or (Ctrl+Shift+I)
- Go to Application tab
- Go to Cookies
- Copy the value of `_U` cookie
- Pass it in .env as COOKIES=_U=value

## AccessToken
 Since Cookies is messy, the lib can generate access token from cookies. Just pass the cookies in .env as COOKIES=_U=value for the first time. The lib will generate access token

 call client.ExportAccessToken() to get the access token
 now you can remove cookies and pass the access token in .env as ACCESS_TOKEN=token

## Working Demo
[@SydAIBot](https://telegram.me/SydAIBot)

## Note
- The lib is still in development, so it may not work as expected

## License
[MIT](  )

### Credits
- [Bing](https://www.bing.com/chat)



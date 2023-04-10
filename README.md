#TODO (critical): close TCP conn after reading response to prevent RCT Signal

# BingAI

BingAI is a Telegram bot that utilizes Bing Sydney AI Lab's [Bing](https://www.bing.com/chat) chatbot to answer your questions and perform simple tasks.

## Usage

To use BingAI, follow these steps:

1. Clone the repository: 
`bash
git clone github.com/amarnathcjd/bingai
`

2. Navigate to the project directory: 
`bash
cd bingai
`

3. Run the project:
`bash
go run main.go
`

## Requirements

Before using BingAI, make sure you have the following:

- Microsoft account with preview access to Bing AI
- Telegram bot token, which you can obtain from [@BotFather](https://telegram.me/BotFather)
- Telegram AppID and AppHash, which you can obtain from [my.telegram.org](https://my.telegram.org)
- Go 1.8 or later

## Get Cookies

To obtain the necessary cookies for BingAI, follow these steps:

1. Go to [Bing](https://www.bing.com/chat) using the latest Edge browser (or set Edge UA in Chrome or Firefox)
2. Log in to your Microsoft account with preview access to Bing AI
3. Open Developer Tools (F12) or (Ctrl+Shift+I)
4. Go to the Application tab
5. Go to Cookies
6. Copy the value of the `_U` cookie
7. Pass it in `.env` as `COOKIES=_U=value`

## AccessToken

Since cookies can be messy, the library can generate an access token from cookies. Simply pass the cookies in `.env` as `COOKIES=_U=value` for the first time. The library will generate an access token.

To obtain the access token, call `client.ExportAccessToken()`. You can then remove the cookies and pass the access token in `.env` as `ACCESS_TOKEN=token`.

## Working Demo

To see a working demo of BingAI, visit [@SydAIBot](https://telegram.me/SydAIBot).

## Note

BingAI is still in development, so it may not work as expected.

## License

BingAI is licensed under the [MIT License](LICENSE).

### Credits

- [Bing](https://www.bing.com/chat)

package main

import (
	"os"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func main() {
	bot, err := linebot.New(os.Getenv("LINE_CHANNEL_SECRET"), os.Getenv("LINE_ACCESS_TOKEN"))
	if err != nil {
		println(err)
	}
	if bot != nil {
		println("OK")
	}
}

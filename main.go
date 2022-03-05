package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func main() {
	bot, err := linebot.New(os.Getenv("LINE_CHANNEL_SECRET"), os.Getenv("LINE_ACCESS_TOKEN"))
	if err != nil {
		println(err.Error())
		return
	}
	if bot != nil {
		println("OK")
	}
	r := gin.Default()
	r.POST("/line", func(c *gin.Context) {
		events, err := bot.ParseRequest(c.Request)
		if err != nil {
			println(err.Error())
		}
		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				var messages []linebot.SendingMessage
				messages = append(messages, linebot.NewTextMessage("OK"))
				_, err = bot.ReplyMessage(event.ReplyToken, messages...).Do()
				if err != nil {
					println(err.Error())
				}
			}
		}
		c.JSON(200, gin.H{
			"message": "OK",
		})
	})
	r.Run()
}

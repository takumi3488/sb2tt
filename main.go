package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"golang.org/x/exp/utf8string"
)

func main() {
	bot, err := linebot.New(os.Getenv("LINE_CHANNEL_SECRET"), os.Getenv("LINE_ACCESS_TOKEN"))
	if err != nil {
		println(err.Error())
		return
	}
	r := gin.Default()

	// LINE メッセージに対する処理
	r.POST("/line", func(c *gin.Context) {
		events, err := bot.ParseRequest(c.Request)
		if err != nil {
			println(err.Error())
		}
		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					text := strings.TrimSpace(message.Text)
					if strings.HasSuffix(text, "シフト管理アプリ「シフトボード」で作成") {
						ss := parse(text, "", time.Now().Local())
						for _, v := range ss {
							_, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(fmt.Sprintf("%s\n%s-%s", v.title, v.start_at, v.end_at))).Do()
							if err != nil {
								println(err.Error())
							}
						}
					}
					if err != nil {
						println(err.Error())
					}
				}
			}
		}
		c.JSON(200, gin.H{})
	})

	// TimeTree Calender Appをカレンダーに追加した際の通知
	r.POST("/timetree", func(c *gin.Context) {
		res := struct {
			Action       string `json:"action"`
			Installation struct {
				ID          string `json:"id"`
				Application struct {
					ID      string `json:"id"`
					Name    string `json:"name"`
					IconURL string `json:"icon_url"`
				} `json:"application"`
				Scopes    []string  `json:"scopes"`
				UpdatedAt time.Time `json:"updated_at"`
				CreatedAt time.Time `json:"created_at"`
			} `json:"installation"`
		}{}
		c.BindJSON(&res)
		if res.Action == "created" {
			bot.PushMessage(os.Getenv("LINE_ADMIN_ID"), linebot.NewTextMessage("INSTALLATION_ID\n"+fmt.Sprint(res.Installation.ID))).Do()
			c.JSON(200, gin.H{"res": res})
		} else {
			c.JSON(400, gin.H{"text": "Internal Error"})
		}
	})
	r.Run()
}

type Schedule struct {
	title    string
	start_at string
	end_at   string
}

func parse(c string, ttl string, now time.Time) []Schedule {
	rows := strings.Split(c, "\n")
	var shedules []Schedule
	for i := 0; ; i++ {
		t := rows[2*i]
		if t == "" {
			break
		}
		n := rows[2*i+1][2:]
		month := t[0:2]
		year := now.Year()
		month_int, _ := strconv.Atoi(month)
		if int(now.Month())-8 > month_int {
			year++
		}
		date := t[3:5]
		tt := utf8string.NewString(t)
		start_at := tt.Slice(8, 13)
		end_at := tt.Slice(16, 21)
		s := new(Schedule)
		s.title = ttl
		if s.title == "" {
			s.title = n
		}
		dt := fmt.Sprintf("%d-%s-%sT", year, month, date)
		e := ":00.000Z"
		s.start_at = fmt.Sprintf("%s%s%s", dt, start_at, e)
		s.end_at = fmt.Sprintf("%s%s%s", dt, end_at, e)
		shedules = append(shedules, *s)
	}
	return shedules
}

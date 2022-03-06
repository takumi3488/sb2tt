package main

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/takumi3488/sb2tt/model"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"golang.org/x/exp/utf8string"
)

func main() {
	bot, err := linebot.New(os.Getenv("LINE_CHANNEL_SECRET"), os.Getenv("LINE_ACCESS_TOKEN"))
	if err != nil {
		panic(err)
	}
	r := gin.Default()

	// LINE イベントに対する処理
	r.POST("/line", func(c *gin.Context) {
		events, err := bot.ParseRequest(c.Request)
		if err != nil {
			println(err.Error())
		}

		for _, event := range events {
			// メッセージ受信時の処理
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					text := strings.TrimSpace(message.Text)

					if strings.HasSuffix(text, "シフト管理アプリ「シフトボード」で作成") {
						// シフトボードからの共有
						var user model.LineUser
						db, err := model.DbOpen()
						if err != nil {
							println(err.Error())
						}
						db.Where(&model.LineUser{UserId: event.Source.UserID}).First(&user)
						ss, err := parse(text, user.DefaultScheduleTitle, time.Now().Local())
						if err != nil {
							bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(fmt.Sprintln(err))).Do()
						}
						for _, v := range ss {
							_, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(fmt.Sprintf("%s\n%s-%s", v.title, v.start_at, v.end_at))).Do()
							if err != nil {
								println(err.Error())
							}
						}
					} else if r := regexp.MustCompile(`\d+`); r.MatchString(text) {
						// Installation id の設定
						db, err := model.DbOpen()
						if err != nil {
							println(err.Error())
							return
						}
						installation_id, _ := strconv.Atoi(text)
						db.Model(&model.LineUser{}).Where("user_id = ?", event.Source.UserID).Update("installation_id", installation_id)
						bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(fmt.Sprintf("Installation ID set to %d.", installation_id))).Do()
					} else if r := regexp.MustCompile(`[^\n]+`); r.MatchString(text) {
						// デフォルトタイトルの設定
						db, err := model.DbOpen()
						if err != nil {
							println(err.Error())
							return
						}
						db.Model(&model.LineUser{}).Where("user_id = ?", event.Source.UserID).Update("default_schedule_title", text)
						bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(fmt.Sprintf("Default scedule title set to %s.", text))).Do()
					}

					if err != nil {
						println(err.Error())
					}
				}
			}

			// フォロー/ブロック解除時の処理
			if event.Type == linebot.EventTypeFollow {
				userId := event.Source.UserID
				db, err := model.DbOpen()
				if err != nil {
					println(err.Error())
					return
				}
				var user model.LineUser
				db.Where(model.LineUser{UserId: userId}).FirstOrCreate(&user)
				if user.InstallationId == 0 {
					bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("Please send your installation id of TimeTree.")).Do()
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
				ID          int `json:"id"`
				Application struct {
					ID      int    `json:"id"`
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
			bot.PushMessage(os.Getenv("LINE_ADMIN_ID"), linebot.NewTextMessage(fmt.Sprintf("INSTALLATION_ID\n%d", res.Installation.ID))).Do()
			c.JSON(200, gin.H{"res": res})
		} else {
			c.JSON(400, gin.H{"text": "Internal Error"})
		}
	})

	// Migrate
	r.POST("/migrate", func(c *gin.Context) {
		model.Migrate()
		c.JSON(200, gin.H{"text": "migrated"})
	})
	if err := r.Run(); err != nil {
		panic(err)
	}
}

type Schedule struct {
	title    string
	start_at string
	end_at   string
}

func parse(c string, ttl string, now time.Time) ([]Schedule, error) {
	var shedules []Schedule
	var err error
	rows := strings.Split(c, "\n")
	flg := false
	for _, row := range rows {
		if strings.HasPrefix(row, "- ") {
			flg = true
			break
		}
	}
	if !flg && ttl == "" {
		return shedules, errors.New("シフトボードでバイト先の表示をONにするか、デフォルトのシフト名を設定してください。")
	}
	for i := 0; ; i++ {
		j := i
		if flg {
			j = 2 * j
		}
		t := strings.TrimSpace(rows[j])
		if t == "" {
			break
		}
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
		if flg {
			s.title = rows[2*i+1][2:]
		}
		dt := fmt.Sprintf("%d-%s-%sT", year, month, date)
		e := ":00.000Z"
		s.start_at = fmt.Sprintf("%s%s%s", dt, start_at, e)
		s.end_at = fmt.Sprintf("%s%s%s", dt, end_at, e)
		shedules = append(shedules, *s)
	}
	return shedules, err
}

package main

import (
	"bytes"
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/Syfaro/telegram-bot-api"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)


const (
	MSGPATTERN = "работ"
	QSOURCE    = "https://ru.wikiquote.org/wiki/%D0%A2%D1%80%D1%83%D0%B4"
	QUOTE      = `\[[0-9]+\]:*\d*`
)

var (
	bot      *tgbotapi.BotAPI
	botToken string
	baseURL  string
	reg      *regexp.Regexp
	)

func initTelegram() {
	var err error

	bot, err = tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Println(err)
		return
	}

	url := baseURL + bot.Token
	_, err = bot.SetWebhook(tgbotapi.NewWebhook(url))
	if err != nil {
		log.Println(err)
	}
}


func webhookHandler(c *gin.Context) {
	defer c.Request.Body.Close()

	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Println(err)
		return
	}

	var update tgbotapi.Update
	err = json.Unmarshal(data, &update)
	if err != nil {
		log.Println(err)
		return
	}

	match, err := regexp.MatchString(MSGPATTERN, strings.ToLower(update.Message.Text))
	if err != nil {
		log.Println(err)
		return
	}

	if match {
		quote, err := getQuote()
		if err != nil {
			quote = "╭∩╮( ͡° ͜ʖ ͡°)╭∩╮"
		}

		response := gin.H{
			"chat_id": update.Message.Chat.ID,
			"reply_to_message_id": update.Message.MessageID,
			"text": quote,
		}
		bts, err := json.Marshal(response)
		if err != nil {
			log.Println(err)
			return
		}

		resp, err := http.DefaultClient.Post("https://api.telegram.org/bot"+ botToken + "/sendMessage", "application/json", bytes.NewReader(bts))
		if resp.Status != "200" {
			log.Println(resp.Status)
			return
		}
	}
}

func getQuote() (string, error) {
	// Request the HTML page.
	res, err := http.Get(QSOURCE)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return "", errors.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", err
	}

	var bandRes string
	count := 0
	sel := doc.Find("html").Find("body").Find("div").Find("div").Find("div").Find("div").Find("div")
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(sel.Length())
	sel.Each(func(i int, s *goquery.Selection) {
		if count == n && s.Index() == 0 && !strings.Contains(s.Text(), "↑") {
			//fmt.Println(s.Text()))
			res := reg.ReplaceAllString(strings.TrimSpace(s.Text()), "${1}")
			bandRes = res
		}
		count++
	})

	return bandRes, nil
}


func main() {
	reg = regexp.MustCompile(QUOTE)
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	botToken = os.Getenv("TELEGRAM_TOKEN")
	if botToken == "" {
		log.Fatal("$TELEGRAM_TOKEN must be set")
	}

	baseURL = os.Getenv("WEBHOOK_URL")
	if baseURL == "" {
		log.Fatal("$WEBHOOK_URL must be set")
	}

	// gin router
	router := gin.New()
	router.Use(gin.Logger())

	// telegram
	initTelegram()
	router.POST("/" + bot.Token, webhookHandler)

	err := router.Run(":" + port)
	if err != nil {
		log.Println(err)
	}
}
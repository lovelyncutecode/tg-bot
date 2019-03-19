package main

import (
	"bytes"
	"encoding/json"
	"github.com/Syfaro/telegram-bot-api"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
)


const MSGPATTERN = "работ"
var (
	bot      *tgbotapi.BotAPI
	botToken string
	baseURL  string
	)

func initTelegram() {
	var err error

	bot, err = tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Println(err)
		return
	}

	// this perhaps should be conditional on GetWebhookInfo()
	// only set webhook if it is not set properly
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

	match, err := regexp.MatchString(MSGPATTERN, update.Message.Text)
	if err != nil {
		log.Println(err)
		return
	}

	// to monitor changes run: heroku logs --tail
	log.Printf("From: %+v Text: %+v\n", update.Message.From, update.Message.Text)
	if match {
		response := gin.H{
			"chat_id": update.Message.Chat.ID,
			"reply_to_message_id": update.Message.MessageID,
			"text": "╭∩╮( ͡° ͜ʖ ͡°)╭∩╮",
		}
		bts, err := json.Marshal(response)
		if err != nil {
			log.Println(err)
			return
		}
	log.Println(baseURL + botToken + "/sendMessage")
		resp, err := http.DefaultClient.Post(baseURL + botToken + "/sendMessage", "application/json", bytes.NewReader(bts))
		if resp.Status != "200" {
			log.Println(resp.Status)
			return
		}
	}
}

func main() {
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
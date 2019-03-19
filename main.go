package main

import (
	"encoding/json"
	"github.com/Syfaro/telegram-bot-api"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"os"
)

//
//var client *http.Client
//
//func main() {
//	//certManager := autocert.Manager{
//	//	Prompt:     autocert.AcceptTOS,
//	//	HostPolicy: autocert.HostWhitelist("example.com"), //Your domain here
//	//	Cache:      autocert.DirCache("certs"),                   //Folder for storing certificates
//	//}
//	//certManager.GetCertificate()
//
//	bot, err := tgbotapi.NewBotAPI("")
//	if err != nil {
//		log.Panic(err)
//	}
//
//	bot.Debug = true
//	log.Printf("Authorized on account %s", bot.Self.UserName)
//
//	port := os.Getenv("PORT")
//	if port == "" {
//		port = "8080"
//		log.Printf("Defaulting to port %s", port)
//	}
//	_, err = bot.SetWebhook(tgbotapi.NewWebhookWithCert("https://www.google.com:" + port + "/" + bot.Token, "cert.pem"))
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	info, err := bot.GetWebhookInfo()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	if info.LastErrorDate != 0 {
//		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
//	}
//
//	updates := bot.ListenForWebhook("/" + bot.Token)
//	go http.ListenAndServeTLS("0.0.0.0:8443", "cert.pem", "key.pem", nil)
//
//	for update := range updates {
//		if update.Message == nil { // ignore any non-Message Updates
//			continue
//		}
//
//		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
//
//		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
//		msg.ReplyToMessageID = update.Message.MessageID
//
//		bot.Send(msg)
//	}
//}
//
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

	bytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Println(err)
		return
	}

	var update tgbotapi.Update
	err = json.Unmarshal(bytes, &update)
	if err != nil {
		log.Println(err)
		return
	}

	// to monitor changes run: heroku logs --tail
	log.Printf("From: %+v Text: %+v\n", update.Message.From, update.Message.Text)
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
		log.Fatal("$TELEGRAM_TOKEN must be set")
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
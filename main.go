package main

import (
	"fmt"
	"lamartire/cryptoratesbot/handler"
	"lamartire/cryptoratesbot/model"
	"lamartire/cryptoratesbot/scheldue"
	"lamartire/cryptoratesbot/service"
	"time"

	"lamartire/cryptoratesbot/db"
	"lamartire/cryptoratesbot/storage"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
)

var (
	tgToken        string
	userStorage    storage.UserStorage
	binanceService service.BinanceService
)

func init() {
	err := godotenv.Load()

	if err != nil {
		panic(fmt.Errorf("can't load env file: %v", err))
	}

	tgToken = os.Getenv("TG_TOKEN")
	dbClientParams := db.ClientParams{
		DBHost:     os.Getenv("DB_HOST"),
		DBName:     os.Getenv("DB_NAME"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
	}
	DB, err := db.ConnectMongoDatabase(dbClientParams)

	if err != nil {
		panic(fmt.Errorf("can't connect database: %v", err))
	}

	usersCollection := DB.Collection("users")
	userStorage = model.UserMongoStorage{Collection: usersCollection}
	binanceService = service.BinanceService{}
}

func main() {
	bot, err := tgbotapi.NewBotAPI(tgToken)

	if err != nil {
		panic(fmt.Errorf("can't initialize bot: %#v", err))
	}

	scheduler := scheldue.Scheduler{
		Bot:                   bot,
		UserStorage:           userStorage,
		CryptoExchangeService: binanceService,
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	if err != nil {
		panic(fmt.Errorf("can't initialize update channel: %v", err))
	}

	botHandler := handler.BotHandler{
		Bot: bot,
	}

	botHandler.AddMiddleware(handler.WithUserMiddleware(userStorage))
	botHandler.AddCommandHandler("start", handler.StartCommand())
	botHandler.AddCommandHandler("add", handler.AddCommand(userStorage, binanceService))
	botHandler.AddCommandHandler("remove", handler.RemoveCommand())
	botHandler.AddCommandHandler("suspend", handler.SuspendCommand(userStorage))
	botHandler.AddCommandHandler("resume", handler.ResumeCommand(userStorage))
	botHandler.AddCommandHandler("interval", handler.IntervalCommand(userStorage))
	botHandler.AddCallbackQueryHander("remove", handler.RemoveCallback(userStorage))

	go scheduler.StartNotifications(time.Second * 5)

	for update := range updates {
		botHandler.ProcessUpdate(update)
	}
}

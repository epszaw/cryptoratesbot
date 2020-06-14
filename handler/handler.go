package handler

import (
	"encoding/json"
	"lamartire/cryptoratesbot/bot"
	"lamartire/cryptoratesbot/storage"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

type Function func(bot bot.Bot, user storage.User, update tgbotapi.Update)

type Middleware func(bot bot.Bot, user storage.User, update tgbotapi.Update) (bot.Bot, storage.User, tgbotapi.Update)

type CallbackQueryPayload struct {
	Method string
	Data   string
}

type BotHandler struct {
	Bot                   bot.Bot
	Middlewares           []Middleware
	CommandsHandlers      map[string]Function
	CallbackQueryHandlers map[string]Function
}

func (b *BotHandler) AddMiddleware(middleware Middleware) {
	b.Middlewares = append(b.Middlewares, middleware)
}

func (b *BotHandler) AddCommandHandler(command string, handler Function) {
	if b.CommandsHandlers == nil {
		b.CommandsHandlers = make(map[string]Function)
	}

	b.CommandsHandlers[command] = handler
}

func (b *BotHandler) AddCallbackQueryHander(method string, handler Function) {
	if b.CallbackQueryHandlers == nil {
		b.CallbackQueryHandlers = make(map[string]Function)
	}

	b.CallbackQueryHandlers[method] = handler
}

func (b *BotHandler) ApplyMidllewares(update tgbotapi.Update) (bot.Bot, storage.User, tgbotapi.Update) {
	var lastUser storage.User

	lastBot := b.Bot
	lastUpdate := update

	for _, middleware := range b.Middlewares {
		lastBot, lastUser, lastUpdate = middleware(lastBot, lastUser, lastUpdate)
	}

	return lastBot, lastUser, lastUpdate
}

func (b *BotHandler) ProcessUpdate(update tgbotapi.Update) {
	bot, user, update := b.ApplyMidllewares(update)

	switch true {
	case update.CallbackQuery != nil:
		var payload CallbackQueryPayload

		if err := json.Unmarshal([]byte(update.CallbackQuery.Data), &payload); err != nil {
			logrus.Errorf("can't parse callback query data: %v", err)
			return
		}

		method := payload.Method
		handler := b.CallbackQueryHandlers[method]

		if handler != nil {
			handler(bot, user, update)
		} else {
			logrus.Errorf("can't process callback query method: %s", method)
		}
	case update.Message != nil && update.Message.IsCommand():
		command := update.Message.Command()
		handler := b.CommandsHandlers[command]

		if handler != nil {
			handler(bot, user, update)
		} else {
			logrus.Errorf("can't process command: %s", command)
		}
	default:
		logrus.Warningf("can't process update - unknown type: %v", update)
	}
}

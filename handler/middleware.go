package handler

import (
	"lamartire/cryptoratesbot/bot"
	"lamartire/cryptoratesbot/storage"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

func WithUserMiddleware(userStorage storage.UserStorage) Middleware {
	return func(bot bot.Bot, user storage.User, update tgbotapi.Update) (bot.Bot, storage.User, tgbotapi.Update) {
		var name string
		var chatID int64

		switch true {
		case update.CallbackQuery != nil:
			name = update.CallbackQuery.From.UserName
			chatID = update.CallbackQuery.Message.Chat.ID
		case update.Message != nil:
			name = update.Message.From.UserName
			chatID = update.Message.Chat.ID
		default:
			logrus.Warningf("can't gen UserName from update: %v", update)
			return bot, user, update
		}

		user, err := userStorage.GetUserByName(name)

		if err == storage.NoResultErr {
			user, err := userStorage.CreateUser(name, chatID)

			if err != nil {
				logrus.Errorf("can't create user: %v", err)
			}

			return bot, user, update
		}

		if err != nil && err != storage.NoResultErr {
			logrus.Errorf("can't get user: %v", err)

			return bot, user, update
		}

		return bot, user, update
	}
}

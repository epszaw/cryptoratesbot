package handler

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"lamartire/cryptoratesbot/bot"
	"lamartire/cryptoratesbot/storage"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func RemoveCallback(userStorage storage.UserStorage) Function {
	return func(bot bot.Bot, user storage.User, update tgbotapi.Update) {
		var payload CallbackQueryPayload

		data := update.CallbackQuery.Data

		if err := json.Unmarshal([]byte(data), &payload); err != nil {
			logrus.Errorf("remove callback error: %v", err)
			return
		}

		name := update.CallbackQuery.From.UserName
		callbackID := update.CallbackQuery.ID
		symbol := payload.Data
		newSymbols := make([]string, 0)

		for _, currentSymbol := range user.Symbols {
			if currentSymbol == symbol {
				continue
			}

			newSymbols = append(newSymbols, currentSymbol)
		}

		user.Symbols = newSymbols

		if err := userStorage.UpdateUserByName(name, user); err != nil {
			logrus.Errorf("remove callback, update user error: %v", err)
			return
		}

		reply := tgbotapi.NewCallback(callbackID, fmt.Sprintf("Пара %s была удалена", symbol))

		if _, err := bot.AnswerCallbackQuery(reply); err != nil {
			logrus.Errorf("remove callbak, answer error: %v", err)
		}
	}
}

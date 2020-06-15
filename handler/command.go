package handler

import (
	"encoding/json"
	"lamartire/cryptoratesbot/bot"
	"lamartire/cryptoratesbot/service"
	"lamartire/cryptoratesbot/storage"
	"lamartire/cryptoratesbot/template"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

func StartCommand() Function {
	return func(bot bot.Bot, user storage.User, update tgbotapi.Update) {
		var message string

		chatID := update.Message.Chat.ID

		message += "/add <символ> <символ> - добавление новой пары"
		message += "\n"
		message += "/remove - удаление существующей пары"
		message += "\n"
		message += "/suspend - остановка уведомлений"
		message += "\n"
		message += "/resume - возобновление уведомлений"
		message += "\n"
		message += "/interval <минуты> - установка интервала уведомлений"
		message += "\n"
		message += "/check - моментальная проверка курса"

		reply := tgbotapi.NewMessage(chatID, message)

		if _, err := bot.Send(reply); err != nil {
			logrus.Errorf("add command: send message %v", err)
			return
		}
	}
}

func AddCommand(userStorage storage.UserStorage, cryptoExchangeService service.CryptoExchangeService) Function {
	return func(bot bot.Bot, user storage.User, update tgbotapi.Update) {
		var message string

		chatID := update.Message.Chat.ID
		name := update.Message.From.UserName
		args := strings.Split(update.Message.CommandArguments(), " ")
		argsLen := len(args)

		if argsLen == 0 {
			message = "Похоже в твоем сообщении нет ни одного символа"
			message += "\n"
			message += "Пришли пару из двух символов. Например /add BTC USDT"
		}

		if argsLen == 1 {
			message = "Видимо в твоем сообщении только один символ"
			message += "\n"
			message += "Пришли пару из двух символов. Например /add BTC USDT"
		}

		if argsLen > 2 {
			message = "Судя по всему, в твоем сообщении больше, чем два символа"
			message += "\n"
			message += "Пришли пару из двух символов. Например /add BTC USDT"
		}

		if message == "" {
			newSymbol := args[0] + args[1]
			newSymbol = strings.ToUpper(newSymbol)

			for _, symbol := range user.Symbols {
				if symbol == newSymbol {
					message = "Эта пара уже была добавлена ранее"
					break
				}
			}

			newSymbols := make([]string, 0)
			newSymbols = append(newSymbols, newSymbol)
			newSymbolsPrices := cryptoExchangeService.GetSymbolsPrices(newSymbols)

			if len(newSymbolsPrices) == 0 {
				message = "Не могу найти цену для этой пары. Попробуй прислать что-то другое, например /add BTC USDT"
			} else {
				user.Symbols = append(user.Symbols, newSymbol)

				message = "Пара была сохранениа. Текущая стоимость:"
				message += "\n"
				message += template.FormatSymbolsPrices(newSymbolsPrices)
			}

			if len(newSymbolsPrices) > 0 {
				err := userStorage.UpdateUserByName(name, user)

				if err != nil {
					logrus.Errorf("add command: update user %v", err)
					return
				}
			}
		}

		reply := tgbotapi.NewMessage(chatID, message)

		if _, err := bot.Send(reply); err != nil {
			logrus.Errorf("add command: send message %v", err)
			return
		}
	}
}

func RemoveCommand() Function {
	return func(bot bot.Bot, user storage.User, update tgbotapi.Update) {
		var reply tgbotapi.Chattable
		var message string
		var buttons [][]tgbotapi.InlineKeyboardButton

		chatID := update.Message.Chat.ID

		if len(user.Symbols) == 0 {
			message = "У тебя нет ни одной наблюдаемой пары"
		}

		if message == "" {
			for _, symbol := range user.Symbols {
				payload := CallbackQueryPayload{
					Method: "remove",
					Data:   symbol,
				}
				jsonPayload, _ := json.Marshal(payload)
				button := tgbotapi.NewInlineKeyboardButtonData(symbol, string(jsonPayload))
				buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(button))
			}

			message = "Выбери пару, чтобы убрать ее из уведомлений"
		}

		reply = tgbotapi.NewMessage(chatID, message)
		sentReply, err := bot.Send(reply)

		if err != nil {
			logrus.Errorf("add command: send message error")
			return
		}

		if len(buttons) > 0 {
			buttonsReply := tgbotapi.NewEditMessageReplyMarkup(
				chatID,
				sentReply.MessageID,
				tgbotapi.NewInlineKeyboardMarkup(buttons...),
			)

			if _, err := bot.Send(buttonsReply); err != nil {
				logrus.Errorf("add command: send buttons error")
			}
		}
	}
}

func SuspendCommand(userStorage storage.UserStorage) Function {
	return func(bot bot.Bot, user storage.User, update tgbotapi.Update) {
		chatID := update.Message.Chat.ID
		name := update.Message.From.UserName

		user.Suspended = true

		if err := userStorage.UpdateUserByName(name, user); err != nil {
			logrus.Errorf("suspend command, update user error: %v", err)
			return
		}

		reply := tgbotapi.NewMessage(chatID, "Уведомления остановлены. Используй /resume, чтобы возобновить их")

		if _, err := bot.Send(reply); err != nil {
			logrus.Errorf("suspend command: send message %v", err)
			return
		}
	}
}

func ResumeCommand(userStorage storage.UserStorage) Function {
	return func(bot bot.Bot, user storage.User, update tgbotapi.Update) {
		chatID := update.Message.Chat.ID
		name := update.Message.From.UserName

		user.Suspended = false

		if err := userStorage.UpdateUserByName(name, user); err != nil {
			logrus.Errorf("resume command, update user error: %v", err)
			return
		}

		reply := tgbotapi.NewMessage(chatID, "Уведомления были возобновлены. Используй /suspend, чтобы остановить их")

		if _, err := bot.Send(reply); err != nil {
			logrus.Errorf("resume command: send message %v", err)
			return
		}
	}
}

func IntervalCommand(userStorage storage.UserStorage) Function {
	return func(bot bot.Bot, user storage.User, update tgbotapi.Update) {
		var message string

		chatID := update.Message.Chat.ID
		name := update.Message.From.UserName
		text := update.Message.CommandArguments()
		interval, err := strconv.Atoi(text)

		if err != nil {
			logrus.Errorf("interval command text parsing error: %v", err)
			return
		}

		if interval <= 0 {
			message = "Полученное значение некорректно. Пришли целое положительное число"
		}

		if message == "" {
			user.Interval = int64(interval)

			if err := userStorage.UpdateUserByName(name, user); err != nil {
				logrus.Errorf("interval command, update user error: %v", err)
				return
			}

			message = "Частота уведомлений была обновлена"
		}

		reply := tgbotapi.NewMessage(chatID, message)

		if _, err := bot.Send(reply); err != nil {
			logrus.Errorf("interval command: send message %v", err)
			return
		}
	}
}

func CheckCommand(exchangeService service.CryptoExchangeService) Function {
	return func(bot bot.Bot, user storage.User, update tgbotapi.Update) {
		var message string

		chatID := update.Message.Chat.ID
		symbols := user.Symbols

		if len(symbols) == 0 {
			message = "У тебя нет ни одной наблюдаемой пары"
		}

		if message == "" {
			prices := exchangeService.GetSymbolsPrices(symbols)
			message = template.FormatSymbolsPrices(prices)
		}

		reply := tgbotapi.NewMessage(chatID, message)

		if _, err := bot.Send(reply); err != nil {
			logrus.Errorf("check command: send message %v", err)
			return
		}
	}
}

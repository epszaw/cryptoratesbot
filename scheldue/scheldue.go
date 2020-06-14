package scheldue

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"lamartire/cryptoratesbot/bot"
	"lamartire/cryptoratesbot/service"
	"lamartire/cryptoratesbot/storage"
	"lamartire/cryptoratesbot/template"
	"time"
)

type Scheduler struct {
	Bot                   bot.Bot
	UserStorage           storage.UserStorage
	CryptoExchangeService service.CryptoExchangeService
}

func (s *Scheduler) StartNotifications(interval time.Duration) {
	ticker := time.NewTicker(interval)

	defer ticker.Stop()

	for range ticker.C {
		now := time.Now().UTC().Unix()
		users, err := s.UserStorage.GetUsersForNotification(now)

		if err == storage.NoResultErr {
			return
		}

		if err != nil {
			logrus.Errorf("scheldue request users error: %v", err)
			return
		}

		for _, user := range users {
			chatID := user.ChatID
			symbols := user.Symbols
			prices := s.CryptoExchangeService.GetSymbolsPrices(symbols)
			reply := tgbotapi.NewMessage(
				chatID,
				template.FormatSymbolsPrices(prices),
			)

			if _, err := s.Bot.Send(reply); err != nil {
				logrus.Errorf("scheldue reply error: %v", err)
				return
			}

			user.LastReply = now

			if err = s.UserStorage.UpdateUserByName(user.Name, user); err != nil {
				logrus.Errorf("scheldue update user error: %v", err)
			}
		}
	}
}

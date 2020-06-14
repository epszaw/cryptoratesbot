package service

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"lamartire/cryptoratesbot/util"
	"net/http"
	"sync"
)

const binanceAPIUrl = "https://api.binance.com/api/v3"

type BinanceService struct {
}

func (b BinanceService) GetSymbolsPrices(symbols []string) []SymbolPrice {
	var wg sync.WaitGroup

	result := make([]SymbolPrice, 0)
	requestBaseURL := binanceAPIUrl + "/ticker/price"

	for _, symbol := range symbols {
		wg.Add(1)

		symbol := symbol

		go func(wg *sync.WaitGroup) {
			defer wg.Done()

			var price SymbolPrice

			requestURL := util.AppendQueryToUrl(requestBaseURL, map[string]string{
				"symbol": symbol,
			})
			res, err := http.Get(requestURL)

			if err != nil {
				logrus.Errorf("binance request error: %v", err)
				return
			}

			body, err := ioutil.ReadAll(res.Body)

			if err != nil {
				logrus.Errorf("binance response parsing error: %v", err)
				return
			}

			if err := res.Body.Close(); err != nil {
				logrus.Errorf("binance response closing error: %v", err)
				return
			}

			if err := json.Unmarshal(body, &price); err != nil {
				logrus.Errorf("binance unmarshal error: %v", err)
				return
			}

			result = append(result, price)
		}(&wg)
	}

	wg.Wait()

	return result
}

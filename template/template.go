package template

import (
	"fmt"
	"lamartire/cryptoratesbot/service"
)

func FormatSymbolsPrices(prices []service.SymbolPrice) string {
	var result string

	for _, price := range prices {
		result += fmt.Sprintf("[%s] %s", price.Symbol, price.Price)
		result += "\n"
	}

	return result
}

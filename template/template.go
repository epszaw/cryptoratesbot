package template

import (
	"fmt"
	"lamartire/cryptoratesbot/service"
	"sort"
	"strings"
)

func FormatSymbolsPrices(prices []service.SymbolPrice) string {
	var result []string

	for _, price := range prices {
		result = append(result, fmt.Sprintf("[%s] %s", price.Symbol, price.Price))
	}

	sort.Strings(result)

	return strings.Join(result, "\n")
}

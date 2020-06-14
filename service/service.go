package service

type SymbolPrice struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

type CryptoExchangeService interface {
	GetSymbolsPrices(symbols []string) []SymbolPrice
}

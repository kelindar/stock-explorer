package finance

import "github.com/aktau/gofinance/util"

// Provider interface allows to query for a particular quote
type Provider interface {
	GetQuote(symbol string) Quote
	GetDividendHistory(symbol string) ([]DividendEntry, error)
}

// Quote represents a single quote for a particular stock symbol
type Quote struct {
	Symbol string
}

// DividendEntry represents a historican entry of a dividend
type DividendEntry struct {
	Symbol string            `json:"Symbol,string"`
	Date   util.YearMonthDay `json:"Date,string"`
	Value  float64           `json:"Dividends,string"`
}

package finance

import "time"

// Provider interface allows to query for a particular quote
type Provider interface {
	GetQuotes(symbols ...string) ([]Quote, error)
	GetDividendHistory(symbol string) ([]DividendEntry, error)
}

// DividendEntry represents a historican entry of a dividend
type DividendEntry struct {
	Symbol string
	Date   time.Time
	Value  float64
}

// Quote represents a single quote for a particular stock symbol
type Quote struct {
	Symbol   string /* e.g.: VEUR.AS, Vanguard dev. europe on Amsterdam */
	Name     string
	Exchange string

	/* last actualization of the results */
	Updated time.Time

	/* volume */
	Volume         int64 // outstanding shares
	AvgDailyVolume int64 // avg amount of shares traded

	/* dividend & related */
	PeRatio          float64   // Price / EPS
	EarningsPerShare float64   // (net income - spec.dividends) / avg.  outstanding shares
	DividendPerShare float64   // total (non-special) dividend payout / total amount of shares
	DividendYield    float64   // annual div. per share / price per share
	DividendExDate   time.Time // last dividend payout date

	/* price & derived */
	Bid, Ask            float64
	Open, PreviousClose float64
	LastTradePrice      float64
	Change              float64

	DayLow, DayHigh   float64
	YearLow, YearHigh float64

	Ma50, Ma200 float64 // 200- and 50-day moving average
}

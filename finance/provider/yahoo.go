package provider

import (
	"database/sql"
	"strconv"
	"strings"
	"time"

	"../"
)

// Provider provides a Query interface implementation
type Provider struct {
}

// NewProvider will create a new Yahoo Provider
func NewProvider() Provider {
	p := Provider{}

	return p
}

// GetQuotes retrieves a quote through a yahoo provider
func (p *Provider) GetQuotes(symbols ...string) ([]finance.Quote, error) {
	db, _ := sql.Open("yql", "||store://datatables.org/alltableswithkeys")
	stmt, err := db.Query(
		"select * from yahoo.finance.quotes where symbol in (?)",
		quoteSymbols(symbols))
	if err != nil {
		return nil, err
	}

	result := make([]finance.Quote, 0, 0)

	for stmt.Next() {
		var data map[string]interface{}

		stmt.Scan(&data)

		quote := finance.Quote{}
		quote.Symbol = readString(data, "Symbol")
		quote.Name = readString(data, "Name")
		quote.Exchange = readString(data, "StockExchange")
		quote.Updated = readDate(data, "LastTradeDate")

		quote.Volume = readInt(data, "Volume")
		quote.AvgDailyVolume = readInt(data, "AverageDailyVolume")

		quote.PeRatio = readFloat(data, "PERatio")
		quote.EarningsPerShare = readFloat(data, "EarningsShare")
		quote.DividendPerShare = readFloat(data, "DividendShare")
		quote.DividendYield = readFloat(data, "DividendYield")
		quote.DividendExDate = readDate(data, "ExDividendDate")

		quote.Bid = readFloat(data, "Bid")
		quote.Ask = readFloat(data, "Ask")

		quote.Open = readFloat(data, "Open")
		quote.PreviousClose = readFloat(data, "PreviousClose")
		quote.LastTradePrice = readFloat(data, "LastTradePriceOnly")
		quote.Change = readFloat(data, "Change")

		quote.DayLow = readFloat(data, "DaysLow")
		quote.DayHigh = readFloat(data, "DaysHigh")
		quote.YearLow = readFloat(data, "YearLow")
		quote.YearHigh = readFloat(data, "YearHigh")

		quote.Ma50 = readFloat(data, "FiftydayMovingAverage")
		quote.Ma200 = readFloat(data, "TwoHundreddayMovingAverage")

		if quote.Name != "" {
			// for each symbol we query, query the financials along with
			// the dividend history.
			quote.Financials, _ = p.GetFinancials(quote.Symbol)
			quote.DividendHistory, _ = p.GetDividendHistory(quote.Symbol)

			// Compute indicators
			quote.DividendFrequency = quote.GetLastYearDividendFrequency()
			quote.Profitability = quote.GetProfitability()
			quote.Growth = quote.GetGrowth()
		}

		result = append(result, quote)
	}

	return result, nil
}

// GetDividendHistory retrieves the dividend history for a particular symbol
func (p *Provider) GetDividendHistory(symbol string) ([]finance.DividendEntry, error) {
	db, _ := sql.Open("yql", "||store://datatables.org/alltableswithkeys")
	stmt, err := db.Query(
		"select * from yahoo.finance.dividendhistory where symbol = ? and startDate = ? and endDate = ?",
		symbol,
		time.Now().AddDate(-3, 0, 0).Format("2006-01-02"),
		time.Now().Format("2006-01-02"))
	if err != nil {
		return nil, err
	}

	result := make([]finance.DividendEntry, 0, 0)

	for stmt.Next() {
		var data map[string]interface{}
		stmt.Scan(&data)

		entry := finance.DividendEntry{}
		entry.Symbol = data["Symbol"].(string)
		entry.Value, _ = strconv.ParseFloat(data["Dividends"].(string), 64)
		date, _ := time.Parse("2006-01-02", data["Date"].(string))
		entry.Date = date

		result = append(result, entry)
	}

	return result, nil
}

// GetFinancials gets the financial information for the symbol
func (p *Provider) GetFinancials(symbol string) ([]finance.Financials, error) {
	m := NewMorningstar()
	return m.GetFinancials(symbol)
}

// Reads an integer
func readInt(data map[string]interface{}, name string) int64 {
	str, ok := data[name].(string)
	if !ok {
		return 0
	}

	result, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0
	}
	return result
}

// Reads a date
func readDate(data map[string]interface{}, name string) time.Time {
	str, ok := data[name].(string)
	if !ok {
		return time.Time{}
	}

	date, err := time.Parse("1/2/2006", str)
	if err != nil {
		return time.Time{}
	}

	return date
}

// Reads a string, safely
func readString(data map[string]interface{}, name string) string {
	str, ok := data[name].(string)
	if !ok {
		return ""
	}
	return str
}

// Reads a float, safely
func readFloat(data map[string]interface{}, name string) float64 {
	str, ok := data[name].(string)
	if !ok {
		return 0
	}

	result, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0
	}
	return result
}

// Format symbols for WHERE clause
func quoteSymbols(symbols []string) string {
	quotedSymbols := mapStr(func(s string) string {
		return `"` + s + `"`
	}, symbols)
	return strings.Join(quotedSymbols, ",")
}

func mapStr(mapping func(string) string, xs []string) []string {
	mxs := make([]string, 0, len(xs))
	for _, s := range xs {
		mxs = append(mxs, mapping(s))
	}
	return mxs
}

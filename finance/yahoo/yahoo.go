package yahoo

import (
	"database/sql"
	"strconv"
	"strings"
	"time"

	"../"
	"github.com/aktau/gofinance/util"
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
		var date time.Time
		stmt.Scan(&data)

		quote := finance.Quote{}
		quote.Symbol = data["Symbol"].(string)
		quote.Name = data["Name"].(string)
		quote.Exchange = data["StockExchange"].(string)

		date, _ = time.Parse("1/2/2006", data["LastTradeDate"].(string))
		quote.Updated = date

		quote.Volume, _ = strconv.ParseInt(data["Volume"].(string), 10, 64)
		quote.AvgDailyVolume, _ = strconv.ParseInt(data["AverageDailyVolume"].(string), 10, 64)

		quote.PeRatio, _ = strconv.ParseFloat(data["PERatio"].(string), 64)
		quote.EarningsPerShare, _ = strconv.ParseFloat(data["EarningsShare"].(string), 64)
		quote.DividendPerShare, _ = strconv.ParseFloat(data["DividendShare"].(string), 64)
		quote.DividendYield, _ = strconv.ParseFloat(data["DividendYield"].(string), 64)

		date, _ = time.Parse("1/2/2006", data["ExDividendDate"].(string))
		quote.DividendExDate = date

		quote.Bid, _ = strconv.ParseFloat(data["Bid"].(string), 64)
		quote.Ask, _ = strconv.ParseFloat(data["Ask"].(string), 64)

		quote.Open, _ = strconv.ParseFloat(data["Open"].(string), 64)
		quote.PreviousClose, _ = strconv.ParseFloat(data["PreviousClose"].(string), 64)
		quote.LastTradePrice, _ = strconv.ParseFloat(data["LastTradePriceOnly"].(string), 64)
		quote.Change, _ = strconv.ParseFloat(data["Change"].(string), 64)

		quote.DayLow, _ = strconv.ParseFloat(data["DaysLow"].(string), 64)
		quote.DayHigh, _ = strconv.ParseFloat(data["DaysHigh"].(string), 64)
		quote.YearLow, _ = strconv.ParseFloat(data["YearLow"].(string), 64)
		quote.YearHigh, _ = strconv.ParseFloat(data["YearHigh"].(string), 64)

		quote.Ma50, _ = strconv.ParseFloat(data["FiftydayMovingAverage"].(string), 64)
		quote.Ma200, _ = strconv.ParseFloat(data["TwoHundreddayMovingAverage"].(string), 64)

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
		time.Now().AddDate(-5, 0, 0).Format("2006-01-02"),
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

// Format symbols for WHERE clause
func quoteSymbols(symbols []string) string {
	quotedSymbols := util.MapStr(func(s string) string {
		return `"` + s + `"`
	}, symbols)
	return strings.Join(quotedSymbols, ",")
}

package scrape

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"../"
	"golang.org/x/net/html"
)

// Morningstar provides a morningstar scrapper
type Morningstar struct {
}

// NewMorningstar will create a new Morningstar scrapper
func NewMorningstar() Morningstar {
	return Morningstar{}
}

// GetFinancials gets the financial information for the symbol
func (p *Morningstar) GetFinancials(symbol string) ([]finance.Financials, error) {

	result := make([]finance.Financials, 0, 0)
	resp, _ := http.Get("http://quotes.morningstar.com/stockq/c-financials?&t=" + symbol)

	b := resp.Body
	defer b.Close()

	z := html.NewTokenizer(b)

	parsed := [][]string{}
	var column []string

	// parse the table into a 2d array
	for tt := z.Next(); tt != html.ErrorToken; tt = z.Next() {
		t := z.Token()

		switch {
		case t.Data == "tr" && t.Type == html.StartTagToken:
			column = []string{}

		case t.Data == "tr" && t.Type == html.EndTagToken:
			parsed = append(parsed, column)

		case t.Data == "td" && t.Type == html.StartTagToken:
			tt = z.Next()
			column = append(column, string(z.Text()))
		}
	}

	offsets := []int{1, 2, 3}

	const y1offset = 1
	const y2offset = 2
	const y3offset = 3

	for _, offset := range offsets {
		date, _ := time.Parse("2006-1", parsed[2][offset])

		f := finance.Financials{}
		f.Symbol = symbol
		f.Date = date

		// process income statement
		f.Income = finance.IncomeStatement{}
		f.Income.Symbol = symbol
		f.Income.Date = date
		f.Income.Revenue = readFloat(parsed[4][offset])
		f.Income.OperatingIncome = readFloat(parsed[5][offset])
		f.Income.NetIncome = readFloat(parsed[6][offset])
		f.Income.EarningsPerShare = readFloat(parsed[7][offset])
		f.Income.DilutedAverageShares = readFloat(parsed[8][offset])

		// process balance sheet
		f.Balance = finance.BalanceSheet{}
		f.Balance.Symbol = symbol
		f.Balance.Date = date
		f.Balance.CurrentAssets = readFloat(parsed[10][offset])
		f.Balance.NonCurrentAssets = readFloat(parsed[11][offset])
		f.Balance.TotalAssets = readFloat(parsed[12][offset])
		f.Balance.CurrentLiabilities = readFloat(parsed[13][offset])
		f.Balance.TotalLiabilities = readFloat(parsed[14][offset])

		// process cash flow sheet
		f.CashFlow = finance.CashFlow{}
		f.CashFlow.Symbol = symbol
		f.CashFlow.Date = date
		f.CashFlow.CashFromOperations = readFloat(parsed[17][offset])
		f.CashFlow.CapitalExpenditures = readFloat(parsed[18][offset])
		f.CashFlow.FreeCashFlow = readFloat(parsed[19][offset])

		result = append(result, f)
	}

	return result, nil
}

// Reads a float, safely
func readFloat(data string) float64 {
	data = strings.Replace(data, ",", "", -1)
	result, err := strconv.ParseFloat(data, 64)
	if err != nil {
		return 0
	}
	return result
}

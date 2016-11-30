package scrape

import (
	"fmt"
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

	if len(parsed) < 3 {
		return nil, fmt.Errorf("morningstar: no data available")
	}

	for _, offset := range offsets {
		date, _ := time.Parse("2006-1", parsed[2][offset])

		f := finance.Financials{}
		f.Symbol = symbol
		f.Date = date

		// process income statement
		f.Income = finance.IncomeStatement{}
		f.Income.Symbol = symbol
		f.Income.Date = date
		f.Income.Revenue = readValue(parsed, "Revenue", offset)
		f.Income.NetInterestIncome = readValue(parsed, "Net Interest Income", offset)
		f.Income.OperatingIncome = readValue(parsed, "Operating Income", offset)
		f.Income.NetIncome = readValue(parsed, "Net Income", offset)
		f.Income.EarningsPerShare = readValue(parsed, "Earnings Per Share", offset)
		f.Income.DilutedAverageShares = readValue(parsed, "Diluted Average Shares", offset)

		// process balance sheet
		f.Balance = finance.BalanceSheet{}
		f.Balance.Symbol = symbol
		f.Balance.Date = date
		f.Balance.NetLoans = readValue(parsed, "Net Loans", offset)
		f.Balance.Deposits = readValue(parsed, "Deposits", offset)
		f.Balance.CurrentAssets = readValue(parsed, "Current Assets", offset)
		f.Balance.NonCurrentAssets = readValue(parsed, "Non Current Assets", offset)
		f.Balance.TotalAssets = readValue(parsed, "Total Assets", offset)
		f.Balance.CurrentLiabilities = readValue(parsed, "Current Liabilities", offset)
		f.Balance.TotalLiabilities = readValue(parsed, "Total Liabilities", offset)

		// process cash flow sheet
		f.CashFlow = finance.CashFlow{}
		f.CashFlow.Symbol = symbol
		f.CashFlow.Date = date
		f.CashFlow.CashFromOperations = readValue(parsed, "Cash From Operations", offset)
		f.CashFlow.CashFromInvesting = readValue(parsed, "Cash From Investing", offset)
		f.CashFlow.CashFromFinancing = readValue(parsed, "Cash From Financing", offset)
		f.CashFlow.CapitalExpenditures = readValue(parsed, "Capital Expenditures", offset)
		f.CashFlow.FreeCashFlow = readValue(parsed, "Free Cash Flow", offset)

		result = append(result, f)
	}

	return result, nil
}

// Finds a value and reads it
func readValue(data [][]string, name string, offset int) float64 {
	for _, row := range data {
		if row == nil || len(row) == 0 {
			continue
		}

		if row[0] == name {
			return readFloat(row[offset])
		}
	}
	return 0
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

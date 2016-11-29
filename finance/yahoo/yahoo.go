package yahoo

import (
	"database/sql"
	"fmt"

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

// GetQuote retrieves a quote through a yahoo provider
func (p *Provider) GetQuote(symbol string) finance.Quote {
	return finance.Quote{}
}

// GetDividendHistory retrieves the dividend history for a particular symbol
func (p *Provider) GetDividendHistory(symbol string) ([]finance.DividendEntry, error) {
	db, _ := sql.Open("yql", "||store://datatables.org/alltableswithkeys")
	stmt, err := db.Query(
		"select * from yahoo.finance.dividendhistory where symbol = ? and startDate = ? and endDate = ?",
		"AAPL",
		"1962-01-01",
		"2016-03-10")
	if err != nil {
		return nil, err
	}

	result := make([]finance.DividendEntry, 0, 0)

	for stmt.Next() {
		var data map[string]interface{}
		stmt.Scan(&data)
		fmt.Printf("%v\n", data)
		//fmt.Printf("%v %v %v %v %v\n", data["Date"], data["Open"], data["High"], data["Low"], data["Close"])
	}

	return result, nil
}

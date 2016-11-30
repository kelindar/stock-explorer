package main

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"

	"./finance/yahoo"
)

func main() {
	provider := yahoo.NewProvider()

	symbols := []string{"WSR"}

	quotes, err := provider.GetQuotes(symbols...)
	if err != nil {
		panic(err)
	}

	for _, quote := range quotes {
		fmt.Printf("%v\n", quote.Name)
		spew.Dump(quote)

		fin, err := provider.GetFinancials(quote.Symbol)
		if err != nil {
			fmt.Println("(no financials found)")
			continue
		}

		for _, f := range fin {
			fmt.Printf("%v\t%v\t%v\n", f.Symbol, f.Date, f.Income.NetIncome)
		}

		/*hist, err := provider.GetDividendHistory(quote.Symbol)
		if err != nil {
			fmt.Println("(no dividend history found)")
			continue
		}

		for _, dividend := range hist {
			fmt.Printf("%v\t%v\t%v\n", dividend.Symbol, dividend.Date, dividend.Value)
		}*/
	}

}

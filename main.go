package main

import (
	"fmt"

	"./finance/yahoo"
)

func main() {
	provider := yahoo.NewProvider()

	symbols := []string{"AAPL", "WSR"}

	quotes, err := provider.GetQuotes(symbols...)
	if err != nil {
		panic(err)
	}

	for _, quote := range quotes {
		fmt.Printf("%v\n", quote.Name)
		//spew.Dump(quote)

		hist, err := provider.GetDividendHistory(quote.Symbol)
		if err != nil {
			panic(err)
		}

		for _, dividend := range hist {
			fmt.Printf("%v\t%v\t%v\n", dividend.Symbol, dividend.Date, dividend.Value)
		}
	}

}

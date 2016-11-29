package main

import "./finance/yahoo"

func main() {
	provider := yahoo.NewProvider()

	hist, err := provider.GetDividendHistory("AAPL")
	if err != nil {
		panic(err)
	}

	for dividend := range hist {
		Printf("%v", dividend.Value)
	}
}

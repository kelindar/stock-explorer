package finance

import (
	"math"
	"time"
)

// GetLastYearDividendFrequency computes the dividend frequency of the past
// year, returns how many dividend payments were made during that year.
func (q *Quote) GetLastYearDividendFrequency() int {
	divs := q.DividendHistory
	if divs == nil {
		return 0
	}

	targetYear := time.Now().Year() - 1
	result := 0

	for _, div := range divs {
		if div.Date.Year() == targetYear {
			result++
		}
	}

	return result
}

// GetRevenueGrowth gets the average revenue growth
func (q *Quote) GetRevenueGrowth() float64 {
	if q.Financials == nil {
		return 0
	}

	y1 := q.Financials[0].Income.Revenue
	y2 := q.Financials[1].Income.Revenue
	y3 := q.Financials[2].Income.Revenue

	c1 := (y1 - y2)
	r1 := math.Abs(c1) / math.Abs(y2)
	if c1 < 0 {
		r1 *= -1
	}

	c2 := (y2 - y3)
	r2 := math.Abs(c2) / math.Abs(y3)
	if c2 < 0 {
		r2 *= -1
	}

	return (r1 + r2) / 2
}

// GetFFOGrowth gets the average cash from operations growth
func (q *Quote) GetFFOGrowth() float64 {
	if q.Financials == nil {
		return 0
	}

	y1 := q.Financials[0].CashFlow.CashFromOperations
	y2 := q.Financials[1].CashFlow.CashFromOperations
	y3 := q.Financials[2].CashFlow.CashFromOperations

	c1 := (y1 - y2)
	r1 := math.Abs(c1) / math.Abs(y2)
	if c1 < 0 {
		r1 *= -1
	}

	c2 := (y2 - y3)
	r2 := math.Abs(c2) / math.Abs(y3)
	if c2 < 0 {
		r2 *= -1
	}

	return (r1 + r2) / 2
}

// GetGrowth estimates the growth rating (A-F)
func (q *Quote) GetGrowth() string {
	return getRating(q.GetRevenueGrowth())
}

// GetProfitability estimates the profitability rating (A-F)
func (q *Quote) GetProfitability() string {
	return getRating(q.GetFFOGrowth())
}

// Gets a rating
func getRating(i float64) string {
	switch {
	case i < 0:
		return "F"
	case i >= 0 && i < 0.05:
		return "E"
	case i >= 0.05 && i < 0.10:
		return "D"
	case i >= 0.10 && i < 0.20:
		return "C"
	case i >= 0.20 && i < 0.30:
		return "B"
	case i >= 0.30:
		return "A"
	}
	return "F"
}

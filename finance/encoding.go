package finance

import (
	"encoding/json"
	"time"
)

var (
	FmtMonthDay      = "Jan _2"    /* yahoo finance often returns these types of dates */
	FmtDayTMonthYear = "_2-Jan-06" /* yahoo finance is quite inconsisten sometimes */
	FmtYearMonthDay  = "2006-01-02"
)

type JsonTime time.Time

// MonthDay is a shortime
type MonthDay JsonTime

// UnmarshalJSON deserializes the JSON
func (jt *MonthDay) UnmarshalJSON(data []byte) error {
	dt := (*JsonTime)(jt)
	return dt.JSONParse(data, FmtMonthDay, FmtDayTMonthYear)
}

// GetTime gets the time from a JSON-encoded time
func (jt *MonthDay) GetTime() time.Time {
	return (time.Time)(*jt)
}

// YearMonthDay is a type for a JsonTime
type YearMonthDay JsonTime

// UnmarshalJSON deserializes the JSON
func (jt *YearMonthDay) UnmarshalJSON(data []byte) error {
	dt := (*JsonTime)(jt)
	return dt.JSONParse(data, FmtYearMonthDay)
}

// GetTime gets the time from a JSON-encoded time
func (jt *YearMonthDay) GetTime() time.Time {
	return (time.Time)(*jt)
}

// JSONParse parses the json
func (dt *JsonTime) JSONParse(data []byte, formats ...string) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	// try all formats, from first to last
	var err error
	for _, format := range formats {
		t, err := time.Parse(format, s)

		// if one works, quit
		if err == nil {
			*dt = (JsonTime)(t)
			return nil
		}
	}

	// all formats failed, return error
	return err
}

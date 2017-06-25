package helpers

import "time"

func ConvertPoloniexDate(val string) time.Time {
	// date example : 2017-06-18 04:31:08
	t, _ := time.Parse("2006-01-02 15:04:05", val)
	return t
}

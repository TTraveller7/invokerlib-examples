package main

import (
	"strconv"
	"time"
)

func SafeParseInt64(s string) int64 {
	res, _ := strconv.ParseInt(s, 10, 64)
	return res
}

func SafeParseFloat64(s string) float64 {
	res, _ := strconv.ParseFloat(s, 64)
	return res
}

func SafeParseTime(s string) time.Time {
	res, _ := time.Parse("2006-01-02 15:04:05.999", s)
	return res
}

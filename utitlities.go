package main

import (
	"strconv"
	"strings"
)

func getDollarValueAsInt(moneyString string) int64 {
	moneyString = strings.TrimSpace(moneyString)
	moneyString = strings.Trim(moneyString, "$")
	if strings.Index(moneyString, ".") == -1 {
		moneyString = moneyString + "00"
	} else {
		moneyString = moneyString[:strings.Index(moneyString, ".")] + moneyString[strings.Index(moneyString, ".")+1:strings.Index(moneyString, ".")+3]
	}
	moneyInt, _ := strconv.ParseInt(moneyString, 10, 64)

	return moneyInt
}

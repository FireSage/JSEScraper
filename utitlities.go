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
		moneyString1 := moneyString[:strings.Index(moneyString, ".")]
		moneyString2 := moneyString[strings.Index(moneyString, ".")+1:] + "00"
		moneyString = moneyString1 + moneyString2[0:2]
	}
	moneyInt, _ := strconv.ParseInt(moneyString, 10, 64)

	return moneyInt
}

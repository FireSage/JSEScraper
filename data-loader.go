package main

import (
	// "fmt"
	// "strconv"
	// "strings"

	"github.com/gocolly/colly"
)

func GetStocksWithUrls(market string, stockList *[]Stock) []Stock {
	url := "https://www.jamstockex.com/trading/trade-quotes/?market=" + market + "-market"

	// stockList := make([]Stock, 0, 120)
	c := colly.NewCollector()

	// search for table with main market stocks in html
	c.OnHTML("table", func(e *colly.HTMLElement) {
		if e.Index == 2 {
			e.ForEach("tbody tr", func(_ int, el *colly.HTMLElement) {
				var s Stock = Stock{
					Ticker: el.ChildText("td:nth-child(2)"),
					JSEUrl: el.ChildAttr("td:nth-child(2) a ", "href"),
				}
				*stockList = append(*stockList, s)
			})
		}
	})

	//attempt to visit url
	c.Visit(url)

	return *stockList
}

// load all stocks from JSE
func loadJSE(history bool, refresh_list bool) []Stock {
	var stockList []Stock = make([]Stock, 0, 120)

	if refresh_list {
		stockList = GetStocksWithUrls(MARKET_COMBINED, &stockList)
	} else {
		read_stocks_from_json_file("stocklist.json", &stockList)
	}

	for i, _ := range stockList {
		stockList[i].loadStock(history)
	}
	write_text_file(get_json_list(stockList), "stocklist.json")

	return stockList
}

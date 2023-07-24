package main

import (
	// "fmt"
	// "strconv"
	// "strings"

	"github.com/gocolly/colly"
)

func GetStocksWithUrls(market string) []Stock {
	url := "https://www.jamstockex.com/trading/trade-quotes/?market=" + market + "-market"
	stockList := make([]Stock, 0, 120)
	// stockUrls := []string{}
	c := colly.NewCollector()

	// search for table with main market stocks in html
	c.OnHTML("table", func(e *colly.HTMLElement) {
		if e.Index == 2 {
			e.ForEach("tbody tr", func(_ int, el *colly.HTMLElement) {
				// stockUrls = append(stockUrls, el.ChildAttr("td:nth-child(2) a ", "href"))
				var s Stock = Stock{
					Ticker: el.ChildText("td:nth-child(2)"),
					JSEUrl: el.ChildAttr("td:nth-child(2) a ", "href"),
				}
				stockList = append(stockList, s)
			})
		}
	})

	//attempt to visit url
	c.Visit(url)

	return stockList
}

// load all stocks from JSE
func loadJSE(history bool) []Stock {
	stockList := GetStocksWithUrls("main")
	// urls := GetStocksWithUrls("main")

	for i, _ := range stockList {
		stockList[i].loadStock(history)
		// fmt.Println(stock)
		// stockList = append(stockList, *stock)
	}

	return stockList
}

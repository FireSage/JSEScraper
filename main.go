package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

type Stock struct {
	Name            string
	Ticker          string
	ShareCount      uint64
	ClosingPrice    float32
	LastTradedPrice float32
	JSEUrl          string
}

func newStock(name string, ticker string, sharecount uint64, closingPrice float32, lastTradedPrice float32) *Stock {
	stock := Stock{
		Name:            name,
		Ticker:          ticker,
		ShareCount:      sharecount,
		ClosingPrice:    closingPrice,
		LastTradedPrice: lastTradedPrice,
	}
	return &stock
}

func GetStocksWithUrls(market string) []string {
	url := "https://www.jamstockex.com/trading/trade-quotes/?market=" + market + "-market"
	var stockUrls []string
	c := colly.NewCollector()

	// search for table with main market stocks in html
	c.OnHTML("table", func(e *colly.HTMLElement) {
		if e.Index == 2 {
			e.ForEach("tr", func(_ int, el *colly.HTMLElement) {
				stockUrls = append(stockUrls, el.ChildAttr("td:nth-child(2) a ", "href"))
			})
		}
	})

	//attempt to visit url
	c.Visit(url)

	return stockUrls
}

func loadStock(url string) *Stock {
	var name string
	var ticker string
	var shareCount uint64
	var lastTradedPrice float32
	var closingPrice float32
	dataCollector := colly.NewCollector()

	dataCollector.OnHTML("h2", func(e *colly.HTMLElement) {
		//fmt.Println(link)
		nameStr := e.Text
		name = (nameStr[2:strings.LastIndex(nameStr, " (")])
		ticker = (nameStr[(strings.LastIndex(nameStr, " (") + 2):strings.LastIndex(nameStr, ")")])

		priceStr := e.DOM.Next().Text()
		closingPrice, err := strconv.ParseFloat(priceStr[strings.Index(priceStr, "$")+1:(strings.Index(priceStr, ".")+2)], 64)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(closingPrice)
		}
	})

	dataCollector.OnHTML(".tw-bg-gray-50.tw-rounded-sm:first-child", func(e *colly.HTMLElement) {

		shareCountString := e.ChildText(".tw-flex.tw-justify-between.tw-bg-gray-50.tw-p-4:nth-child(5) span:nth-child(2)")
		shareCount, err := strconv.ParseInt(strings.ReplaceAll(shareCountString[0:strings.LastIndex(shareCountString, "\n u")], ",", ""), 10, 64)
		if err != nil {
			fmt.Println(shareCount)
		} else {
			fmt.Println(err)
		}
		// stockList += ")\n"

	})

	//attempt to visit url
	dataCollector.Visit(url)

	return newStock(name, ticker, shareCount, closingPrice, lastTradedPrice)
}
func loadJSE() []Stock {
	var stockList []Stock
	urls := GetStocksWithUrls("main")

	for _, url := range urls {
		// TODO:: prevent empty first element
		stock := loadStock(url)
		stockList = append(stockList, *stock)
	}

	return stockList
}

func SayHello() string {
	return "Hello, World"
}

func main() {
	loadJSE()
	// loadJSE("https://www.jamstockex.com/trading/trade-quotes/")
	// fmt.Println(loadJSE("https://www.jamstockex.com/trading/trade-quotes/"))
}

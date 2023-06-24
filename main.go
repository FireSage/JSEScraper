package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

type Stock struct {
	Name             string
	Ticker           string
	ShareCount       uint64
	ClosingPrice     int64
	LastTradedPrice  int64
	LastTradedVolume uint64
	JSEUrl           string
}

func newStock(name string, ticker string, sharecount uint64, closingPrice int64, lastTradedPrice int64, lastTradedVolume uint64) *Stock {
	stock := Stock{
		Name:             name,
		Ticker:           ticker,
		ShareCount:       sharecount,
		ClosingPrice:     closingPrice,
		LastTradedPrice:  lastTradedPrice,
		LastTradedVolume: lastTradedVolume,
	}
	return &stock
}

func GetStocksWithUrls(market string) []string {
	url := "https://www.jamstockex.com/trading/trade-quotes/?market=" + market + "-market"
	stockUrls := []string{}
	c := colly.NewCollector()

	// search for table with main market stocks in html
	c.OnHTML("table", func(e *colly.HTMLElement) {
		if e.Index == 2 {
			e.ForEach("tbody tr", func(_ int, el *colly.HTMLElement) {
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
	var lastTradedVolume uint64
	var lastTradedPrice int64
	var closingPrice int64

	dataCollector := colly.NewCollector()

	// get name, ticker, close price
	dataCollector.OnHTML("h2", func(e *colly.HTMLElement) {
		//fmt.Println(link)
		nameStr := e.Text
		name = (nameStr[2:strings.LastIndex(nameStr, " (")])
		ticker = (nameStr[(strings.LastIndex(nameStr, " (") + 2):strings.LastIndex(nameStr, ")")])

		priceStr := e.DOM.Next().Text()
		priceStr = priceStr[strings.Index(priceStr, "$")+1 : (strings.Index(priceStr, ".") + 3)]
		closingPrice = getDollarValueAsInt(priceStr) //, err := strconv.ParseFloat(priceStr[strings.Index(priceStr, "$")+1:(strings.Index(priceStr, ".")+3)], 64)
		// if err != nil {
		// fmt.Println(err)
		// } else {
		fmt.Println(closingPrice)
		// }
	})

	// get share count
	dataCollector.OnHTML(".tw-bg-gray-50.tw-rounded-sm:first-child", func(e *colly.HTMLElement) {
		var err error
		shareCountString := e.ChildText(".tw-flex.tw-justify-between.tw-bg-gray-50.tw-p-4:nth-child(5) span:nth-child(2)")
		shareCount, err = strconv.ParseUint(strings.ReplaceAll(shareCountString[0:strings.LastIndex(shareCountString, "\n u")], ",", ""), 10, 64)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(shareCount)
		}
		// stockList += ")\n"

	})

	dataCollector.OnHTML(".tw-container > div:nth-child(2) .tw-justify-between span:nth-child(2)", func(e *colly.HTMLElement) {
		switch e.Index {
		case 4:
			lastTradedVolumeString := e.Text
			lastTradedVolume, _ = strconv.ParseUint(strings.ReplaceAll(lastTradedVolumeString, ",", ""), 10, 64)
		case 5:
			lastTradedPrice = getDollarValueAsInt(e.Text) //strconv.ParseFloat(e.Text[strings.Index(e.Text, "$")+1:(strings.Index(e.Text, ".")+3)], 64)
			// TODO:: continue here

		}
	})

	//attempt to visit url
	dataCollector.Visit(url)
	// TODO:: Throw error or retry when stock fails to load, fix
	return newStock(name, ticker, shareCount, closingPrice, lastTradedPrice, lastTradedVolume)
}
func loadJSE() []Stock {
	stockList := make([]Stock, 0, 120)
	urls := GetStocksWithUrls("main")

	for _, url := range urls {
		stock := loadStock(url)
		stockList = append(stockList, *stock)
	}

	return stockList
}

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

func main() {
	loadJSE()
	// loadJSE("https://www.jamstockex.com/trading/trade-quotes/")
	// fmt.Println(loadJSE("https://www.jamstockex.com/trading/trade-quotes/"))
}

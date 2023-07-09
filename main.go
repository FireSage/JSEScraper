package main

// "github.com/gocolly/colly"

/*
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
*/
func loadJSE() []Stock {
	stockList := GetStocksWithUrls("main")
	// urls := GetStocksWithUrls("main")

	for i, _ := range stockList {
		stockList[i].loadStock()
		// fmt.Println(stock)
		// stockList = append(stockList, *stock)
	}

	return stockList
}

func main() {
	loadJSE()
	// loadJSE("https://www.jamstockex.com/trading/trade-quotes/")
	// fmt.Println(loadJSE("https://www.jamstockex.com/trading/trade-quotes/"))
}

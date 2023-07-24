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
	history          []PriceHistory
}

type PriceHistory struct {
	Price  uint16
	Volume uint64
	Date   uint64
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

func (s *Stock) loadStock(history bool) {
	dataCollector := colly.NewCollector()

	// get name, ticker, close price
	dataCollector.OnHTML("h2", func(e *colly.HTMLElement) {
		nameStr := e.Text
		s.Name = (nameStr[2:strings.LastIndex(nameStr, " (")])
		s.Ticker = (nameStr[(strings.LastIndex(nameStr, " (") + 2):strings.LastIndex(nameStr, ")")])

		priceStr := e.DOM.Next().Text()
		priceStr = priceStr[strings.Index(priceStr, "$")+1 : (strings.Index(priceStr, ".") + 3)]
		s.ClosingPrice = getDollarValueAsInt(priceStr)
	})

	// get share count
	dataCollector.OnHTML(".tw-bg-gray-50.tw-rounded-sm:first-child", func(e *colly.HTMLElement) {
		var err error
		shareCountString := e.ChildText(".tw-flex.tw-justify-between.tw-bg-gray-50.tw-p-4:nth-child(5) span:nth-child(2)")
		s.ShareCount, _ = strconv.ParseUint(strings.ReplaceAll(shareCountString[0:strings.LastIndex(shareCountString, "\n u")], ",", ""), 10, 64)
		if err != nil {
			fmt.Println(err)
		}
	})

	dataCollector.OnHTML(".tw-container > div:nth-child(2) .tw-justify-between span:nth-child(2)", func(e *colly.HTMLElement) {
		switch e.Index {
		case 4:
			lastTradedVolumeString := e.Text
			s.LastTradedVolume, _ = strconv.ParseUint(strings.ReplaceAll(lastTradedVolumeString, ",", ""), 10, 64)
		case 5:
			s.LastTradedPrice = getDollarValueAsInt(e.Text) //strconv.ParseFloat(e.Text[strings.Index(e.Text, "$")+1:(strings.Index(e.Text, ".")+3)], 64)
			// TODO:: continue here

		}
	})

	// load history if flag is true
	if history {
		dataCollector.OnHTML("script", func(e *colly.HTMLElement) {
			if e.Index == 12 {
				price_text := e.Text[strings.Index(e.Text, "name: 'Price',"):]
				volume_text := e.Text[strings.Index(e.Text, "name: 'Volume',"):]
				price_text = price_text[strings.Index(price_text, "[[")+2 : strings.Index(price_text, "]],")]
				volume_text = volume_text[strings.Index(volume_text, "[[")+2 : strings.Index(volume_text, "]],")]
				s.processPriceHistory(&price_text, &volume_text)
			}
		})
	}

	//attempt to visit url
	dataCollector.Visit(s.JSEUrl)
}

// converts price history string to array of price history
func (s *Stock) processPriceHistory(priceHistory *string, volumeeHistory *string) {
	priceHist := strings.Split(*priceHistory, "],[")
	volHist := strings.Split(*volumeeHistory, "],[")
	var length int = len(volHist)
	for i := 0; i < length; i++ {
		date, _ := strconv.ParseUint(priceHist[i][0:strings.Index(priceHist[i], ",")], 10, 64)
		price := uint16(getDollarValueAsInt(priceHist[i][strings.Index(priceHist[i], ",")+1:]))
		volume, _ := strconv.ParseUint(volHist[i][strings.Index(volHist[i], ",")+1:], 10, 64)
		var hist PriceHistory
		hist.Date = date
		hist.Price = price
		hist.Volume = volume
		s.history = append(s.history, hist)
	}
}

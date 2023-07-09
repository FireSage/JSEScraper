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

func (s *Stock) loadStock() {
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

	//attempt to visit url
	dataCollector.Visit(s.JSEUrl)
}

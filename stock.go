package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

var MARKET_COMBINED string = "combined"
var MARKET_MAIN string = "main"
var MARKET_JUNIOR string = "junior"
var MARKET_USD string = "usd-equities"
var MARKET_BOND string = "bond"

type Stock struct {
	Name             string
	Ticker           string
	ShareCount       uint64
	ClosingPrice     int64
	LastTradedPrice  int64
	LastTradedVolume uint64
	JSEUrl           string
	History          []PriceHistory
}

type PriceHistory struct {
	OpenPrice  uint16
	HighPrice  uint16
	LowPrice   uint16
	ClosePrice uint16
	Volume     uint64
	Date       uint64
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
		if s.Name == "" {
			nameStr := e.Text
			s.Name = (nameStr[2:strings.LastIndex(nameStr, " (")])
			s.Ticker = (nameStr[(strings.LastIndex(nameStr, " (") + 2):strings.LastIndex(nameStr, ")")])
		}

		priceStr := e.DOM.Next().Text()
		priceStr = priceStr[strings.Index(priceStr, "$")+1 : (strings.Index(priceStr, ".") + 3)]
		s.ClosingPrice = getDollarValueAsInt(priceStr)
	})

	// get share count
	dataCollector.OnHTML(".tw-bg-gray-50.tw-rounded-sm:first-child", func(e *colly.HTMLElement) {
		if s.ShareCount == 0 {
			var err error
			shareCountString := e.ChildText(".tw-flex.tw-justify-between.tw-bg-gray-50.tw-p-4:nth-child(5) span:nth-child(2)")
			s.ShareCount, _ = strconv.ParseUint(strings.ReplaceAll(shareCountString[0:strings.LastIndex(shareCountString, "\n u")], ",", ""), 10, 64)
			if err != nil {
				fmt.Println(err)
			}
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
	// TODO: load corporate actions
	if history {
		dataCollector.OnHTML("script", func(e *colly.HTMLElement) {
			if e.Index == 12 {
				if !strings.Contains(e.Text, "type: 'candlestick',\n name: ''") {
					price_text := e.Text[strings.Index(e.Text, "type: 'candlestick',"):]
					volume_text := e.Text[strings.Index(e.Text, "name: 'Volume',"):]
					price_text = price_text[strings.Index(price_text, "[[")+2 : strings.Index(price_text, "]],")]
					volume_text = volume_text[strings.Index(volume_text, "[[")+2 : strings.Index(volume_text, "]],")]
					s.processPriceHistory(&price_text, &volume_text)
				}
			}
		})
	}

	//attempt to visit url
	dataCollector.Visit(s.JSEUrl)

	// write_text_file(s.get_json_string(), s.Ticker)
}

// converts price history string to array of price history
func (s *Stock) processPriceHistory(priceHistory *string, volumeeHistory *string) {
	priceHist := strings.Split(*priceHistory, "],[")
	volHist := strings.Split(*volumeeHistory, "],[")
	var lastRecordedHistory uint64 = 0
	if len(s.History) > 0 {
		lastRecordedHistory = s.History[len(s.History)-1].Date
	}
	var length int = len(volHist)
	for i := 0; i < length; i++ {
		priceHistoryData := strings.Split(priceHist[i], ",")
		date, _ := strconv.ParseUint(priceHistoryData[0], 10, 64)

		if lastRecordedHistory < date {
			openprice := uint16(getDollarValueAsInt(priceHistoryData[1]))
			highprice := uint16(getDollarValueAsInt(priceHistoryData[2]))
			LowPrice := uint16(getDollarValueAsInt(priceHistoryData[3]))
			closePrice := uint16(getDollarValueAsInt(priceHistoryData[4]))
			volume, _ := strconv.ParseUint(volHist[i][strings.Index(volHist[i], ",")+1:], 10, 64)
			var hist PriceHistory
			hist.Date = date
			hist.OpenPrice = openprice
			hist.HighPrice = highprice
			hist.LowPrice = LowPrice
			hist.ClosePrice = closePrice
			hist.Volume = volume
			s.History = append(s.History, hist)
		}

	}
}

func (s *Stock) get_json_string() string {
	stock_string, _ := json.Marshal(s)
	return string(stock_string)
}

func get_json_list(stocks []Stock) string {
	stock_string, _ := json.Marshal(stocks)
	return string(stock_string)
}

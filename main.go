package main 

import (
	"fmt"

	"github.com/gocolly/colly"
)

var stockList string

func loadJSE( url string ) string{
	// initialise collector
	c :=colly.NewCollector()
	dataCollector := c.Clone()

	// search for tables in html
	c.OnHTML("table", func(e *colly.HTMLElement) {
		if e.Index == 1 {
			// fmt.Println(e.Index)

			e.ForEach("tr", func(_ int, el *colly.HTMLElement){
				//fmt.Println(el.ChildText("td:nth-child(2)"))
				link := el.ChildAttr("td:nth-child(2) a ", "href")
				//fmt.Println(link)
				dataCollector.Visit(link)
			})
		}
		

	})
	
	dataCollector.OnHTML("h2", func(e *colly.HTMLElement){
		//fmt.Println(e.Text)
		stockList += ("(" + e.Text)
		//fmt.Println("found")
		//fmt.Println(e.Request.URL)
	})

	dataCollector.OnHTML(".tw-bg-gray-50.tw-rounded-sm:first-child", func(e *colly.HTMLElement){
		//fmt.Println(e.Text)
		fmt.Println(e.ChildText(".tw-flex.tw-justify-between.tw-bg-gray-50.tw-p-4:nth-child(5)"))
		
		stockList += ","
		stockList += e.ChildText(".tw-flex.tw-justify-between.tw-bg-gray-50.tw-p-4:nth-child(5)")
		stockList += ")\n"

		fmt.Println(stockList)
	})
	//attempt to visit url
	c.Visit(url)
	
	return stockList
}

func SayHello() string {
	return "Hello, World"
}

func main(){
	// loadJSE("https://www.jamstockex.com/trading/trade-quotes/")

	fmt.Println(loadJSE("https://www.jamstockex.com/trading/trade-quotes/"))
}

package main 

import (
	"fmt"
	"strings"
	"strconv"
	"github.com/gocolly/colly"
)

// type Stock struct{
// 	Name string
// 	Ticker string
// 	ShareCount uint
// } 
// func newStock(string name, string ticker uint sharecount) *Stock{
// 	stock := Stock{
// 		Name: name
// 		Ticker: ticker
// 		ShareCount: shareCount
// 	}
// 	return &stock
// }


func loadJSE( url string ) string{
	var stockList string
	// initialise collector
	c :=colly.NewCollector()
	dataCollector := c.Clone()
	// stock := Stock

	// search for table with main market stocks in html
	c.OnHTML("table", func(e *colly.HTMLElement) {
		if e.Index == 2 {

			e.ForEach("tr", func(_ int, el *colly.HTMLElement){
				//fmt.Println(el.ChildText("td:nth-child(2)"))
				link := el.ChildAttr("td:nth-child(2) a ", "href")
				//fmt.Println(link)
				dataCollector.Visit(link)
			})
		}
		
	})
	
	dataCollector.OnHTML("h2", func(e *colly.HTMLElement){
		//fmt.Println(link)
		nameStr := e.Text
		name := (nameStr[0:strings.LastIndex(nameStr, " (")])
		ticker := (nameStr[(strings.LastIndex(nameStr, " (")+2): strings.LastIndex(nameStr, ")")])
		stockList += ("(" + name)
		stockList += ("," + ticker)
		//fmt.Println("found")
		//fmt.Println(e.Request.URL)
		fmt.Println(stockList);
	})

	dataCollector.OnHTML(".tw-bg-gray-50.tw-rounded-sm:first-child", func(e *colly.HTMLElement){
		//fmt.Println(e.Text)
		// fmt.Println(e.ChildText(".tw-flex.tw-justify-between.tw-bg-gray-50.tw-p-4:nth-child(5)"))
		
		// fmt.Println(nameStr);

		stockList += ", "
		shareCountString := e.ChildText(".tw-flex.tw-justify-between.tw-bg-gray-50.tw-p-4:nth-child(5) span:nth-child(2)")
		shareCount, err := strconv.ParseInt( strings.ReplaceAll(shareCountString[0: strings.LastIndex(shareCountString, "\n u")], ",", ""), 10, 64 )
		if err == nil{
			fmt.Println(shareCount)
		}else{
			fmt.Println(err)
		}
		// stockList += ")\n"

	})

	//attempt to visit url
	c.Visit(url)
	
	return stockList
}

func SayHello() string {
	return "Hello, World"
}

func main(){
	loadJSE("https://www.jamstockex.com/trading/trade-quotes/")

	// fmt.Println(loadJSE("https://www.jamstockex.com/trading/trade-quotes/"))
}

package main

import "flag"

func main() {
	history := flag.Bool("h", false, "Load history")
	refresh_list := flag.Bool("r", false, "Load history")

	flag.Parse()
	loadJSE(*history, *refresh_list)
	// loadJSE("https://www.jamstockex.com/trading/trade-quotes/")
	// fmt.Println(loadJSE("https://www.jamstockex.com/trading/trade-quotes/"))
}

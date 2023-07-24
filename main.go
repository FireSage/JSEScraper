package main

import (
	"flag"
)

func main() {
	history := flag.Bool("h", false, "Load history")

	flag.Parse()
	// fmt.Println()
	loadJSE(*history)
	// loadJSE("https://www.jamstockex.com/trading/trade-quotes/")
	// fmt.Println(loadJSE("https://www.jamstockex.com/trading/trade-quotes/"))
}

package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/gocarina/gocsv"
)

const usage = `Usage:
  -v, --verbose verbose output
  -h, --help prints help information 
`

// // for testing, comment out when not needed
// func main() {
// 	holdings := getShareholdingQtr("512529", "FY22Q4")
// 	JsonToCsv(holdings, "testdata/test.json")
// }

func main() {
	// fmt.Println("##### Marketview v0.1 #####")
	var flagFetchSymbols, flagViewSymbols, flagFetchShareholding, flagViewShareHolding bool
	flag.BoolVar(&flagFetchSymbols, "fetch-sym", false, "Fetch symbols from BSE/NSE Servers")
	flag.BoolVar(&flagFetchSymbols, "S", false, "Fetch symbols from BSE/NSE Servers")
	flag.BoolVar(&flagViewSymbols, "view-sym", false, "View symbols in csv format")
	flag.BoolVar(&flagViewSymbols, "s", false, "View symbols in csv format")
	flag.BoolVar(&flagFetchShareholding, "fetch-holding", false, "Fetch company's 7 yrs shareholding from BSE Servers")
	flag.BoolVar(&flagFetchShareholding, "H", false, "Fetch company's 7 yrs shareholding from BSE Servers")
	flag.BoolVar(&flagViewShareHolding, "view-holding", false, "View Shareholding in csv format")
	flag.BoolVar(&flagViewShareHolding, "h", false, "View Shareholding in csv format")

	_, data_dir := GetConfig()
	symbols_data_dir := path.Join(data_dir, "symbols")
	os.MkdirAll(symbols_data_dir, 0o700)
	db_filepath := path.Join(symbols_data_dir, "data.db")
	db := GetDB(db_filepath)

	flag.Usage = func() { fmt.Print(usage) }
	flag.Parse()

	if flagFetchSymbols {
		FetchSymbols(db)
		return
	} else if flagViewSymbols {
		mappings := GetSymbolsMappingDB(db)
		err := gocsv.MarshalFile(mappings, os.Stdout)
		CheckErr(err)
		return
	} else if flagFetchShareholding {
	} else if flagViewShareHolding {
	}

	fmt.Println("Error! No arguments specified.")
	fmt.Print(usage)
	os.Exit(1)
}

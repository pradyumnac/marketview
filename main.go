package main

import (
	"encoding/csv"
	"flag"
	"log"
	"os"
	"path"

	"github.com/gocarina/gocsv"
)

// const usage = `Usage:
//   -v, --verbose verbose output
//   -h, --help prints help information
// `

// // for testing, comment out when not needed
// func main() {
// 	holdings := FetchRecentShareholdings("500209", 28, nil) // nill db will skip db step
// 	size := fmt.Sprintf("%d", len(holdings.Holdings))
// 	StructToJson(holdings, "testdata/500209-"+size+"-14012023.json")
// 	fmt.Printf("%s records written", size)
// }

func main() {
	// fmt.Println("##### Marketview v0.1 #####")

	var flagFetchSymbols, flagViewSymbols bool
	var flagFetchShareholding, flagViewShareHolding string
	flag.BoolVar(&flagFetchSymbols, "fetch-sym", false, "Fetch symbols from BSE/NSE Servers")
	flag.BoolVar(&flagFetchSymbols, "S", false, "Fetch symbols from BSE/NSE Servers")
	flag.BoolVar(&flagViewSymbols, "view-sym", false, "View symbols in csv format")
	flag.BoolVar(&flagViewSymbols, "s", false, "View symbols in csv format")
	flag.StringVar(&flagFetchShareholding, "fetch-holding", "", "Fetch company's 7 yrs shareholding from BSE Servers")
	flag.StringVar(&flagFetchShareholding, "H", "", "Fetch company's 7 yrs shareholding from BSE Servers")
	flag.StringVar(&flagViewShareHolding, "view-holding", "", "View Shareholding in csv format")
	flag.StringVar(&flagViewShareHolding, "h", "", "View Shareholding in csv format")
	// flag.Usage = func() { fmt.Print(usage) }
	flag.Parse()

	// General Config
	_, data_dir := GetConfig()
	// symbolsDataDir := path.Join(data_dir, "symbols")
	// os.MkdirAll(symbolsDataDir, 0o700)
	dbFilepath := path.Join(data_dir, "data.db")
	db := GetDB(dbFilepath)

	// // ############## Testing code here ##################
	// holdings := ViewRecentShareholdingsDb("500209", db)
	// log.Print(holdings)
	// os.Exit(1)
	// // ############# Testing Code Ends ###################

	// Controller logic

	if flagFetchSymbols {
		FetchSymbols(db)
		return
	} else if flagViewSymbols {

		mappings := GetSymbolsMappingDB(db)
		err := gocsv.MarshalCSVWithoutHeaders(mappings, csv.NewWriter(os.Stdout))
		CheckErr(err)
		return
	} else if len(flagFetchShareholding) > 0 {
		log.Printf("Fetching script id: %s", flagFetchShareholding)
		// os.Exit(1)
		FetchRecentShareholdings(flagFetchShareholding, 28, db)
		return
	} else if len(flagViewShareHolding) > 0 {
		// log.Println("Not Implemented yet")
		// os.Exit(1)
		holdings := ViewRecentShareholdingsDb(flagViewShareHolding, db)
		// log.Print(holdings)
		csv.NewWriter(os.Stdout).WriteAll(holdings)
		return
	}
	// No arguments matched
	log.Println("Error! No arguments specified.")
	flag.Usage()
	os.Exit(1)
}

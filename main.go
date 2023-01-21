package main

import (
	"fmt"
	"os"
	"path"
)

const usage = `Usage:
  -v, --verbose verbose output
  -h, --help prints help information 
`

// // for testing, comment out when not needed
// func main() {
// 	holdings := FetchRecentShareholdings("500209", 28, nil) // nill db will skip db step
// 	size := fmt.Sprintf("%d", len(holdings.Holdings))
// 	StructToJson(holdings, "testdata/500209-"+size+"-14012023.json")
// 	fmt.Printf("%s records written", size)
// }

func main() {
	fmt.Println("##### Marketview v0.1 #####")
	// 	var flagFetchSymbols, flagViewSymbols, flagFetchShareholding, flagViewShareHolding bool
	// 	flag.BoolVar(&flagFetchSymbols, "fetch-sym", false, "Fetch symbols from BSE/NSE Servers")
	// 	flag.BoolVar(&flagFetchSymbols, "S", false, "Fetch symbols from BSE/NSE Servers")
	// 	flag.BoolVar(&flagViewSymbols, "view-sym", false, "View symbols in csv format")
	// 	flag.BoolVar(&flagViewSymbols, "s", false, "View symbols in csv format")
	// 	flag.BoolVar(&flagFetchShareholding, "fetch-holding", false, "Fetch company's 7 yrs shareholding from BSE Servers")
	// 	flag.BoolVar(&flagFetchShareholding, "H", false, "Fetch company's 7 yrs shareholding from BSE Servers")
	// 	flag.BoolVar(&flagViewShareHolding, "view-holding", false, "View Shareholding in csv format")
	// 	flag.BoolVar(&flagViewShareHolding, "h", false, "View Shareholding in csv format")
	// 	flag.Usage = func() { fmt.Print(usage) }
	// 	flag.Parse()
	// 	// General Config
	_, data_dir := GetConfig()
	symbolsDataDir := path.Join(data_dir, "symbols")
	os.MkdirAll(symbolsDataDir, 0o700)
	dbFilepath := path.Join(symbolsDataDir, "data.db")
	db := GetDB(dbFilepath)
	// ############## Testing code here ##################
	bseScripId := "500209"
	noOfQtrs := 2
	FetchRecentShareholdings(bseScripId, noOfQtrs, db)
	// ############# Testing Code Ends ###################
	// // Controller logc
	//
	//	if flagFetchSymbols {
	//		FetchSymbols(db)
	//		return
	//	} else if flagViewSymbols {
	//
	//		mappings := GetSymbolsMappingDB(db)
	//		err := gocsv.MarshalFile(mappings, os.Stdout)
	//		CheckErr(err)
	//		return
	//	} else if flagFetchShareholding {
	//
	// } else if flagViewShareHolding {
	// }
	// // No arguments matched
	// fmt.Println("Error! No arguments specified.")
	// fmt.Print(usage)
	// os.Exit(1)
}

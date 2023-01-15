package main

import "fmt"

const usage = `Usage:
  -v, --verbose verbose output
  -h, --help prints help information 
`

// for testing, comment out when not needed
func main() {
	holdings := FetchRecentShareholdings("500209", 28, nil) // nill db will skip db step
	size := fmt.Sprintf("%d", len(holdings.Holdings))
	StructToJson(holdings, "testdata/500209-"+size+"-14012023.json")
	fmt.Printf("%s records written", size)
}

//func main() {
//	fmt.Println("##### Marketview v0.1 #####")
//	// 	var flagFetchSymbols, flagViewSymbols, flagFetchShareholding, flagViewShareHolding bool
//	// 	flag.BoolVar(&flagFetchSymbols, "fetch-sym", false, "Fetch symbols from BSE/NSE Servers")
//	// 	flag.BoolVar(&flagFetchSymbols, "S", false, "Fetch symbols from BSE/NSE Servers")
//	// 	flag.BoolVar(&flagViewSymbols, "view-sym", false, "View symbols in csv format")
//	// 	flag.BoolVar(&flagViewSymbols, "s", false, "View symbols in csv format")
//	// 	flag.BoolVar(&flagFetchShareholding, "fetch-holding", false, "Fetch company's 7 yrs shareholding from BSE Servers")
//	// 	flag.BoolVar(&flagFetchShareholding, "H", false, "Fetch company's 7 yrs shareholding from BSE Servers")
//	// 	flag.BoolVar(&flagViewShareHolding, "view-holding", false, "View Shareholding in csv format")
//	// 	flag.BoolVar(&flagViewShareHolding, "h", false, "View Shareholding in csv format")
//	// 	flag.Usage = func() { fmt.Print(usage) }
//	// 	flag.Parse()
//	// 	// General Config
//	_, data_dir := GetConfig()
//	symbols_data_dir := path.Join(data_dir, "symbols")
//	os.MkdirAll(symbols_data_dir, 0o700)
//	db_filepath := path.Join(symbols_data_dir, "data.db")
//	db := GetDB(db_filepath)
//	// ############## Testing code here ##################
//	bse_scrip_id := "500209"
//	noOfQtrs := 28
//	FetchRecentShareholdings(bse_scrip_id, noOfQtrs, db)
//	// ############# Testing Code Ends ###################
//	// // Controller logc
//	//
//	//	if flagFetchSymbols {
//	//		FetchSymbols(db)
//	//		return
//	//	} else if flagViewSymbols {
//	//
//	//		mappings := GetSymbolsMappingDB(db)
//	//		err := gocsv.MarshalFile(mappings, os.Stdout)
//	//		CheckErr(err)
//	//		return
//	//	} else if flagFetchShareholding {
//	//
//	// } else if flagViewShareHolding {
//	// }
//	// // No arguments matched
//	// fmt.Println("Error! No arguments specified.")
//	// fmt.Print(usage)
//	// os.Exit(1)
//}

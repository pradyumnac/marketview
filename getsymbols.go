package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gocarina/gocsv"
	"gorm.io/gorm"
)

// User agent string sent with headers for performing requests
var (
	_USERAGENT_STRING = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36"
)

var (
	CONFIG_DIR string
	DATA_DIR   string
)

func FetchRes(url string) []byte {
	spaceClient := http.Client{
		Timeout: time.Second * 5,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", _USERAGENT_STRING)

	res, getErr := spaceClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	resp, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	// response in []byte
	return resp
}

func FetchBseSymbols() []BseSymbol {
	url := "https://api.bseindia.com/BseIndiaAPI/api/ListofScripData/w?Group=&Scripcode=&industry=&segment=Equity&status=Active"

	resp := FetchRes(url)
	var scrips []BseSymbol
	jsonErr := json.Unmarshal(resp, &scrips)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return scrips
}

func FetchNseSymbols() []NseSymbol {
	url := "https://www1.nseindia.com/content/equities/EQUITY_L.csv"

	resp := FetchRes(url)
	var scrips []NseSymbol
	if err := gocsv.UnmarshalBytes(resp, &scrips); err != nil { // Load clients from file
		panic(err)
	}
	return scrips
}

// Fetches both bse && nse symbols from respective servers
// Builds mapping of NSE/BSE data using ISIN ( BuildBseNseSymbolMaps )
// Saves this data ascsv file in symbols_data_dir
func FetchSymbols(db *gorm.DB) []SymbolsMapping {
	symbols_bse := FetchBseSymbols()
	// bse_symbols_filepath := path.Join(symbols_data_dir, "bse.csv")
	// SaveBseSymbolstoCsv(symbols_bse, bse_symbols_filepath)
	SaveBseSymbols(symbols_bse, db)

	symbols_nse := FetchNseSymbols()
	// nse_symbols_filepath := path.Join(symbols_data_dir, "nse.csv")
	// SaveNseSymbolstoCsv(symbols_nse, nse_symbols_filepath)
	SaveNseSymbols(symbols_nse, db)

	// build mappings
	mappings := BuildBseNseSymbolMaps(db)
	SaveMappings(mappings, db)

	return mappings
}

// Returns all BSE<->NSE symbo mapings based on isin from db
func GetSymbolsMappingDB(db *gorm.DB) []SymbolsMapping {
	var mappings []SymbolsMapping

	// Connect to database
	if err := db.Find(&mappings).Error; err != nil {
		log.Fatal("Unable to fetch mappings.")
	}

	return mappings
}

// Build a mapping of nse and bse symbols using iisin_number
// For records where either is missing, the corresponding column is marked as empty
// Pass this data to Save SymbolsMappng function to save to db
// To use this data from db, call GetSymbolsMapping()
func BuildBseNseSymbolMaps(db *gorm.DB) []SymbolsMapping {
	// a list of isins for which maping s build
	isin_visited := make(map[string]bool)
	var mappings []SymbolsMapping

	var symbols_bse []BseSymbol
	if err := db.Find(&symbols_bse).Error; err != nil {
		log.Fatal("Unable to fetch BSE symbols from DB.")
	}

	var symbols_nse []NseSymbol
	if err := db.Find(&symbols_nse).Error; err != nil {
		log.Fatal("Unable to fetch mappings.")
	}

	for _, symbol := range symbols_bse {
		var nsesymbol NseSymbol
		if _, exist := isin_visited[symbol.ISIN]; !exist {
			var nsecd string
			if err := db.First(&nsesymbol, "isin= ?", symbol.ISIN).Error; err == nil {
				// record found
				nsecd = nsesymbol.ScripCd
			} else {
				// record missing
				nsecd = ""
			}
			mapping := SymbolsMapping{
				ISIN:      symbol.ISIN,
				ScripName: symbol.ScripName,
				BseCd:     symbol.ScripCd,
				BseId:     symbol.ScripId,
				NseCd:     nsecd,
				Industry:  symbol.Industry,
				Group:     symbol.Group,
			}
			mappings = append(mappings, mapping)
		}
		isin_visited[symbol.ISIN] = true
	}

	for _, symbol := range symbols_nse {
		var bsesymbol BseSymbol
		if _, exist := isin_visited[symbol.ISIN]; !exist {
			// exists in nse but not in bse
			var bsecd, bseid, industry, group string
			if err := db.First(&bsesymbol, "isin= ?", symbol.ISIN).Error; err == nil {
				// record found
				bsecd = bsesymbol.ScripCd
				bseid = bsesymbol.ScripId
				industry = bsesymbol.Industry
				group = bsesymbol.Group
			} else {
				// record missing
				bsecd = ""
				bseid = ""
				industry = ""
				group = ""
			}
			mapping := SymbolsMapping{
				ISIN:      symbol.ISIN,
				ScripName: symbol.ScripName,
				BseCd:     bsecd,
				BseId:     bseid,
				NseCd:     symbol.ScripCd,
				Industry:  industry,
				Group:     group,
			}
			mappings = append(mappings, mapping)
		}
		isin_visited[symbol.ISIN] = true
	}

	return mappings
}

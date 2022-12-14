package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"time"

	"github.com/gocarina/gocsv"
)

// User agent string sent with headers for performing requests
var (
	_USERAGENT_STRING = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36"
)

var (
	CONFIG_DIR string
	DATA_DIR   string
)

type BseSymbol struct {
	SCRIP_CD    string
	Scrip_Name  string
	Status      string
	GROUP       string
	FACE_VALUE  string
	ISIN_NUMBER string
	INDUSTRY    string
	Scrip_id    string
	Segment     string
	NSURL       string
	Issuer_Name string
	Mktcap      string
}

type NseSymbol struct {
	SCRIP_CD    string `csv:"SYMBOL"`
	Scrip_Name  string `csv:"NAME OF COMPANY"`
	Status      string `csv:""`
	GROUP       string `csv:""`
	FACE_VALUE  string `csv:""`
	ISIN_NUMBER string `csv:"ISIN NUMBER"`
	INDUSTRY    string `csv:""`
	Scrip_id    string `csv:""`
	Segment     string `csv:""`
	NSURL       string `csv:""`
	Issuer_Name string `csv:""`
	Mktcap      string `csv:""`
}

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
// Saves this data ascsv file in symbols_data_dir
func getSymbols(symbols_data_dir string) ([]BseSymbol, []NseSymbol) {
	symbols_bse := FetchBseSymbols()
	bse_symbols_filepath := path.Join(symbols_data_dir, "bse.csv")
	SaveBseSymbolstoCsv(symbols_bse, bse_symbols_filepath)

	symbols_nse := FetchNseSymbols()
	nse_symbols_filepath := path.Join(symbols_data_dir, "nse.csv")
	SaveNseSymbolstoCsv(symbols_nse, nse_symbols_filepath)
	return symbols_bse, symbols_nse
}

package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
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

type BseScrip struct {
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

type NseScrip struct {
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

func FetchBSE() []BseScrip {
	url := "https://api.bseindia.com/BseIndiaAPI/api/ListofScripData/w?Group=&Scripcode=&industry=&segment=Equity&status=Active"

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

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	var scrips []BseScrip
	jsonErr := json.Unmarshal([]byte(body), &scrips)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return scrips
}

func FetchNSE() []NseScrip {
	url := "https://www1.nseindia.com/content/equities/EQUITY_L.csv"

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

	// body, readErr := ioutil.ReadAll(res.Body)
	// if readErr != nil {
	// 	log.Fatal(readErr)
	// }

	var scrips []NseScrip
	if err := gocsv.Unmarshal(res.Body, &scrips); err != nil { // Load clients from file
		panic(err)
	}

	return scrips
}

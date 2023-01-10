package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// User agent string sent with headers for performing requests
var (
	_USERAGENT_STRING = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36"
)

var (
	CONFIG_DIR string
	DATA_DIR   string
)

type Scrip struct {
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

func GetConfig() (string, string) {
	USER_HOME_DIR, _ := os.UserHomeDir()
	USER_CONFIG_DIR := path.Join(USER_HOME_DIR, ".config")
	USER_DATA_DIR := path.Join(USER_HOME_DIR, ".local", "share")

	CONFIG_DIR = path.Join(USER_CONFIG_DIR, "marketview")
	DATA_DIR := path.Join(USER_DATA_DIR, "marketview")

	os.MkdirAll(filepath.Dir(CONFIG_DIR), 0o700)
	os.MkdirAll(filepath.Dir(DATA_DIR), 0o700)

	return CONFIG_DIR, DATA_DIR
}

func FetchBSE() []Scrip {
	url := "https://api.bseindia.com/BseIndiaAPI/api/ListofScripData/w?Group=&Scripcode=&industry=&segment=Equity&status=Active"

	spaceClient := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
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

	var scrips []Scrip
	jsonErr := json.Unmarshal([]byte(body), &scrips)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return scrips
}

func FetchNSE() []Scrip {
	url := "https://www1.nseindia.com/content/equities/EQUITY_L.csv"

	spaceClient := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
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

	var scrips []Scrip
	jsonErr := json.Unmarshal([]byte(body), &scrips)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return scrips
}

func StructToCSV(scrips []Scrip, csvFilePath string) {
	f, err := os.Create(csvFilePath)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	// _, err = f.WriteString("SCRIP_CD, Scrip_Name, Status, GROUP, FACE_VALUE, ISIN_NUMBER, INDUSTRY, Scrip_id, Segment, NSURL, Issuer_Name,  Mktcap\r\n")
	_, err = f.WriteString("Scrip_id, SCRIP_CD, ISIN_NUMBER, Scrip_Name, NSURL, INDUSTRY, GROUP, FACE_VALUE,  Issuer_Name,  Mktcap\r\n")
	if err != nil {
		log.Fatal(err)
	}

	for _, sym := range scrips {
		f.WriteString(sym.Scrip_id + ", " + sym.SCRIP_CD + ", " + sym.ISIN_NUMBER + ", " + sym.Scrip_Name + ", " + sym.NSURL + ", " + sym.INDUSTRY + ", " + sym.GROUP + ", " + sym.FACE_VALUE + ", " + strings.Replace(sym.Issuer_Name, ",", "", -1) + ", " + sym.Mktcap + "\r\n")
		if err != nil {
			log.Fatal(err)
		}
	}

	f.Sync()
}

func SaveCsv(csvData []byte, csvFilePath string) {
	f, err := os.Create(csvFilePath)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	_, err = f.Write(csvData)
	f.Sync()
}

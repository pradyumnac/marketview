package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/dnlo/struct2csv"
)

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

func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Convert any Struct (JSON Compatible) to JSON
// TODO: remove ny and implement type constraints
func JsonToCsv(structForJson any, jsonFilePath string) {
	jsonStr, err := json.MarshalIndent(structForJson, "", "  ")
	CheckErr(err)

	f, err := os.Create(jsonFilePath)
	CheckErr(err)
	defer f.Close()

	if _, err := f.Write(jsonStr); err != nil {
		log.Fatalf("Unable to write json data tofile")
	}
}

// Convert any Struct Compatible type to csv
// TODO: remove ny and implement type constraints
func StructToCsv(structForCsv any, csvFilePath string) {
	buff := &bytes.Buffer{}
	w := struct2csv.NewWriter(buff)
	err := w.WriteStructs(structForCsv)
	if err != nil {
		log.Fatalf("Unable to create csv for NSE Symbols")
	}
	SaveBytesCsv(buff.Bytes(), csvFilePath)
}

// Save Structure for NSe symbols to csv file
func SaveNseSymbolstoCsv(symbols []NseSymbol, csvFilePath string) {
	buff := &bytes.Buffer{}
	w := struct2csv.NewWriter(buff)
	err := w.WriteStructs(symbols)
	if err != nil {
		log.Fatalf("Unable to create csv for NSE Symbols")
	}
	SaveBytesCsv(buff.Bytes(), csvFilePath)
}

// Saves Structure for BSE Symbols to CSV
func SaveBseSymbolstoCsv(symbols []BseSymbol, csvFilePath string) {
	buff := &bytes.Buffer{}
	w := struct2csv.NewWriter(buff)
	err := w.WriteStructs(symbols)
	if err != nil {
		log.Fatalf("Unable to create csv for BSE Symbols")
	}
	SaveBytesCsv(buff.Bytes(), csvFilePath)
}

// Saves csv data in bytes to csv file
func SaveBytesCsv(csvAsBytes []byte, csvFilePath string) {
	f, err := os.Create(csvFilePath)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	_, err = f.Write(csvAsBytes)
	f.Sync()
}

// FOR BSESHAREHOLDING SERVERS:Returns the BSE qtrid for last 7 years. Latest qtr will be the last completed qtr as of today
func getLatestQtrId() int {
	fy, qtr := getYyQtr()
	return fy*5 + qtr - 2
}

// Returns Qtrstring based on curent datetime
func getLatestQtrString() string {
	fy, qtr := getYyQtr()
	return fmt.Sprintf("FY%dQ%d", fy, qtr)
}

// FOR BSESHAREHOLDING SERVERS: Qtr String- "FY20Q1" to qtrid - 116 as accepted by bse shareholding page
func getQtridFromQtrstring(qtr_string string) int {
	var year, qtr int
	fmt.Sscanf(qtr_string, "FY%dQ%d", &year, &qtr)

	if year < 1 && year > 99 {
		log.Fatalf("Invalid year: %d", year)
	}

	if qtr < 1 && qtr > 4 {
		log.Fatalf("Invalid qtr: %d", qtr)
	}

	// back calculated formula
	return 5*year + qtr - 2
}

// Get current financial year and completed quarter
// qtrToId := map[string]int{"FY16Q3": 88, "FY16Q4": 89, "FY17Q1": 90, "FY17Q2": 91, "FY17Q3": 92, "FY17Q4": 93, "FY18Q1": 94, "FY18Q2": 95, "FY18Q3": 96, "FY18Q4": 97, "FY19Q1": 98, "FY19Q2": 99, "FY19Q3": 100, "FY19Q4": 101, "FY20Q1": 102, "FY20Q2": 103, "FY20Q3": 104, "FY20Q4": 105, "FY21Q1": 106, "FY21Q2": 107, "FY21Q3": 108, "FY21Q4": 109, "FY22Q1": 110, "FY22Q2": 111, "FY22Q3": 112, "FY22Q4": 113, "FY23Q1": 114, "FY23Q2": 115, "FY23Q3": 116, "FY23Q4": 117, "FY24Q4": 118, "FY24Q1": 119, "FY24Q2": 120, "FY24Q3": 121}
func getYyQtr() (int, int) {
	today := time.Now()
	year := int(today.Year()) % 100 // 2 digit year
	month := int(today.Month())
	var qtr int

	// dally a date
	switch month {
	case 1, 2, 3:
		// For dec, FY(thisyr)Q3
		qtr = 3
	case 4, 5, 6:
		qtr = 4
	case 7, 8, 9:
		qtr = 1
		year += 1
	case 10, 11, 12:
		qtr = 2
		year += 1
	}

	return year, qtr
}

// FOR BSESHAREHOLDING SERVERS: returns a slice of qtr ids ( last numOfQtrs )
func getLastNQtrids(lastQtrId int, noOfQtrs int) []int {
	firstQtrId := lastQtrId - noOfQtrs + 1
	r := make([]int, lastQtrId-firstQtrId+1)
	for i := range r {
		r[i] = firstQtrId + i
	}
	return r
}

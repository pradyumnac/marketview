package main

import (
	"bytes"
	"log"
	"os"
	"path"
	"path/filepath"

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

func SaveNseSymbolstoCsv(symbols []NseSymbol, csvFilePath string) {
	buff := &bytes.Buffer{}
	w := struct2csv.NewWriter(buff)
	err := w.WriteStructs(symbols)
	if err != nil {
		log.Fatalf("Unable to create csv for NSE Symbols")
	}
	SaveBytesCsv(buff.Bytes(), csvFilePath)
}

func SaveBseSymbolstoCsv(symbols []BseSymbol, csvFilePath string) {
	buff := &bytes.Buffer{}
	w := struct2csv.NewWriter(buff)
	err := w.WriteStructs(symbols)
	if err != nil {
		log.Fatalf("Unable to create csv for BSE Symbols")
	}
	SaveBytesCsv(buff.Bytes(), csvFilePath)
}

func SaveBytesCsv(csvAsBytes []byte, csvFilePath string) {
	f, err := os.Create(csvFilePath)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	_, err = f.Write(csvAsBytes)
	f.Sync()
}

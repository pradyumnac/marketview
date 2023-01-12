package main

import (
	"os"
	"path"

	"github.com/gocarina/gocsv"
)

func main() {
	_, data_dir := GetConfig()
	symbols_data_dir := path.Join(data_dir, "symbols")
	os.MkdirAll(symbols_data_dir, 0o700)
	db_filepath := path.Join(symbols_data_dir, "data.db")
	db := GetDB(db_filepath)

	// FetchSymbols(db)

	mappings := GetSymbolsMappingDB(db)
	err := gocsv.MarshalFile(mappings, os.Stdout)
	CheckErr(err)
}

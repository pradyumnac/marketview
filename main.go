package main

import (
	"os"
	"path"
)

func main() {
	_, data_dir := GetConfig()
	symbols_dir := path.Join(data_dir, "symbols")

	scrips := FetchBSE()
	os.MkdirAll(symbols_dir, 0o700)
	bse_symbols_filepath := path.Join(symbols_dir, "bse.csv")
	StructToCSV(scrips, bse_symbols_filepath)

	scrips = FetchNSE()
	nse_symbols_filepath := path.Join(symbols_dir, "bse.csv")
	StructToCSV(scrips, nse_symbols_filepath)
}

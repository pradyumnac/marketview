package main

import (
	"path"
)

func main() {
	_, data_dir := GetConfig()
	symbols_dir := path.Join(data_dir, "symbols")

	// bse_scrips := FetchBSE()
	// os.MkdirAll(symbols_dir, 0o700)
	// bse_symbols_filepath := path.Join(symbols_dir, "bse.csv")
	// StructToCSV(bse_scrips, bse_symbols_filepath)

	nse_scrips := FetchNSE()
	nse_symbols_filepath := path.Join(symbols_dir, "bse.csv")
	StructToCSV(nse_scrips, nse_symbols_filepath)
}

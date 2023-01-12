package main

import (
	"os"
	"path"
)

func main() {
	_, data_dir := GetConfig()
	symbols_data_dir := path.Join(data_dir, "symbols")
	os.MkdirAll(symbols_data_dir, 0o700)

	getSymbols(symbols_data_dir)

}

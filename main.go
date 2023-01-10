package main

func main() {
	scrips := FetchBSE()
	StructToCSV(scrips, "../data/symbols/bse.csv")

	// scrips = FetchNSE()
	// StructToCSV(scrips, "../data/symbols/bse.csv")
}

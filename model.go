package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// This table contains the comapny name, bse and nse code mapped using isin_number
type SymbolsMapping struct {
	gorm.Model
	NseCd       string
	BseCd       string
	BseId       string
	ISIN        string
	ScripName   string
	Industry    string
	Group       string
	MarketCap   string
	DateListing string // will help in identifying new additions
}

type BseSymbol struct {
	gorm.Model
	ScripCd    string
	ScripName  string
	Status     string
	Group      string
	FaceValue  string
	ISIN       string
	Industry   string
	ScripId    string
	Segment    string
	NsUrl      string
	IssuerName string
	Mktcap     string
}

type NseSymbol struct {
	gorm.Model
	ScripCd     string `csv:"SYMBOL"`
	ScripName   string `csv:"NAME OF COMPANY"`
	Series      string `csv:"SERIES"`
	DateListing string `csv:"DATE OF LISTING"`
	PaidupValue string `csv:"PAID UP VALUE"`
	MarketLot   string `csv:"MARKET_LOT"`
	ISIN        string `csv:"ISIN NUMBER"`
	FaceValue   string `csv:"FACE VALUE"`
}

// Saves a slice of NseSymbol s to DB
func SaveBseSymbols(symbols []BseSymbol, db *gorm.DB) {
	// Delete all entries
	db.Model(&BseSymbol{}).Delete(&BseSymbol{})

	// Insert data afresh
	db.CreateInBatches(symbols, 1000)
}

// Saves a slice of NseSymbol s to DB
func SaveNseSymbols(symbols []NseSymbol, db *gorm.DB) {
	// Deleteall entries
	db.Model(&NseSymbol{}).Delete(&NseSymbol{})

	// Insert data a fresh
	db.CreateInBatches(symbols, 1000)
}

func GetDB(db_path string) *gorm.DB {
	// Connect to database
	db, err := gorm.Open(sqlite.Open(db_path), &gorm.Config{})
	CheckErr(err)

	// Migrate the schema
	db.AutoMigrate(&NseSymbol{}, &BseSymbol{}, &SymbolsMapping{})

	return db
}

// Returns a record set of all BSE<->NSE symbo mapings based on isin from db
func GetSymbolsMapping(symbols_bse []BseSymbol, symbols_nse []NseSymbol) []SymbolsMapping {

	var mappings []SymbolsMapping

	return mappings
}

// Buld a mapping of nse and bse symbols using iisin_number
// For records where either is missing, the corresponding column is marked as empty
// Pass this data to Save SymbolsMappng function to save to db
// To use this data from db, call GetSymbolsMapping()
func BuildBseNseSymbolMaps(symbols_bse []BseSymbol, symbols_nse []NseSymbol, db *gorm.DB) []SymbolsMapping {
	// a list of isins for which maping s build
	var isin_visited map[string]bool
	var mappings []SymbolsMapping

	for _, symbol := range symbols_bse {
		var nsesymbol NseSymbol
		if _, exist := isin_visited[symbol.ISIN]; !exist {
			var nsecd, datelisting string
			if err := db.First(&nsesymbol, "isin= ?", symbol.ISIN); err == nil {
				// record found
				nsecd = nsesymbol.ScripCd
				datelisting = nsesymbol.DateListing
			} else {
				// record missing
				nsecd = ""
				datelisting = ""
			}
			mapping := SymbolsMapping{
				ISIN:        symbol.ISIN,
				ScripName:   symbol.ScripName,
				BseCd:       symbol.ScripCd,
				BseId:       symbol.ScripId,
				NseCd:       nsecd,
				Industry:    symbol.Industry,
				Group:       symbol.Group,
				MarketCap:   symbol.Mktcap,
				DateListing: datelisting,
			}
			mappings = append(mappings, mapping)
		}
		isin_visited[symbol.ISIN] = true
	}

	for _, symbol := range symbols_nse {
		var bsesymbol BseSymbol
		if _, exist := isin_visited[symbol.ISIN]; !exist {

			//exists in nse but not in bse
			var bsecd, bseid, industry, group, marketcap string
			if err := db.First(&bsesymbol, "isin= ?", symbol.ISIN); err == nil {
				// record found
				bsecd = bsesymbol.ScripCd
				bseid = bsesymbol.ScripId
				industry = bsesymbol.Industry
				group = bsesymbol.Group
				marketcap = bsesymbol.Mktcap
			} else {
				// record missing
				bsecd = ""
				bseid = ""
				industry = ""
				group = ""
				marketcap = ""
			}
			mapping := SymbolsMapping{
				ISIN:      symbol.ISIN,
				ScripName: symbol.ScripName,
				BseCd:     bsecd,
				BseId:     bseid,
				NseCd:     symbol.ScripCd,
				Industry:  industry,
				Group:     group,
				MarketCap: marketcap,
			}
			mappings = append(mappings, mapping)
		}
		isin_visited[symbol.ISIN] = true
	}

	return mappings
}

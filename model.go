package main

import (
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

// This table contains the company name, bse and nse code mapped using isin_number
type SymbolsMapping struct {
	NseCd     string
	BseCd     string
	BseId     string
	ISIN      string `gorm:"primaryKey"`
	ScripName string
	Industry  string
	Group     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type BseSymbol struct {
	ScripCd    string `json:"scrip_cd"`
	ScripName  string `json:"scrip_name"`
	Status     string `json:"status"`
	Group      string `json:"group"`
	FaceValue  string `json:"face_value"`
	ISIN       string `gorm:"primaryKey" json:"isin_number"`
	Industry   string `json:"industry"`
	ScripId    string `json:"scrip_id"`
	Segment    string `json:"segment"`
	NsUrl      string `json:"nsurl"`
	IssuerName string `json:"issuer_name"`
	MktCap     string `json:"mktcap"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type NseSymbol struct {
	ScripCd     string `csv:"SYMBOL"`
	ScripName   string `csv:"NAME OF COMPANY"`
	Series      string `csv:"SERIES"`
	DateListing string `csv:"DATE OF LISTING"`
	PaidupValue string `csv:"PAID UP VALUE"`
	MarketLot   string `csv:"MARKET_LOT"`
	ISIN        string `csv:"ISIN NUMBER" gorm:"primaryKey"`
	FaceValue   string `csv:"FACE VALUE"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Saves a slice of NseSymbol s to DB
func SaveBseSymbols(symbols []BseSymbol, db *gorm.DB) {
	// Delete all entries
	// db.Model(&BseSymbol{}).Delete(&BseSymbol{})

	// Insert data afresh
	db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "isin"}}, // key colume
		DoNothing: true,
	}).CreateInBatches(symbols, 1000)
}

// Saves a slice of NseSymbol s to DB
func SaveNseSymbols(symbols []NseSymbol, db *gorm.DB) {
	// Deleteall entries
	// db.Model(&NseSymbol{}).Delete(&NseSymbol{})

	// Insert data a fresh
	db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "isin"}}, // key colume
		DoNothing: true,
	}).CreateInBatches(&symbols, 1000)
}

// Saves a slice of NseSymbol s to DB
func SaveMappings(mappings []SymbolsMapping, db *gorm.DB) {
	// Deleteall entries
	db.Model(&SymbolsMapping{}).Delete(&SymbolsMapping{})

	// Insert data a fresh
	db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "isin"}}, // key colume
		DoNothing: true,
	}).CreateInBatches(mappings, 1000)
}

func GetDB(db_path string) *gorm.DB {
	// Connect to database
	db, err := gorm.Open(sqlite.Open(db_path), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
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

// Build a mapping of nse and bse symbols using iisin_number
// For records where either is missing, the corresponding column is marked as empty
// Pass this data to Save SymbolsMappng function to save to db
// To use this data from db, call GetSymbolsMapping()
func BuildBseNseSymbolMaps(symbols_bse []BseSymbol, symbols_nse []NseSymbol, db *gorm.DB) []SymbolsMapping {
	// a list of isins for which maping s build
	isin_visited := make(map[string]bool)
	var mappings []SymbolsMapping

	for _, symbol := range symbols_bse {
		var nsesymbol NseSymbol
		if _, exist := isin_visited[symbol.ISIN]; !exist {
			var nsecd string
			if err := db.First(&nsesymbol, "isin= ?", symbol.ISIN).Error; err == nil {
				// record found
				nsecd = nsesymbol.ScripCd
			} else {
				// record missing
				nsecd = ""
			}
			mapping := SymbolsMapping{
				ISIN:      symbol.ISIN,
				ScripName: symbol.ScripName,
				BseCd:     symbol.ScripCd,
				BseId:     symbol.ScripId,
				NseCd:     nsecd,
				Industry:  symbol.Industry,
				Group:     symbol.Group,
			}
			mappings = append(mappings, mapping)
		}
		isin_visited[symbol.ISIN] = true
	}

	for _, symbol := range symbols_nse {
		var bsesymbol BseSymbol
		if _, exist := isin_visited[symbol.ISIN]; !exist {
			// exists in nse but not in bse
			var bsecd, bseid, industry, group string
			if err := db.First(&bsesymbol, "isin= ?", symbol.ISIN).Error; err == nil {
				// record found
				bsecd = bsesymbol.ScripCd
				bseid = bsesymbol.ScripId
				industry = bsesymbol.Industry
				group = bsesymbol.Group
			} else {
				// record missing
				bsecd = ""
				bseid = ""
				industry = ""
				group = ""
			}
			mapping := SymbolsMapping{
				ISIN:      symbol.ISIN,
				ScripName: symbol.ScripName,
				BseCd:     bsecd,
				BseId:     bseid,
				NseCd:     symbol.ScripCd,
				Industry:  industry,
				Group:     group,
			}
			mappings = append(mappings, mapping)
		}
		isin_visited[symbol.ISIN] = true
	}

	return mappings
}

package main

import (
	"log"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

// These Struct can be saved as csv
type CsvCompatibleType interface {
	ShareholdingLineItem
}

// These Struct can be saved as json
type JsonCompatibleType interface {
	ShareholdingQtr | ShareholdingLineItem | Shareholdings
}

// ##################################### Generic ####################################

// FB factory
func GetDB(dbPath string) *gorm.DB {
	// Connect to database
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	CheckErr(err)

	// Migrate the schema
	db.AutoMigrate(&NseSymbol{}, &BseSymbol{}, &SymbolsMapping{}, &Shareholdings{}, &ShareholdingQtr{}, &ShareholdingLineItem{})

	return db
}

// ###################################### Symbols ####################################

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

// Saves the bse/nse mapping struct to DB
func SaveMappings(mappings []SymbolsMapping, db *gorm.DB) {
	// Deleteall entries
	db.Model(&SymbolsMapping{}).Delete(&SymbolsMapping{})

	// Insert data a fresh
	result := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "isin"}}, // key colume
		DoNothing: true,
	}).CreateInBatches(mappings, 1000)

	log.Printf("Added %d mappings", result.RowsAffected)
}

// ####################### ShareHolding #################################

// this stores a single line item of shareholding data
type ShareholdingLineItem struct {
	gorm.Model
	TypeCd            string `json:"type_cd"`
	TypeName          string `json:"type_name"`
	QtrId             string `json:"qtrid"`
	BseScripId        string `json:"bse_scrip_id"`
	CategoryName      string `json:"category"`
	HolderCount       string `json:"holder_count"`
	NoOfShares        string `json:"no_of_shares"`
	PctHolding        string `json:"pct_holding"`
	ShareholdingQtrID uint
}

// This struct can store a company's shareholding for a quarter
type ShareholdingQtr struct {
	gorm.Model
	BseScripId       string                 `json:"bse_scrip_id"`
	QtrId            string                 `json:"qtrid"`
	OverviewHoldings []ShareholdingLineItem `json:"overview"`
	PublicHoldings   []ShareholdingLineItem `json:"public"`
	DiiHoldings      []ShareholdingLineItem `json:"dii"`
	FiiHoldings      []ShareholdingLineItem `json:"fii"`
	PromoterHoldings []ShareholdingLineItem `json:"promoter"`
	ShareholdingsID  uint
}

// This struct stores shareholdings of a company across quarters
type Shareholdings struct {
	gorm.Model
	BseScripId string            `gorm:"unique" json:"bse_scrip_id"`
	Holdings   []ShareholdingQtr `json:"holdings"`
}

// Gets a company's shareholdings from database
func GetRecentShareholdingsDb(holdings Shareholdings, db *gorm.DB) {
}

// Stores a slice of Shareholding struct to database
func SaveRecentShareholdingsDb(holdings Shareholdings, db *gorm.DB) {
	// Insert data a fresh
	result := db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&holdings)

	if result.Error != nil {
		log.Fatalf("Error: Unable to save record to database: %s", result.Error)
	} else {
		log.Printf("Added %d shareholdings for %s", result.RowsAffected, holdings.BseScripId)
	}
}

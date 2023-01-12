package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type BseSymbol struct {
	gorm.Model
	SCRIP_CD    string
	Scrip_Name  string
	Status      string
	GROUP       string
	FACE_VALUE  string
	ISIN_NUMBER string
	INDUSTRY    string
	Scrip_id    string
	Segment     string
	NSURL       string
	Issuer_Name string
	Mktcap      string
}

type NseSymbol struct {
	gorm.Model
	Symbol       string `csv:"SYMBOL"`
	Nanme        string `csv:"NAME OF COMPANY"`
	Series       string `csv:"SERIES"`
	Date_listing string `csv:"DATE OF LISTING"`
	Paidup_value string `csv:"PAID UP VALUE"`
	Market_lot   string `csv:"MARKET_LOT"`
	Isin_number  string `csv:"ISIN NUMBER"`
	Face_value   string `csv:"FACE VALUE"`
}

// Saves a slice of NseSymbol s to DB
func SaveBseSymbols(symbols []BseSymbol, db *gorm.DB) {
	db.CreateInBatches(symbols, 1000)
}

// Saves a slice of NseSymbol s to DB
func SaveNseSymbols(symbols []NseSymbol, db *gorm.DB) {
	db.CreateInBatches(symbols, 1000)
}

func GetDB(db_path string) *gorm.DB {
	// Connect to database
	db, err := gorm.Open(sqlite.Open(db_path), &gorm.Config{})
	CheckErr(err)

	// Migrate the schema
	db.AutoMigrate(&NseSymbol{}, &BseSymbol{})

	return db
}

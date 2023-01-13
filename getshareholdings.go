package main

import (
	"fmt"
	"log"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type ShareholdingCategory struct {
	CategoryName string
	HolderCount  string
	NoOfShares   string
	PctHolding   string
}

type Shareholding struct {
	categories []ShareholdingCategory
}

// Gets QTR STring based on curent datetime
// Warning: This function return the last completed qtr
// (( For whch results/documents will be avalable
func getLatestQtrString() string {
	today := time.Now()
	year := int(today.Year()) % 100 // 2 digit year
	month := int(today.Month())
	var qtr int

	// dally a date
	switch month {
	case 1, 2, 3:
		// For dec, FY(thisyr)Q3
		qtr = 3
	case 4, 5, 6:
		qtr = 4
	case 7, 8, 9:
		qtr = 1
		year += 1
	case 10, 11, 12:
		qtr = 2
		year += 1
	}

	return fmt.Sprintf("FY%dQ%d", year, qtr)
}

// Qtr String- "FY20Q1" to qtrid - 116 as accepted by bse shareholding page
func getShareholdingQtrId(qtr_string string) int {
	// qtrToId := map[string]int{"FY16Q3": 88, "FY16Q4": 89, "FY17Q1": 90, "FY17Q2": 91, "FY17Q3": 92, "FY17Q4": 93, "FY18Q1": 94, "FY18Q2": 95, "FY18Q3": 96, "FY18Q4": 97, "FY19Q1": 98, "FY19Q2": 99, "FY19Q3": 100, "FY19Q4": 101, "FY20Q1": 102, "FY20Q2": 103, "FY20Q3": 104, "FY20Q4": 105, "FY21Q1": 106, "FY21Q2": 107, "FY21Q3": 108, "FY21Q4": 109, "FY22Q1": 110, "FY22Q2": 111, "FY22Q3": 112, "FY22Q4": 113, "FY23Q1": 114, "FY23Q2": 115, "FY23Q3": 116, "FY23Q4": 117, "FY24Q4": 118, "FY24Q1": 119, "FY24Q2": 120, "FY24Q3": 121}
	// return qtrToId[qtr_string]

	var year, qtr int
	fmt.Sscanf(qtr_string, "FY%dQ%d", &year, &qtr)

	if year < 1 && year > 99 {
		log.Fatalf("Invalid year: %d", year)
	}

	if qtr < 1 && qtr > 4 {
		log.Fatalf("Invalid qtr: %d", qtr)
	}

	// back calculated formula
	return 5*year + qtr - 2
}

// parse shareholding category page
func ParseCategory(doc *goquery.Document) []ShareholdingCategory {
	var categories []ShareholdingCategory
	// Find the review items
	doc.Find("#tdData > table > tbody > tr:nth-child(5) > td > table > tr").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		categories = append(categories, ShareholdingCategory{
			CategoryName: s.Find("td.TTRow_left").Text(),
			HolderCount:  s.Find("td.TTRow_Right:nth-child(1)").Text(),
			NoOfShares:   s.Find("td.TTRow_Right:nth-child(2)").Text(),
			PctHolding:   s.Find("td.TTRow_Right:nth-child(3)").Text(),
		})
	})

	return categories
}

// get shareholding data for a bse listed company
// earliest available  qtr string for infy : 88: FY16Q3 - for dec2015
func getShareholding(bse_scrip_id string, qtr_string string) Shareholding {
	// https://www.bseindia.com/corporates/shpSecurities.aspx?scripcd=500209&qtrid=115.00&Flag=New

	holding := Shareholding{}

	qtrid := getShareholdingQtrId(qtr_string)
	url_string := "https://www.bseindia.com/corporates/shpSecurities.aspx?scripcd=%s&qtrid=%d.00&Flag=New"
	url := fmt.Sprintf(url_string, bse_scrip_id, qtrid)

	res := FetchResBody(url)
	defer res.Close()

	doc, err := goquery.NewDocumentFromReader(res)
	if err != nil {
		log.Fatal(err)
	}
	holding.categories = ParseCategory(doc)

	return holding
}

// returns sharehokding for the co for latest published qtr
func getLatestShareholding(bse_scrip_id string) Shareholding {
	return getShareholding(bse_scrip_id, getLatestQtrString())
}

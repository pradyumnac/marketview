package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

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

// Get shareholding breakup for Non Promoters
func getPublicShareholding(bse_scrip_id string, qtrid int, qtr_string string) ([]ShareholdingLineItem, []ShareholdingLineItem, []ShareholdingLineItem) {
	var public_holdings []ShareholdingLineItem
	var dii_holdings []ShareholdingLineItem
	var fii_holdings []ShareholdingLineItem

	// qtrid := getShareholdingQtrId(qtr_string)
	url_string := "https://www.bseindia.com/corporates/shpPublicShareholder.aspx?scripcd=%s&qtrid=%d"
	url := fmt.Sprintf(url_string, bse_scrip_id, qtrid)

	res := FetchResBody(url)
	defer res.Close()

	doc, err := goquery.NewDocumentFromReader(res)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("#tdData > table > tbody > tr:nth-child(3) > td > table > tbody > tr").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		left_column := s.Find("td.TTRow_left")
		if left_column.Length() > 0 {
			right_column := s.Find("td.TTRow_right")
			typename := "Public" // TODO: may be mf, DII, FII, public
			no_of_shares := strings.Replace(right_column.Eq(1).Text(), ",", "", -1)
			if len(no_of_shares) > 0 {
				public_holdings = append(public_holdings, ShareholdingLineItem{
					TypeCd:       "1",
					TypeName:     typename,
					BseScripId:   bse_scrip_id,
					Qtr:          qtr_string,
					CategoryName: left_column.Text(),
					HolderCount:  strings.Replace(right_column.Eq(0).Text(), ",", "", -1),
					NoOfShares:   no_of_shares,
					PctHolding:   strings.Replace(right_column.Eq(3).Text(), ",", "", -1),
				})
			}
		}
	})
	return public_holdings, dii_holdings, fii_holdings
}

// Get Promoter & Promoter group shareholding breakup
func getPromoterShareholding(bse_scrip_id string, qtrid int, qtr_string string) []ShareholdingLineItem {
	var promoter_holdings []ShareholdingLineItem

	// qtrid := getShareholdingQtrId(qtr_string)
	url_string := "https://www.bseindia.com/corporates/shpPromoterNGroup.aspx?scripcd=%s&qtrid=%d"
	url := fmt.Sprintf(url_string, bse_scrip_id, qtrid)

	res := FetchResBody(url)
	defer res.Close()

	doc, err := goquery.NewDocumentFromReader(res)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("#tdData > table > tbody > tr:nth-child(3) > td > table > tbody > tr").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		left_column := s.Find("td.TTRow_left")
		if left_column.Length() > 0 {
			right_column := s.Find("td.TTRow_right")
			typename := "Promoter"
			no_of_shares := strings.Replace(right_column.Eq(1).Text(), ",", "", -1)
			if len(no_of_shares) > 0 {
				promoter_holdings = append(promoter_holdings, ShareholdingLineItem{
					TypeCd:       "1",
					TypeName:     typename,
					BseScripId:   bse_scrip_id,
					Qtr:          qtr_string,
					CategoryName: left_column.Text(),
					HolderCount:  strings.Replace(right_column.Eq(0).Text(), ",", "", -1),
					NoOfShares:   no_of_shares,
					PctHolding:   strings.Replace(right_column.Eq(3).Text(), ",", "", -1),
				})
			}
		}
	})
	return promoter_holdings
}

// parse shareholding category page
func ParseCategory(bse_scrip_id string, qtr_string string, doc *goquery.Document) []ShareholdingLineItem {
	var categories []ShareholdingLineItem
	// Find the review items
	doc.Find("#tdData > table > tbody > tr:nth-child(5) > td > table > tbody > tr").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title

		left_column := s.Find("td.TTRow_left")
		if left_column.Length() > 0 {
			right_column := s.Find("td.TTRow_right")
			typename := "overview"
			no_of_shares := strings.Replace(right_column.Eq(1).Text(), ",", "", -1)
			if len(no_of_shares) > 0 {
				categories = append(categories, ShareholdingLineItem{
					TypeCd:       "1",
					TypeName:     typename,
					BseScripId:   bse_scrip_id,
					Qtr:          qtr_string,
					CategoryName: left_column.Text(),
					HolderCount:  strings.Replace(right_column.Eq(0).Text(), ",", "", -1),
					NoOfShares:   no_of_shares,
					PctHolding:   strings.Replace(right_column.Eq(3).Text(), ",", "", -1),
				})
			}
		}
	})
	return categories
}

// get shareholding data for a bse listed company
// earliest available  qtr string for infy : 88: FY16Q3 - for dec2015
func getShareholdingQtr(bse_scrip_id string, qtr_string string) ShareholdingQtr {
	// https://www.bseindia.com/corporates/shpSecurities.aspx?scripcd=500209&qtrid=115.00&Flag=New

	holdingQtr := ShareholdingQtr{}

	qtrid := getShareholdingQtrId(qtr_string)
	url_string := "https://www.bseindia.com/corporates/shpSecurities.aspx?scripcd=%s&qtrid=%d.00&Flag=New"
	url := fmt.Sprintf(url_string, bse_scrip_id, qtrid)

	res := FetchResBody(url)
	defer res.Close()

	doc, err := goquery.NewDocumentFromReader(res)
	if err != nil {
		log.Fatal(err)
	}
	overview_holdings := ParseCategory(bse_scrip_id, qtr_string, doc)
	promoter_holdings := getPromoterShareholding(bse_scrip_id, qtrid, qtr_string)
	public_holdings, dii_holdings, fii_holdings := getPublicShareholding(bse_scrip_id, qtrid, qtr_string)

	// log.Println(overview_holdings)
	// log.Println(promoter_holdings)
	// log.Println(public_holdings)

	holdingQtr.QtrString = qtr_string
	holdingQtr.BseScripId = bse_scrip_id
	holdingQtr.OverviewHoldings = overview_holdings
	holdingQtr.PromoterHoldings = promoter_holdings
	holdingQtr.PublicHoldings = public_holdings
	holdingQtr.DiiHoldings = dii_holdings
	holdingQtr.FiiHoldings = fii_holdings

	// Do some parsing to get  a single holding structure
	return holdingQtr
}

// returns sharehokding for the co for latest published qtr
func getLatestShareholding(bse_scrip_id string) ShareholdingQtr {
	return getShareholdingQtr(bse_scrip_id, getLatestQtrString())
}

func GetRecentShareholdings(bse_scrip_id string) ShareHoldings {
	companyShareHoldings := make(ShareHoldings)

	return companyShareHoldings
}

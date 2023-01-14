package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Get shareholding breakup for Non Promoters
func getPublicShareholding(bse_scrip_id string, qtrid_string string) ([]ShareholdingLineItem, []ShareholdingLineItem, []ShareholdingLineItem) {
	var public_holdings []ShareholdingLineItem
	var dii_holdings []ShareholdingLineItem
	var fii_holdings []ShareholdingLineItem

	// qtrid := getShareholdingQtrId(qtr_string)
	url_string := "https://www.bseindia.com/corporates/shpPublicShareholder.aspx?scripcd=%s&qtrid=%d"
	url := fmt.Sprintf(url_string, bse_scrip_id, qtrid_string)

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
			no_of_shares := strings.Trim(strings.Replace(right_column.Eq(1).Text(), ",", "", -1), " ")
			pct_holding := strings.Trim(strings.Replace(right_column.Eq(5).Text(), ",", "", -1), " ")
			holder_count := strings.Trim(strings.Replace(right_column.Eq(0).Text(), ",", "", -1), " ")
			if len(no_of_shares) > 0 {
				public_holdings = append(public_holdings, ShareholdingLineItem{
					TypeCd:       "2",
					TypeName:     typename,
					BseScripId:   bse_scrip_id,
					QtrId:        qtrid_string,
					CategoryName: left_column.Text(),
					HolderCount:  holder_count,
					NoOfShares:   no_of_shares,
					PctHolding:   pct_holding,
				})
			}
		}
	})
	return public_holdings, dii_holdings, fii_holdings
}

// Get Promoter & Promoter group shareholding breakup
func getPromoterShareholding(bse_scrip_id string, qtrid_string string) []ShareholdingLineItem {
	var promoter_holdings []ShareholdingLineItem

	// qtrid := getShareholdingQtrId(qtr_string)
	url_string := "https://www.bseindia.com/corporates/shpPromoterNGroup.aspx?scripcd=%s&qtrid=%s"
	url := fmt.Sprintf(url_string, bse_scrip_id, qtrid_string)

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
			no_of_shares := strings.Trim(strings.Replace(right_column.Eq(1).Text(), ",", "", -1), " ")
			pct_holding := strings.Trim(strings.Replace(right_column.Eq(3).Text(), ",", "", -1), " ")
			holder_count := strings.Trim(strings.Replace(right_column.Eq(0).Text(), ",", "", -1), " ")
			if len(no_of_shares) > 0 {
				promoter_holdings = append(promoter_holdings, ShareholdingLineItem{
					TypeCd:       "2",
					TypeName:     typename,
					BseScripId:   bse_scrip_id,
					QtrId:        qtrid_string,
					CategoryName: left_column.Text(),
					HolderCount:  holder_count,
					NoOfShares:   no_of_shares,
					PctHolding:   pct_holding,
				})
			}
		}
	})
	return promoter_holdings
}

// parse shareholding category page
func ParseOverview(bse_scrip_id string, qtrid_string string, doc *goquery.Document) []ShareholdingLineItem {
	var categories []ShareholdingLineItem
	// Find the review items
	doc.Find("#tdData > table > tbody > tr:nth-child(5) > td > table > tbody > tr").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title

		left_column := s.Find("td.TTRow_left")
		if left_column.Length() > 0 {
			right_column := s.Find("td.TTRow_right")
			typename := "overview"
			no_of_shares := strings.Trim(strings.Replace(right_column.Eq(1).Text(), ",", "", -1), " ")
			pct_holding := strings.Trim(strings.Replace(right_column.Eq(3).Text(), ",", "", -1), " ")
			holder_count := strings.Trim(strings.Replace(right_column.Eq(0).Text(), ",", "", -1), " ")
			if len(no_of_shares) > 0 {
				categories = append(categories, ShareholdingLineItem{
					TypeCd:       "1",
					TypeName:     typename,
					BseScripId:   bse_scrip_id,
					QtrId:        qtrid_string,
					CategoryName: left_column.Text(),
					HolderCount:  holder_count,
					NoOfShares:   no_of_shares,
					PctHolding:   pct_holding,
				})
			}
		}
	})
	return categories
}

// get shareholding data for a bse listed company
// earliest available  qtr string for infy : 88: FY16Q3 - for dec2015
func getShareholdingQtr(bse_scrip_id string, qtrid_string string) ShareholdingQtr {
	// https://www.bseindia.com/corporates/shpSecurities.aspx?scripcd=500209&qtrid=115.00&Flag=New

	holdingQtr := ShareholdingQtr{}

	// qtrid := getShareholdingQtrId(qtr_string)
	url_string := "https://www.bseindia.com/corporates/shpSecurities.aspx?scripcd=%s&qtrid=%s.00&Flag=New"
	url := fmt.Sprintf(url_string, bse_scrip_id, qtrid_string)

	res := FetchResBody(url)
	defer res.Close()

	doc, err := goquery.NewDocumentFromReader(res)
	if err != nil {
		log.Fatal(err)
	}

	// Initialized struct

	holdingQtr.QtrId = qtrid_string
	holdingQtr.BseScripId = bse_scrip_id

	// check if shareholding data is avilable for the quarter
	if doc.Find("#tdData > table > tbody > tr:nth-child(5) > td > table > tbody > tr").Length() > 0 {
		overview_holdings := ParseOverview(bse_scrip_id, qtrid_string, doc)
		promoter_holdings := getPromoterShareholding(bse_scrip_id, qtrid_string)
		public_holdings, dii_holdings, fii_holdings := getPublicShareholding(bse_scrip_id, qtrid_string)

		holdingQtr.OverviewHoldings = overview_holdings
		holdingQtr.PromoterHoldings = promoter_holdings
		holdingQtr.PublicHoldings = public_holdings
		holdingQtr.DiiHoldings = dii_holdings
		holdingQtr.FiiHoldings = fii_holdings
	}

	// Do some parsing to get  a single holding structure
	return holdingQtr
}

// returns sharehokding for the co for latest published qtr
func getLatestShareholding(bse_scrip_id string) ShareholdingQtr {
	qtrid_string := fmt.Sprintf("%d", getLatestQtrId())
	return getShareholdingQtr(bse_scrip_id, qtrid_string)
}

// Get the comapny share holding data for last 7 years
func GetRecentShareholdings(bse_scrip_id string, noOfQtrs int) ShareHoldings {
	companyShareHoldings := make(ShareHoldings)

	qtrIds := getLastNQtrids(getLatestQtrId(), noOfQtrs)
	for qtrid := range qtrIds {
		// apply wait group here
		qtrid_string := fmt.Sprintf("%d", qtrIds[qtrid])
		companyShareHoldings[qtrid_string] = getShareholdingQtr(bse_scrip_id, qtrid_string)
	}

	// waitgroup wait

	return companyShareHoldings
}

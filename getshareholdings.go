package main

import (
	"fmt"
	"time"
)

// MMM-YYYY to qtrid accedped by bse shareholding page
func getQtrId(qtr string) int {
	firstDataDate, _ := time.Parse("Jan 2006", "Dec 2015")
	qtrDate, _ := time.Parse("Jan 2006", qtr)
	dateDiff := int(qtrDate.Sub(firstDataDate).Hours() / 24 / 91)

	// 88: dec 2015 - date with first data in bse
	// before this, some scrips have data (like acc)
	// but page format different
	return dateDiff + 88
}

// get shareholding data for a bse listed company
// bse_scrip_id: 6 digit bse numerical cod
// qtr(string): format - MMM-YYYY
// earliest avaikable ffor infy : 88: dec-2015
func getShareholding(bse_scrip_id string, qtr string) {
	// https://www.bseindia.com/corporates/shpSecurities.aspx?scripcd=500209&qtrid=115.00&Flag=New

	qtrid := getQtrId(qtr)
	url_string := "https://www.bseindia.com/corporates/shpSecurities.aspx?scripcd=%s&qtrid=%d.00&Flag=New"
	url := fmt.Sprintf(url_string, bse_scrip_id, qtrid)

	response := FetchRes(url)
}

// returns sharehokding for the co for latest published qtr
func getLatestShareholding(bse_scrip_id string) {
}

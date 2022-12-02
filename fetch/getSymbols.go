package symbols

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func FetchBSE() []Scrip {
	url := "https://api.bseindia.com/BseIndiaAPI/api/ListofScripData/w?Group=&Scripcode=&industry=&segment=Equity&status=Active"

	spaceClient := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")

	res, getErr := spaceClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	// fmt.Println(string(body))

	// scrips := []scrip{}
	var scrips []Scrip
	jsonErr := json.Unmarshal([]byte(body), &scrips)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return scrips
}

func FetchNSE() []byte {
	url := "https://www1.nseindia.com/content/equities/EQUITY_L.csv"

	spaceClient := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")

	res, getErr := spaceClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	return body
}

// func main() {
// 	scrips := fetchBSE()
// 	structToCSV(scrips, "../data/symbols/bse.csv")

// 	csvString := fetchNSE()
// 	saveCsv(csvString, "../data/symbols/bse.csv")
// }

package symbols

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func fetchBSE(csvFilePath string) {
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
	var scrips []scrip
	jsonErr := json.Unmarshal([]byte(body), &scrips)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	structToCSV(scrips, csvFilePath)
}

func fetchNSE(csvFilePath string) {
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

	f, err := os.Create(csvFilePath)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	_, err = f.Write(body)
	f.Sync()
}

func main() {
	fetchBSE("../data/symbols/bse.csv")
	fetchNSE("../data/symbols/nse.csv")
}

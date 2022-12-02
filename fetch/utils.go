package symbols

import (
	"log"
	"os"
	"strings"
)

type Scrip struct {
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

func structToCSV(scrips []Scrip, csvFilePath string) {
	f, err := os.Create(csvFilePath)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	// _, err = f.WriteString("SCRIP_CD, Scrip_Name, Status, GROUP, FACE_VALUE, ISIN_NUMBER, INDUSTRY, Scrip_id, Segment, NSURL, Issuer_Name,  Mktcap\r\n")
	_, err = f.WriteString("Scrip_id, SCRIP_CD, ISIN_NUMBER, Scrip_Name, NSURL, INDUSTRY, GROUP, FACE_VALUE,  Issuer_Name,  Mktcap\r\n")
	if err != nil {
		log.Fatal(err)
	}

	for _, sym := range scrips {
		f.WriteString(sym.Scrip_id + ", " + sym.SCRIP_CD + ", " + sym.ISIN_NUMBER + ", " + sym.Scrip_Name + ", " + sym.NSURL + ", " + sym.INDUSTRY + ", " + sym.GROUP + ", " + sym.FACE_VALUE + ", " + strings.Replace(sym.Issuer_Name, ",", "", -1) + ", " + sym.Mktcap + "\r\n")
		if err != nil {
			log.Fatal(err)
		}
	}

	f.Sync()
}

func saveCsv(csvData []byte, csvFilePath string) {
	f, err := os.Create(csvFilePath)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	_, err = f.Write(csvData)
	f.Sync()
}

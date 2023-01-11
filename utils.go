package main

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func GetConfig() (string, string) {
	USER_HOME_DIR, _ := os.UserHomeDir()
	USER_CONFIG_DIR := path.Join(USER_HOME_DIR, ".config")
	USER_DATA_DIR := path.Join(USER_HOME_DIR, ".local", "share")

	CONFIG_DIR = path.Join(USER_CONFIG_DIR, "marketview")
	DATA_DIR := path.Join(USER_DATA_DIR, "marketview")

	os.MkdirAll(filepath.Dir(CONFIG_DIR), 0o700)
	os.MkdirAll(filepath.Dir(DATA_DIR), 0o700)

	return CONFIG_DIR, DATA_DIR
}

func StructToCSV(scrips []NseScrip, csvFilePath string) {
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

func SaveCsv(csvData []byte, csvFilePath string) {
	f, err := os.Create(csvFilePath)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	_, err = f.Write(csvData)
	f.Sync()
}

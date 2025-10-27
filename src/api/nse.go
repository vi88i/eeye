package api

import (
	"archive/zip"
	"bytes"
	"eeye/src/constants"
	"eeye/src/models"
	"eeye/src/utils"
	"fmt"
	"log"
	"time"

	"github.com/gocarina/gocsv"
)

// Downloads the zip file in memory
// Extracts the CSV file (only file in zip) and unmarshals it into a slice of NSEStockData
func getBhavcopyData(zipFileName string) ([]models.NSEStockData, error) {
	empty := utils.EmptySlice[models.NSEStockData]()

	// download zip file
	log.Printf("getBhavcopyData: Trying to fetch %v", zipFileName)
	resp, err := NseClient.
		R().
		Get(zipFileName)
	if err != nil {
		return empty, fmt.Errorf("getBhavcopyData: %w", err)
	}

	// Get a reader for zip file
	r, err := zip.NewReader(bytes.NewReader(resp.Body()), int64(len(resp.Body())))
	if err != nil {
		return empty, fmt.Errorf("getBhavcopyData: %w", err)
	}

	// Get the first (and only) file in the zip
	if len(r.File) == 0 {
		return empty, fmt.Errorf("getBhavcopyData: zip file is empty")
	}

	rc, err := r.File[0].Open()
	if err != nil {
		return empty, fmt.Errorf("getBhavcopyData: %w", err)
	}
	defer func() {
		_ = rc.Close()
	}()

	stocks := make([]models.NSEStockData, 0, constants.NumOfStocks)
	if err := gocsv.Unmarshal(rc, &stocks); err != nil {
		return empty, fmt.Errorf("getBhavcopyData: %w", err)
	}

	return stocks, nil
}

// DownloadLatestBhavcopy attempts to download the latest NSE Bhavcopy zip file
// and extract the bhavcopy data CSV. It tries up to 10 previous days
// if the latest file is not available (e.g., weekends, holidays).
func DownloadLatestBhavcopy() ([]models.NSEStockData, string, error) {
	const (
		Probes = 10
	)

	var (
		now   = time.Now().Local()
		i     = 0
		empty = utils.EmptySlice[models.NSEStockData]()
	)
	for i < Probes {
		var (
			// Format: BhavCopy_NSE_CM_0_0_0_YYYYMMDD_F_0000.csv.zip
			dateStr        = now.Format("20060102")
			zipFileName    = fmt.Sprintf("BhavCopy_NSE_CM_0_0_0_%s_F_0000.csv.zip", dateStr)
			lastTradingDay = now.Format("2006-01-02")
		)

		stocks, err := getBhavcopyData(zipFileName)
		if err == nil && len(stocks) > 0 {
			log.Printf("DownloadLatestBhavcopy: Last trading day is %v", lastTradingDay)
			return stocks, lastTradingDay, nil
		}

		log.Printf("DownloadLatestBhavcopy: %v\n", err)
		now = now.AddDate(0, 0, -1)
		i++
	}

	return empty, "", fmt.Errorf("DownloadLatestBhavcopy: no valid data found")
}

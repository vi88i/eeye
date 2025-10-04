package api

import (
	"archive/zip"
	"bytes"
	"eeye/src/constants"
	"eeye/src/models"
	"eeye/src/utils"
	"fmt"
	"log"
	"path"
	"time"

	"github.com/gocarina/gocsv"
)

// Downloads the zip file in memory
// Extracts the target CSV file and unmarshals it into a slice of NSEStockData
func getMarketCapData(zipFileName, marketCapDataFileName string) ([]models.NSEStockData, error) {
	empty := utils.EmptySlice[models.NSEStockData]()

	// download zip file
	resp, err := NseClient.
		R().
		Get(zipFileName)
	if err != nil {
		return empty, fmt.Errorf("getMarketCapData: %w", err)
	}

	// Get a reader for zip file
	r, err := zip.NewReader(bytes.NewReader(resp.Body()), int64(len(resp.Body())))
	if err != nil {
		return empty, fmt.Errorf("getMarketCapData: %w", err)
	}

	// traverse the zip file and find target file
	for _, f := range r.File {
		if !f.FileInfo().IsDir() && path.Base(f.Name) == marketCapDataFileName {
			rc, err := f.Open()
			if err != nil {
				return nil, fmt.Errorf("getMarketCapData: open file: %w", err)
			}
			defer func() {
				_ = rc.Close()
			}()

			stocks := make([]models.NSEStockData, 0, constants.NumOfStocks)
			if err := gocsv.Unmarshal(rc, &stocks); err != nil {
				return empty, fmt.Errorf("getMarketCapData: %w", err)
			}

			return stocks, nil
		}
	}

	return empty, nil
}

// DownloadLatestBhavcopy attempts to download the latest NSE Bhavcopy zip file
// and extract the market capitalization data CSV. It tries up to 10 previous days
// if the latest file is not available (e.g., weekends, holidays).
func DownloadLatestBhavcopy() ([]models.NSEStockData, string, error) {
	const (
		PadFmt = "%02d"
		Probes = 10
	)

	var (
		now   = time.Now().Local()
		i     = 0
		empty = utils.EmptySlice[models.NSEStockData]()
	)
	for i < Probes {
		var (
			year                  = fmt.Sprintf(PadFmt, now.Year()%2000)
			month                 = fmt.Sprintf(PadFmt, now.Month())
			day                   = fmt.Sprintf(PadFmt, now.Day())
			zipFileName           = fmt.Sprintf("PR%v%v%v.zip", day, month, year)
			marketCapDataFileName = fmt.Sprintf("MCAP%v%v%v.csv", day, month, now.Year())
			lastTradingDay        = fmt.Sprintf("%v-%v-%v", now.Year(), month, day)
		)

		stocks, err := getMarketCapData(zipFileName, marketCapDataFileName)
		if err == nil && len(stocks) > 0 {
			return stocks, lastTradingDay, nil
		}

		log.Printf("DownloadLatestBhavcopy: %v\n", err)
		now = now.AddDate(0, 0, -1)
		i++
	}

	return empty, "", fmt.Errorf("DownloadLatestBhavcopy: no valid data found")
}

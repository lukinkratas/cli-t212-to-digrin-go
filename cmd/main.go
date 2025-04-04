package main

import (
	"fmt"
	"time"
	"slices"
	"os"

	"github.com/joho/godotenv"
	utils "github.com/lukinkratas/cli-t212-to-digrin-go/internal/utils"
)


func GetInputDt() string {

	var currentDt time.Time = time.Now()
	var previousMonthDt time.Time = currentDt.AddDate(0, -1, 0)
	var previousMonthDtStr string = previousMonthDt.Format("2006-01")

	var inputDtStr string
	fmt.Println("Reporting Year Month in \"YYYY-mm\" format: ")
	fmt.Printf("Or confirm default \"%v\" by ENTER.\n", previousMonthDtStr)
	fmt.Scanln(&inputDtStr)

	if inputDtStr == "" {
		inputDtStr = previousMonthDtStr
	}

	return inputDtStr

}


func main() {

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	var bucketName string = os.Getenv("BUCKET_NAME")

	var inputDtStr string = GetInputDt()

	var inputDt time.Time
	inputDt, err = time.Parse("2006-01", inputDtStr)
	if err != nil {
		panic(err)
	}

	var fromDt time.Time = utils.GetFirstDayOfMonth(inputDt)
	var toDt time.Time = utils.GetFirstDayOfNextMonth(inputDt)

	fmt.Printf("  fromDt: %v\n", fromDt)
	fmt.Printf("  toDt: %v\n", toDt)

	var createdReportId int

	for {
		
		createdReportId = utils.CreateExport(fromDt, toDt)

		if createdReportId != 0 {
			break
		}

		time.Sleep(10 * time.Second)

	}
	
	// // createdReportId Mock Up
	// createdReportId = 1594033

	fmt.Printf("  createdReportId: %v\n", createdReportId)

	var downloadLink string

	for {

		var reportsList []utils.Report = utils.FetchReports()

		// report list is empty
		if len(reportsList) == 0 {
			time.Sleep(10 * time.Second)
			continue
		}
		
		// if report list is not empty
		var createdReport utils.Report

		// reverse order for loop, cause latest export is expected to be at the end
		slices.Reverse(reportsList)

		for _, report := range reportsList {

			if report.Id == createdReportId {
		        createdReport = report
				break
			}

		}

		if createdReport.Status == "Finished" {
			downloadLink = createdReport.DownloadLink
			break
		}

	}

	// // downloadLink Mock Up
	// downloadLink = "https://tzswiy3zk5dms05cfeo.s3.eu-central-1.amazonaws.com/from_2025-03-01_to_2025-04-01_MTc0MzU4MDY0MDE0Mw.csv?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Date=20250402T075723Z&X-Amz-SignedHeaders=host&X-Amz-Expires=604799&X-Amz-Credential=AKIARJCCZCDEKCUWYOXG%2F20250402%2Feu-central-1%2Fs3%2Faws4_request&X-Amz-Signature=857a3b30cb532fdc0d52137a8af7602cbdfd84f597de0c74f61727403c71be3c"

	fmt.Printf("  downloadLink: %v\n", downloadLink)

	var t212CsvEncoded []byte = utils.DownloadReport(downloadLink)

	var fileName string = fmt.Sprintf("%s.csv", inputDtStr)

	var keyName string = fmt.Sprintf("t212/%s", fileName)
	utils.S3PutObject(t212CsvEncoded, bucketName, keyName)

	var csvRecords []utils.CsvRecord = utils.ReadCsv(t212CsvEncoded)

	csvRecords = utils.TransformCsv(csvRecords)

	utils.SaveCsv(csvRecords, fileName)

	var digrinCsvEncoded []byte = utils.EncodeCsv(csvRecords)
	keyName = fmt.Sprintf("digrin/%s", fileName)
	utils.S3PutObject(digrinCsvEncoded, bucketName, keyName)

}

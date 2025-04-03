package main

import (
	"fmt"
	"time"
	"net/http"
	"io/ioutil"
	"os"
	"bytes"
	"strings"
	// "slices"
	"encoding/json"
	"encoding/csv"

	"github.com/joho/godotenv"
	"github.com/gocarina/gocsv"
	utils "github.com/lukinkratas/cli-t212-to-digrin-go/internal/utils"
)

const bucketName string = "t212-to-digrin-test"

type Payload struct {
	DataIncluded map[string]bool `json:"dataIncluded"`
	TimeFrom     string          `json:"timeFrom"` // Use string for RFC3339 formatted time
	TimeTo       string          `json:"timeTo"`
}

type Export struct {
	ReportId int `json:"reportId"`
}

type Report struct {
	Id           int             `json:"reportId"`
	TimeFrom     string          `json:"timeFrom"`
	TimeTo       string          `json:"timeTo"`
	DataIncluded map[string]bool `json:"dataIncluded"`
	Status       string          `json:"status"`
	DownloadLink string          `json:"downloadLink"`
}


type CsvRow struct {
	Action                               string `csv:"Action"`
	Time                                 string `csv:"Time"`
	ISIN                                 string `csv:"ISIN"`
	Ticker                               string `csv:"Ticker"`
	Name                                 string `csv:"Name"`
	Notes                                string `csv:"Notes"`
	Id                                   string `csv:"ID"`
	NoOfShares                           float64 `csv:"No. of shares"`
	PricePerShare                        float64 `csv:"Price / share"`
	CurrencyPricePerShare                string `csv:"Currency (Price / share)"`
	ExchangeRate                         string `csv:"Exchange rate"`
	CurrencyResult                       string `csv:"Currency (Result)"`
	Total                                float64 `csv:"Total"`
	CurrencyTotal                        string `csv:"Currency (Total)"`
	WithholdingTax                       float64 `csv:"Withholding tax"`
	CurrencyWithholdingTax               string `csv:"Currency (Withholding tax)"`
	CurrencyConversionFromAmount         float64 `csv:"Currency conversion from amount"`
	CurrencyCurrencyConversionFromAmount string `csv:"Currency (Currency conversion from amount)"`
	CurrencyConversionToAmount           float64 `csv:"Currency conversion to amount"`
	CurrencyCurrencyConversionToAmount   string `csv:"Currency (Currency conversion to amount)"`
	CurrencyConversionFee                float64 `csv:"Currency conversion fee"`
	CurrencyCurrencyConversionFee        string `csv:"Currency (Currency conversion fee)"`
	FrenchTransactionTax                 float64 `csv:"French transaction tax"`
	CurrencyFrenchTransactionTax         string `csv:"Currency (French transaction tax)"`
}

func GetInputDt() string {

	var currentDt time.Time = time.Now()
	var defaultDt time.Time = currentDt.AddDate(0, -1, 0)
	var defaultDtStr string = defaultDt.Format("2006-01")

	var inputDtStr string
	fmt.Println("Reporting Year Month in \"YYYY-mm\" format: ")
	fmt.Printf("Or confirm default \"%v\" by ENTER.\n", defaultDtStr)
	fmt.Scanln(&inputDtStr)

	if inputDtStr == "" {
		inputDtStr = defaultDtStr
	}

	return inputDtStr

}

func GetFirstDayOfMonth(dt time.Time) time.Time {
	return time.Date(dt.Year(), dt.Month(), 1, 0, 0, 0, 0, time.UTC)
}

func GetFirstDayOfNextMonth(dt time.Time) time.Time {
	var nextMonthDt time.Time = dt.AddDate(0, 1, 0) // works even for Jan and Dec
	return time.Date(nextMonthDt.Year(), nextMonthDt.Month(), 1, 0, 0, 0, 0, time.UTC)
}

func CreateExport(fromDt time.Time, toDt time.Time) int {

	const url string = "https://live.trading212.com/api/v0/history/exports"

	dataIncluded := map[string]bool{
		"includeDividends":    true,
		"includeInterest":     true,
		"includeOrders":       true,
		"includeTransactions": true,
	}
	
	payload := Payload{
		DataIncluded: dataIncluded,
		TimeFrom: fromDt.Format(time.RFC3339),
		TimeTo:   toDt.Format(time.RFC3339),
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		panic(err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", os.Getenv("T212_API_KEY"))

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	
	fmt.Printf("  Response Status: %v\n", response.Status)

	if response.Status != "200 OK" {
		return 0
	}
	
	defer response.Body.Close()
	reponseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	var reponseBody Export
	err = json.Unmarshal(reponseBytes, &reponseBody)
	if err != nil {
		panic(err)
	}

	return reponseBody.ReportId
}

func FetchReports() []Report {

	const url string = "https://live.trading212.com/api/v0/history/exports"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Add("Authorization", os.Getenv("T212_API_KEY"))
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	fmt.Printf("  Response Status: %v\n", response.Status)

	if response.Status != "200 OK" {
		return nil
	}

	defer response.Body.Close()
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	var reponseBody []Report
	err = json.Unmarshal(responseBytes, &reponseBody)
	if err != nil {
		panic(err)
	}

	return reponseBody
}


func DownloadReport(downloadLink string) []byte {

	req, err := http.NewRequest("GET", downloadLink, nil)
	if err != nil {
		panic(err)
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	fmt.Printf("  Response Status: %v\n", response.Status)

	if response.Status != "200 OK" {
		return nil
	}

	defer response.Body.Close()
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	return responseBytes
}


// func Transform(){
// 	// # Read input CSV
//     // report_df = pd.read_csv(StringIO(df_bytes.decode('utf-8')))

//     // # Filter out blacklisted tickers
//     // report_df = report_df[~report_df['Ticker'].isin(TICKER_BLACKLIST)]
//     // report_df = report_df[report_df['Action'].isin(['Market buy', 'Market sell'])]

//     // # Apply the mapping to the ticker column
//     // report_df['Ticker'] = report_df['Ticker'].apply(map_ticker)

//     // # convert dtypes
//     // return report_df.convert_dtypes()
// }

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	var inputDtStr string = GetInputDt()

	// var inputDt time.Time
	// inputDt, err = time.Parse("2006-01", inputDtStr)
	// if err != nil {
	// 	panic(err)
	// }

	// var fromDt time.Time = GetFirstDayOfMonth(inputDt)
	// var toDt time.Time = GetFirstDayOfNextMonth(inputDt)

	// fmt.Printf("  fromDt: %v\n", fromDt)
	// fmt.Printf("  toDt: %v\n", toDt)

	var createdReportId int

	// for {
		
	// 	createdReportId = CreateExport(fromDt, toDt)

	// 	if createdReportId != 0 {
	// 		break
	// 	}

	// 	time.Sleep(30 * time.Second)

	// }
	
	// createdReportId Mock Up
	createdReportId = 1594033

	fmt.Printf("  createdReportId: %v\n", createdReportId)

	var downloadLink string

	// for {

	// 	var reportsList []Report
	// 	reportsList = FetchReports()

	// 	// report list is empty
	// 	if len(reportsList) == 0 {
	// 		time.Sleep(60 * time.Second)
	// 		continue
	// 	}
		
	// 	// if report list is not empty
	// 	var createdReport Report

	// 	// reverse order for loop, cause latest export is expected to be at the end
	// 	for report in slices.Reverse(reportsList) {

	// 		if report.Id == createdReportId {
	// 	        createdReport = report // is this needed?
	// 			break
	// 		}

	// 	}

	// 	if report.Status == "Finished" {
	// 		downloadLink = report.DownloadLink
	// 		break
	// 	}

	// }

	// downloadLink Mock Up
	downloadLink = "https://tzswiy3zk5dms05cfeo.s3.eu-central-1.amazonaws.com/from_2025-03-01_to_2025-04-01_MTc0MzU4MDY0MDE0Mw.csv?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Date=20250402T075723Z&X-Amz-SignedHeaders=host&X-Amz-Expires=604799&X-Amz-Credential=AKIARJCCZCDEKCUWYOXG%2F20250402%2Feu-central-1%2Fs3%2Faws4_request&X-Amz-Signature=857a3b30cb532fdc0d52137a8af7602cbdfd84f597de0c74f61727403c71be3c"

	fmt.Printf("  downloadLink: %v\n", downloadLink)

	var t212Bytes []byte
	t212Bytes = DownloadReport(downloadLink)
	// fmt.Printf("  string(t212Bytes): %v\n", string(t212Bytes))

	var fileName string
	fileName = fmt.Sprintf("%s.csv", inputDtStr)

	var keyName string
	keyName = fmt.Sprintf("t212/%s", fileName)
	utils.S3PutObject(t212Bytes, bucketName, keyName)

	// Read the CSV file into a slice of Record structs
	var csvRows []CsvRow
	err = gocsv.UnmarshalBytes(t212Bytes, &csvRows)
	if err != nil {
		panic(err)
	}
 
	// Print the records
	for idx, csvRow := range csvRows {
		fmt.Printf("  gocsv %v %v\n", idx, csvRow)
	}

	// Write the CSV data locally
	csvFile, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()
 
	writer := csv.NewWriter(csvFile)
	defer writer.Flush()
	
	var rows []string
	rows = strings.Split(string(t212Bytes), "\n")
	for _, row := range rows {
		writer.Write(strings.Split(row, ","))
	}

	keyName = fmt.Sprintf("digrin/%s", fileName)
	utils.S3PutObject(t212Bytes, bucketName, keyName)

}

package main

import (
	"fmt"
	"time"
	"net/http"
	"io/ioutil"
	"bytes"
	"os"
	"encoding/json"

	"github.com/joho/godotenv"
)

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

func GetInputDt() string {

	var CurrentDt time.Time = time.Now()
	var DefaultDt time.Time = CurrentDt.AddDate(0, -1, 0)
	var DefaultDtStr string = DefaultDt.Format("2006-01")

	var InputDtStr string
	fmt.Println("Reporting Year Month in \"YYYY-mm\" format: ")
	fmt.Printf("Or confirm default \"%v\" by ENTER.\n", DefaultDtStr)
	fmt.Scanln(&InputDtStr)

	if InputDtStr == "" {
		InputDtStr = DefaultDtStr
	}

	return InputDtStr

}

func GetFirstDayOfMonth(Dt time.Time) time.Time {
	return time.Date(Dt.Year(), Dt.Month(), 1, 0, 0, 0, 0, time.UTC)
}

func GetFirstDayOfNextMonth(Dt time.Time) time.Time {
	var NextMonthDt time.Time = Dt.AddDate(0, 1, 0) // works even for Jan and Dec
	return time.Date(NextMonthDt.Year(), NextMonthDt.Month(), 1, 0, 0, 0, 0, time.UTC)
}

func CreateExport(FromDt time.Time, ToDt time.Time) int {

	const url string = "https://live.trading212.com/api/v0/history/exports"

	DataIncluded := map[string]bool{
		"includeDividends":    true,
		"includeInterest":     true,
		"includeOrders":       true,
		"includeTransactions": true,
	}
	
	payloadData := Payload{
		DataIncluded: DataIncluded,
		TimeFrom: FromDt.Format(time.RFC3339),
		TimeTo:   ToDt.Format(time.RFC3339),
	}

	payloadBytes, err := json.Marshal(payloadData)
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

	var reponseData Export
	err = json.Unmarshal(reponseBytes, &reponseData)
	if err != nil {
		panic(err)
	}

	return reponseData.ReportId
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

	var reponseData []Report
	err = json.Unmarshal(responseBytes, &reponseData)
	if err != nil {
		panic(err)
	}

	return reponseData
}


func DownloadReport(DownloadLink string) string {

	req, err := http.NewRequest("GET", DownloadLink, nil)
	if err != nil {
		panic(err)
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	fmt.Printf("  Response Status: %v\n", response.Status)

	if response.Status != "200 OK" {
		return ""
	}

	defer response.Body.Close()
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	// var reponseData string
	// err = json.Unmarshal(responseBytes, &reponseData)
	// if err != nil {
	// 	panic(err)
	// }

	return string(responseBytes)
}


func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	var InputDtStr string = GetInputDt()

	var InputDt time.Time
	InputDt, err = time.Parse("2006-01", InputDtStr)
	if err != nil {
		panic(err)
	}

	var FromDt time.Time = GetFirstDayOfMonth(InputDt)
	var ToDt time.Time = GetFirstDayOfNextMonth(InputDt)

	var CreatedReportId int

	// for {
		
	// 	CreatedReportId = CreateExport(FromDt, ToDt)

	// 	if CreatedReportId != 0 {
	// 		break
	// 	}

	// 	time.Sleep(30 * time.Second)

	// }
	
	// CreatedReportId Mock Up
	CreatedReportId = 1594033

	fmt.Printf("  CreatedReportId: %v\n", CreatedReportId)

	var DownloadLink string

	// for {

	// 	var ReportsList []Report
	// 	ReportsList = FetchReports()

	// 	// report list is empty
	// 	if len(ReportsList) == 0 {
	// 		time.Sleep(60 * time.Second)
	// 		continue
	// 	}
		
	// 	// if report list is not empty
	// 	var Report Report

	// 	// reverse order for loop, cause latest export is expected to be at the end
	// 	for i := len(ReportsList) - 1; i >= 0; i-- {

	// 		if ReportsList[i].Id == CreatedReportId {
	// 			Report = ReportsList[i]
	// 			break
	// 		}

	// 	}

	// 	if Report.Status == "Finished" {
	// 		DownloadLink = Report.DownloadLink
	// 		break
	// 	}

	// }

	// DownloadLink Mock Up
	DownloadLink = "https://tzswiy3zk5dms05cfeo.s3.eu-central-1.amazonaws.com/from_2025-03-01_to_2025-04-01_MTc0MzU4MDY0MDE0Mw.csv?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Date=20250402T075723Z&X-Amz-SignedHeaders=host&X-Amz-Expires=604799&X-Amz-Credential=AKIARJCCZCDEKCUWYOXG%2F20250402%2Feu-central-1%2Fs3%2Faws4_request&X-Amz-Signature=857a3b30cb532fdc0d52137a8af7602cbdfd84f597de0c74f61727403c71be3c"

	fmt.Printf("  DownloadLink: %v\n", DownloadLink)

	csv := DownloadReport(DownloadLink)
	fmt.Printf("  csv: %v\n", csv)

}

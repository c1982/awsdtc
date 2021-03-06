package main

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

//go:embed pages/index.html
var IndexPage string

func RunPage(address string) {
	http.HandleFunc("/datatransfers", func(w http.ResponseWriter, r *http.Request) {
		billdate := r.URL.Query().Get("date")
		if billdate == "" {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%s", errors.New("date parameter cannot be empty"))
			return
		}

		start, end, err := GetDates(billdate)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%s", err)
			return
		}

		report, regionalreport, err := GenerateData(start, end, "MONTHLY")
		if err != nil {
			log.Printf("error: %s, start: %s, end: %s", err, start, end)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%s", err)
			return
		}

		data := struct {
			Costs   []UsageItem
			Regions []RegionalUsage
			Groups  []string
		}{
			Costs:   report,
			Regions: regionalreport,
			Groups:  regions.GroupByRegion(),
		}

		rspdata, err := json.Marshal(data)

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s", rspdata)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s", IndexPage)
	})

	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func GetDates(billdate string) (startdate, enddate string, err error) {
	start, err := time.Parse("2006|January", billdate)
	if err != nil {
		return startdate, enddate, err
	}
	end := start.AddDate(0, 1, -1)
	return start.Format("2006-01-02"), end.Format("2006-01-02"), nil
}

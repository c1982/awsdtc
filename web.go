package main

import (
	"fmt"
	"net/http"
	"text/template"
)

func RunPage() {
	index := template.Must(template.ParseFiles("./pages/index.html"))

	http.HandleFunc("/datatransfers", func(w http.ResponseWriter, r *http.Request) {

	})

	http.HandleFunc("/regionaldatatransfers", func(w http.ResponseWriter, r *http.Request) {

	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		start := "2021-05-01"
		end := "2021-05-10"

		report, regionalreport, err := GenerateData(start, end, "MONTHLY")
		if err != nil {
			fmt.Fprintf(w, "%s", err)
			return
		}

		data := struct {
			Costs   []UsageItem
			Regions []RegionalUsage
		}{
			Costs:   report,
			Regions: regionalreport,
		}

		index.Execute(w, data)
	})

	http.ListenAndServe(":8000", nil)
}

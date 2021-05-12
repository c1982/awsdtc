package main

import (
	"fmt"
	"net/http"
	"text/template"
)

func RunPage() {
	index := template.Must(template.ParseFiles("./pages/index.html"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		start := "2021-05-01"
		end := "2021-05-12"
		granularity := "MONTHLY"

		report, err := GenerateData(start, end, granularity)
		if err != nil {
			fmt.Fprintf(w, "%s", err)
			return
		}

		data := struct {
			Costs      []UsageItem
			AwsRegions Regions
		}{
			Costs:      report,
			AwsRegions: regions,
		}

		index.Execute(w, data)
	})

	http.ListenAndServe(":8000", nil)
}

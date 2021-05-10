package main

import (
	"fmt"
)

func main() {
	start := "2021-05-01"
	end := "2021-05-10"
	granularity := "MONTHLY"

	result, err := GenerateData(start, end, granularity)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Report:", result)
}

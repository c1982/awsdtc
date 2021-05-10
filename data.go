package main

import (
	"errors"
	"log"
	"regexp"
)

type UsageItem struct {
	BlendedCost       string
	BlendedCostUnit   string
	UsageQuantity     string
	UsageQuantityUnit string
	SourceRegion      string
	DestinationRegion string
	TransferDirection string
}

func GenerateData(start, end, granularity string) (items []UsageItem, err error) {
	output, err := GetCostAndUsage(start, end, granularity)
	if err != nil {
		return items, err
	}

	items = []UsageItem{}
	for i := 0; i < len(output.ResultsByTime); i++ {
		for g := 0; g < len(output.ResultsByTime[i].Groups); g++ {
			group := output.ResultsByTime[i].Groups[g]
			name := group.Keys[0]
			src, dst, direction, err := ParseUsageType(*name)
			if err != nil {
				log.Println("parse error: ", *name)
				continue
			}

			item := UsageItem{
				DestinationRegion: dst,
				SourceRegion:      src,
				TransferDirection: direction,
			}
			item.BlendedCost = *group.Metrics["BlendedCost"].Amount
			item.BlendedCostUnit = *group.Metrics["BlendedCost"].Unit
			item.UsageQuantity = *group.Metrics["UsageQuantity"].Amount
			item.UsageQuantityUnit = *group.Metrics["UsageQuantity"].Unit

			items = append(items, item)
		}
	}

	return items, nil
}

func ParseUsageType(text string) (src, dst, direction string, err error) {
	rg, err := regexp.Compile(`([A-Z0-9].+)\-([A-Z0-9].+)\-AWS\-(.+)\-Bytes`)
	if err != nil {
		return src, dst, direction, err
	}

	groups := rg.FindStringSubmatch(text)
	if len(groups) != 4 {
		return src, dst, direction, errors.New("cannot explode groups from usage type")
	}

	return groups[1], groups[2], groups[3], nil
}

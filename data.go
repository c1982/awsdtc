package main

import (
	"errors"
	"regexp"
)

type UsageItem struct {
	Name              string
	BlendedCost       string
	BlendedCostUnit   string
	UsageQuantity     string
	UsageQuantityUnit string
	SourceRegion      AwsRegion
	DestinationRegion AwsRegion
	TransferDirection string
}

func GenerateDataMock(start, end, granularity string) (items []UsageItem, err error) {
	return []UsageItem{
		{
			Name:              "",
			DestinationRegion: regions.GetByBillName("EUC1"),
			SourceRegion:      regions.GetByBillName("EUN1"),
			TransferDirection: "In",
			BlendedCost:       "0.01",
			BlendedCostUnit:   "USD",
			UsageQuantity:     "5.004",
			UsageQuantityUnit: "GB",
		},
		{
			Name:              "",
			DestinationRegion: regions.GetByBillName("EUC1"),
			SourceRegion:      regions.GetByBillName("EU"),
			TransferDirection: "Out",
			BlendedCost:       "0.41",
			BlendedCostUnit:   "USD",
			UsageQuantity:     "6.004",
			UsageQuantityUnit: "GB",
		},
		{
			Name:              "",
			DestinationRegion: regions.GetByBillName("EUC1"),
			SourceRegion:      regions.GetByBillName("AFS1"),
			TransferDirection: "Out",
			BlendedCost:       "0.41",
			BlendedCostUnit:   "USD",
			UsageQuantity:     "6.004",
			UsageQuantityUnit: "GB",
		},
	}, nil
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
				continue
			}

			item := UsageItem{
				Name:              *name,
				DestinationRegion: regions.GetByBillName(dst),
				SourceRegion:      regions.GetByBillName(src),
				TransferDirection: direction,
				BlendedCost:       *group.Metrics["BlendedCost"].Amount,
				BlendedCostUnit:   *group.Metrics["BlendedCost"].Unit,
				UsageQuantity:     *group.Metrics["UsageQuantity"].Amount,
				UsageQuantityUnit: *group.Metrics["UsageQuantity"].Unit,
			}
			items = append(items, item)
		}
	}

	return items, nil
}

func ParseUsageType(text string) (src, dst, direction string, err error) {
	//TODO: Cover these types EUN1-DataTransfer-In-Bytes, EUN1-DataTransfer-Out-Bytes,EUN1-DataTransfer-Regional-Bytes
	rg, err := regexp.Compile(`([A-Z0-9].+)\-([A-Z0-9].+)\-AWS\-(.+)\-Bytes`)
	if err != nil {
		return src, dst, direction, err
	}

	groups := rg.FindStringSubmatch(text)
	if len(groups) != 4 {
		return src, dst, direction, errors.New("cannot explode groups from 'usage type'")
	}

	return groups[1], groups[2], groups[3], nil
}

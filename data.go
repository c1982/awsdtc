package main

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/service/costexplorer"
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
	Opacity           string
}

type RegionalUsage struct {
	Name                 string
	Region               AwsRegion
	DataTransferOut      UsageItem
	DataTransferIn       UsageItem
	DataTransferRegional UsageItem
	UsagePercent         string
}

func GenerateDataMock(start, end, granularity string) (items []UsageItem, regionalusages []RegionalUsage, err error) {
	return []UsageItem{
			{
				Name:              "EUC1-EUN1",
				DestinationRegion: regions.GetByBillName("EUC1"),
				SourceRegion:      regions.GetByBillName("EUN1"),
				TransferDirection: "In",
				BlendedCost:       "0.01",
				BlendedCostUnit:   "USD",
				UsageQuantity:     "5.004",
				UsageQuantityUnit: "GB",
				Opacity:           "1",
			},
			{
				Name:              "EUC1-EU",
				DestinationRegion: regions.GetByBillName("EUC1"),
				SourceRegion:      regions.GetByBillName("EU"),
				TransferDirection: "Out",
				BlendedCost:       "0.41",
				BlendedCostUnit:   "USD",
				UsageQuantity:     "6.004",
				UsageQuantityUnit: "GB",
				Opacity:           "1",
			},
			{
				Name:              "EUC1-AFS1",
				DestinationRegion: regions.GetByBillName("EUC1"),
				SourceRegion:      regions.GetByBillName("AFS1"),
				TransferDirection: "Out",
				BlendedCost:       "0.41",
				BlendedCostUnit:   "USD",
				UsageQuantity:     "6.004",
				UsageQuantityUnit: "GB",
				Opacity:           "1",
			},
		}, []RegionalUsage{

			{"Frankfurt",
				regions.GetByBillName("EUC1"),
				UsageItem{
					Name:              "EUC1-AFS1",
					DestinationRegion: regions.GetByBillName("EUC1"),
					SourceRegion:      regions.GetByBillName("AFS1"),
					TransferDirection: "Out",
					BlendedCost:       "0.41",
					BlendedCostUnit:   "USD",
					UsageQuantity:     "6.004",
					UsageQuantityUnit: "GB",
					Opacity:           "1",
				},
				UsageItem{
					Name:              "EUC1-AFS1",
					DestinationRegion: regions.GetByBillName("EUC1"),
					SourceRegion:      regions.GetByBillName("AFS1"),
					TransferDirection: "Out",
					BlendedCost:       "0.41",
					BlendedCostUnit:   "USD",
					UsageQuantity:     "6.004",
					UsageQuantityUnit: "GB",
					Opacity:           "1",
				},
				UsageItem{
					Name:              "EUC1-AFS1",
					DestinationRegion: regions.GetByBillName("EUC1"),
					SourceRegion:      regions.GetByBillName("AFS1"),
					TransferDirection: "Out",
					BlendedCost:       "0.41",
					BlendedCostUnit:   "USD",
					UsageQuantity:     "6.004",
					UsageQuantityUnit: "GB",
					Opacity:           "1",
				},
				"0",
			},
			{"N.California",
				regions.GetByBillName("USW1"),
				UsageItem{
					Name:              "USW1-AFS1",
					DestinationRegion: regions.GetByBillName("USW1"),
					SourceRegion:      regions.GetByBillName("AFS1"),
					TransferDirection: "Out",
					BlendedCost:       "0.41",
					BlendedCostUnit:   "USD",
					UsageQuantity:     "10.0",
					UsageQuantityUnit: "GB",
					Opacity:           "1",
				},
				UsageItem{
					Name:              "USW1-AFS1",
					DestinationRegion: regions.GetByBillName("USW1"),
					SourceRegion:      regions.GetByBillName("AFS1"),
					TransferDirection: "Out",
					BlendedCost:       "0.41",
					BlendedCostUnit:   "USD",
					UsageQuantity:     "12.004",
					UsageQuantityUnit: "GB",
					Opacity:           "1",
				},
				UsageItem{
					Name:              "USW1-AFS1",
					DestinationRegion: regions.GetByBillName("USW1"),
					SourceRegion:      regions.GetByBillName("AFS1"),
					TransferDirection: "Out",
					BlendedCost:       "0.41",
					BlendedCostUnit:   "USD",
					UsageQuantity:     "12.004",
					UsageQuantityUnit: "GB",
					Opacity:           "1",
				},
				"0",
			},
		}, nil
}

func GenerateData(start, end, granularity string) (usages []UsageItem, regionalusages []RegionalUsage, err error) {
	output, err := GetCostAndUsage(start, end, granularity)
	if err != nil {
		return usages, regionalusages, err
	}

	tmpregionalusages := CreateRegions()
	usages = []UsageItem{}
	for i := 0; i < len(output.ResultsByTime); i++ {
		for g := 0; g < len(output.ResultsByTime[i].Groups); g++ {
			group := output.ResultsByTime[i].Groups[g]
			name := group.Keys[0]

			regionID, transfertype, usage, regional, _ := CreateRegionalUsage(group)
			if regional {
				regionalusage := tmpregionalusages[regionID]
				regionalusage.Name = *name
				regionalusage.Region = regions.GetByBillName(regionID)

				if transfertype == "In" {
					regionalusage.DataTransferIn = usage
				}

				if transfertype == "Out" {
					regionalusage.DataTransferOut = usage
				}

				if transfertype == "Regional" {
					regionalusage.DataTransferRegional = usage
				}
				tmpregionalusages[regionID] = regionalusage
			}

			src, dst, direction, err := ParseUsageType(*name)
			if err != nil {
				continue
			}

			lineOpacity := "0.1"
			if !strings.HasPrefix(*group.Metrics["UsageQuantity"].Amount, "0.0") {
				lineOpacity = *group.Metrics["UsageQuantity"].Amount
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
				Opacity:           lineOpacity,
			}
			usages = append(usages, item)
		}
	}

	calculatedusages := CalculatetUsagePercents(tmpregionalusages)
	return usages, calculatedusages, nil
}

func CreateRegions() (r map[string]RegionalUsage) {
	r = map[string]RegionalUsage{}
	for i := 0; i < len(regions); i++ {
		r[regions[i].A1] = RegionalUsage{
			Name:         regions[i].Location,
			Region:       regions[i],
			UsagePercent: "0",
			DataTransferRegional: UsageItem{
				UsageQuantity:     "0",
				UsageQuantityUnit: "GB",
				BlendedCost:       "0",
				BlendedCostUnit:   "USD",
			},
			DataTransferOut: UsageItem{
				UsageQuantity:     "0",
				UsageQuantityUnit: "GB",
				BlendedCost:       "0",
				BlendedCostUnit:   "USD",
			},
			DataTransferIn: UsageItem{
				UsageQuantity:     "0",
				UsageQuantityUnit: "GB",
				BlendedCost:       "0",
				BlendedCostUnit:   "USD",
			},
		}
	}
	return r
}

func CreateRegionalUsage(group *costexplorer.Group) (regionID, transfertype string, usage UsageItem, regional bool, err error) {
	//Patterns: EUN1-DataTransfer-In-Bytes, EUN1-DataTransfer-Out-Bytes,EUN1-DataTransfer-Regional-Bytes
	regional = false
	text := group.Keys[0]

	regionID, transfertype, err = ParseRegionalUsageType(*text)
	if err != nil {
		return regionID, transfertype, usage, regional, err
	}

	regional = true
	usage = UsageItem{
		Name:              *text,
		DestinationRegion: regions.GetByBillName(regionID),
		SourceRegion:      regions.GetByBillName(regionID),
		TransferDirection: transfertype,
		BlendedCost:       *group.Metrics["BlendedCost"].Amount,
		BlendedCostUnit:   *group.Metrics["BlendedCost"].Unit,
		UsageQuantity:     *group.Metrics["UsageQuantity"].Amount,
		UsageQuantityUnit: *group.Metrics["UsageQuantity"].Unit,
	}

	return regionID, transfertype, usage, regional, nil
}

func ParseRegionalUsageType(text string) (regionID, transfertype string, err error) {
	virginiaRegion := "USE1"

	if text == "DataTransfer-In-Bytes" {
		return virginiaRegion, "In", nil
	}

	if text == "DataTransfer-Out-Bytes" {
		return virginiaRegion, "Out", nil
	}

	if text == "DataTransfer-Regional-Bytes" {
		return virginiaRegion, "Regional", nil
	}

	rg, err := regexp.Compile(`(.+)\-DataTransfer\-(.+)-Bytes`)
	if err != nil {
		return regionID, transfertype, err
	}

	groups := rg.FindStringSubmatch(text)
	if len(groups) != 3 {
		return regionID, transfertype, errors.New("cannot explode groups from 'DataTransfer'")
	}

	return groups[1], groups[2], nil
}

func ParseUsageType(text string) (src, dst, direction string, err error) {
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

func CalculatetUsagePercents(usages map[string]RegionalUsage) []RegionalUsage {
	total := GetTotalUsageQuantity(usages)
	tmpmap := []RegionalUsage{}
	for _, v := range usages {
		f, err := strconv.ParseFloat(v.DataTransferRegional.UsageQuantity, 64)
		if err != nil {
			log.Println("UsageQuantity cannot parse to float", err, "UsageQuantity:", v.DataTransferRegional.UsageQuantity)
			continue
		}

		percent := f / total * 100
		if percent < 5.0 {
			percent = 7
		}

		if percent > 50 {
			percent = 40
		}

		v.UsagePercent = fmt.Sprintf("%f", percent)
		tmpmap = append(tmpmap, v)
	}

	sort.Slice(tmpmap, func(i, j int) bool {
		return tmpmap[i].Name > tmpmap[j].Name
	})
	return tmpmap
}

func GetTotalUsageQuantity(usages map[string]RegionalUsage) float64 {
	var total float64
	for _, v := range usages {
		f, err := strconv.ParseFloat(v.DataTransferRegional.UsageQuantity, 64)
		if err != nil {
			log.Println("UsageQuantity cannot parse to float", err, "UsageQuantity:", v.DataTransferRegional.UsageQuantity)
			continue
		}
		total += f
	}

	return total
}

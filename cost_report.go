package main

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/costexplorer"
)

type AwsRegion struct {
	Region   string
	Location string
	Code     string
	A1       string
}

type Regions []AwsRegion

func (r Regions) GetByBillName(billname string) AwsRegion {

	for i := 0; i < len(r); i++ {
		if r[i].A1 == billname {
			return r[i]
		}
	}

	return AwsRegion{}
}

var (
	regions = Regions{
		{"Asia Pacific", "Hong Kong", "ap-east-1", "APE1"},
		{"Asia Pacific", "Tokyo", "ap-northeast-1", "APN1"},
		{"Asia Pacific", "Seoul", "ap-northeast-2", "APN2"},
		{"Asia Pacific", "Osaka", "ap-northeast-3", "APN3"},
		{"Asia Pacific", "Singapore", "ap-southeast-1", "APS1"},
		{"Asia Pacific", "Sydney", "ap-southeast-2", "APS2"},
		{"Asia Pacific", "Mumbai", "ap-south-1", "APS3"},
		{"Canada", "Montreal", "ca-central-1", "CAN1"},
		{"Africa", "Cape Town", "af-south-1", "CPT"},
		{"Europe", "Stockholm", "eu-north-1", "EUN1"},
		{"Europe", "Frankfurt", "eu-central-1", "EUC1"},
		{"Europe", "Ireland", "eu-west-1", "EU"},
		{"Europe", "London", "eu-west-2", "EUW2"},
		{"Europe", "Paris", "eu-west-3", "EUW3"},
		{"Europe", "Milan", "eu-south-1", "EUS1"},
		{"Middle East", "Bahrain", "me-south-1", "MES1"},
		{"South America", "São Paulo", "sa-east-1", "SAE1"},
		{"AWS GovCloud", "US-West", "us-gov-west-1", "UGW1"},
		{"AWS GovCloud", "US-East", "us-gov-east-1", "UGE1"},
		{"US East", "N. Virginia", "us-east-1", "USE1"},
		{"US East", "Ohio", "us-east-2", "USE2"},
		{"US West", "N. California", "us-west-1", "USW1"},
		{"US West", "Oregon", "us-west-2", "USW2"},
	}

	awsMetrics = aws.StringSlice([]string{
		"BlendedCost",
		"UnblendedCost",
		"UsageQuantity",
	})

	awsSession *session.Session
	awsCostSvc *costexplorer.CostExplorer
)

func init() {
	var err error
	awsSession, err = session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	if err != nil {
		log.Fatal("AWS session cannot be started: ", err)
	}

	awsCostSvc = costexplorer.New(awsSession)
}

func GetCostAndUsage(start, end, granularity string) (output *costexplorer.GetCostAndUsageOutput, err error) {
	output, err = awsCostSvc.GetCostAndUsage(&costexplorer.GetCostAndUsageInput{
		TimePeriod: &costexplorer.DateInterval{
			Start: aws.String(start),
			End:   aws.String(end),
		},
		Granularity: aws.String(granularity),
		GroupBy: []*costexplorer.GroupDefinition{
			{
				Type: aws.String("DIMENSION"),
				Key:  aws.String("USAGE_TYPE"),
			},
		},
		// Filter: &costexplorer.Expression{
		// 	Dimensions: &costexplorer.DimensionValues{
		// 		Key:    aws.String("USAGE_TYPE"),
		// 		Values: aws.StringSlice([]string{"APE1-MES1-AWS-In-Bytes"}),
		// 	},
		// },
		Metrics: awsMetrics,
	})

	return output, err
}

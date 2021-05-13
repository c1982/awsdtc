package main

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/costexplorer"
)

const (
	defaultAWSRegion = "us-east-1"
)

type AwsRegion struct {
	Region    string
	Location  string
	Code      string
	A1        string
	Latitude  string
	Longitude string
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
		{"Asia Pacific", "Hong Kong", "ap-east-1", "APE1", "22.302711", "114.177216"},
		{"Asia Pacific", "Tokyo", "ap-northeast-1", "APN1", "35.689487", "139.691711"},
		{"Asia Pacific", "Seoul", "ap-northeast-2", "APN2", "37.566536", "126.977966"},
		{"Asia Pacific", "Osaka", "ap-northeast-3", "APN3", "34.693737", "135.502167"},
		{"Asia Pacific", "Singapore", "ap-southeast-1", "APS1", "1.352083", "103.819839"},
		{"Asia Pacific", "Sydney", "ap-southeast-2", "APS2", "-33.868820", "151.209290"},
		{"Asia Pacific", "Mumbai", "ap-south-1", "APS3", "19.0760", "72.8777"},
		{"Canada", "Montreal", "ca-central-1", "CAN1", "45.5017", "-73.5673"},
		{"Africa", "Cape Town", "af-south-1", "AFS1", "-33.9249", "18.4241"},
		{"Europe", "Stockholm", "eu-north-1", "EUN1", "59.3293", "18.0686"},
		{"Europe", "Frankfurt", "eu-central-1", "EUC1", "50.1109", "8.6821"},
		{"Europe", "Ireland", "eu-west-1", "EU", "53.1424", "-7.6921"},
		{"Europe", "London", "eu-west-2", "EUW2", "51.5074", "0.1278"},
		{"Europe", "Paris", "eu-west-3", "EUW3", "48.8566", "2.3522"},
		{"Europe", "Milan", "eu-south-1", "EUS1", "45.4642", "9.1900"},
		{"Middle East", "Bahrain", "me-south-1", "MES1", "26.0667", "50.5577"},
		{"South America", "São Paulo", "sa-east-1", "SAE1", "-23.5505", "-46.6333"},
		{"US East", "N. Virginia", "us-east-1", "USE1", "37.4316", "-78.6569"},
		{"US East", "Ohio", "us-east-2", "USE2", "40.4173", "-82.9071"},
		{"US East", "Miami", "us-east-1-mia", "MIA1", "25.761681", "-80.191788"},
		{"US East", "Houston", "us-east-1-iah", "IAH1", "29.7604", "-95.3698"},
		{"US East", "Boston", "us-east-1-bos", "BOS1", "42.35843", "-71.05977"},
		{"US West", "N. California", "us-west-1", "USW1", "36.7783", "-122.08"},
		{"US West", "Oregon", "us-west-2", "USW2", "44.000000", "-120.500000"},
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
		Region: aws.String(defaultAWSRegion)},
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
		//TODO: delete when the time comes
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

# The AWS Data Transfer Cost Explorer

The AWS Data Transfer Cost Explorer tool analyzes the billed Data Transfer items in your AWS account and presents them visualized on a map.

### Motivation;

We have a continious cost optimization case on AWS. Especially the Data Transfer tab on the Bills screen is quite long and it takes a long time to understand which areas are used more.

Another need is to catch unusual Data Transfers in our infrastructure. For example, in our infrastructure, it is not possible to get traffic from Tokyo to Sao Paulo, but thanks to this vehicle, we can see and solve it. 

General:
![](./ss-explorer.png)

Filtered:
![](./ss-frankfurt.png)

## Installation

* [awsdtc](https://github.com/c1982/awsdtc) 

## Configuration

* set AWS credentials in `~/.aws/credentials` file

```ini
[default]
aws_access_key_id = A******************U
aws_secret_access_key = WD/**********************************MA
```

or

* set AWS credentials system environment

```bash
export AWS_ACCESS_KEY_ID=A******************U
export AWS_SECRET_ACCESS_KEY=WD/**********************************MA
```

## Policy

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "aws-portal:ViewUsage",
                "aws-portal:ViewBilling",
                "cur:DescribeReportDefinitions",
            ],
            "Resource": "*"
        }
    ]
}
```

## Running

1. Download Binary for your OS
2. Run awsdtc executable

## Usage

TODO: will complete
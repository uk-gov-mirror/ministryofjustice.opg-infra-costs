package unblendedcosts

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/costexplorer"
)

// CostAndUsageInput returns a correctly formatted input struct for SDK usage
func (c *AWSCosts) CostAndUsageInput() *costexplorer.GetCostAndUsageInput {

	return &costexplorer.GetCostAndUsageInput{
		Granularity: aws.String(c.Granularity),
		TimePeriod: &costexplorer.DateInterval{
			Start: aws.String(c.Start.Format("2006-01-02")),
			End:   aws.String(c.End.Format("2006-01-02")),
		},
		Metrics: []*string{
			aws.String(c.CostType),
		},
		GroupBy: []*costexplorer.GroupDefinition{
			{
				Type: aws.String("DIMENSION"),
				Key:  aws.String("SERVICE"),
			},
		},
	}

}

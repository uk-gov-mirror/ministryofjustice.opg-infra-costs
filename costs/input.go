package costs

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/costexplorer"
)

// CostAndUsageInput takes the time & granularity details and generates an Input object for the
// request to cost explorer
func CostAndUsageInput(
	start time.Time,
	end time.Time,
	granularity string,
	costType string) *costexplorer.GetCostAndUsageInput {

	return &costexplorer.GetCostAndUsageInput{
		Granularity: aws.String(granularity),
		TimePeriod: &costexplorer.DateInterval{
			Start: aws.String(start.Format("2006-01-02")),
			End:   aws.String(end.Format("2006-01-02")),
		},
		Metrics: []*string{
			aws.String(costType),
		},
		GroupBy: []*costexplorer.GroupDefinition{
			{
				Type: aws.String("DIMENSION"),
				Key:  aws.String("SERVICE"),
			},
		},
	}
}

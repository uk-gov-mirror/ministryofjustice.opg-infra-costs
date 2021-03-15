package unblendedcosts

import (
	"opg-infra-costs/session"

	"github.com/aws/aws-sdk-go/service/costexplorer"
)

// AssumedSession returns an authorised session based on the account configured
func (c *AWSCosts) AssumedSession() (*costexplorer.CostExplorer, error) {

	sess, err := session.Assumed(c.Account.Arn(), c.Account.Region)
	if err != nil {
		return &costexplorer.CostExplorer{}, err
	}
	return costexplorer.New(sess), nil
}

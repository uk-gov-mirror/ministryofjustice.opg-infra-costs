package costs

import (
	"opg-infra-costs/accounts"
	"opg-infra-costs/session"

	"github.com/aws/aws-sdk-go/service/costexplorer"
)

// Client returns a costexplorer using an assumed session
func Client(
	account accounts.Account) (*costexplorer.CostExplorer, error) {

	sess, err := session.Assumed(account.Arn(), account.Region)
	if err != nil {
		return nil, err
	}
	return costexplorer.New(sess), nil
}

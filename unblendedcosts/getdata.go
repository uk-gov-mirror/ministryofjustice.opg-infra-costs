package unblendedcosts

import (
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/costexplorer"
)

// requestAndResponse provides request and repsonse objects
// for the sdk
func (c *AWSCosts) requestAndResponse(session *costexplorer.CostExplorer) (*request.Request, *costexplorer.GetCostAndUsageOutput) {
	// generate a suitable sdk
	sdkInput := c.CostAndUsageInput()
	return session.GetCostAndUsageRequest(sdkInput)

}

// Fetch generates the request and fetches teh raw response from AWS
// via an authenticated assumed session
func (c *AWSCosts) Fetch() error {
	// get an authenticated session for this account
	session, err := c.AssumedSession()
	if err != nil {
		return err
	}
	//
	request, response := c.requestAndResponse(session)
	// send the request
	err = request.Send()
	if err != nil {
		return err
	}
	// set the raw response
	c.RawResult = response
	return nil
}

// CostData makes the api call, converts the data structures and returns
// new struct
func (c *AWSCosts) CostData() (CostData, error) {
	// fetch the data
	err := c.Fetch()
	if err != nil {
		return CostData{}, err
	}
	// convert to local format
	rows, err := c.ConvertToCostRows()
	return CostData{Entries: rows}, err

}

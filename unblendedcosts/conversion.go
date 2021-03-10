package unblendedcosts

import (
	"strconv"
	"strings"
)

// ConvertToCostRows takes the RawResult and converts to the
// CostRow format for local use
func (c *AWSCosts) ConvertToCostRows() ([]CostRow, error) {

	var resultsCosts []CostRow
	filterByService := c.ServiceFilter
	l := len(c.ServiceFilter)
	// convert the raw api result to a slice of CostRow
	for _, results := range c.RawResult.ResultsByTime {
		startDate := *results.TimePeriod.Start
		for _, groups := range results.Groups {
			for _, metrics := range groups.Metrics {

				fVal, e := strconv.ParseFloat(*metrics.Amount, 64)
				if e != nil {
					return nil, e
				}
				r := CostRow{
					Date:    startDate,
					Service: *groups.Keys[0],
					Cost:    fVal,
					Account: c.Account,
				}

				r.ServiceNameCorrection()

				// if there is no filter, or if the filter contained in the service name
				if l == 0 || (l > 0 && strings.Contains(strings.ToUpper(r.Service), strings.ToUpper(filterByService))) {
					resultsCosts = append(resultsCosts, r)
				}
			}
		}
	}

	return resultsCosts, nil
}

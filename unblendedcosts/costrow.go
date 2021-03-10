package unblendedcosts

import (
	"fmt"
	"strings"
)

func (c *CostRow) asMap() map[string]string {
	return map[string]string{
		"Date":    c.Date,
		"Service": c.Service,
		"Cost":    fmt.Sprintf("%f", c.Cost),
	}
}

// Get helper
func (c *CostRow) Get(prop string) string {
	// pass down to account
	if strings.Contains(prop, "Account.") {
		return c.Account.Get(strings.ReplaceAll(prop, "Account.", ""))
	} else if strings.Contains(prop, "Meta.") {
		meta := c.Meta
		return meta[strings.ReplaceAll(prop, "Meta.", "")]
	}
	return c.asMap()[prop]
}

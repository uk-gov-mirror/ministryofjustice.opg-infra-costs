package unblendedcosts

// ServiceNameCorrection updates names of services
// AWS changed the name in cost explorer, so this makes them all match
// to keep spreadsheets / tables tidy
func (r *CostRow) ServiceNameCorrection() {

	switch name := r.Service; name {
	case "Amazon EC2 Container Service":
		r.Service = "Amazon Elastic Container Service"
	}

}

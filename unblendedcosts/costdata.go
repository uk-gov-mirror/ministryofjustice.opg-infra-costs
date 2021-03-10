package unblendedcosts

// Total returns grand total of .Entries
func (cd *CostData) Total() float64 {
	total := 0.0
	for _, r := range cd.Entries {
		total = total + r.Cost
	}
	return total
}

// GroupByKeys returns a new CostData re-organised to be grouped by the key array passed
func (cd *CostData) GroupByKeys(keys []string) CostData {
	costs := CostData{}
	mapped := cd.GroupByKeysMap(keys)
	// add to the costs - range does not give order
	for _, row := range mapped {
		costs.Entries = append(costs.Entries, row)
	}
	return costs
}

func (cd *CostData) GroupByKeysMap(keys []string) map[string]CostRow {
	mapped := make(map[string]CostRow)
	//loop over entries and group costs
	for _, c := range cd.Entries {
		found := CostRow{}
		ok := false
		foundFloat := 0.0
		key := GenerateGroupKey(keys, &c)
		// if found, then sum
		if found, ok = mapped[key]; ok {
			foundFloat = found.Cost
		}
		sum := foundFloat + c.Cost
		found = c
		found.Cost = sum
		mapped[key] = found
	}
	return mapped
}

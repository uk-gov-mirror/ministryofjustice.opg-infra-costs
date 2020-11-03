package costs

// IntersectingMaps returns all keys from a that exist in b
func IntersectingMaps(a map[string]CostRow, b map[string]CostRow) []string {
	inter := make([]string, 0)
	for k, _ := range a {
		for i, _ := range b {
			if k == i {
				inter = append(inter, k)
			}
		}
	}
	return inter
}

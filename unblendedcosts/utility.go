package unblendedcosts

import (
	"fmt"
	"sort"
	"strings"
)

// GenerateGroupKey generates a combined string to use as a single level key for maps
// based on several fields
func GenerateGroupKey(keys []string, cr *CostRow) string {
	key := ""
	// sort the keys, reverse order
	sort.Sort(sort.Reverse(sort.StringSlice(keys)))
	for _, k := range keys {
		t := strings.Trim(cr.Get(k), " ")
		t = strings.ToUpper(t)
		key = key + fmt.Sprintf("%v:%v||", k, t)
	}
	return strings.Trim(key, "||")
}

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

// Increased looks for values that have increased between matching services
func Increased(
	intersect []string,
	mapA map[string]CostRow,
	mapB map[string]CostRow,
	increaseP int,
	baseCost float64,
) CostData {
	var increasedCosts CostData

	for _, k := range intersect {
		row := mapA[k]
		costA := mapA[k].Cost
		costB := mapB[k].Cost

		if costA > baseCost || costB > baseCost {
			costDiff := costB - costA
			//look for %
			increase := int(costDiff / (costA / 100))
			// store meta data as map for display
			row.Meta = map[string]string{
				"A":    fmt.Sprintf("%.2f", costA),
				"B":    fmt.Sprintf("%.2f", costB),
				"Diff": fmt.Sprintf("%.2f", costDiff),
				"P":    fmt.Sprintf("%d", increase),
			}

			if increase > increaseP {
				increasedCosts.Entries = append(increasedCosts.Entries, row)
			}
		}
	}
	return increasedCosts
}

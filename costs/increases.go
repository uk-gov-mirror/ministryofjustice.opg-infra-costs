package costs

import "fmt"

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

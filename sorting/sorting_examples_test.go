package sorting_test

import (
	"fmt"
	"sort"

	"github.com/chengxiaoer/geomGo"
	"github.com/chengxiaoer/geomGo/sorting"
)

func ExampleNewFlatCoordSorting2D() {
	// Some description
	coords := []float64{1, 0, 0, 1, 2, 2, 2, -2, -1, 0}
	sort.Sort(sorting.NewFlatCoordSorting2D(geom.XY, coords))
	fmt.Println(coords)
	// Output: [-1 0 0 1 1 0 2 -2 2 2]
}

func ExampleNewFlatCoordSorting() {
	coords := []float64{1, 0, 0, 1, 2, 2, -1, 0}
	isLess := func(c1, c2 []float64) bool {
		return c1[0] < c2[0]
	}
	sort.Sort(sorting.NewFlatCoordSorting(geom.XY, coords, isLess))
	fmt.Println(coords)
	// Output: [-1 0 0 1 1 0 2 2]
}

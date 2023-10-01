package check_points

import (
	"fmt"
	"github.com/agfy/doom_environment/image_comparer"
	"testing"
)

func TestCheckPoints(t *testing.T) {
	checkPoints, err := NewCheckPoints("loc_1_lvl_1/")
	if err != nil {
		t.Error(err)
	}

	testPoints, err := NewCheckPoints("test_points/")
	if err != nil {
		t.Error(err)
	}

	for _, testPoint := range testPoints.Points {
		fmt.Printf("%d ", testPoint.Score)
		for _, checkPoint := range checkPoints.Points {
			eq, err := image_comparer.AreImagesEqual(testPoint.Img, checkPoint.Img)
			if err != nil {
				t.Error(err)
			}
			if eq {
				fmt.Printf(" equal %d ", checkPoint.Score)
			}
		}
		print("\n")
	}
}

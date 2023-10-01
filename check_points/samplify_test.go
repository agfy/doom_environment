package check_points

import (
	"fmt"
	"github.com/agfy/doom_environment/image_comparer"
	"testing"
)

func TestSamplify(t *testing.T) {
	checkPoints, err := NewCheckPoints("loc_1_lvl_1/sample_1/")
	if err != nil {
		t.Error(err)
	}

	testPoints, err := NewCheckPoints("loc_1_lvl_1/test_points/")
	if err != nil {
		t.Error(err)
	}

	for _, samples := range []int{2, 4, 6, 8, 10} {
		samplifiedTestPoints := make([]CheckPoint, len(testPoints.Points))
		for i, testPoint := range testPoints.Points {
			samplifiedTestPoints[i] = CheckPoint{Img: image_comparer.Samplify(testPoint.Img, samples), Score: testPoint.Score}
		}

		samplifiedCheckPoints := make([]CheckPoint, len(checkPoints.Points))
		for i, checkPoint := range checkPoints.Points {
			samplifiedCheckPoints[i] = CheckPoint{Img: image_comparer.Samplify(checkPoint.Img, samples), Score: checkPoint.Score}
			name := fmt.Sprintf("loc_1_lvl_1/sample_%d/doom_%d.jpg", samples, i)
			err = image_comparer.Save(name, samplifiedCheckPoints[i].Img)
			if err != nil {
				t.Error(err)
			}
		}

		for _, testPoint := range samplifiedTestPoints {
			fmt.Printf("%d ", testPoint.Score)
			for _, checkPoint := range samplifiedCheckPoints {
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
		print("\n")
	}

}

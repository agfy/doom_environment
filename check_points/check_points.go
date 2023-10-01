package check_points

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"
)

type CheckPoint struct {
	Img   image.Image
	Score int
}

type CheckPoints struct {
	Points []CheckPoint
}

func NewCheckPoints(dir string) (CheckPoints, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return CheckPoints{}, err
	}
	points := make([]CheckPoint, len(files))
	for i, _ := range files {
		fileName := fmt.Sprintf("%sdoom_%d.jpg", dir, i)
		f, err := os.Open(fileName)
		if err != nil {
			return CheckPoints{}, err
		}
		defer f.Close()
		img, err := jpeg.Decode(f)

		if err != nil {
			return CheckPoints{}, err
		}

		//words := strings.Split(file.Name(), "_")
		//if len(words) < 2 {
		//	return CheckPoints{}, errors.New("failed to parse filename " + file.Name())
		//}
		//strScore := strings.TrimSuffix(words[len(words)-1], ".jpg")
		//score, err := strconv.Atoi(strScore)
		//if err != nil {
		//	return CheckPoints{}, errors.New("failed to atoi" + strScore)
		//}

		points[i] = CheckPoint{Img: img, Score: i}
	}

	return CheckPoints{Points: points}, nil
}

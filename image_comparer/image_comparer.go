package image_comparer

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"math"
	"os"
)

var threshHold = 30.0

const equalThreshHold = 2000
const searchWindowPercent = 40

func AbsDiff(left, right uint32) uint32 {
	if left > right {
		return left - right
	} else {
		return right - left
	}
}

func AreImagesEqual(img1, img2 image.Image) (bool, error) {
	if img1.Bounds() != img2.Bounds() {
		return false, errors.New("not equal image size")
	}

	searchWindowPixels := img1.Bounds().Size().X * searchWindowPercent / 100
	equalPixels := 0
	for x := 0; x < img1.Bounds().Size().X; x++ {
		for y := 0; y < img1.Bounds().Size().Y; y++ {
			r1, g1, b1, _ := img1.At(x, y).RGBA()
			for i := searchWindowPixels * -1; i < searchWindowPixels; i++ {
				if x+i <= 0 || x+i >= img1.Bounds().Size().X {
					continue
				}
				r2, g2, b2, _ := img2.At(x+i, y).RGBA()
				if AbsDiff(r1, r2) < equalThreshHold && AbsDiff(g1, g2) < equalThreshHold && AbsDiff(b1, b2) < equalThreshHold {
					equalPixels++
					break
				}
			}

		}
	}
	equalPixelPercent := equalPixels * 100 / (img1.Bounds().Size().Y * img1.Bounds().Size().X)
	//fmt.Printf("equal pixel percent %d ", equalPixelPercent)
	if equalPixelPercent > 75 {
		return true, nil
	}
	return false, nil
}

func AreImagesEqual1(img1, img2 image.Image) (bool, error) {
	if img1.Bounds() != img2.Bounds() {
		return false, errors.New("not equal image size")
	}
	psnr := PSNR(img1, img2)
	fmt.Printf("%.2f ", psnr)
	if psnr > threshHold {
		return true, nil
	}
	return false, nil
}

func PSNR(img1, img2 image.Image) float64 {
	var se uint32
	for x := 0; x < img1.Bounds().Size().X; x++ {
		for y := 0; y < img1.Bounds().Size().Y; y++ {
			r1, g1, b1, _ := img1.At(x, y).RGBA()
			r2, g2, b2, _ := img2.At(x, y).RGBA()
			se += (r1-r2)*(r1-r2) + (g1-g2)*(g1-g2) + (b1-b2)*(b1-b2)
		}
	}

	mse := float64(se) / float64(3*img1.Bounds().Size().X*img1.Bounds().Size().Y)
	return 10.0 * math.Log10(256.0*256.0/mse)
}

func Save(name string, img image.Image) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()

	opt := jpeg.Options{
		Quality: 90,
	}
	err = jpeg.Encode(f, img, &opt)
	//err = png.Encode(f, img)
	if err != nil {
		return err
	}
	return nil
}

func Samplify(img image.Image, sample int) image.Image {
	width := img.Bounds().Max.X
	height := img.Bounds().Max.Y
	rect := image.Rectangle{Max: image.Point{X: width / sample, Y: height / sample}}
	sampledImg := image.NewRGBA(rect)

	for x := 0; x < width; x += sample {
		for y := 0; y < height; y += sample {
			pixelMap := make(map[color.Color]int)
			for xOffset := 0; xOffset < sample; xOffset++ {
				for yOffset := 0; yOffset < sample; yOffset++ {
					pixelMap[img.At(x+xOffset, y+yOffset)]++
				}
			}

			maxValue := 0
			mostFrequentColour := color.Color(color.RGBA{})
			for k, v := range pixelMap {
				if v > maxValue {
					maxValue = v
					mostFrequentColour = k
				}
			}
			sampledImg.Set(x/sample, y/sample, mostFrequentColour)
		}
	}
	return sampledImg
}

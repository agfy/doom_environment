package doom_environment

import (
	"image"
)

type Observation struct {
	Image image.Image
	Score int
}

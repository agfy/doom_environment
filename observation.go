package doom_environment

import (
	"github.com/go-vgo/robotgo"
)

type Observation struct {
	Bitmap robotgo.Bitmap
	Score  int
}

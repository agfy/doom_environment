package doom_environment

import (
	"testing"
)

func TestPlay(t *testing.T) {
	env, err := Create(1, 6)
	if err != nil {
		t.Error("failed to create environment", err)
	}
	defer env.Close()

	for {
		score, err := env.GetScore(0)
		//input := env.GetObservation(0)
		//println(len(input))
		if err != nil {
			t.Error("failed to get score", err)
		}
		println(score)
	}
}

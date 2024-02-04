package doom_environment

import (
	"testing"
	"time"
)

func TestScore(t *testing.T) {
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

func TestPlay(t *testing.T) {
	env, err := Create(1, 6)
	if err != nil {
		t.Error("failed to create environment", err)
	}
	defer env.Close()

	err = env.Reset()
	if err != nil {
		t.Error("failed to reset environment", err)
	}
	stepDuration := 100 * time.Millisecond

	for i := 0; i < 10; i++ {
		err = env.Step([]bool{true, false, false, false, false, false, false, false}, 0)
		if err != nil {
			t.Error(err)
		}
		time.Sleep(stepDuration)
	}

	for i := 0; i < 10; i++ {
		err = env.Step([]bool{false, true, false, false, false, false, false, false}, 0)
		if err != nil {
			t.Error(err)
		}
		time.Sleep(stepDuration)
	}

	for i := 0; i < 10; i++ {
		err = env.Step([]bool{true, false, true, false, false, false, false, false}, 0)
		if err != nil {
			t.Error(err)
		}
		time.Sleep(stepDuration)
	}

	for i := 0; i < 10; i++ {
		err = env.Step([]bool{false, true, false, true, false, false, false, false}, 0)
		if err != nil {
			t.Error(err)
		}
		time.Sleep(stepDuration)
	}

	err = env.Step([]bool{true, false, false, false, false, false, false, false}, 0)
	if err != nil {
		t.Error(err)
	}
	time.Sleep(stepDuration)

	err = env.Step([]bool{false, true, false, false, false, false, false, false}, 0)
	if err != nil {
		t.Error(err)
	}
	time.Sleep(stepDuration)

	err = env.Step([]bool{false, false, false, false, false, false, false, false}, 0)
	time.Sleep(stepDuration)
	env.Close()
}

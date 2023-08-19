package doom_environment

import (
	"strconv"
	"testing"
	"time"
)

func TestPlay(t *testing.T) {
	env, err := Create(1)
	if err != nil {
		t.Error("failed to create environment", err)
	}
	defer env.Close()

	var obs *Observation
	for i := 0; i < 1000; i++ {
		obs = env.GetObservation(0)
		err = env.Save("doom_"+strconv.Itoa(i)+".jpg", obs.Image)
		if err != nil {
			t.Error(err)
		}

		time.Sleep(time.Second)
	}
}

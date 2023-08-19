package main

import (
	"strconv"
	"time"

	"github.com/agfy/doom_environment"
)

func main() {
	env, err := doom_environment.Create(1)
	if err != nil {
		println("failed to create environment", err.Error())
		return
	}
	defer env.Close()

	err = env.Start()
	if err != nil {
		println("failed to start environment", err.Error())
		return
	}

	var obs *doom_environment.Observation
	for i := 0; i < 1000; i++ {
		obs = env.GetObservation(0)
		err = env.Save("doom_"+strconv.Itoa(i)+".jpg", obs.Image)
		if err != nil {
			println("failed to save observation", err.Error())
			return
		}

		time.Sleep(time.Second)
	}
}

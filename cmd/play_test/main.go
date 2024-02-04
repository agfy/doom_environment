package main

import (
	"fmt"
	"github.com/agfy/doom_environment"
	"time"
)

func main() {
	env, err := doom_environment.Create(1, 6)
	if err != nil {
		_ = fmt.Errorf("failed to create environment %v", err)
	}

	err = env.Reset()
	if err != nil {
		fmt.Printf("failed to reset environment, %v", err)
	}
	stepDuration := 500 * time.Millisecond

	for i := 0; i < 10; i++ {
		err = env.Step([]bool{true, false, false, false, false, false, false, false}, 0)
		if err != nil {
			fmt.Printf("failed to reset environment, %v", err)
		}
		time.Sleep(stepDuration)
	}

	for i := 0; i < 10; i++ {
		err = env.Step([]bool{false, true, false, false, false, false, false, false}, 0)
		if err != nil {
			fmt.Printf("failed to reset environment, %v", err)
		}
		time.Sleep(stepDuration)
	}

	for i := 0; i < 10; i++ {
		err = env.Step([]bool{true, false, true, false, false, false, false, false}, 0)
		if err != nil {
			fmt.Printf("failed to reset environment, %v", err)
		}
		time.Sleep(stepDuration)
	}

	for i := 0; i < 10; i++ {
		err = env.Step([]bool{false, true, false, true, false, false, false, false}, 0)
		if err != nil {
			fmt.Printf("failed to reset environment, %v", err)
		}
		time.Sleep(stepDuration)
	}

	//for {
	//	score, err := env.GetScore(0)
	//	if err != nil {
	//		_ = fmt.Errorf("failed to get score %v", err)
	//	}
	//	println(score)
	//}

	env.Close()
}

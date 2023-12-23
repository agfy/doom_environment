package main

import (
	"fmt"
	"github.com/agfy/doom_environment"
)

func main() {
	env, err := doom_environment.Create(1, 0)
	if err != nil {
		_ = fmt.Errorf("failed to create environment %v", err)
	}
	defer env.Close()

	for {
		score, err := env.GetScore(0)
		if err != nil {
			_ = fmt.Errorf("failed to get score %v", err)
		}
		println(score)
	}
}

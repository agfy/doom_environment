package main

import (
	"github.com/agfy/doom_environment"
	"time"
)

func main() {
	env, err := doom_environment.Create(1, 2)
	if err != nil {
		println("failed to create environment", err.Error())
		return
	}
	defer env.Close()

	err = env.Reset()
	if err != nil {
		println("failed to start environment", err.Error())
		return
	}

	//	var obs *doom_environment.Observation
	for i := 0; i < 1000; i++ {
		_ = env.GetImage(0)
		//err = env.Save("doom_"+strconv.Itoa(i)+".jpg", bitmap)
		//if err != nil {
		//	println("failed to save observation", err.Error())
		//	return
		//}
		//bitmap := robotgo.CaptureImg(640, 337, 640, 514)
		//err = env.Save("doom_"+strconv.Itoa(i)+".jpg", bitmap)
		// use `defer robotgo.FreeBitmap(bit)` to free the bitmap
		//defer robotgo.FreeBitmap(bitmap)
		//fmt.Println("...", bitmap)

		//err = robotgo.SavePng(bitmap, "doom_"+strconv.Itoa(i)+".png")
		//if err != nil {
		//	println("failed to start environment", err.Error())
		//	return
		//}

		time.Sleep(time.Second)
	}
}

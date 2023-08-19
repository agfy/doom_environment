package doom_environment

import (
	"errors"
	"fmt"
	"github.com/go-vgo/robotgo"
	"image"
	"image/jpeg"
	"os"
	"os/exec"
	"sync"
	"time"
)

var windowName = "prboom-plus"

type DoomEnvironment struct {
	pids  []int32
	mutex sync.Mutex
}

func Create(numberOfWindows int) (*DoomEnvironment, error) {
	for i := 0; i < numberOfWindows; i++ {
		go func() {
			cmd := exec.Command("prboom-plus", "doom1")
			err := cmd.Run()
			if err != nil {
				fmt.Println("an error occurred while starting doom", err.Error())
			}
		}()
	}

	numberOfTries := 10
	for i := 0; i < numberOfTries; i++ {
		pids, err := robotgo.FindIds(windowName)
		if err != nil {
			return nil, err
		}
		if len(pids) == numberOfWindows {
			time.Sleep(time.Second)
			return &DoomEnvironment{pids: pids}, nil
		}

		time.Sleep(time.Second)
	}

	return nil, errors.New("number of process not equal numberOfWindows")
}

func (e *DoomEnvironment) Start() error {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	for _, pid := range e.pids {
		time.Sleep(100 * time.Millisecond)
		err := robotgo.KeyTap("enter")
		if err != nil {
			return err
		}

		err = robotgo.ActivePID(pid)
		if err != nil {
			return err
		}

		time.Sleep(100 * time.Millisecond)
		err = robotgo.KeyTap("enter")
		if err != nil {
			return err
		}
		err = robotgo.KeyTap("esc")
		if err != nil {
			return err
		}
		time.Sleep(100 * time.Millisecond)
		err = robotgo.KeyTap("enter")
		if err != nil {
			return err
		}
		time.Sleep(100 * time.Millisecond)
		err = robotgo.KeyTap("enter")
		if err != nil {
			return err
		}
		time.Sleep(100 * time.Millisecond)
		err = robotgo.KeyTap("up")
		if err != nil {
			return err
		}
		time.Sleep(100 * time.Millisecond)
		err = robotgo.KeyTap("up")
		if err != nil {
			return err
		}
		time.Sleep(100 * time.Millisecond)
		err = robotgo.KeyTap("enter")
		if err != nil {
			return err
		}

		//e.GetObservation(i)
	}
	return nil
}

func (e *DoomEnvironment) Step(act, env int) (*Observation, error) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	strAction, exist := GetAction(act)
	if !exist {
		fmt.Println("action not in action space")
	}

	err := e.Act(strAction, env)
	if err != nil {
		fmt.Println(err)
	}

	obs := e.GetObservation(env)

	return obs, nil
}

func (e *DoomEnvironment) GetObservation(env int) *Observation {
	x, y, w, h := robotgo.GetBounds(e.pids[env])
	img := robotgo.CaptureImg(x-10, y-8, w, h-3)

	return &Observation{Image: img}
}

func (e *DoomEnvironment) Save(name string, img image.Image) error {
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

func (e *DoomEnvironment) Act(action string, env int) error {
	err := robotgo.ActivePID(e.pids[env])
	if err != nil {
		return err
	}

	time.Sleep(100 * time.Millisecond)
	err = robotgo.KeyTap(action)
	if err != nil {
		return err
	}

	return nil
}

func (e *DoomEnvironment) Reset() {

}

func (e *DoomEnvironment) Record() {

}

func (e *DoomEnvironment) Close() {
	for _, pid := range e.pids {
		err := robotgo.Kill(pid)
		if err != nil {
			fmt.Println("an error occurred while killing doom", err.Error())
		}
	}
}

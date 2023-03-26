package environment

import (
	"doom-environment/action"
	"doom-environment/observation"
	"errors"
	"fmt"
	"github.com/go-vgo/robotgo"
	"github.com/vcaesar/bitmap"
	"os/exec"
	"strconv"
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
				fmt.Println("an error occurred.", err.Error())
			}
		}()
	}
	time.Sleep(time.Second)

	pids, err := robotgo.FindIds(windowName)
	if err != nil {
		return nil, err
	}
	if len(pids) != numberOfWindows {
		return nil, errors.New("number of process not equal numberOfWindows")
	}

	return &DoomEnvironment{pids: pids}, nil
}

func (e *DoomEnvironment) Start() error {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	for i, pid := range e.pids {
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

		e.GetObservation(i)
	}
	return nil
}

func (e *DoomEnvironment) Step(act, env int) (observation.Observation, error) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	strAction, exist := action.GetAction(act)
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

func (e *DoomEnvironment) GetObservation(env int) observation.Observation {

	x, y, w, h := robotgo.GetBounds(e.pids[env])
	println(x, y, w, h)
	bit := robotgo.CaptureScreen(x-10, y-8, w, h-2)
	bitMap := robotgo.ToBitmap(bit)
	defer robotgo.FreeBitmap(bit)
	bitmap.Save(bit, strconv.Itoa(env)+"_test_1.png")

	return observation.Observation{Bitmap: bitMap}
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

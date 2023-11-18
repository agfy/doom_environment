package doom_environment

import (
	"errors"
	"fmt"
	"github.com/agfy/doom_environment/check_points"
	"github.com/agfy/doom_environment/image_comparer"
	"github.com/go-vgo/robotgo"
	"image"
	"os/exec"
	"strconv"
	"sync"
	"time"
)

const (
	windowName = "prboom-plus"
	width      = 640
	height     = 514
)

type DoomEnvironment struct {
	samples     int
	pids        []int32
	mutex       sync.Mutex
	checkPoints check_points.CheckPoints
	maxScores   []int
}

func Create(numberOfWindows, samples int) (*DoomEnvironment, error) {
	for i := 0; i < numberOfWindows; i++ {
		go func() {
			cmd := exec.Command("prboom-plus", "doom1")
			err := cmd.Run()
			if err != nil {
				fmt.Println("an error occurred while starting doom", err.Error())
			}
		}()
	}

	checkPoints, err := check_points.NewCheckPoints("check_points/loc_1_lvl_1/sample_" + strconv.Itoa(samples) + "/")
	if err != nil {
		fmt.Println("failed to create checkpoints", err.Error())
	}

	numberOfTries := 10
	for i := 0; i < numberOfTries; i++ {
		pids, err := robotgo.FindIds(windowName)
		if err != nil {
			return nil, err
		}
		if len(pids) == numberOfWindows {
			time.Sleep(time.Second)
			return &DoomEnvironment{
				checkPoints: checkPoints,
				pids:        pids,
				maxScores:   make([]int, numberOfWindows),
				samples:     samples,
			}, nil
		}

		time.Sleep(time.Second)
	}

	return nil, errors.New("number of process not equal numberOfWindows")
}

func (e *DoomEnvironment) GetInputNeuronNumber() int {
	return width / e.samples
}

func (e *DoomEnvironment) GetOutputNeuronNumber() int {
	return height / e.samples
}

func (e *DoomEnvironment) Reset() error {
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

func (e *DoomEnvironment) Step(act, env int) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	strAction, exist := GetAction(act)
	if !exist {
		return errors.New("action not in action space")
	}

	err := e.Act(strAction, env)
	if err != nil {
		return err
	}

	return nil
}

func (e *DoomEnvironment) GetObservation(env int) image.Image {
	//x, y, w, h := robotgo.GetBounds(e.pids[env]) causes x11 error "Maximum number of clients reached"
	//img := robotgo.CaptureImg(x-10, y-8, w, h-3)
	img := robotgo.CaptureImg(640, 337, 640, 514)

	return img
}

func (e *DoomEnvironment) GetScore(env int) (int, error) {
	obs := e.GetObservation(env)
	for _, checkPoint := range e.checkPoints.Points {
		eq, err := image_comparer.AreImagesEqual(image_comparer.Samplify(obs, e.samples), checkPoint.Img)
		if err != nil {
			return 0, err
		}
		if eq && checkPoint.Score > e.maxScores[env] {
			e.maxScores[env] = checkPoint.Score
		}
	}

	return e.maxScores[env], nil
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

package doom_environment

import (
	"errors"
	"fmt"
	"github.com/agfy/doom_environment/check_points"
	"github.com/agfy/doom_environment/image_comparer"
	"github.com/go-vgo/robotgo"
	"image"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"time"
)

var ActionSpace = map[int]string{
	0: robotgo.KeyW,  //move forward
	1: robotgo.KeyS,  //move backward
	2: robotgo.KeyA,  //strafe left
	3: robotgo.KeyD,  //strafe right
	4: robotgo.Lctrl, //fire
	5: robotgo.Space, //use
	6: robotgo.Left,  //left arrow
	7: robotgo.Right, //right arrow
}

const (
	windowName     = "prboom-plus"
	width          = 640
	height         = 514
	maxColourValue = 0xFFFF
)

type DoomEnvironment struct {
	samples         int
	pids            []int
	mutex           sync.Mutex
	checkPoints     check_points.CheckPoints
	maxScores       []int
	numberOfWindows int
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

	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	checkPoints, err := check_points.NewCheckPoints(basepath + "/check_points/loc_1_lvl_1/sample_" + strconv.Itoa(samples) + "/")
	if err != nil {
		return nil, fmt.Errorf("failed to create checkpoints %v", err)
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
				checkPoints:     checkPoints,
				pids:            pids,
				maxScores:       make([]int, numberOfWindows),
				samples:         samples,
				numberOfWindows: numberOfWindows,
			}, nil
		}

		time.Sleep(time.Second)
	}

	return nil, errors.New("number of process not equal numberOfWindows")
}

func (e *DoomEnvironment) GetInputNeuronNumber() int {
	return width * height / e.samples
}

func (e *DoomEnvironment) GetOutputNeuronNumber() int {
	return len(ActionSpace)
}

func (e *DoomEnvironment) GetWindowNumber() int {
	return e.numberOfWindows
}

func (e *DoomEnvironment) Reset() error {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	for _, pid := range e.pids {
		time.Sleep(100 * time.Millisecond)
		err := robotgo.KeyTap("enter", pid)
		if err != nil {
			return err
		}

		time.Sleep(100 * time.Millisecond)
		err = robotgo.KeyTap("enter", pid)
		if err != nil {
			return err
		}
		err = robotgo.KeyTap("esc", pid)
		if err != nil {
			return err
		}
		time.Sleep(100 * time.Millisecond)
		err = robotgo.KeyTap("enter", pid)
		if err != nil {
			return err
		}
		time.Sleep(100 * time.Millisecond)
		err = robotgo.KeyTap("enter", pid)
		if err != nil {
			return err
		}
		time.Sleep(100 * time.Millisecond)
		err = robotgo.KeyTap("up", pid)
		if err != nil {
			return err
		}
		time.Sleep(100 * time.Millisecond)
		err = robotgo.KeyTap("up", pid)
		if err != nil {
			return err
		}
		time.Sleep(100 * time.Millisecond)
		err = robotgo.KeyTap("enter", pid)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *DoomEnvironment) Step(act, env int) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	strAction, exist := e.GetAction(act)
	if !exist {
		return errors.New("action not in action space")
	}

	err := e.Act(strAction, env)
	if err != nil {
		return err
	}

	return nil
}

func (e *DoomEnvironment) GetImage(env int) image.Image {
	//x, y, w, h := robotgo.GetBounds(e.pids[env]) causes x11 error "Maximum number of clients reached"
	//img := robotgo.CaptureImg(x-10, y-8, w, h-3)
	img := robotgo.CaptureImg(640, 337, width, height)

	return img
}

func (e *DoomEnvironment) GetObservation(env int) []int {
	img := e.GetImage(env)
	bounds := img.Bounds()
	result := make([]int, 3*bounds.Max.Y*bounds.Max.X)
	var r, g, b uint32
	for x := 0; x < bounds.Max.X; x++ {
		for y := 0; y < bounds.Max.Y; y++ {
			r, g, b, _ = img.At(x, y).RGBA()
			result[3*(x+y*width)] = int(9 * r / maxColourValue)
			result[3*(x+y*width)+1] = int(9 * g / maxColourValue)
			result[3*(x+y*width)+2] = int(9 * b / maxColourValue)
		}
	}

	return result
}

func (e *DoomEnvironment) GetScore(env int) (int, error) {
	obs := e.GetImage(env)
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
	err := robotgo.KeyTap(action, e.pids[env])
	if err != nil {
		return err
	}

	return nil
}

func (e *DoomEnvironment) Record() {
	println("useless commit")
}

func (e *DoomEnvironment) Close() {
	for _, pid := range e.pids {
		err := robotgo.Kill(pid)
		if err != nil {
			fmt.Println("an error occurred while killing doom", err.Error())
		}
	}
}

func (e *DoomEnvironment) GetAction(action int) (string, bool) {
	v, ok := ActionSpace[action]
	if !ok {
		return "", false
	}
	return v, true
}

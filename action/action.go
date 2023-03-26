package action

import "github.com/go-vgo/robotgo"

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

func GetAction(action int) (string, bool) {
	v, ok := ActionSpace[action]
	if !ok {
		return "", false
	}
	return v, true
}

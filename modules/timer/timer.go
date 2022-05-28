package timer

import (
	"fmt"
	timer "github.com/tiniyo/timer"
)


var VerifyTimer *timer.TiniyoTimer

func Timer(cb func(data interface{}) error) *timer.TiniyoTimer {
	if VerifyTimer !=  nil{
		fmt.Println("Timer is not null returning previously created timer")
		return VerifyTimer
	}
	var NewTimer timer.TiniyoTimer
	NewTimer.InitializeTiniyoTimer(cb)
	NewTimer.Run()
	VerifyTimer = &NewTimer
	return VerifyTimer
}

package main

import (
	"fmt"
	"time"

	timermgr "github.com/Lei2050/timer-manager"
)

func main() {
	timerMgr := timermgr.New()
	canceler := timerMgr.NewCanceler()
	var needCancel timermgr.TimerID
	tick := time.NewTicker(time.Millisecond * 100)
	timerMgr.AddTimer(time.Second, func(a ...any) {
		fmt.Printf("1 second timer, args:%+v, %+v\n", a, time.Now())
	}, nil, 1, "hello")
	timerMgr.AddTimer(time.Second*2, func(a ...any) {
		fmt.Printf("2 second timer, args:%+v, %+v\n", a, time.Now())
	}, nil, 2, "hello")
	timerMgr.AddTimer(time.Second*3, func(a ...any) {
		fmt.Printf("3 second timer, args:%+v, %+v\n", a, time.Now())
	}, nil, 3, "hello")

	timerMgr.AddTimer(time.Second*100, func(a ...any) {
		fmt.Printf("100 second timer, args:%+v, %+v\n", a, time.Now())
	}, canceler, 100, "hello")

	timerMgr.AddTimer(time.Second*4, func(a ...any) {
		fmt.Printf("4 second timer, args:%+v, %+v\n    going to add repeat timer...\n", a, time.Now())

		timerMgr.AddRepeatTimer(time.Second*1, func(a ...any) {
			fmt.Printf("1 second repeat timer, args:%+v, %+v\n", a, time.Now())
		}, canceler, 1, "repeat")
		timerMgr.AddRepeatTimer(time.Second*2, func(a ...any) {
			fmt.Printf("2 second repeat timer, args:%+v, %+v\n", a, time.Now())
		}, canceler, 2, "repeat")
		needCancel = timerMgr.AddRepeatTimer(time.Second*3, func(a ...any) {
			fmt.Printf("3 second repeat timer, args:%+v, %+v\n", a, time.Now())
		}, nil, 3, "repeat")

		fmt.Printf("OHHHHHH\n")

	}, nil, "add repeat timers")

	timerMgr.AddTimer(time.Second*60, func(a ...any) {
		fmt.Printf("60 second timer, going to cancel, args:%+v, %+v\n", a, time.Now())
		canceler.Cancel()
		timerMgr.CancelTimer(needCancel)
		canceler.Cancel()
		timerMgr.CancelTimer(needCancel)
		canceler.Cancel()
		timerMgr.CancelTimer(needCancel)
	}, canceler, 60, "hello")

	for c := range tick.C {
		timerMgr.Tick(c)
	}
}

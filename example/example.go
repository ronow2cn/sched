/*
* @Author: huang
* @Date:   2017-09-07 16:00:19
* @Last Modified by:   huang
* @Last Modified time: 2017-09-07 16:31:08
 */
package main

import (
	"fmt"
	"github.com/ronow2cn/sched/asyncop"
	"github.com/ronow2cn/sched/loop"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func OnSignal(f func(os.Signal)) {
	go func() {
		c := make(chan os.Signal, 8)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

		for s := range c {
			f(s)
			if s != syscall.SIGHUP {
				break
			}
		}
	}()
}

func main() {
	quit := make(chan int)
	OnSignal(func(s os.Signal) {
		fmt.Println("shutdown signal received ...")
		close(quit)
	})

	start()
	<-quit
	stop()
}

func start() {
	fmt.Println("start time ", time.Now())
	//begin asyncop and loop
	asyncop.Start()
	loop.Run()

	handler()
}

func handler() {
	//loop thread
	loop.Push(func() {
		fmt.Println("loop push...")
	})

	//asyncop background thread
	asyncop.Push(func() {
		fmt.Println("asyncop push...")
	}, nil)

	//loop timer func
	timerTestLoop(time.Now())
}

func timerTestLoop(ts time.Time) {
	nextTime := nextTimerTest(ts)

	loop.SetTimeOut(nextTime, func() {
		fmt.Println("Timer func in loop", time.Now())
		timerTestLoop(nextTime)
	})
}

//get next timer call time
func nextTimerTest(ts time.Time) time.Time {
	return ts.Add(time.Duration(3) * time.Second)
}

func stop() {
	//end asyncop and loop
	asyncop.Stop()
	loop.Stop()

	asyncop.Close()
}

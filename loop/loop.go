/*
* @Author: huang
* @Date:   2017-09-07 11:09:50
* @Last Modified by:   huang
* @Last Modified time: 2017-09-07 16:41:19
 */
package loop

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// ============================================================================
const (
	MAXLOOPCHANNUM = 100000 //set max q cache by yourself
)

var (
	q      = make(chan func(), MAXLOOPCHANNUM)
	timerq = NewTimerQueue()
	quit   = make(chan int)
	wg     sync.WaitGroup

	numHandled int32
)

// ============================================================================
// listen func by loop.Push
func Run() {
	wg.Add(1)

	go loopFunc()
	go loopTimer()
}

func Stop() {
	close(quit)
	wg.Wait()
}

func Push(f func()) {
	defer func() {
		if err := recover(); err != nil {
			// ignore EPIPE
		}
	}()

	q <- f
}

//set timer to callback func f at time 'ts'
func SetTimeOut(ts time.Time, f func()) *Timer {
	return timerq.SetTimeOut(ts, f)
}

//cancel haven't called timer
func CancelTimer(t *Timer) {
	timerq.Cancel(t)
}

//change timer called time
func UpdateTimer(t *Timer, ts time.Time) {
	timerq.Update(t, ts)
}

func QLen() int32 {
	return int32(len(q))
}

func NumHandled() int32 {
	return atomic.SwapInt32(&numHandled, 0)
}

// ============================================================================
//loop thread, all loop push and timer func are running here
func loopFunc() {
	defer wg.Done()

	for f := range q {
		safeExecute(f)
		atomic.AddInt32(&numHandled, 1)
	}
}

//handle timer func, it will pushed to loop thread when it is expired,
func loopTimer() {
	defer close(q)

	for {
		select {
		case <-quit:
			return

		default:
			Push(func() {
				now := time.Now()
				for timerq.Expire(now) {
				}
			})

			time.Sleep(100 * time.Millisecond)
		}
	}
}

//call func in safe , cause recover
func safeExecute(f func()) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("critical exception:%v", err)
			fmt.Println(Callstack())
		}
	}()

	f()
}

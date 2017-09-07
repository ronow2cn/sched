/*
* @Author: huang
* @Date:   2017-09-07 11:29:56
* @Last Modified by:   huang
* @Last Modified time: 2017-09-07 15:06:08
 */
package loop

import (
	"sync"
	"testing"
	"time"
)

func TestSetTimeOut(T *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	Run() //begin loop packpkg

	t := time.Now()
	T.Log("begin time", t)

	SetTimeOut(t.Add(time.Duration(2)*time.Second), func() {
		T.Log("timeout", time.Now())

		defer wg.Done()
	})

	wg.Wait()

	T.Log("end")
}

func TestUpdateTimer(T *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	t := time.Now()
	T.Log("begin time", t)

	tid := SetTimeOut(t.Add(time.Duration(10)*time.Second), func() {
		T.Log("timeout", time.Now())
		defer wg.Done()
	})

	UpdateTimer(tid, t.Add(time.Duration(1)*time.Second))

	wg.Wait()
	T.Log("end")
}

func TestCancelTimer(T *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	t := time.Now()
	T.Log("begin time", t)

	tid := SetTimeOut(t.Add(time.Duration(2)*time.Second), func() {
		T.Log("timeout1", time.Now())
		defer wg.Done()
	})

	SetTimeOut(t.Add(time.Duration(3)*time.Second), func() {
		T.Log("timeout2", time.Now())
		defer wg.Done()
	})

	SetTimeOut(t.Add(time.Duration(1)*time.Second), func() {
		T.Log("cancel 1", time.Now())
		CancelTimer(tid)
	})

	wg.Wait()
	T.Log("end time", time.Now())
}

func TestPush(T *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	Push(func() {
		T.Log("push func")
		wg.Done()
	})

	Stop() //end loop packpkg
	wg.Wait()
}

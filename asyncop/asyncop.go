/*
* @Author: huang
* @Date:   2017-09-07 11:09:11
* @Last Modified by:   huang
* @Last Modified time: 2017-09-07 16:34:24
 */
package asyncop

import (
	"github.com/ronow2cn/sched/loop"
	"sync"
)

// ============================================================================

var (
	q    = make(chan *asyncOPT, 100000)
	quit = make(chan int)
	wg   sync.WaitGroup
)

// ============================================================================

type asyncOPT struct {
	op func() // run in background thread
	cb func() // run in logic thread //!!NOTICE!!::cb func() can't call asyncop.push func again
}

// ============================================================================
// listen chan data to callback
func Start() {
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			select {
			case <-quit:
				return

			case aop := <-q:
				aop.op()
				if aop.cb != nil {
					loop.Push(aop.cb)
				}
			}
		}
	}()
}

//stop listening by push
func Stop() {
	close(quit)
	wg.Wait()
}

//handle remainind asyncOPT
func Close() {
	close(q)
	for aop := range q {
		aop.op()
		if aop.cb != nil {
			aop.cb() //here is why cb func() can't call asyncop.push func again, because q is closed
		}
	}
}

//add asyncOPT func() to chan q, and callback after
func Push(op func(), cb func()) {
	defer func() {
		if err := recover(); err != nil {
			// ignore EPIPE
		}
	}()

	q <- &asyncOPT{op, cb}
}

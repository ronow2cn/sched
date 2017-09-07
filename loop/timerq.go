/*
* @Author: huang
* @Date:   2017-09-07 11:05:27
* @Last Modified by:   huang
* @Last Modified time: 2017-09-07 11:32:02
 */
package loop

import (
	"container/heap"
	"time"
)

// ============================================================================

type Timer struct {
	ts    time.Time // expire timestamp
	f     func()    // timer callback
	index int       // heap index
	valid bool      // is valid timer?
}

func (self *Timer) Time() time.Time {
	return self.ts
}

// ============================================================================
// heap interface implement, 5 methods
type timerArray []*Timer

func (self timerArray) Len() int {
	return len(self)
}

func (self timerArray) Less(i, j int) bool {
	return self[i].ts.Before(self[j].ts)
}

func (self timerArray) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
	self[i].index, self[j].index = i, j
}

func (self *timerArray) Push(v interface{}) {
	t := v.(*Timer)
	t.index = len(*self)
	t.valid = true
	*self = append(*self, t)
}

func (self *timerArray) Pop() interface{} {
	arr := *self
	n := len(arr)
	t := arr[n-1]
	t.index = -1
	t.valid = false
	*self = arr[:n-1]

	return t
}

// ============================================================================

type TimerQueue struct {
	q timerArray
}

func NewTimerQueue() *TimerQueue {
	tq := &TimerQueue{}
	heap.Init(&tq.q)
	return tq
}

func (self *TimerQueue) SetTimeOut(ts time.Time, f func()) *Timer {
	t := &Timer{
		ts: ts,
		f:  f,
	}

	heap.Push(&self.q, t)
	return t
}

func (self *TimerQueue) Cancel(t *Timer) {
	if t != nil && t.valid && t.index >= 0 && t.index < len(self.q) {
		heap.Remove(&self.q, t.index)
	}
}

func (self *TimerQueue) Update(t *Timer, ts time.Time) {
	if t != nil && t.valid && t.index >= 0 && t.index < len(self.q) {
		t.ts = ts
		heap.Fix(&self.q, t.index)
	}
}

func (self *TimerQueue) Expire(now time.Time) bool {
	if self.q.Len() == 0 {
		return false
	}

	t := self.q[0]

	if t.ts.After(now) {
		return false
	}

	heap.Pop(&self.q)
	t.f()

	return true
}

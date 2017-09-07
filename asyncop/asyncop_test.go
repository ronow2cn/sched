/*
* @Author: huang
* @Date:   2017-09-07 14:41:17
* @Last Modified by:   huang
* @Last Modified time: 2017-09-07 15:00:25
 */
package asyncop

import (
	"github.com/ronow2cn/sched/loop"
	"sync"
	"testing"
)

func TestPush(T *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	loop.Run()
	Start()

	Push(func() {
		T.Log("push background")
	}, func() {
		T.Log("push loop")
		wg.Done()
	})

	wg.Wait()
	T.Log("end")
}

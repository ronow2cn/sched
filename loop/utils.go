/*
* @Author: huang
* @Date:   2017-09-07 11:42:15
* @Last Modified by:   huang
* @Last Modified time: 2017-09-07 11:42:43
 */
package loop

import (
	"runtime/debug"
)

func Callstack() string {
	return string(debug.Stack())
}

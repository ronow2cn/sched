##一种主逻辑协程+后台协程的设计方案

###主逻辑协程----loop

loop包，提供一种处理逻辑协程的方案。并包含一种定时器Timer，可以设定某个时间触发回调函数，达到定时运行某个函数的目的。

**启动时在初始化位置调用函数`loop.Run()`,用来监听需要处理的函数。**

**将需要在该`loop`协程运行的函数f抛到`loop`协程里面运行:**

    loop.Push(f func())
例：

    //loop 协程
    loop.Push(func() {
        fmt.Println("loop push...")
    })

**`loop`协程还支持`timer`定时器函数回调：**

    //在ts时间运行函数f, f函数在loop线程safeExecute()中调用
    loop.SetTimeOut(ts time.Time, f func())

例:

    //需要定时重复运行的函数   
    func timerTestLoop(ts time.Time) {
        nextTime := nextTimerTest(ts)
        loop.SetTimeOut(nextTime, func() {
            fmt.Println("Timer func in loop", time.Now())
            timerTestLoop(nextTime)
        })
    }
    //获取下一个运行的时间
    func nextTimerTest(ts time.Time) time.Time {
        return ts.Add(time.Duration(3) * time.Second)
    }

`timer`对于还没到时间的函数：
可以用函数`func UpdateTimer(t *Timer, ts time.Time)`, 更改定时器触发时间。
可以用函数`func CancelTimer(t *Timer)`, 取消这个定时器。


**程序结束时调用函数`loop.Stop()`。**

***

###后台协程----asyncop

asyncop包，提供一个后台协程方案。在主loop协程中，可以将一些loop不相干的逻辑放到asyncop后台协程处理。

**启动时在初始化位置调用函数`asyncop.Start()`,用来监听需要处理的函数。**

**将需要在该`asyncop`协程运行的函数f抛到`asyncop`协程里面运行:**

    func Push(op func(), cb func()) 
    op:需要在该协程运行的函数；
    cb:在op函数运行结束之后，push回loop协程运行的函数,不传则设置为nil。


例：

    //asyncop background thread
    asyncop.Push(func() {
        fmt.Println("asyncop push...")
    }, nil)

    asyncop.Push(func() {
        fmt.Println(" run in asyncop ...")
    }, func(){
        fmt.Println("run in loop...")
    })

**程序结束时调用`asyncop.Stop()`和`asyncop.Close()`**

`asyncop.Stop()`：停止`Push`加入的函数。
`asyncop.Close()`: 处理`chan q`中还未处理完的。

**注意：**

cb函数不允许有asyncop.Push操作(防止在Close处理剩余cb时,存在asyncop.Push操作，因为此时Push不了)





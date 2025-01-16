package gofaketime

/*
   #include <time.h>
*/
import "C"

import (
	"sync"
	"time"

	"bou.ke/monkey"
)

/*
通过猴子补丁替换time.Now()的方式来支持faketime;
*/

var lockerNow = sync.Mutex{}

func getTimeFromC() int64 {
	ch := make(chan int64, 1)
	go func() {
		var t C.time_t
		ch <- int64(C.time(&t)) // 将时间存储在 t 中
	}()
	return <-ch
}
func fakeTime() time.Time {
	lockerNow.Lock()
	defer lockerNow.Unlock()
	return time.Unix(getTimeFromC(), 0)
}

type FakeTime struct {
	faker *monkey.PatchGuard
}

func NewFakeTime() *FakeTime {
	return &FakeTime{faker: monkey.Patch(time.Now, fakeTime)}
}

func (f *FakeTime) Close() {
	f.faker.Unpatch()
}

func (f *FakeTime) Restore() {
	f.faker.Restore()
}

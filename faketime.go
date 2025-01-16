package gofaketime

/*
   #include <time.h>
*/
import "C"

import (
	"runtime"
	"sync"
	"time"

	"bou.ke/monkey"
)

var lockerNow = sync.Mutex{}

func fakeTime() time.Time {
	lockerNow.Lock()
	defer lockerNow.Unlock()
	// 绑定到当前系统线程，避免切换到 g0
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	return time.Unix(int64(C.time(nil)), 0)

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

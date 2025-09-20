//go:build ios

package application

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework UIKit
#import <Foundation/Foundation.h>
#import <dispatch/dispatch.h>

extern void dispatchOnMainThreadCallback(unsigned int);

static void dispatchOnMainThread(unsigned int id) {
	dispatch_async(dispatch_get_main_queue(), ^{
		dispatchOnMainThreadCallback(id);
	});
}

static bool onMainThread() {
	return [NSThread isMainThread];
}
*/
import "C"
import "fmt"

func (a *iosApp) isOnMainThread() bool {
	return bool(C.onMainThread())
}

func (a *iosApp) dispatchOnMainThread(id uint) {
	fmt.Printf("🔵 [mainthread_ios.go] dispatchOnMainThread(id=%d)\n", id)
	C.dispatchOnMainThread(C.uint(id))
}

//export dispatchOnMainThreadCallback
func dispatchOnMainThreadCallback(callbackID C.uint) {
	fmt.Printf("🔵 [mainthread_ios.go] dispatchOnMainThreadCallback(id=%d)\n", callbackID)
	mainThreadFunctionStoreLock.RLock()
	id := uint(callbackID)
	fn := mainThreadFunctionStore[id]
	if fn == nil {
		Fatal("dispatchCallback called with invalid id: %v", id)
	}
	delete(mainThreadFunctionStore, id)
	mainThreadFunctionStoreLock.RUnlock()
	fn()
}

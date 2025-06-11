//go:build (android || darwin) && !cli

package main

/*
#include <stdint.h>
*/
import "C"
import (
	"sync"
	"unsafe"

	"github.com/freefly-systems/mavp2p-library/router"
)

var (
	rtr *router.Program
	mu  sync.Mutex
)

//export StartRouter
func StartRouter(argc C.int, argv **C.char) C.int {
	mu.Lock()
	defer mu.Unlock()

	if rtr != nil {
		return -1 // already running
	}

	// Convert C args to Go string slice
	args := make([]string, int(argc))
	ptrSize := unsafe.Sizeof(uintptr(0))
	for i := 0; i < int(argc); i++ {
		cstr := *(**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(argv)) + uintptr(i)*ptrSize))
		args[i] = C.GoString(cstr)
	}

	p, err := router.NewProgram(args)
	if err != nil {
		return -2
	}
	rtr = p
	return 0
}

//export StopRouter
func StopRouter() {
	mu.Lock()
	defer mu.Unlock()
	if rtr != nil {
		rtr.Close()
		rtr = nil
	}
}

//export IsRouterRunning
func IsRouterRunning() C.int {
	mu.Lock()
	defer mu.Unlock()
	if rtr != nil {
		return 1 // running
	}
	return 0 // not running
}

func main() {}

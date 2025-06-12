//go:build (android || darwin) && !cli

package main

/*
#include <stdint.h>
*/
import "C"
import (
	"fmt"
	"os"
	"sync"
	"unsafe"

	"github.com/freefly-systems/mavp2p-library/router"
)

var (
	rtr *router.Program
	mu  sync.Mutex
)

//export TestConnection
func TestConnection() C.int {
	return 42 // Simple test - just return a known value
}

//export AddNumbers
func AddNumbers(a C.int, b C.int) C.int {
	return a + b // Another simple test
}

//export StartDefaultRouter
func StartDefaultRouter() C.int {
	// Log to Android logcat
	fmt.Fprintf(os.Stderr, "Starting default router\n")

	mu.Lock()
	defer mu.Unlock()

	if rtr != nil {
		fmt.Fprintf(os.Stderr, "Router already running\n")
		return -1 // already running
	}

	// Default endpoints: UDP server on 14550 and TCP server on 5760
	args := []string{
		"udps:0.0.0.0:14550",
		"udpc:localhost:6001",
	}

	fmt.Fprintf(os.Stderr, "Initializing router with endpoints: %v\n", args)

	p, err := router.NewProgram(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start router: %v\n", err)
		return -2
	}

	fmt.Fprintf(os.Stderr, "Router started successfully\n")
	rtr = p
	return 0
}

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

package main

import (
	"flag"
	"github.com/pzsz/glutils"
	"runtime"
)

func main() {
	flag.Parse()

	runtime.LockOSThread()

	appstatemanager := glutils.GetManager()

	state := NewMCPlayAppState()

	appstatemanager.Setup(state, "Marching cubes")
	appstatemanager.RunLoop()
	appstatemanager.Destroy()
}

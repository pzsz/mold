package main

import (
	"flag"
	"github.com/pzsz/glutils"
	"github.com/pzsz/marchingcubes/rendering"
	"runtime"
)

func main() {
	flag.Parse()

	runtime.LockOSThread()

	appstatemanager := glutils.GetManager()

	state := rendering.NewMCPlayAppState()

	appstatemanager.Setup(state, "Marching cubes")
	appstatemanager.RunLoop()
	appstatemanager.Destroy()
}

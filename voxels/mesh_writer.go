package voxels

import (

)

type MeshWriter interface {
	StartVertex() int
	AddPosition(x,y,z float32)
	AddNormal(x,y,z float32)
	AddColour(r,g,b,a byte)
	AddTexCoord(u, v float32)
	AddIndice3(a,b,c int)
}

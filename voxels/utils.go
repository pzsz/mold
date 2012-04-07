package voxels

import (
	"math"
	v "github.com/pzsz/lin3dmath"
	)

func DrawSphere(store VoxelField, x,y,z float32, radius float32, power int) {
	startX,endX := int(x-radius),int(x+radius+1)
	startY,endY := int(y-radius),int(y+radius+1)
	startZ,endZ := int(z-radius),int(z+radius+1)

	op := func (ix,iy,iz int) int {
		distX := float32(ix) - x
		distY := float32(iy) - y
		distZ := float32(iz) - z
		dist := float32(math.Sqrt(float64(distX*distX+
			distY*distY + distZ*distZ)))
		if dist > radius {
			return 0
		}

		return int(float32(power) * (radius - dist) / radius);
	}

	store.AddValue(v.Vector3i{startX, startY, startZ},
		v.Vector3i{endX, endY, endZ}, op)
}

func DrawGround(store VoxelField, level int) {
	sizeCube := store.Size()
	
	op := func (ix,iy,iz int) int {
		if iy < level {
			return 255
		}
		return 0
	}

	store.AddValue(sizeCube.Start, sizeCube.End, op)
}
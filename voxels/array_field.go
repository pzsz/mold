package voxels

import (
	v "github.com/pzsz/lin3dmath"
)

type ArrayVoxelField struct {
	sizeX, sizeY, sizeZ                int
	translateX, translateY, translateZ int
	data                               []byte
}

func CreateArrayVoxelField(x, y, z, tx, ty, tz int) *ArrayVoxelField {
	return &ArrayVoxelField{
	sizeX: x,
	sizeY: y,
	sizeZ: z,
	translateX: tx,
	translateY: ty,
	translateZ: tz,
	data:  make([]byte, x*y*z)}
}

func (self *ArrayVoxelField) Size() (v.Boxi) {
	return v.Boxi{
		v.Vector3i{self.translateX, self.translateY, self.translateZ},
		v.Vector3i{self.translateX+self.sizeX, self.translateY+self.sizeY, self.translateZ+self.sizeZ}}
}

func (self *ArrayVoxelField) GetValue(x, y, z int) int {
	x -= self.translateX
	y -= self.translateY
	z -= self.translateZ

	if x < 0 || y < 0 || z < 0 {
		return 0
	}
	if x >= self.sizeX || y >= self.sizeY || z >= self.sizeZ {
		return 0
	}

	return int(self.data[x+y*self.sizeX+z*self.sizeX*self.sizeY])
}

func (self *ArrayVoxelField) AddValue(from, to v.Vector3i, eval VoxelFunc) {
	from.SubIP(v.Vector3i{self.translateX, self.translateY, self.translateZ})
	to.SubIP(v.Vector3i{self.translateX, self.translateY, self.translateZ})

	for z := from.Z; z < to.Z; z++ {
		for y := from.Y; y < to.Y; y++ {
			for x := from.X; x < to.X; x++ {
				if x < 0 || y < 0 || z < 0 ||
					x >= self.sizeX || y >= self.sizeY || z >= self.sizeZ {
					continue
				}

				id := x + y*self.sizeX + z*self.sizeX*self.sizeY
				v := int(self.data[id]) + eval(x+self.translateX, y+self.translateY, z+self.translateZ)

				if v < 0 {
					v = 0
				} else if v > 255 {
					v = 255
				}
				self.data[id] = byte(v)
			}
		}
	}
}

func (self *ArrayVoxelField) SetValue(from, to v.Vector3i, eval VoxelFunc) {
	from = from.Sub(v.Vector3i{self.translateX, self.translateY, self.translateZ})
	to = to.Sub(v.Vector3i{self.translateX, self.translateY, self.translateZ})

	for z := from.Z; z < to.Z; z++ {
		for y := from.Y; y < to.Y; y++ {
			for x := from.X; x < to.X; x++ {
				if x < 0 || y < 0 || z < 0 ||
					x >= self.sizeX || y >= self.sizeY || z >= self.sizeZ {
					continue
				}

				id := x + y*self.sizeX + z*self.sizeX*self.sizeY
				v := eval(x, y, z)
				if v < 0 {
					v = 0
				} else if v > 255 {
					v = 255
				}
				self.data[id] = byte(v)
			}
		}
	}
}

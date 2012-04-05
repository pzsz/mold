package voxels

import (
	v "github.com/pzsz/lin3dmath"
	)

type VoxelFunc func (int, int, int) int

// Interface to voxels storage.
// Where voxels are a scalar field with ISOFIELD==127
type VoxelField interface {
	Size() (int,int,int)
	GetValue(x,y,z int) (int)

	AddValue(from, to v.Vector3i, eval VoxelFunc)
	SetValue(from, to v.Vector3i, eval VoxelFunc)
}

// Return gradient vector of voxels at given position
func GetVoxelGradient(storage VoxelField, pos v.Vector3i) (ret v.Vector3f) {
	xl := float32(storage.GetValue(pos.X-1, pos.Y, pos.Z))
	xr := float32(storage.GetValue(pos.X+1, pos.Y, pos.Z))

	yl := float32(storage.GetValue(pos.X, pos.Y-1, pos.Z))
	yr := float32(storage.GetValue(pos.X, pos.Y+1, pos.Z))

	zl := float32(storage.GetValue(pos.X, pos.Y, pos.Z-1))
	zr := float32(storage.GetValue(pos.X, pos.Y, pos.Z+1))

	return v.Vector3f{xr-xl / 2.0, yr-yl / 2.0, zr-zl / 2.0}
}

// Useful abstraction over 8 Voxels formed into cube
// Voxels id ordering:
//   7----6
//  /|   /|
// 4----5 |
// | |  | |
// | 3--|-2
// |/   |/
// 0----1
type VoxelCube struct {
	VoxelField VoxelField
	Pos    v.Vector3i
}

// Get Position of given voxel
func (self *VoxelCube) GetPosition(id int) (v.Vector3i) {
	switch id {
	case 0:
		return self.Pos.Add(v.Vector3i{0,0,0});
	case 1:
		return self.Pos.Add(v.Vector3i{1,0,0});
	case 2:
		return self.Pos.Add(v.Vector3i{1,0,1});
	case 3:
		return self.Pos.Add(v.Vector3i{0,0,1});
	case 4:
		return self.Pos.Add(v.Vector3i{0,1,0});
	case 5:
		return self.Pos.Add(v.Vector3i{1,1,0});
	case 6:
		return self.Pos.Add(v.Vector3i{1,1,1});
	case 7:
		return self.Pos.Add(v.Vector3i{0,1,1});
	}
	panic("Unknown voxel id: ")
}

// Get Position of given voxel
func (self *VoxelCube) GetValue(id int) (int) {
	var x,y,z = self.Pos.X,self.Pos.Y,self.Pos.Z;

	switch id {
	case 0:
		return self.VoxelField.GetValue(x,y,z);			
	case 1:
		return self.VoxelField.GetValue(x+1,y,z);
	case 2:
		return self.VoxelField.GetValue(x+1,y,z+1);
	case 3:
		return self.VoxelField.GetValue(x,y,z+1);
	case 4:
		return self.VoxelField.GetValue(x,y+1,z);
	case 5:
		return self.VoxelField.GetValue(x+1,y+1,z);
	case 6:
		return self.VoxelField.GetValue(x+1,y+1,z+1);
	case 7:
		return self.VoxelField.GetValue(x,y+1,z+1);
	}
	panic("Unknown voxel id: ")
}

func (self *VoxelCube) GetGradient(id int) (v.Vector3f) {
	pos := self.GetPosition(id)
	return GetVoxelGradient(self.VoxelField, pos)
}

type ArrayVoxelField struct {
	sizeX, sizeY, sizeZ int;
	data []byte;
}

func CreateArrayVoxelField(x,y,z int) (*ArrayVoxelField) {
	return &ArrayVoxelField{
	sizeX: x,
	sizeY: y,
	sizeZ: z,
	data: make([]byte, x*y*z)}
}

func (self *ArrayVoxelField) Size() (int,int,int) {
	return self.sizeX, self.sizeY, self.sizeZ;
}

func (self *ArrayVoxelField) GetValue(x,y,z int) (int) {
	if x < 0 || y < 0 || z < 0 { return 0 }
	if x >= self.sizeX || y >= self.sizeY || z >= self.sizeZ { return 0 }

	return int(self.data[x + y * self.sizeX + z * self.sizeX*self.sizeY]);
}

func (self *ArrayVoxelField) AddValue(from, to v.Vector3i, eval VoxelFunc) {
	for z:=from.Z; z < to.Z; z++ {
		for y:=from.Y; y < to.Y; y++ {
			for x:=from.X; x < to.X; x++ {
				if x < 0 || y < 0 || z < 0 || x >= self.sizeX || y >= self.sizeY || z >= self.sizeZ { 
					continue
				}

				id := x + y * self.sizeX + z * self.sizeX*self.sizeY
				v := int(self.data[id]) + eval(x,y,z)

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
	for z:=from.Z; z < to.Z; z++ {
		for y:=from.Y; y < to.Y; y++ {
			for x:=from.X; x < to.X; x++ {				
				if x < 0 || y < 0 || z < 0 || x >= self.sizeX || y >= self.sizeY || z >= self.sizeZ { 
					continue
				}

				id := x + y * self.sizeX + z * self.sizeX*self.sizeY
				v := eval(x,y,z)
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

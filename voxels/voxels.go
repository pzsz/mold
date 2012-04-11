package voxels

import (
	v "github.com/pzsz/lin3dmath"
)

type VoxelFunc func(int, int, int) int

// Interface to voxels storage.
// Where voxels are a scalar field with ISOFIELD==127
type VoxelField interface {
	Size() (v.Boxi)
	GetValue(x, y, z int) int

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

	return v.Vector3f{(xr - xl)/2.0, (yr - yl)/2.0, (zr - zl)/2.0}
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
	Pos        v.Vector3i
}
// Get Position of given voxel
func (self *VoxelCube) GetPosition(id int) v.Vector3i {
	switch id {
	case 0:
		return self.Pos.Add(v.Vector3i{0, 0, 0})
	case 1:
		return self.Pos.Add(v.Vector3i{1, 0, 0})
	case 2:
		return self.Pos.Add(v.Vector3i{1, 0, 1})
	case 3:
		return self.Pos.Add(v.Vector3i{0, 0, 1})
	case 4:
		return self.Pos.Add(v.Vector3i{0, 1, 0})
	case 5:
		return self.Pos.Add(v.Vector3i{1, 1, 0})
	case 6:
		return self.Pos.Add(v.Vector3i{1, 1, 1})
	case 7:
		return self.Pos.Add(v.Vector3i{0, 1, 1})
	}
	panic("Unknown voxel id: ")
}

// Get Position of given voxel
func (self *VoxelCube) GetValue(id int) int {
	var x, y, z = self.Pos.X, self.Pos.Y, self.Pos.Z

	switch id {
	case 0:
		return self.VoxelField.GetValue(x, y, z)
	case 1:
		return self.VoxelField.GetValue(x+1, y, z)
	case 2:
		return self.VoxelField.GetValue(x+1, y, z+1)
	case 3:
		return self.VoxelField.GetValue(x, y, z+1)
	case 4:
		return self.VoxelField.GetValue(x, y+1, z)
	case 5:
		return self.VoxelField.GetValue(x+1, y+1, z)
	case 6:
		return self.VoxelField.GetValue(x+1, y+1, z+1)
	case 7:
		return self.VoxelField.GetValue(x, y+1, z+1)
	}
	panic("Unknown voxel id: ")
}

func (self *VoxelCube) GetGradient(id int) v.Vector3f {
	pos := self.GetPosition(id)
	return GetVoxelGradient(self.VoxelField, pos)
}

func (s *VoxelCube) Interpolate(x, y, z float32) int {
	bottom := float32(s.GetValue(0))*(1-x)*(1-y) +
		float32(s.GetValue(1))*(x)*(1-y) +
		float32(s.GetValue(2))*(x)*(y) +
		float32(s.GetValue(3))*(1-x)*(y)

	top := float32(s.GetValue(4))*(1-x)*(1-y) +
		float32(s.GetValue(5))*(x)*(1-y) +
		float32(s.GetValue(6))*(x)*(y) +
		float32(s.GetValue(7))*(1-x)*(y)
	return int(bottom*z + top*(1-z))
}

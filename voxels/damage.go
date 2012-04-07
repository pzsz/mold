package voxels

import (
	v "github.com/pzsz/lin3dmath"
	)

type DamageFunc func(v.Boxi)

// Voxel field wrapper
type DamageWrapper struct {
	Field       VoxelField
	DamageFunc  DamageFunc
}

func NewDamageWrapper(Field VoxelField, DamageFunc DamageFunc) *DamageWrapper {
	return &DamageWrapper{Field, DamageFunc}
}

func (s *DamageWrapper) Size() (v.Boxi) {
	return s.Field.Size()
}

func (s *DamageWrapper) GetValue(x,y,z int) (int) {
	return s.Field.GetValue(x, y, z)
}

func (s *DamageWrapper) AddValue(from, to v.Vector3i, eval VoxelFunc) {
	s.Field.AddValue(from, to, eval)	
	s.DamageFunc(v.Boxi{from, to})
}

func (s *DamageWrapper) SetValue(from, to v.Vector3i, eval VoxelFunc) {
	s.Field.SetValue(from, to, eval)
	s.DamageFunc(v.Boxi{from, to})
}

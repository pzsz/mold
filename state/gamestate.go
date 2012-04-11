package state

import (
	"github.com/pzsz/marchingcubes/voxels"
	"github.com/pzsz/marchingcubes/wobject"
	)

type GameState struct {
	VoxelField    *voxels.DamageWrapper
	ObjectManager *wobject.WObjectManager

}

func NewGameState() *GameState {
	return &GameState{
		VoxelField: voxels.NewDamageWrapper(voxels.CreateArrayVoxelField(
			256, 128, 256,
			-128, -64, -128), nil),
	ObjectManager: wobject.NewWObjectManager(nil)}
}

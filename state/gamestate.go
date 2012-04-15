package state

import (
	"github.com/pzsz/mold/voxels"
	"github.com/pzsz/mold/wobject"
	v "github.com/pzsz/lin3dmath"
	)

type GameState struct {
	VoxelField    *voxels.DamageWrapper
	ObjectManager *wobject.WObjectManager

	Player          *wobject.WObject
}

func NewGameState() *GameState {
	state := &GameState{
		VoxelField: voxels.NewDamageWrapper(voxels.CreateArrayVoxelField(
			256, 128, 256,
			-128, -64, -128), nil),
	ObjectManager: wobject.NewWObjectManager()}

	state.ObjectManager.VoxelField = state.VoxelField

	return state
}

func (s *GameState) UpdatePlayerCtrl(moveDir v.Vector3f, horAngle v.Angle) {
	move := wobject.FindSimpleMoveWModule(s.Player)
	move.Speed = moveDir

	s.Player.Rotation = v.QuaternionFromAngle(
		v.Vector3f{0,1,0}, 
		horAngle)
	
}

func (s *GameState) CreatePlayer() {
	wo := wobject.NewWObject(2, 1)
	wo.Modules[0] = wobject.NewSimpleMoveWModule()
	wo.Modules[1] = wobject.NewSimpleCollisionWModule()
	
	s.Player = wo
	s.ObjectManager.RegisterObject(s.Player, v.Vector3f{0,10,0}, nil)
}
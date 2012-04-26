package state

import (
	v "github.com/pzsz/lin3dmath"
	"github.com/pzsz/mold/bulletbridge"
	"github.com/pzsz/mold/voxels"
	"github.com/pzsz/mold/wobject"
)

type GameState struct {
	VoxelField    *voxels.DamageWrapper
	ObjectManager *wobject.WObjectManager

	Player *wobject.WObject

	VoxelsMesh *VoxelsMesh
	Physics    *bulletbridge.BBWorld
}

func NewGameState() *GameState {
	state := &GameState{
		VoxelField: voxels.NewDamageWrapper(voxels.CreateArrayVoxelField(
			256, 128, 256,
			-128, -64, -128), nil),
		ObjectManager: wobject.NewWObjectManager()}

	state.ObjectManager.DataField = state
	state.ObjectManager.VoxelField = state.VoxelField

	state.Physics = bulletbridge.NewBBWorld()

	state.VoxelsMesh = NewVoxelsMesh(
		state.VoxelField,
		VoxelsMeshConfig{
			BlockArraySize: v.Vector3i{16, 16, 16},
			BlockSize:      v.Vector3i{8, 8, 8},
		},
		state.Physics)

	voxels.DrawPerlin(state.VoxelField)
	state.VoxelsMesh.RefreshMesh()

	return state
}

func GetGameState(manager *wobject.WObjectManager) *GameState {
	return manager.DataField.(*GameState)
}

func (s *GameState) UpdatePlayerCtrl(moveDir v.Vector3f, horAngle v.Angle) {
	move := wobject.FindSimpleMoveWModule(s.Player)
	move.Speed = moveDir

	s.Player.Rotation = v.QuaternionFromAngle(
		v.Vector3f{0, 1, 0},
		horAngle)

}

func (s *GameState) CreatePlayer() {
	wo := wobject.NewWObject("player", 2, 1)
	wo.Modules[0] = wobject.NewSimpleMoveWModule()
	wo.Modules[1] = wobject.NewSimpleCollisionWModule()

	s.Player = wo
	s.ObjectManager.RegisterObject(s.Player, v.Vector3f{0, 10, 0}, nil)
}

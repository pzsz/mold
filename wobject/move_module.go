package wobject

import (
	v "github.com/pzsz/lin3dmath"
)

type SimpleMoveWModule struct {
	Object        *WObject
	Speed         v.Vector3f
	SurfaceNormal v.Vector3f
	Gravity       float32
}

func NewSimpleMoveWModule() *SimpleMoveWModule {
	return &SimpleMoveWModule{
		SurfaceNormal: v.Vector3f{0, 1, 0},
		Gravity:       -2,
	}
}

func FindSimpleMoveWModule(Object *WObject) *SimpleMoveWModule {
	for i := 0; i < len(Object.Modules); i++ {
		cast, ok := Object.Modules[i].(*SimpleMoveWModule)
		if ok {
			return cast
		}
	}
	return nil

}

func (self *SimpleMoveWModule) Setup(ob *WObject) {
	self.Object = ob
}

func (self *SimpleMoveWModule) InitNew() {

}

func (self *SimpleMoveWModule) Process(time_step float32) {
	move := self.Speed.Mul(time_step)
	move_quat := v.QuaternionFromAngle(self.SurfaceNormal,
		self.Object.Rotation.GetAngle())
	rot_move := move_quat.Mul3f(move)

//	rot_move.Y += self.Gravity * time_step

	self.Object.MovePosition(rot_move.X, rot_move.Y, rot_move.Z)
}

package wobject

import (
	v "github.com/pzsz/lin3dmath"
	"github.com/pzsz/mold/voxels"
)

type SimpleCollisionWModule struct {
	Object        *WObject
	Speed         v.Vector3f
	SurfaceNormal v.Vector3f
	Gravity       float32
}

func NewSimpleCollisionWModule() *SimpleCollisionWModule {
	return &SimpleCollisionWModule{
		SurfaceNormal: v.Vector3f{0, 1, 0},
		Gravity:       -3,
	}
}

func FindSimpleCollisionWModule(Object *WObject) *SimpleCollisionWModule {
	for i := 0; i < len(Object.Modules); i++ {
		cast, ok := Object.Modules[i].(*SimpleCollisionWModule)
		if ok {
			return cast
		}
	}
	return nil

}

func (self *SimpleCollisionWModule) Setup(ob *WObject) {
	self.Object = ob
}

func (self *SimpleCollisionWModule) InitNew() {

}

func (self *SimpleCollisionWModule) Process(time_step float32) {
	pos := self.Object.Position
	
	voxelField := self.Object.Manager.VoxelField

	x0 := v.FastFloor32(pos.X)
	y0 := v.FastFloor32(pos.Y)
	z0 := v.FastFloor32(pos.Z)

	cube := voxels.VoxelCube{voxelField, v.Vector3i{int(x0), int(y0), int(z0)}}
	val := cube.InterpolateValue(pos.X - x0, pos.Y - y0, pos.Z - z0)

	if val > 128 {
		//self.Object.MovePosition(0, 5*time_step, 0)		
	}
}

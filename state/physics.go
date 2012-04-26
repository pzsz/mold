package state

import (
	v "github.com/pzsz/lin3dmath"
	"github.com/pzsz/mold/bulletbridge"
	"github.com/pzsz/mold/wobject"
	"fmt"
)

type PhysicsSphereWModule struct {
	Object       *wobject.WObject
	RBody        *bulletbridge.BBRigidBody 
}

func NewPhysicsSphereWModule() *PhysicsSphereWModule {
	return &PhysicsSphereWModule{}
}

func FindPhysicsSphereWModule(Object *wobject.WObject) *PhysicsSphereWModule {
	for i := 0; i < len(Object.Modules); i++ {
		cast, ok := Object.Modules[i].(*PhysicsSphereWModule)
		if ok {
			return cast
		}
	}
	return nil

}

func (self *PhysicsSphereWModule) Setup(ob *wobject.WObject) {
	self.Object = ob
	state := GetGameState(self.Object.Manager)
	self.RBody = state.Physics.NewRigidBody(bulletbridge.NewCShapeSphere(1), 1, 
		v.Vector3f{0,10,0})
}

func (self *PhysicsSphereWModule) InitNew() {

}

func (self *PhysicsSphereWModule) Process(time_step float32) {
	self.Object.SetPositionv3f(self.RBody.GetPosition())
	fmt.Printf("%v\n",self.Object.Position)
}

func CreateBall(pos v.Vector3f) *wobject.WObject {
	wo := wobject.NewWObject("ball",1, 1)
	wo.Modules[0] = NewPhysicsSphereWModule()

	return wo
}

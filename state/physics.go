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

type CharacterControlerWModule struct {
	Object       *wobject.WObject
	Controler    *bulletbridge.BBCharacterControler
}

func NewCharacterControlerWModule() *CharacterControlerWModule {
	return &CharacterControlerWModule{}
}


func FindCharacterControlerWModule(Object *wobject.WObject) *CharacterControlerWModule {
	for i := 0; i < len(Object.Modules); i++ {
		cast, ok := Object.Modules[i].(*CharacterControlerWModule)
		if ok {
			return cast
		}
	}
	return nil
}


func (self *CharacterControlerWModule) Setup(ob *wobject.WObject) {
	self.Object = ob
	state := GetGameState(self.Object.Manager)
	self.Controler = state.Physics.NewCharacterControler(1, 3, self.Object.Position)
}

func (self *CharacterControlerWModule) InitNew() {

}

func (self *CharacterControlerWModule) Process(time_step float32) {
	self.Object.SetPositionv3f(self.Controler.GetPosition())
}

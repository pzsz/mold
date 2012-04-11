package wobject

import (
	v "github.com/pzsz/lin3dmath"
)

const FRICTION_ACC = 2
const MIN_SPEED = 0.01
const COLLISION_SPEED_DAMPING = 0.99
const COLLISION_SPEED_IMPULSE_DAMPING = 0.5
const GRAVITY = 4

type MoveWModule struct {
	Object       *WObject
	Acc          v.Vector3f
	Speed        v.Vector3f
	SpeedMax     float32
	SpeedImpulse v.Vector3f
	Immobilized  bool
	SpeedFactor  float32
}

func NewMoveWModule(speedmax float32) *MoveWModule {
	return &MoveWModule{SpeedMax: speedmax, SpeedFactor: 1}
}

func FindMoveWModule(Object *WObject) *MoveWModule {
	for i := 0; i < len(Object.Modules); i++ {
		cast, ok := Object.Modules[i].(*MoveWModule)
		if ok {
			return cast
		}
	}
	return nil

}

func (self *MoveWModule) Setup(ob *WObject) {
	self.Object = ob
}

func (self *MoveWModule) InitNew() {

}

func (self *MoveWModule) Process(time_step float32) {
	if self.Immobilized {
		return
	}

	delta_accx := self.Acc.X * time_step
	delta_accy := self.Acc.Y * time_step
	delta_accz := self.Acc.Z * time_step

	self.Speed.X += delta_accx
	self.Speed.Y += delta_accy
	self.Speed.Z += delta_accz

	speed_len := self.Speed.Len()
	self.Speed.NormalizeIP()

	max_speed := self.SpeedMax * self.SpeedFactor

	if speed_len > max_speed {
		speed_len = max_speed
	}
	if speed_len < MIN_SPEED {
		speed_len = 0
		self.Speed.ZeroIP()
	} else {
		speed_len -= FRICTION_ACC * time_step
	}

	self.Speed.MulIP(speed_len)

	delta_x := self.Speed.X * time_step
	delta_y := self.Speed.Y * time_step
	delta_z := self.Speed.Z * time_step
	
	self.Object.MovePosition(delta_x, delta_y, delta_z)
}

func (self *MoveWModule) Stop() {
	self.Acc.X = 0
	self.Acc.Y = 0
	self.Acc.Z = 0

	self.Speed = v.Vector3f{0, 0, 0}
	self.SpeedImpulse = v.Vector3f{0, 0, 0}
}

func (self *MoveWModule) SetAcc(power float32, direction v.Degree) {
	self.Acc.X = direction.X() * power
	self.Acc.Y = direction.Y() * power
}

func (self *MoveWModule) SetAccDirect(x, y float32) {
	self.Acc.X = x
	self.Acc.Y = y
}

func (self *MoveWModule) AddImpulse(im v.Vector3f) {
	self.SpeedImpulse.AddIP(im)
}

func (self *MoveWModule) Immobilize(im bool) {
	self.Immobilized = im
	if im {
		self.Stop()
	}
}

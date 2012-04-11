package wobject

import (

)

type HealthWModule struct {
	object        *WObject
	health        int
	HealthMax     int
	DamageDefence int
}

type ICauseHealthDeathListener interface {
	OnCauseHealthDeath(ob *WObject)
}

func NewHealthWModule(health int) *HealthWModule {
	return &HealthWModule{health: health, HealthMax: health}
}

func FindHealthWModule(object *WObject) *HealthWModule {
	for i := 0; i < len(object.Modules); i++ {
		cast, ok := object.Modules[i].(*HealthWModule)
		if ok {
			return cast
		}
	}
	return nil
}

func DealDamage(object *WObject, damage int, origin *WObject) bool {
	hm := FindHealthWModule(object)
	if hm != nil {
		hm.Damage(damage, origin)
		return true
	}
	return false
}

func (self *HealthWModule) Setup(ob *WObject) {
	self.object = ob
}

func (self *HealthWModule) InitNew() {

}

func (self *HealthWModule) Process(time_step float32) {

}

func (self *HealthWModule) HealthRatio() float32 {
	return float32(self.health) / float32(self.HealthMax)
}

func (self *HealthWModule) SetHealthMax(v int) {
	self.HealthMax = v
	if self.health > self.HealthMax {
		self.health = self.HealthMax
	}
}

func (self *HealthWModule) GetHealth() int {
	return self.health
}

func (self *HealthWModule) GetHealthMax() int {
	return self.HealthMax
}

func (self *HealthWModule) SetFullHealth() {
	self.health = self.HealthMax
}

func (self *HealthWModule) Damage(val int, origin *WObject) {
	if val < 0 {
		val += self.DamageDefence
		if val >= 0 {
			return
		}
	}

	self.health -= val
	if self.health <= 0 {
		self.object.Die()

		if origin != nil {
			for _, mod := range origin.Modules {
				listener, ok := mod.(ICauseHealthDeathListener)
				if ok {
					listener.OnCauseHealthDeath(self.object)
				}
			}
		}
	}
	if self.health > self.HealthMax {
		self.health = self.HealthMax
	}
}

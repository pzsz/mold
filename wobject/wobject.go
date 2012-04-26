package wobject

import (
	v "github.com/pzsz/lin3dmath"
)

type IWObjectOwner interface {
	GetName() string
}

type WObjectType string

type WObjectCriterionFunc func(*WObject) bool

type WModule interface {
	InitNew()
	Setup(ob *WObject)
	Process(time_step float32)
}

type WDestroyableModule interface {
	WModule
	Destroy()
}

type WModuleRenderer interface {
	WModule
	Render()
}

type WObject struct {
	Id         int32
	Name       string
	Position   v.Vector3f
	Rotation   v.Quaternion
	LookVector v.Vector3f
	Size       float32
	Alive      bool
	Manager    *WObjectManager
	Owner      IWObjectOwner

	Modules        []WModule
	RendererModule WModuleRenderer
}

func (self *WObject) setup(world *WObjectManager,
	id int32) {
	self.Manager = world
	self.Id = id
	self.Alive = true
}

func (self *WObject) initNew(position v.Vector3f,
	owner IWObjectOwner,
	renderer WModuleRenderer) {

	self.Position = position
	self.Owner = owner

	// Append renderer
	self.RendererModule = renderer
	if renderer != nil {
		self.Modules = append(self.Modules, renderer)
	}

	// Setup and init modules
	for i := 0; i < len(self.Modules); i++ {
		self.Modules[i].Setup(self)
	}
	for i := 0; i < len(self.Modules); i++ {
		self.Modules[i].InitNew()
	}

	// Register in bucket
	//	bx, by := self.Manager.Buckets.GetBucketCoord(self.Position)
	//	self.Manager.Buckets.AddWObject(self, 
	//		bx, by)
}

func (self *WObject) initRestored(renderer WModuleRenderer) {
	// Append renderer
	self.RendererModule = renderer
	if renderer != nil {
		self.Modules = append(self.Modules, renderer)
	}

	// Setup modules
	for i := 0; i < len(self.Modules); i++ {
		self.Modules[i].Setup(self)
	}

	// Register in bucket
	//	bx, by := self.Manager.Buckets.GetBucketCoord(self.Position)
	//	self.Manager.Buckets.AddWObject(self, 
	//		bx, by)
}

func (self *WObject) attachRenderer(renderer WModuleRenderer) {
	if self.RendererModule != nil {
		panic("Wtf man, renderer already attached")
	}
	self.RendererModule = renderer
	self.Modules = append(self.Modules, renderer)
	renderer.Setup(self)
}

/** Called by WObjectManager after finding out that object died */
func (self *WObject) destroy() {
	for i := 0; i < len(self.Modules); i++ {
		dmod, cast := self.Modules[i].(WDestroyableModule)
		if cast {
			dmod.Destroy()
		}
	}

	//	bx, by := self.Manager.Buckets.GetBucketCoord(self.Position)
	//	self.Manager.Buckets.RemoveWObject(self, 
	//		bx, by)
}

func (self *WObject) Process(time_step float32) {
	for i := 0; i < len(self.Modules); i++ {
		self.Modules[i].Process(time_step)
		if !self.Alive {
			break
		}
	}
}

func (self *WObject) Die() {
	self.Alive = false
}

func NewWObject(name string, modulesNum int, size float32) *WObject {
	return &WObject{
		Name:    name,
		Size:    size,
		Modules: make([]WModule, modulesNum)}
}

func (self *WObject) GetDistance2(other *WObject) (vec v.Vector3f, dist2 float32) {
	vec = other.Position.Sub(self.Position)
	dist2 = vec.Len2()
	return
}

func (self *WObject) GetDistance(other *WObject) (vec v.Vector3f, dist2 float32) {
	vec = other.Position.Sub(self.Position)
	dist2 = vec.Len()
	return
}

func (self *WObject) MovePosition(x, y, z float32) {
	nx := self.Position.X + x
	ny := self.Position.Y + y
	nz := self.Position.Z + z
	self.SetPosition(nx, ny, nz)
}

func (self *WObject) SetPosition(x, y, z float32) {
	self.Position = v.Vector3f{x, y, z}
}

func (self *WObject) SetPositionv3f(p v.Vector3f) {
	self.Position = p
}

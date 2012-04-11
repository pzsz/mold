package wobject

import (
	"container/list"
	v "github.com/pzsz/lin3dmath"
)

type IWObjectWorld interface {
	GetWorldSizeF() (float32, float32)
	GetWorldSizeI() (int, int)
}

type WObjectManager struct {
	Objects          *list.List	
	id_gen           int32
	RendererFactory  WModuleRendererFactory
	DataField        interface{}
}

func NewWObjectManager(DataField interface{}) (*WObjectManager) {
//	buckets := NewWObjectBuckets()

	return &WObjectManager{ 
	Objects: list.New(),
//	Buckets : buckets,
	DataField: DataField}
}

/* Register new object */
func (self *WObjectManager) RegisterObject(pObject *WObject, position v.Vector3f, owner IWObjectOwner) {
	self.Objects.PushBack(pObject)
	
	// Create renderer
	var renderer WModuleRenderer
	if self.RendererFactory != nil {
		renderer = self.RendererFactory.CreateWObjectRenderer(pObject)
	}

	pObject.setup(self, self.id_gen)
	self.id_gen+=1

	pObject.initNew(position, owner, renderer)
}

/* Register loaded object - there will be no new-object initialization performed */
func (self *WObjectManager) RestoreObject(pObject *WObject) {
	self.Objects.PushBack(pObject)
	
	// Create renderer
	var renderer WModuleRenderer
	if self.RendererFactory != nil {
		renderer = self.RendererFactory.CreateWObjectRenderer(pObject)
	}

	pObject.setup(self, self.id_gen)
	self.id_gen+=1
	
	pObject.initRestored(renderer)
}

func (self *WObjectManager) SetRenderer(RendererFactory WModuleRendererFactory) {
	self.RendererFactory = RendererFactory
	self.RendererFactory.SetupManager(self)

	// Add renderers to objects
	for e := self.Objects.Front(); e != nil; e = e.Next() {
		var object *WObject = e.Value.(*WObject)
		renderer := self.RendererFactory.CreateWObjectRenderer(object)
		object.attachRenderer(renderer)
	}
}

func (self *WObjectManager) Process(time_step float32) {
	for e := self.Objects.Front(); e != nil; {
		var object *WObject = e.Value.(*WObject)
		object.Process(time_step)
		if !object.Alive {
			// Remove object
			next_e := e.Next()
			object.destroy()
			self.Objects.Remove(e)
			e = next_e
		} else {
			e = e.Next() 
		}
	}
}

func (self *WObjectManager) NormalizePosition(p v.Vector2f) v.Vector2f {
	if p.X < 0 {p.X = 0}
	if p.Y < 0 {p.Y = 0}

	/*w,h := self.GetWorldSizeF()
	if p.X >= w {p.X = w - 0.01}
	if p.Y >= h {p.Y = h - 0.01}*/

	return p
}

func (self *WObjectManager) DestroyAllObjects() {
	for e := self.Objects.Front(); e != nil; e = e.Next() {
		var object *WObject = e.Value.(*WObject)
		object.Die()
		object.destroy()
	}
	self.Objects = list.New()
}
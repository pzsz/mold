package state

import (
	v "github.com/pzsz/lin3dmath"
	"github.com/pzsz/glutils"
	"github.com/pzsz/mold/wobject"
)

type SimpleRendererWModule struct {
	Object         *wobject.WObject
	Mesh           *glutils.MeshBuffer
	Context        *RenderContext
}

func (s *SimpleRendererWModule)	InitNew() {

}

func (s *SimpleRendererWModule)	Setup(ob *wobject.WObject) {
	s.Object = ob
}

func (s *SimpleRendererWModule)	Process(time_step float32) {
	
}

func (s *SimpleRendererWModule)	Render() {
	rop := glutils.NewSimpleRenderOp(false, s.Mesh)
	pos := s.Object.Position
	matrix := v.MatrixTranslate(pos.X, pos.Y, pos.Z)
	rop.Render(s.Context.Camera, matrix)
}


type RenderContext struct {
	BoxMesh       *glutils.MeshBuffer
	Camera        *glutils.Camera
}

func NewRenderContext(Camera *glutils.Camera) *RenderContext {
	context := &RenderContext{Camera: Camera}
	context.BoxMesh = glutils.BuildCubeBuffer(v.Vector3f{0.5, 0.5, 0.5})
	
	return context
}

func (s *RenderContext) SetupManager(man *wobject.WObjectManager) {
	
}

func (s *RenderContext) CreateWObjectRenderer(obj *wobject.WObject) wobject.WModuleRenderer {
	if obj.Name == "ball" {
		return &SimpleRendererWModule{Context: s, Mesh:s.BoxMesh}
	}
	return nil
}


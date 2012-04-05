package rendering

import (
	//	"github.com/pzsz/gl"
	"github.com/banthar/Go-SDL/sdl"
	"github.com/pzsz/glutils"
	v "github.com/pzsz/lin3dmath"
	"github.com/pzsz/marchingcubes/voxels"

)

type MCPlayAppState struct {
	Manager   *glutils.AppStateManager
	Camera    *glutils.Camera
	Voxels    voxels.VoxelField
	Renderer  *VoxelsRenderer

	shader    *glutils.ShaderProgram
	bindFunc  func(prog *glutils.ShaderProgram)

	pos       v.Vector3f
	moveDir   v.Vector3f
}

func NewMCPlayAppState() *MCPlayAppState {
	return &MCPlayAppState{}
}

func (self *MCPlayAppState) Setup(manager *glutils.AppStateManager) {
	self.Manager = manager

	self.Camera = glutils.NewCamera(glutils.GetViewport())
	self.Camera.SetFrustrumProjection(60, 0.1, 100)

	storage := voxels.CreateArrayVoxelField(1024, 1024, 128)
	self.Voxels = storage
	voxels.DrawSphere(storage, 6, 6, 0, 6, 250)
	voxels.DrawSphere(storage, 20, 6, 0, 6, 180)

	var err error
	if self.shader, err = glutils.GetProgram(
		"shaders/blob.vertex", 
		"shaders/blob.fragment"); err != nil {
		panic(err.Error())
	}

	self.bindFunc = func(prog *glutils.ShaderProgram) {
		prog.GetUniform("light0_direction").Uniform3f(1, 1, -1)
	}

	self.Renderer = NewVoxelsRenderer(storage,
		VoxelsRendererConfig{
	            BlockArraySize: v.Vector3i{5,5,5},
                    BlockSize: v.Vector3i{8,8,8},
	})

	self.Renderer.RefreshMesh()
}

func (self *MCPlayAppState) OnViewportResize(x, y float32) {

}

func (self *MCPlayAppState) Destroy() {

}

func (self *MCPlayAppState) Pause() {

}

func (self *MCPlayAppState) Resume() {

}

func (self *MCPlayAppState) Process(time_step float32) {
	glutils.Clear()

	self.pos.AddIP(self.moveDir.Mul(time_step))

	self.Renderer.SetCenter(self.pos)

	self.Camera.SetModelview(
		self.pos.X, self.pos.Y, self.pos.Z+32,
		self.pos.X, self.pos.Y, self.pos.Z,
		0, 1, 0)

	self.Renderer.Render(self.Camera, 
		self.shader,
		self.bindFunc)

	sdl.GL_SwapBuffers()
}

func (self *MCPlayAppState) OnKeyDown(key *sdl.Keysym) {
	switch(key.Sym) {
	case sdl.K_w:
		self.moveDir.Y = 10
		break;
	case sdl.K_s:
		self.moveDir.Y = -10
		break
	case sdl.K_a:
		self.moveDir.X = -10
		break
	case sdl.K_d:
		self.moveDir.X = 10
		break
	}
}

func (self *MCPlayAppState) OnKeyUp(key *sdl.Keysym) {
	switch(key.Sym) {
	case sdl.K_w:
		self.moveDir.Y = 0
		break
	case sdl.K_s:
		self.moveDir.Y = 0
		break
	case sdl.K_a:
		self.moveDir.X = 0
		break
	case sdl.K_d:
		self.moveDir.X = 0
		break
	}
}

func (self *MCPlayAppState) OnMouseMove(x, y float32) {

}

func (self *MCPlayAppState) OnMouseClick(x, y float32, button int, down bool) {
	mpos := self.Camera.ScreenToPlaneXY(x, y, 0)
	voxels.DrawSphere(self.Voxels, mpos.X, mpos.Y, 0, 6, 100)
	self.Renderer.RefreshMesh()
}

func (self *MCPlayAppState) OnSdlEvent(event *sdl.Event) {

}

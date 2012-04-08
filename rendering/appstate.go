package rendering

import (
	//	"github.com/pzsz/gl"
	"github.com/banthar/Go-SDL/sdl"
	"github.com/pzsz/glutils"
	v "github.com/pzsz/lin3dmath"
	"github.com/pzsz/marchingcubes/voxels"
	"fmt"
)

type MCPlayAppState struct {
	Manager   *glutils.AppStateManager
	Camera    *glutils.Camera
	Controller *glutils.FpsController
	Voxels    voxels.VoxelField
	Renderer  *VoxelsRenderer

	shader    *glutils.ShaderProgram
	bindFunc  func(prog *glutils.ShaderProgram)

	pos       v.Vector3f
	moveDir   v.Vector3f

	lastX, lastY  float32    

	frameStats glutils.FrameStats
}

func NewMCPlayAppState() *MCPlayAppState {
	return &MCPlayAppState{}
}

func (self *MCPlayAppState) Setup(manager *glutils.AppStateManager) {
	self.Manager = manager

	self.Camera = glutils.NewCamera(glutils.GetViewport())

	self.Controller = glutils.NewFpsController(self.Camera)
	self.Controller.Pos = v.Vector3f{0, 0, -30}
	self.Camera.SetFrustrumProjection(60, 0.1, 100)

	storage := voxels.NewDamageWrapper(voxels.CreateArrayVoxelField(
		1024, 128, 1024,
		-512, -64, -512), nil)

	self.Voxels = storage

	var err error
	if self.shader, err = glutils.GetProgram(
		"shaders/blob.vertex", 
		"shaders/blob.fragment"); err != nil {
		panic(err.Error())
	}

	self.bindFunc = func(prog *glutils.ShaderProgram) {
		prog.GetUniform("light0_direction").Uniform3f(1, -1, -1)
	}

	self.Renderer = NewVoxelsRenderer(storage,
		VoxelsRendererConfig{
	            BlockArraySize: v.Vector3i{32,8,32},
                    BlockSize: v.Vector3i{8,8,8},
	})

	voxels.DrawWave(storage)

	voxels.DrawSphere(storage, 0, 0, 0, 6, 250)

	voxels.DrawSphere(storage, 10, 0, 0, 6, 250)

	self.Renderer.RefreshMesh()

	sdl.ShowCursor(0)
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

	self.Controller.MoveBy(self.moveDir.Y * time_step, 
		self.moveDir.X * time_step)

	self.Renderer.SetCenter(self.Controller.Pos)

	self.Controller.SetupCamera()

	self.Renderer.Render(self.Camera, 
		self.shader,
		self.bindFunc)

	sdl.GL_SwapBuffers()

	if self.frameStats.FrameFinished(int64(time_step * 1000)) {
		fmt.Printf("%v\n", &self.frameStats)
	}
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
	dx := self.lastX - x 
	dy := self.lastY - y

	if dx < 50 && dx > -50 && dy < 50 && dy > -50 {
		self.Controller.RotateBy(dx * -0.005, dy * -0.005)
	}

	self.lastX = x 
	self.lastY = y

}

func (self *MCPlayAppState) OnMouseClick(x, y float32, button int, down bool) {
	mpos := self.Camera.ScreenToSphere(x,y, 10)
	//mpos := self.Camera.ScreenToPlaneXY(x, y, 0)

//	voxels.DrawSphere(self.Voxels, mpos.X, mpos.Y, 0, 6, 100)
	voxels.DrawSphere(self.Voxels, mpos.X, mpos.Y, mpos.Z, 4, 100)
	
	//voxels.DrawSphere(self.Voxels, 0, 0, 0, 6, 100)
	self.Renderer.RefreshMesh()
}

func (self *MCPlayAppState) OnSdlEvent(event *sdl.Event) {

}

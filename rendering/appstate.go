package rendering

import (
	//	"github.com/pzsz/gl"
	"github.com/banthar/Go-SDL/sdl"
	"github.com/pzsz/glutils"
	v "github.com/pzsz/lin3dmath"
	"github.com/pzsz/mold/voxels"
	"github.com/pzsz/mold/state"
	"fmt"
)

type MCPlayAppState struct {
	Manager       *glutils.AppStateManager
	moveDir       v.Vector3f

	frameStats    glutils.FrameStats

	gameState       *state.GameState

	Camera          *glutils.Camera
	Controller      *glutils.FpsController
	VoxelsRenderer  *VoxelsRenderer
	shader          *glutils.ShaderProgram
	bindFunc        func(prog *glutils.ShaderProgram)
}

func NewMCPlayAppState() *MCPlayAppState {
	return &MCPlayAppState{}
}

func (self *MCPlayAppState) Setup(manager *glutils.AppStateManager) {
	self.Manager = manager

	self.Camera = glutils.NewCamera(glutils.GetViewport())
	self.Camera.SetFrustrumProjection(60, 0.1, 200)

	self.Controller = glutils.NewFpsController(self.Camera)
	self.Controller.Pos = v.Vector3f{0, 0, -30}

	self.gameState = state.NewGameState()

	var err error
	if self.shader, err = glutils.GetProgram(
		"shaders/blob.vertex", 
		"shaders/blob.fragment"); err != nil {
		panic(err.Error())
	}
	self.bindFunc = func(prog *glutils.ShaderProgram) {
		prog.GetUniform("light0_direction").Uniform3f(0, -1, 0)
	}

	self.VoxelsRenderer = NewVoxelsRenderer(
		self.gameState.VoxelField,
		VoxelsRendererConfig{
	            BlockArraySize: v.Vector3i{16,16,16},
                    BlockSize: v.Vector3i{8,8,8},
	})

	voxels.DrawPerlin(self.gameState.VoxelField)
	voxels.DrawSphere(self.gameState.VoxelField, 0, 0, 0, 6, 250)
	voxels.DrawSphere(self.gameState.VoxelField, 10, 0, 0, 6, 250)
	self.VoxelsRenderer.RefreshMesh()

	self.gameState.CreatePlayer()

	sdl.ShowCursor(1)
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

	self.gameState.UpdatePlayerCtrl(
		self.moveDir, v.Angle(self.Controller.HorAxis))

	self.gameState.ObjectManager.Process(time_step)

	self.Controller.Pos = self.gameState.Player.Position

	self.Controller.SetupCamera()

	self.VoxelsRenderer.SetCenter(
		self.Controller.Pos)

	self.VoxelsRenderer.Render(self.Camera, 
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
		self.moveDir.Z = -10
		break;
	case sdl.K_s:
		self.moveDir.Z = 10
		break
	case sdl.K_a:
		self.moveDir.X = -10
		break
	case sdl.K_d:
		self.moveDir.X = 10
		break
	case sdl.K_q:
		self.moveDir.Y = 10
		break
	case sdl.K_e:
		self.moveDir.Y = -10
		break

	}
}

func (self *MCPlayAppState) OnKeyUp(key *sdl.Keysym) {
	switch e := key.Sym; {
	case e == sdl.K_w || e == sdl.K_s:
		self.moveDir.Z = 0
		break
	case e == sdl.K_a || e == sdl.K_d:
		self.moveDir.X = 0
		break
	case e == sdl.K_q || e == sdl.K_e:
		self.moveDir.Y = 0
		break
	}
}

func (self *MCPlayAppState) OnMouseMove(x, y, dx, dy float32) {	
	self.Controller.RotateBy(dx * 0.005, dy * 0.005)
}

func (self *MCPlayAppState) OnMouseClick(x, y float32, button int, down bool) {
	vvector := self.Controller.GetViewVector().Mul(10)
	mpos := self.Controller.Pos.Add(vvector)

	voxels.DrawSphere(self.gameState.VoxelField, mpos.X, mpos.Y, mpos.Z, 4, 100)
	
	self.VoxelsRenderer.RefreshMesh()
}

func (self *MCPlayAppState) OnSdlEvent(event *sdl.Event) {

}

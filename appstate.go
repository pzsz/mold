package main

import (
	"github.com/banthar/Go-SDL/sdl"
	"github.com/pzsz/glutils"
	v "github.com/pzsz/lin3dmath"
	"github.com/pzsz/mold/voxels"
	"github.com/pzsz/mold/state"
	"github.com/pzsz/mold/wobject"
	"fmt"
)

type MCPlayAppState struct {
	Manager       *glutils.AppStateManager
	moveDir       v.Vector3f

	frameStats    glutils.FrameStats

	gameState       *state.GameState

	Camera          *glutils.Camera
	Controller      *glutils.FpsController

	shader          *glutils.ShaderProgram
	bindFunc        func(prog *glutils.ShaderProgram)
}

func NewMCPlayAppState() *MCPlayAppState {
	return &MCPlayAppState{}
}

func (self *MCPlayAppState) Setup(manager *glutils.AppStateManager) {
	self.Manager = manager

	self.Camera = glutils.NewCamera(glutils.GetViewport())
	self.Camera.SetFrustrumProjection(60, 0.01, 200)

	self.Controller = glutils.NewFpsController(self.Camera)
	self.Controller.Pos = v.Vector3f{0, 0, -30}

	self.gameState = state.NewGameState()
	
	self.gameState.ObjectManager.SetRenderer(state.NewRenderContext(self.Camera))

	var err error
	if self.shader, err = glutils.GetProgram(
		"shaders/blob.vertex", 
		"shaders/blob.fragment"); err != nil {
		panic(err.Error())
	}
	self.bindFunc = func(prog *glutils.ShaderProgram) {
		prog.GetUniform("light0_direction").Uniform3f(0, -1, 0)
	}

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

	self.gameState.Physics.Process(time_step)

	self.gameState.UpdatePlayerCtrl(
		self.moveDir.Mul(0.2), v.Angle(self.Controller.HorAxis))


	self.gameState.ObjectManager.Process(time_step)


	self.Controller.Pos = self.gameState.Player.Position

	self.Controller.SetupCamera()

	self.gameState.VoxelsMesh.SetCenter(
		self.Controller.Pos)

	self.gameState.VoxelsMesh.Render(self.Camera, 
		self.shader,
		self.bindFunc)

	self.gameState.ObjectManager.ForAllObjects(func(ob *wobject.WObject) {
		if ob.RendererModule != nil {
			ob.RendererModule.Render()
		}
	})


	sdl.GL_SwapBuffers()

	if self.frameStats.FrameFinished(int64(time_step * 1000)) {
		fmt.Printf("%v\n", &self.frameStats)
	}
}

func (self *MCPlayAppState) OnKeyDown(key *sdl.Keysym) {
	switch(key.Sym) {
	case sdl.K_w:
		self.moveDir.Z = -1
		break;
	case sdl.K_s:
		self.moveDir.Z = 1
		break
	case sdl.K_a:
		self.moveDir.X = -1
		break
	case sdl.K_d:
		self.moveDir.X = 1
		break
	case sdl.K_q:
		self.moveDir.Y = 1
		break
	case sdl.K_e:
		self.moveDir.Y = -1
		break
	case sdl.K_l:
		ball := state.CreateBall(v.Vector3f{0,10,0})
		self.gameState.ObjectManager.RegisterObject(ball,  v.Vector3f{0, 10, 0}, nil)
		break;
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

	voxels.DrawSphere(self.gameState.VoxelField, mpos.X, mpos.Y, mpos.Z, 2, 100)
	
	self.gameState.VoxelsMesh.RefreshMesh()
}

func (self *MCPlayAppState) OnSdlEvent(event *sdl.Event) {

}

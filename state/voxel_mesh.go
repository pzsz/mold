package state

import (
	"fmt"
	"github.com/pzsz/glutils"
	v "github.com/pzsz/lin3dmath"
	"github.com/pzsz/mold/bulletbridge"
	"github.com/pzsz/mold/voxels"
)

type VMBState int

const (
	VMB_DIRTY     = VMBState(0)
	VMB_BUILDING  = VMBState(1)
	VMB_OK        = VMBState(2)
	VMB_DESTROYED = VMBState(3)
)

type VoxelsMeshBlock struct {
	Mesh          *glutils.MeshBuffer
	Position      v.Vector3i
	Size          v.Vector3i
	State         VMBState
	CollisionMesh *bulletbridge.BBStaticMesh
}

func (s *VoxelsMeshBlock) SetMesh(Mesh *glutils.MeshBuffer, CollisionMesh *bulletbridge.BBStaticMesh) {
	if s.Mesh != nil {
		s.Destroy()
	}
	s.Mesh = Mesh
	s.CollisionMesh = CollisionMesh
	if s.Mesh != nil {
		// Create VBO
		s.Mesh.CopyArraysToVBO()
		// Create collision
	}
}

func (s *VoxelsMeshBlock) Destroy() {
	if s.Mesh != nil {
		s.Mesh.Destroy()
		s.Mesh = nil
	}
	if s.CollisionMesh != nil {
		s.CollisionMesh.Destroy()
		s.CollisionMesh = nil
	}
}

type VMBUpdate struct {
	BlockGen int
	Mesh     *glutils.MeshBuffer
	WorldPos v.Vector3i
}

type VoxelsMeshConfig struct {
	BlockArraySize v.Vector3i
	BlockSize      v.Vector3i
}

type VoxelsMesh struct {
	voxelField    *voxels.DamageWrapper
	blockArray    []VoxelsMeshBlock
	blockStart    v.Vector3i
	blockArrayGen int

	renderSize v.Vector3i

	config VoxelsMeshConfig

	vbmUpdateChannel chan VMBUpdate
	baUpdateChannel  chan []VoxelsMeshBlock

	Physics *bulletbridge.BBWorld
}

func NewVoxelsMesh(voxelField *voxels.DamageWrapper,
	config VoxelsMeshConfig,
	Physics *bulletbridge.BBWorld) *VoxelsMesh {
	ret := &VoxelsMesh{
		voxelField:       voxelField,
		config:           config,
		renderSize:       config.BlockArraySize.Mul3I(config.BlockSize),
		vbmUpdateChannel: make(chan VMBUpdate, 100),
		baUpdateChannel:  make(chan []VoxelsMeshBlock, 5),
	        Physics:          Physics,
	}

	ret.voxelField.DamageFunc = func(b v.Boxi) {
		ret.UpdateArea(b)
	}
	ret.newRangesArray()

	return ret
}

func (s *VoxelsMesh) UpdatePoint(point v.Vector3i) {
	b := point.Div3I(s.config.BlockSize)
	b.Sub(s.blockStart)

	baSize := s.config.BlockArraySize
	if b.X < 0 || b.X >= baSize.X ||
		b.Y < 0 || b.Y >= baSize.Y ||
		b.Z < 0 || b.Z >= baSize.Z {
		return
	}
	id := b.X + b.Y*baSize.X + b.Z*baSize.X*baSize.Y
	s.blockArray[id].State = VMB_DIRTY
}

func (s *VoxelsMesh) UpdateArea(damage v.Boxi) {
	//	block_start := damage.Start.Div3I(s.config.BlockSize).Sub(s.blockStart)
	//	block_end := damage.End.Div3I(s.config.BlockSize).Sub(s.blockStart)

	damage = damage.GrowBy(1)
	block_start := damage.Start.Sub(s.blockStart.Mul3I(s.config.BlockSize)).Div3I(s.config.BlockSize)
	block_end := damage.End.Sub(s.blockStart.Mul3I(s.config.BlockSize)).Div3I(s.config.BlockSize)

	updamage := v.Boxi{block_start, block_end}

	baSize := s.config.BlockArraySize

	for z := updamage.Start.Z; z <= updamage.End.Z; z++ {
		for y := updamage.Start.Y; y <= updamage.End.Y; y++ {
			for x := updamage.Start.X; x <= updamage.End.X; x++ {
				if x < 0 || x >= baSize.X ||
					y < 0 || y >= baSize.Y ||
					z < 0 || z >= baSize.Z {
					continue
				}
				id := x + y*baSize.X + z*baSize.X*baSize.Y
				s.blockArray[id].State = VMB_DIRTY
			}
		}
	}
}

func (s *VoxelsMesh) GetBlock(x, y, z int) *VoxelsMeshBlock {
	b := s.config.BlockArraySize
	return &s.blockArray[x+y*b.X+z*b.X*b.Y]
}

func (s *VoxelsMesh) SetCenter(pos v.Vector3f) {
	units_half_size := s.config.BlockArraySize.Mul3I(s.config.BlockSize).DivI(2)
	newBlockStart := pos.Sub(units_half_size.To3F()).To3I().Div3I(s.config.BlockSize)

	if newBlockStart != s.blockStart {
		dif := newBlockStart.Sub(s.blockStart)
		fmt.Printf("New translate start %v -> %v, dif %v\n", s.blockStart, newBlockStart, dif)
		s.translateBlockArray(dif.X, dif.Y, dif.Z)
		s.RefreshMesh()
	}
}

func (s *VoxelsMesh) newRangesArray() {
	size := s.config.BlockArraySize
	array_size := size.X * size.Y * size.Z

	s.blockArray = make([]VoxelsMeshBlock, array_size)
	for id, _ := range s.blockArray {
		z := id / (size.X * size.Y)
		xy := (id - size.X*size.Y*z)
		y := xy / size.X
		x := xy % size.X

		pos := s.blockStart.Mul3I(s.config.BlockSize)
		pos.AddIP(v.Vector3i{x, y, z}.Mul3I(s.config.BlockSize))

		s.initVMB(&s.blockArray[id], pos)
	}
}

func (s *VoxelsMesh) initVMB(vbm *VoxelsMeshBlock, position v.Vector3i) {
	vbm.Position = position
	vbm.Size = s.config.BlockSize
	vbm.State = VMB_DIRTY
}

func (s *VoxelsMesh) translateBlockArray(tx, ty, tz int) {
	size := s.config.BlockArraySize
	array_size := size.X * size.Y * size.Z

	s.blockStart.X += tx
	s.blockStart.Y += ty
	s.blockStart.Z += tz

	newBlockArray := make([]VoxelsMeshBlock, array_size)
	for z := 0; z < size.Z; z++ {
		for y := 0; y < size.Y; y++ {
			for x := 0; x < size.X; x++ {
				id := x + y*size.X + z*size.X*size.Y
				ox, oy, oz := x+tx, y+ty, z+tz
				out_of_band := (ox < 0 || ox >= size.X ||
					oy < 0 || oy >= size.Y ||
					oz < 0 || oz >= size.Z)

				if !out_of_band {
					oid := ox + oy*size.X + oz*size.X*size.Y
					newBlockArray[id] = s.blockArray[oid]
				} else {
					pos := s.blockStart.Mul3I(s.config.BlockSize)
					pos.AddIP(v.Vector3i{x, y, z}.Mul3I(s.config.BlockSize))
					s.initVMB(&newBlockArray[id], pos)
				}
			}
		}
	}

	if tx != 0 {
		rmx_start := 0
		rmx_end := 0
		if tx > 0 {
			rmx_start = 0
			rmx_end = tx
		} else if tx < 0 {
			rmx_start = size.X + tx
			rmx_end = size.X
		}
		destroyMeshes(rmx_start, rmx_end,
			0, size.Y,
			0, size.Z,
			s.blockArray, size)
	}

	if ty != 0 {
		rmy_start := 0
		rmy_end := 0
		if ty > 0 {
			rmy_start = 0
			rmy_end = ty
		} else if tx < 0 {
			rmy_start = size.Y + ty
			rmy_end = size.Y
		}
		destroyMeshes(0, size.X,
			rmy_start, rmy_end,
			0, size.Z,
			s.blockArray, size)
	}

	if tz != 0 {
		rmz_start := 0
		rmz_end := 0
		if tz > 0 {
			rmz_start = 0
			rmz_end = tz
		} else if tx < 0 {
			rmz_start = size.Z + tz
			rmz_end = size.Z
		}
		destroyMeshes(0, size.X,
			0, size.Y,
			rmz_start, rmz_end,
			s.blockArray, size)
	}

	s.blockArray = newBlockArray
	s.blockArrayGen += 1
}

func destroyMeshes(sx, ex, sy, ey, sz, ez int, array []VoxelsMeshBlock, size v.Vector3i) {
	fmt.Printf("Destroing X %v-%v Y %v-%v Z %v-%v\n", sx, ex, sy, ey, sz, ez)
	for x := sx; x < ex; x++ {
		for z := sz; z < ez; z++ {
			for y := sy; y < ey; y++ {
				id := x + y*size.X + z*size.X*size.Y
				blck := &array[id]
				if (blck.State == VMB_DIRTY ||
					blck.State == VMB_OK) &&
					array[id].Mesh != nil {
					array[id].Destroy()
					array[id].State = VMB_DESTROYED
				}
			}
		}
	}
}

// Regenerate meshes of block
func (s *VoxelsMesh) RefreshMesh() {
	blockArrayGeneration := s.blockArrayGen

	for id, _ := range s.blockArray {
		block := &s.blockArray[id]

		switch block.State {
		case VMB_DIRTY:
			go func() {
				mesh_builder := glutils.NewMeshBuilder()

				voxels.BuildMeshRange(
					s.voxelField,
					block.Position,
					block.Position.Add(block.Size),
					mesh_builder)

				if !mesh_builder.IsEmpty() {
					mesh := glutils.NewMeshBuffer(
						0, 0, glutils.RENDER_POLYGONS,
						glutils.BUF_NORMAL)

					mesh_builder.Finalize(false, mesh)
					s.vbmUpdateChannel <- VMBUpdate{
						blockArrayGeneration,
						mesh,
						block.Position}
				} else {
					s.vbmUpdateChannel <- VMBUpdate{
						blockArrayGeneration,
						nil,
						block.Position}
				}
			}()
			block.State = VMB_BUILDING
			break
		}
	}
}

func (s *VoxelsMesh) Process() {
	baSize := s.config.BlockArraySize
	for {
		select {
		case update := <-s.vbmUpdateChannel:
			b := update.WorldPos.Div3I(s.config.BlockSize).Sub(s.blockStart)

			if b.X < 0 || b.X >= baSize.X ||
				b.Y < 0 || b.Y >= baSize.Y ||
				b.Z < 0 || b.Z >= baSize.Z {
				continue
			}

			id := b.X + b.Y*baSize.X + b.Z*baSize.X*baSize.Y
			block := &s.blockArray[id]

			if update.Mesh != nil {
				varray, iarray := update.Mesh.GetArrays()
				collision_mesh := s.Physics.NewStaticMesh(
					update.Mesh.CalcVertexSize(),
					varray, iarray)
				
				block.SetMesh(update.Mesh, collision_mesh)
			} else {
				block.SetMesh(nil, nil)
			}
			block.State = VMB_OK
		default:
			return
		}
	}
}

func (s *VoxelsMesh) Render(camera *glutils.Camera,
	program *glutils.ShaderProgram, bindfunc func(prog *glutils.ShaderProgram)) {

	s.Process()

	for id, _ := range s.blockArray {
		block := &s.blockArray[id]
		if block.Mesh != nil {
			rop := glutils.NewShaderRenderOp(
				false, program, bindfunc, block.Mesh)
			rop.Render(camera, v.MatrixOne())
		}
	}
}

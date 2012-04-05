package rendering

import (
	"github.com/pzsz/glutils"
	v "github.com/pzsz/lin3dmath"
	"github.com/pzsz/marchingcubes/voxels"
	"fmt"
	)

type VoxelsBlockMesh struct {
	Mesh      *glutils.MeshBuffer
	Position  v.Vector3i
	Size      v.Vector3i
	Dirty     bool
}

func (s *VoxelsBlockMesh) UpdateMesh(vox_storage voxels.VoxelField) {
	var builder *glutils.MeshBuilder
	if s.Mesh != nil {
		builder = glutils.ReuseMeshBuilder(s.Mesh)
	} else {
		builder = glutils.NewMeshBuilder(0, 0, glutils.RENDER_POLYGONS, 
			glutils.BUF_NORMAL | glutils.BUF_COLOUR, true)
	}
	
	voxels.BuildMeshRange(vox_storage, s.Position, s.Position.Add(s.Size),
		builder)	
}

type VoxelsRendererConfig struct {
	BlockArraySize  v.Vector3i
	BlockSize       v.Vector3i
}

type VoxelsRenderer struct {
	storage       voxels.VoxelField
	blockArray    []VoxelsBlockMesh
	blockStart    v.Vector3i 


	config        VoxelsRendererConfig
}

func NewVoxelsRenderer(storage voxels.VoxelField, config VoxelsRendererConfig) *VoxelsRenderer {
	ret := &VoxelsRenderer{
	storage: storage,
        config: config,
	}
	ret.newRangesArray()
	
	return ret
}

func (s *VoxelsRenderer) UpdatePoint(point v.Vector3i) {
	b := point.Div3I(s.config.BlockSize)
	b.Sub(s.blockStart)
	
	baSize := s.config.BlockArraySize
	if b.X < 0 || b.X >= baSize.X ||
		b.Y < 0 || b.Y >= baSize.Y ||
		b.Z < 0 || b.Z >= baSize.Z {
		return
	}
	id :=b.X + b.Y*baSize.X + b.Z*baSize.X*baSize.Y	
	s.blockArray[id].Dirty = true
}

func (s *VoxelsRenderer) UpdateLargeArea(damage v.Box3i) {
	block_start := damage.Start.Div3I(s.config.BlockSize)
	block_end := damage.End.Div3I(s.config.BlockSize)

	block_damage := v.Box3i{block_start, block_end}
	render_cube := v.Box3i{s.blockStart, s.blockStart.Add(s.config.BlockArraySize)}

	baSize := s.config.BlockArraySize
	found, updamage := block_damage.Intersection(render_cube)
	if !found {
		return
	}
	for z := updamage.Start.Z; z<updamage.End.Z; z++ {
	for y := updamage.Start.Y; z<updamage.End.Y; y++ {
	for x := updamage.Start.X; z<updamage.End.X; x++ {
		id :=x + y*baSize.X + z*baSize.X*baSize.Y
		s.blockArray[id].Dirty = true
	}
	}
	}
}

func (s *VoxelsRenderer) GetBlock(x, y, z int) *VoxelsBlockMesh {
	b := s.config.BlockArraySize
	return &s.blockArray[x + y*b.X + z*b.X*b.Y]
}

func (s *VoxelsRenderer) SetCenter(pos v.Vector3f) {
	units_half_size := s.config.BlockArraySize.Mul3I(s.config.BlockSize).DivI(2)
	newBlockStart := pos.Sub(units_half_size.To3F()).To3I().Div3I(s.config.BlockSize)

	if newBlockStart != s.blockStart {
		dif := newBlockStart.Sub(s.blockStart)
		s.translateBlockArray(dif.X, dif.Y, dif.Z)
		s.blockStart = newBlockStart
		s.RefreshMesh()

		fmt.Printf("New translate start %v, dif %v\n", s.blockStart, dif)
	}
}

func (s *VoxelsRenderer) newRangesArray() {
	size := s.config.BlockArraySize
	array_size := size.X*size.Y*size.Z

	s.blockArray = make([]VoxelsBlockMesh, array_size)
	for id, _ := range s.blockArray {
		block := &s.blockArray[id]

		block.Size = s.config.BlockSize
		block.Position = s.blockStart.Mul3I(s.config.BlockSize)
		
		z := id / (size.X*size.Y)
		xy := (id-size.X*size.Y*z)
		y := xy / size.X
		x := xy % size.X

		block.Position.AddIP(v.Vector3i{x, y, z}.Mul3I(s.config.BlockSize))
		block.Dirty = true
	}
}

func (s *VoxelsRenderer) translateBlockArray(tx, ty, tz int) {
	size := s.config.BlockArraySize
	array_size := size.X*size.Y*size.Z

	s.blockStart.X += tx
	s.blockStart.Y += ty
	s.blockStart.Z += tz

	newBlockArray := make([]VoxelsBlockMesh, array_size)
	for z:=0; z<size.Z; z++ {
		for y:=0; y<size.Y; y++ {
			for x:=0; x<size.X; x++ {
				id := x + y*size.X + z*size.X*size.Y
				ox,oy,oz := x+tx,y+ty,z+tz
				out_of_band := (ox < 0 || ox >= size.X ||
					oy < 0 || oy >= size.Y ||
					oz < 0 || oz >= size.Z)

				if !out_of_band {
					oid := ox + oy*size.X + oz*size.X*size.Y
					newBlockArray[id] = s.blockArray[oid]
				} else {
					block := &newBlockArray[id]
					block.Size = s.config.BlockSize
					block.Position = s.blockStart.Mul3I(s.config.BlockSize)
					block.Position.AddIP(v.Vector3i{x, y, z}.Mul3I(s.config.BlockSize))
					block.Dirty = true
				}
			}
		}	
	}	

	s.blockArray = newBlockArray
}


func (s *VoxelsRenderer) RefreshMesh() {
	for id, _ := range s.blockArray {
		block := &s.blockArray[id]
		if !block.Dirty {
			continue
		}

		if block.Mesh == nil {
			block.Mesh = glutils.NewMeshBuffer(
				0, 0, 
				glutils.RENDER_POLYGONS,
				glutils.BUF_NORMAL)
			block.Mesh.AllocBuffers()
		}

		mesh_builder := glutils.ReuseMeshBuilder(block.Mesh)

		voxels.BuildMeshRange(
			s.storage, 
			block.Position,
			block.Position.Add(block.Size),
			mesh_builder)

		if !mesh_builder.IsEmpty() {
			block.Mesh = mesh_builder.Finalize()
		} else {
			block.Mesh.Destroy()
			block.Mesh = nil
		}
	}
}

func (s *VoxelsRenderer) Render(camera *glutils.Camera,
	program *glutils.ShaderProgram, bindfunc func(prog *glutils.ShaderProgram)) {

	for id, _ := range s.blockArray {
		block := &s.blockArray[id]
		if block.Mesh != nil {
			rop := glutils.NewShaderRenderOp(
				false, program, bindfunc, block.Mesh)
			rop.Render(camera, v.MatrixOne())
		}
	}
}

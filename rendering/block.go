package rendering

import (
	"github.com/pzsz/glutils"
	v "github.com/pzsz/lin3dmath"
	"github.com/pzsz/marchingcubes/voxels"
	)

type VoxelsBlockMesh struct {
	Mesh      *glutils.MeshBuffer
	Position  v.Vector3i
	Size      v.Vector3i
	Dirty     bool
}

func (s *VoxelsBlockMesh) UpdateMesh(vox_storage voxels.VoxelsStorage) {
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
	storage       voxels.VoxelsStorage
	blockArray    []VoxelsBlockMesh
	blockStart    v.Vector3i 

	config        VoxelsRendererConfig
}

func NewVoxelsRenderer(storage voxels.VoxelsStorage, config VoxelsRendererConfig) *VoxelsRenderer {
	return &VoxelsRenderer{
	storage: storage,
        config: config,
	}
}

func (s *VoxelsRenderer) GetBlock(x, y, z int) *VoxelsBlockMesh {
	b := s.config.BlockArraySize
	return &s.blockArray[x + y*b.X + z*b.X*b.Y]
}

func (s *VoxelsRenderer) CreateRangesArray() {
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
		block.Position.AddIP(v.Vector3i{x, y, z})
	}
}

func (s *VoxelsRenderer) RefreshMesh() {
	for id, _ := range s.blockArray {
		block := &s.blockArray[id]

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
			block.Size,
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

		rop := glutils.NewShaderRenderOp(
			false, program, bindfunc, block.Mesh)
		rop.Render(camera, v.MatrixTranslate(
			float32(block.Position.X),
			float32(block.Position.Y),
			float32(block.Position.Z)))
	}
}

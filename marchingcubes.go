package marchingcubes

import (
	"math"
	v "github.com/pzsz/lin3dmath"
)

type VoxelCube struct {
	Voxels Voxels
	Pos    v.Vector3i
}

func (self *VoxelCube) GetPosition(id int) (v.Vector3i) {
	switch id {
	case 0:
		return self.Pos.Add(v.Vector3i{0,0,0});
	case 1:
		return self.Pos.Add(v.Vector3i{1,0,0});
	case 2:
		return self.Pos.Add(v.Vector3i{1,0,1});
	case 3:
		return self.Pos.Add(v.Vector3i{0,0,1});
	case 4:
		return self.Pos.Add(v.Vector3i{0,1,0});
	case 5:
		return self.Pos.Add(v.Vector3i{1,1,0});
	case 6:
		return self.Pos.Add(v.Vector3i{1,1,1});
	case 7:
		return self.Pos.Add(v.Vector3i{0,1,1});
	}
	return v.Vector3i{-1,-1,-1};
}

func (self *VoxelCube) GetValue(id int) (int) {
	var x,y,z = self.Pos.X,self.Pos.Y,self.Pos.Z;

	switch id {
	case 0:
		return self.Voxels.GetValue(x,y,z);			
	case 1:
		return self.Voxels.GetValue(x+1,y,z);
	case 2:
		return self.Voxels.GetValue(x+1,y,z+1);
	case 3:
		return self.Voxels.GetValue(x,y,z+1);
	case 4:
		return self.Voxels.GetValue(x,y+1,z);
	case 5:
		return self.Voxels.GetValue(x+1,y+1,z);
	case 6:
		return self.Voxels.GetValue(x+1,y+1,z+1);
	case 7:
		return self.Voxels.GetValue(x,y+1,z+1);
	}
	return 255;
}

func SampleVoxels(vstore Voxels, mesh *Mesh) {
	sizeX,sizeY,sizeZ := vstore.Size()
	for z:=0;z < sizeZ-1;z++ {
		for y:=0;y < sizeY-1;y++ {
			for x:=0;x < sizeX-1;x++ {
				cube := VoxelCube{vstore, v.Vector3i{x,y,z}}
				SampleVoxelCube(&cube, mesh)
			}
		}
	}
}

func SampleVoxelCube(cube *VoxelCube, mesh *Mesh) {
	var cubeindex uint;
	cubeindex = 0;

	var vertlist [12]v.Vector3f;

	if cube.GetValue(0) < ISOLEVEL {cubeindex |= 1;}
	if cube.GetValue(1) < ISOLEVEL {cubeindex |= 2;}
	if cube.GetValue(2) < ISOLEVEL {cubeindex |= 4;}
	if cube.GetValue(3) < ISOLEVEL {cubeindex |= 8;}
	if cube.GetValue(4) < ISOLEVEL {cubeindex |= 16;}
	if cube.GetValue(5) < ISOLEVEL {cubeindex |= 32;}
	if cube.GetValue(6) < ISOLEVEL {cubeindex |= 64;}
	if cube.GetValue(7) < ISOLEVEL {cubeindex |= 128;}

	/* Cube is entirely in/out of the surface */
	if edgeTable[cubeindex] == 0 {
		return;
	}
	
	/* Find the vertices where the surface intersects the cube */
	if edgeTable[cubeindex] & 1 != 0{
		vertlist[0] = VertexInterp(cube,0,1);
	}
	if edgeTable[cubeindex] & 2 != 0 {
		vertlist[1] = VertexInterp(cube,1,2);
	}
	if edgeTable[cubeindex] & 4 != 0 {
		vertlist[2] = VertexInterp(cube,2,3);
	}
	if edgeTable[cubeindex] & 8 != 0 {
		vertlist[3] = VertexInterp(cube,3,0);
	}
	if edgeTable[cubeindex] & 16 != 0 {
		vertlist[4] = VertexInterp(cube,4,5);
	}
	if edgeTable[cubeindex] & 32 != 0 {
		vertlist[5] = VertexInterp(cube,5,6);
	}
	if edgeTable[cubeindex] & 64 != 0 {
		vertlist[6] = VertexInterp(cube,6,7);
	}
	if edgeTable[cubeindex] & 128 != 0 {
		vertlist[7] = VertexInterp(cube,7,4);
	}
	if edgeTable[cubeindex] & 256 != 0 {
		vertlist[8] = VertexInterp(cube,0,4);
	}
	if edgeTable[cubeindex] & 512 != 0 {
		vertlist[9] = VertexInterp(cube,1,5);
	}
	if edgeTable[cubeindex] & 1024 != 0 {
		vertlist[10] = VertexInterp(cube,2,6);
	}
	if edgeTable[cubeindex] & 2048 != 0 {
		vertlist[11] = VertexInterp(cube,3,7);
	}

	/* Create the triangle */
	for i:=0;triTable[cubeindex][i]!=-1;i+=3 {
		vertIds := triTable[cubeindex];
		
		v1 := mesh.AddVertex( &Vertex{pos: vertlist[vertIds[i]], normal: v.Vector3f{0,0,0}});
		v2 := mesh.AddVertex( &Vertex{pos: vertlist[vertIds[i+1]], normal: v.Vector3f{0,0,0}});
		v3 := mesh.AddVertex( &Vertex{pos: vertlist[vertIds[i+2]], normal: v.Vector3f{0,0,0}});
		mesh.AddPolygon(CreatePolygon(v1,v2,v3));
	}
}

/*
   Linearly interpolate the position where an isosurface cuts
   an edge between two vertices, each with their own scalar value
*/

func VertexInterp(cube *VoxelCube,a,b int) (v.Vector3f) {
	var mu float32 = 0;
	var p v.Vector3f;
	var valp1 int = int(cube.GetValue(a));
	var valp2 int = int(cube.GetValue(b));
	p1 := cube.GetPosition(a);
	p2 := cube.GetPosition(b);

	if math.Abs(float64(ISOLEVEL-valp1)) < 0.00001 { 
		return p1.ToF();
	}
	if math.Abs(float64(ISOLEVEL-valp2)) < 0.00001 {
		return p2.ToF()
	}
	if math.Abs(float64(valp1-valp2)) < 0.00001 {
		return p1.ToF();
	}
	mu = float32(ISOLEVEL - valp1) / float32(valp2 - valp1);
	p.X = float32(p1.X) + float32(mu * float32(p2.X - p1.X));
	p.Y = float32(p1.Y) + float32(mu * float32(p2.Y - p1.Y));
	p.Z = float32(p1.Z) + float32(mu * float32(p2.Z - p1.Z));

	return p;
}
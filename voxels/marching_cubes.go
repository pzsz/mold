package voxels

import (
	"math"
	v "github.com/pzsz/lin3dmath"
)

// Make mesh from 'cube' range inside voxel storage
func BuildMeshRange(vstore VoxelsStorage, start, end v.Vector3i, mesh MeshWriter) {
//	sizeX,sizeY,sizeZ := vstore.Size()

	for z:=start.Z; z < end.Z-1;z++ {
		for y:=start.Y; y < end.Y-1;y++ {
			for x:=start.X;x < end.X-1;x++ {
				BuildMeshCube(&VoxelCube{vstore, v.Vector3i{x,y,z}}, mesh)
			}
		}
	}
}

// Make mesh from one 'cube' inside voxel storage
func BuildMeshCube(cube *VoxelCube, mesh MeshWriter) {
	var cubeindex uint;
	cubeindex = 0;

	var pos_list [12]v.Vector3f;
	var normal_list [12]v.Vector3f;

	if cube.GetValue(0) < ISOLEVEL {cubeindex |= 1;}
	if cube.GetValue(1) < ISOLEVEL {cubeindex |= 2;}
	if cube.GetValue(2) < ISOLEVEL {cubeindex |= 4;}
	if cube.GetValue(3) < ISOLEVEL {cubeindex |= 8;}
	if cube.GetValue(4) < ISOLEVEL {cubeindex |= 16;}
	if cube.GetValue(5) < ISOLEVEL {cubeindex |= 32;}
	if cube.GetValue(6) < ISOLEVEL {cubeindex |= 64;}
	if cube.GetValue(7) < ISOLEVEL {cubeindex |= 128;}

	/* Cube is entirely in/out of the surface */
	if EDGE_TABLE[cubeindex] == 0 {
		return;
	}
	
	/* Find the vertices where the surface intersects the cube */
	if EDGE_TABLE[cubeindex] & 1 != 0{
		pos_list[0] = VertexPositionInterp(cube,0,1);
		normal_list[0] = VertexNormalGradientInterp(cube, 0, 1)
	}
	if EDGE_TABLE[cubeindex] & 2 != 0 {
		pos_list[1] = VertexPositionInterp(cube,1,2);
		normal_list[1] = VertexNormalGradientInterp(cube, 1, 2)
	}
	if EDGE_TABLE[cubeindex] & 4 != 0 {
		pos_list[2] = VertexPositionInterp(cube,2,3);
		normal_list[2] = VertexNormalGradientInterp(cube, 2, 3)
	}
	if EDGE_TABLE[cubeindex] & 8 != 0 {
		pos_list[3] = VertexPositionInterp(cube,3,0);
		normal_list[3] = VertexNormalGradientInterp(cube, 3, 0)
	}
	if EDGE_TABLE[cubeindex] & 16 != 0 {
		pos_list[4] = VertexPositionInterp(cube,4,5);
		normal_list[4] = VertexNormalGradientInterp(cube, 4, 5)
	}
	if EDGE_TABLE[cubeindex] & 32 != 0 {
		pos_list[5] = VertexPositionInterp(cube,5,6);
		normal_list[5] = VertexNormalGradientInterp(cube, 5, 6)
	}
	if EDGE_TABLE[cubeindex] & 64 != 0 {
		pos_list[6] = VertexPositionInterp(cube,6,7);
		normal_list[6] = VertexNormalGradientInterp(cube, 6, 7)
	}
	if EDGE_TABLE[cubeindex] & 128 != 0 {
		pos_list[7] = VertexPositionInterp(cube,7,4);
		normal_list[7] = VertexNormalGradientInterp(cube, 7, 4)
	}
	if EDGE_TABLE[cubeindex] & 256 != 0 {
		pos_list[8] = VertexPositionInterp(cube,0,4);
		normal_list[8] = VertexNormalGradientInterp(cube, 0, 4)
	}
	if EDGE_TABLE[cubeindex] & 512 != 0 {
		pos_list[9] = VertexPositionInterp(cube,1,5);
		normal_list[9] = VertexNormalGradientInterp(cube, 1, 5)
	}
	if EDGE_TABLE[cubeindex] & 1024 != 0 {
		pos_list[10] = VertexPositionInterp(cube,2,6);
		normal_list[10] = VertexNormalGradientInterp(cube, 2, 6)
	}
	if EDGE_TABLE[cubeindex] & 2048 != 0 {
		pos_list[11] = VertexPositionInterp(cube,3,7);
		normal_list[11] = VertexNormalGradientInterp(cube, 3, 7)
	}

	/* Create the triangle */
	for i:=0;TRI_TABLE[cubeindex][i]!=-1;i+=3 {
		vertIds := TRI_TABLE[cubeindex];

		pos0 := pos_list[vertIds[i]]
		pos1 := pos_list[vertIds[i+1]]
		pos2 := pos_list[vertIds[i+2]]

		normal0 := normal_list[vertIds[i]]
		normal1 := normal_list[vertIds[i+1]]
		normal2 := normal_list[vertIds[i+2]]

		v1 := mesh.StartVertex()
		mesh.AddPosition(pos0.X, pos0.Y, pos0.Z)
		mesh.AddNormal(normal0.X, normal0.Y, normal0.Z)

		v2 := mesh.StartVertex()
		mesh.AddPosition(pos1.X, pos1.Y, pos1.Z)
		mesh.AddNormal(normal1.X, normal1.Y, normal1.Z)

		v3 := mesh.StartVertex()
		mesh.AddPosition(pos2.X, pos2.Y, pos2.Z)
		mesh.AddNormal(normal2.X, normal2.Y, normal2.Z)

		mesh.AddIndice3(v3,v2,v1)
	}
}

//   Linearly interpolate the position where an isosurface cuts
//   an edge between two vertices, each with their own scalar value
func VertexPositionInterp(cube *VoxelCube,a,b int) (p v.Vector3f) {
	var mu float32 = 0
	var valp1 int = cube.GetValue(a)
	var valp2 int = cube.GetValue(b)
	p1 := cube.GetPosition(a)
	p2 := cube.GetPosition(b)

	if ISOLEVEL-valp1 == 0 { 
		p = p1.ToF()
		return 
	}
	if ISOLEVEL-valp2 == 0 {
		p = p2.ToF()
		return
	}
	if math.Abs(float64(valp1-valp2)) < 2 {
		p = p1.ToF()
		return
	}
	mu = float32(ISOLEVEL - valp1) / float32(valp2 - valp1)
	p.X = float32(p1.X) + float32(mu * float32(p2.X - p1.X))
	p.Y = float32(p1.Y) + float32(mu * float32(p2.Y - p1.Y))
	p.Z = float32(p1.Z) + float32(mu * float32(p2.Z - p1.Z))
	return 
}

func VertexNormalGradientInterp(cube *VoxelCube,a,b int) (v.Vector3f) {
	var mu float32 = 0;
	var valp1 int = int(cube.GetValue(a));
	var valp2 int = int(cube.GetValue(b));
	p1 := cube.GetGradient(a)
	p2 := cube.GetGradient(b)

	if ISOLEVEL-valp1 == 0 { 
		return p1
	}
	if ISOLEVEL-valp2 == 0 {
		return p2
	}
	if math.Abs(float64(valp1-valp2)) < 2 {
		return p1
	}
	mu = float32(ISOLEVEL - valp1) / float32(valp2 - valp1)
	return p1.Add(p2.Sub(p1).Mul(mu)).Normalize()
}

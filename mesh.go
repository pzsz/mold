package marchingcubes

import (
	v "github.com/pzsz/lin3dmath"
)

type Vertex struct {
	pos    v.Vector3f;
	normal v.Vector3f;
}

type Polygon [3]int

type Mesh struct {
	Verts  []Vertex;
	Polys  []Polygon;
}

func CreateMesh() (*Mesh) {
	return new(Mesh);
}

func CreatePolygon(a,b,c int) (*Polygon) {
	return &Polygon{a,b,c};
}

func (m *Mesh) AddVertex(v *Vertex) (int) {
	m.Verts = append(m.Verts, *v);
	return len(m.Verts)-1;
}

func (m *Mesh) AddPolygon(v *Polygon) {
	m.Polys = append(m.Polys, *v);
}


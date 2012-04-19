package bulletbrideg

// #cgo CFLAGS: -I../bulletbridge_cpp
// #cgo LDFLAGS: -lstdc++ -L../bulletbridge_cpp -lbb -L/usr/local/lib -lBulletSoftBody -lBulletDynamics -lBulletCollision -lLinearMath
// #include <bbridge.h>
import "C"
import "unsafe"

type BBWorld struct {
	cptr *C.BB_World
}

type BBStaticMesh struct {
	cptr *C.BB_StaticMesh
	vertexArray   []byte
	indiceArray   []byte
}

func NewBBWorld() *BBWorld {
	ptr := C.BB_NewWorld()
	return &BBWorld{ptr}
}

func (s *BBWorld) Destroy() {
	C.BB_DestroyWorld(s.cptr)
}

func (s *BBWorld) NewStaticMesh(vertexSize int, vertexArray, indiceArray []byte) *BBStaticMesh {
	
	va_ptr := unsafe.Pointer(&vertexArray[0])
	ia_ptr := unsafe.Pointer(&indiceArray[0])

	ptr := C.BB_NewStaticMesh(s.cptr, C.int(vertexSize), 
		(*C.char)(va_ptr), C.int(len(vertexArray)), 
		(*C.char)(ia_ptr), C.int(len(indiceArray)))
	
	return &BBStaticMesh{
		ptr,
		vertexArray,
		indiceArray}
}

func (s *BBStaticMesh) Destroy() {
	C.BB_DestroyStaticMesh(s.cptr)
}
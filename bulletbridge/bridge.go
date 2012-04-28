package bulletbridge

// #cgo CFLAGS: -I../bulletbridge_cpp
// #cgo LDFLAGS: -lstdc++ -L../bulletbridge_cpp -lbb -L/usr/local/lib -lBulletSoftBody -lBulletDynamics -lBulletCollision -lLinearMath
// #include <bbridge.h>
import "C"
import "unsafe"
import v "github.com/pzsz/lin3dmath"

type BBWorld struct {
	cptr *C.BB_World
}

type BBStaticMesh struct {
	cptr *C.BB_StaticMesh
	vertexArray   []byte
	indiceArray   []byte
}

type BBRigidBody struct {
	cptr *C.BB_RBody
}

type BBCharacterControler struct {
	cptr *C.BB_CharacterControler
}

type BBCShape struct {
	cobj C.BB_CShape
}

func NewBBWorld() *BBWorld {
	ptr := C.BB_NewWorld()
	return &BBWorld{ptr}
}

func (s *BBWorld) Destroy() {
	C.BB_DestroyWorld(s.cptr)
}

func (s *BBWorld) Process(time_step float32) {
	C.BB_ProcessWorld(s.cptr, C.float(time_step))
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

func toCVector(vec v.Vector3f) (cvec C.BB_Vector3) {
	cvec.x = C.float(vec.X)
	cvec.y = C.float(vec.Y)
	cvec.z = C.float(vec.Z)
	return
}

func fromCVector(cvec C.BB_Vector3) (v.Vector3f) {
	return v.Vector3f{
		float32(cvec.x),
		float32(cvec.y),
		float32(cvec.z)}
}

func NewCShapeSphere(radius float32) BBCShape {
	cobj := C.BB_NewCShapeSphere(C.float(radius))
	return BBCShape{cobj}
}

func (s *BBWorld) NewRigidBody(shape BBCShape, mass float32, pos v.Vector3f) *BBRigidBody {	
	ptr := C.BB_NewRBody(s.cptr, shape.cobj, C.float(mass), toCVector(pos))
	return &BBRigidBody{ptr}
}

func (s *BBRigidBody) GetPosition() v.Vector3f {
	return fromCVector(C.BB_GetPositionRBody(s.cptr))
}

func (s *BBWorld) NewCharacterControler(width, height float32, pos v.Vector3f) *BBCharacterControler {	
	ptr := C.BB_NewCharacterControler(s.cptr, C.float(height), C.float(width), toCVector(pos))
	return &BBCharacterControler{ptr}
}

func (s *BBCharacterControler) GetPosition() v.Vector3f {
	return fromCVector(C.BB_GetPositionCharacterControler(s.cptr))
}

func (s *BBCharacterControler) SetWalkDirection(walk v.Vector3f)  {
	C.BB_SetWalkDirection(s.cptr, toCVector(walk))
}

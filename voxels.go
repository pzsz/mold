package marchingcubes

type Voxels interface {
	Size() (int,int,int)
	GetValue(x,y,z int) (int)
	SetValue(x,y,z int, v int)
}

type ByteVoxels struct {
	sizeX, sizeY, sizeZ int;
	data []byte;
}

func CreateByteVoxels(x,y,z int) (*ByteVoxels) {
	return &ByteVoxels{
	sizeX: x,
	sizeY: y,
	sizeZ: z,
	data: make([]byte, x*y*z)}
}

func (self *ByteVoxels) Size() (int,int,int) {
	return self.sizeX, self.sizeY, self.sizeZ;
}

func (self *ByteVoxels) GetValue(x,y,z int) (int) {
	if x < 0 || y < 0 || z < 0 { return 0; }
	if x >= self.sizeX || y >= self.sizeY || z >= self.sizeZ { return 0;}

	return int(self.data[x + y * self.sizeX + z * self.sizeX*self.sizeY]);
}

func (self *ByteVoxels) SetValue(x,y,z int, v int) {
	if x < 0 || y < 0 || z < 0 { return; }
	if x >= self.sizeX || y >= self.sizeY || z >= self.sizeZ { return;}

	self.data[x + y * self.sizeX + z * self.sizeX*self.sizeY] = byte(v);
}

	

	

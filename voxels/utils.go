package voxels

import (
	"math"
	)

func DrawSphere(store VoxelsStorage, x,y,z float32, radius float32, power int) {
	startX,endX := int(x-radius),int(x+radius+1)
	startY,endY := int(y-radius),int(y+radius+1)
	startZ,endZ := int(z-radius),int(z+radius+1)

	for ix:=startX;ix < endX;ix++ {
		for iy:=startY;iy < endY;iy++ {
			for iz:=startZ;iz < endZ;iz++ {
				var distX float32 = float32(ix) - x;
				var distY float32 = float32(iy) - y;
				var distZ float32 = float32(iz) - z;
				dist := float32(math.Sqrt(float64(distX*distX+
					distY*distY + distZ*distZ)));
				if dist > radius {
					continue;
				}
				
				var v int = int(store.GetValue(ix, iy, iz));
				v += int(float32(power) * (radius - dist) / radius);
				if v > 255 {v = 255;}
				if v < 0 {v = 0;}
				
				store.SetValue(ix, iy, iz, v);
			}
		}
	}
}

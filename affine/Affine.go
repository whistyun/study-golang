package affine

import (
	"math"

	"gonum.org/v1/gonum/mat"
)

func Rotate(xy, yz, zx float64) mat.Matrix {
	s := math.Sin(xy)
	c := math.Cos(xy)
	matxy := mat.NewDense(4, 4, []float64{
		c, -s, 0, 0,
		s, c, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	})

	s = math.Sin(yz)
	c = math.Cos(yz)
	matyz := mat.NewDense(4, 4, []float64{
		1, 0, 0, 0,
		0, c, -s, 0,
		0, s, c, 0,
		0, 0, 0, 1,
	})

	s = math.Sin(zx)
	c = math.Cos(zx)
	matzx := mat.NewDense(4, 4, []float64{
		c, 0, s, 0,
		0, 1, 0, 0,
		-s, 0, c, 0,
		0, 0, 0, 1,
	})

	ans := mat.NewDense(4, 4, nil)
	ans.Product(matxy, matyz, matzx)
	return ans
}

func Translate(x, y, z float64) mat.Matrix {
	return mat.NewDense(4, 4, []float64{
		1, 0, 0, x,
		0, 1, 0, y,
		0, 0, 1, z,
		0, 0, 0, 1,
	})
}

func Scale(x, y, z float64) mat.Matrix {
	return mat.NewDense(4, 4, []float64{
		x, 0, 0, 0,
		0, y, 0, 0,
		0, 0, z, 0,
		0, 0, 0, 1,
	})
}

package examples

import (
	"fmt"
	"github.com/vitsensei/gosac/ransac"
	"math"
	"math/rand"
	"os"
)

func CreateRandomCloud(size int) ransac.Population {
	pts := make([]ransac.Data, size, size)
	firstHalfSize := int(size / 2)
	secondHalfSize := size - firstHalfSize

	for i := 0; i < firstHalfSize; i++ {
		pts[i] = ransac.Data {
			NumericFields: []float64{rand.NormFloat64() + 1, rand.NormFloat64() + 1},
		}
	}

	for i := 0; i < secondHalfSize; i++ {
		pts[i + firstHalfSize] = ransac.Data {
			NumericFields: []float64{rand.NormFloat64() - 1, rand.NormFloat64() - 1},
		}
	}

	return ransac.Population{Data: pts}
}

// The line is represented by the following formula:
// ax + by + c = 0
type Line struct {
	a float64
	b float64
	c float64
	squareRootAB float64
}

func (l *Line) ToFile(path string) {
	f, err := os.Create(path)
	must(err)

	_, _ = f.WriteString(fmt.Sprintf("%f, %f, %f\n", l.a, l.b, l.c))
}

func GetXY(data ransac.Data) (float64, float64) {
	return data.NumericFields[0], data.NumericFields[1]
}

func (l *Line) FitData(p ransac.Population, ind []int) {
	// Assuming our Line will have the equation: ax + by + c = 0,
	// we will have the following equation:
	// a = y1 - y2
	// b = x2 - x1
	// c = (y2 - y1) * x1 + (x1 - x2) * y1

	x1, y1 := GetXY(p.Data[ind[0]])
	x2, y2 := GetXY(p.Data[ind[1]])
	if len(ind) >= 2 {
		l.a = y1 - y2
		l.b = x2 - x1
		l.c = (y2 - y1) * x1 + (x1 - x2) * y1
		l.squareRootAB = math.Sqrt(l.a * l.a + l.b * l.b)
	}
}

func (l Line) SingleLoss(p ransac.Population, ind int) float64 {
	x1, y1 := GetXY(p.Data[ind])
	return math.Abs(l.a * x1 + l.b * y1 + l.c) / l.squareRootAB
}

func (l Line) Loss(p ransac.Population, ind []int) float64 {
	//rmsLoss := 0.0
	//
	//var dist float64
	//for _, i := range ind {
	//	x1, y1 := GetXY(p.Data[i])
	//	dist = math.Abs(l.a * x1 + l.b * y1 + l.c) / l.squareRootAB
	//	rmsLoss += dist * dist
	//}
	//
	//return math.Sqrt(rmsLoss / float64(len(ind)))
	return 1 / float64(len(ind))
}

func (l Line) Clone() ransac.Model {
	m := &l
	return m
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
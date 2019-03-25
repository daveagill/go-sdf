package sdf

import "math"

type point struct {
	x, y int
}

func (p point) dstSq(q point) float64 {
	return float64((p.x-q.x)*(p.x-q.x) + (p.y-q.y)*(p.y-q.y))
}

func (p point) nearest(pts []point) (*point, float64) {
	minDst2 := math.MaxFloat64
	var nearest *point

	for i := range pts {
		pt := &pts[i]
		dst2 := pt.dstSq(p)
		if dst2 < minDst2 {
			minDst2 = dst2
			nearest = pt
		}
	}

	return nearest, math.Sqrt(minDst2)
}

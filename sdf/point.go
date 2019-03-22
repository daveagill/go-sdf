package sdf

import "math"

type point struct {
	x, y int
}

func (p point) dstSq(q point) float64 {
	return float64((p.x-q.x)*(p.x-q.x) + (p.y-q.y)*(p.y-q.y))
}

func (p point) dstFromPts(pts []point) float64 {
	dst := math.MaxFloat64

	for _, pt := range pts {
		dst = math.Min(dst, p.dstSq(pt))
	}

	return math.Sqrt(dst)
}

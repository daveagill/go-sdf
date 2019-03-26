package sdf

import (
	"errors"
	"image"
	"image/color"
	"math"
)

// SDF models a rectangular & discretized Signed Distance Field
type SDF struct {
	Field  []float64
	Width  int
	Height int
}

// New returns a zeroed SDF of the given size
func New(w, h int) *SDF {
	return &SDF{
		Field:  make([]float64, w*h),
		Width:  w,
		Height: h,
	}
}

// At returns the field value at the given coordinate
func (sdf *SDF) At(x, y int) float64 {
	return sdf.Field[y*sdf.Width+x]
}

// Set writes the field value at the given coordinate
func (sdf *SDF) Set(x, y int, v float64) {
	sdf.Field[y*sdf.Width+x] = v
}

// DisplacementField is a vectorized Signed-Distance-Field where each field value is associated
// with its nearest boundary point.
type DisplacementField struct {
	*SDF
	boundaryPts []point
}

// NearestBoundaryAt returns X,Y coordinate of the nearest boundary point from the given point
func (df *DisplacementField) NearestBoundaryAt(x, y int) (int, int) {
	pt := df.boundaryPts[y*df.Width+x]
	return pt.x, pt.y
}

// Calculate a new DisplacementField from the given Stencil
func Calculate(s Stencil) *DisplacementField {
	w, h := s.Size()
	df := DisplacementField{
		New(w, h),
		make([]point, w*h),
	}

	pts := findBoundaries(s)

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			pt, dst := point{x, y}.nearest(pts)

			// use -ve sign if we are inside and +ve if outside
			if s.Within(x, y) {
				dst = -dst
			}

			df.Set(x, y, dst)
			df.boundaryPts[y*w+x] = *pt
		}
	}

	return &df
}

func findBoundaries(s Stencil) []point {
	boundaryPts := []point{}

	w, h := s.Size()
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			// boundaries are any points within the stencil that have adjacent points outside the stencil
			if s.Within(x, y) {
				lftTransparent := x == 0 || !s.Within(x-1, y)
				rgtTransparent := x == w-1 || !s.Within(x+1, y)
				topTransparent := y == 0 || !s.Within(x, y-1)
				botTransparent := y == h-1 || !s.Within(x, y+1)

				if lftTransparent || rgtTransparent || topTransparent || botTransparent {
					boundaryPts = append(boundaryPts, point{x, y})
				}
			}
		}
	}

	return boundaryPts
}

// Draw returns an 8-bit grayscale representation of a Signed-Distance-Field
func (sdf *SDF) Draw() *image.Gray {
	gray := image.NewGray(image.Rect(0, 0, sdf.Width, sdf.Height))

	for y := 0; y < sdf.Height; y++ {
		for x := 0; x < sdf.Width; x++ {
			// clamp field distance to a range [-127, 127] and then map that to [0, 255]
			dst := sdf.At(x, y)
			clamped := math.Max(-127, math.Min(127, dst))
			mapped := uint8(clamped + 127)

			col := color.Gray{mapped}
			gray.Set(x, y, col)
		}
	}

	return gray
}

// Lerp returns the linear interpolation between two SDFs, weighted by t in range [0, 1]
func Lerp(a *SDF, b *SDF, t float64) (*SDF, error) {
	if a.Width != b.Width {
		return nil, errors.New("SDF a and SDF b must have matching width")
	}
	if a.Height != b.Height {
		return nil, errors.New("SDF a and SDF b must have matching height")
	}

	ret := New(a.Width, a.Height)
	for i := range a.Field {
		ret.Field[i] = a.Field[i] + (b.Field[i]-a.Field[i])*t
	}

	return ret, nil
}

package sdf

import (
	"errors"
	"image"
	"image/color"
	"math"
)

const (
	// OpaqueAlpha is an alpha-threshold so fully-opaque pixels will form the boundary of the SDF
	OpaqueAlpha = uint16(math.MaxUint16)
	// HalfAlpha is an alpha-threshold so 50% opaque pixels will form the boundary of the SDF
	HalfAlpha = uint16(math.MaxUint16 / 2)
)

// Stencil defines a binary surface where pixels are either inside or outside the stencil
type Stencil interface {
	// Within predicates whether the given coordinate is inside or outside of the stencil surface
	Within(x, y int) bool
	// Size returns the width and height of the Stencil
	Size() (int, int)
}

// ImageAlphaStencil implements a Stencil where the alpha channel of an image is thresholded
// against the Alpha value
type ImageAlphaStencil struct {
	Image image.Image
	Alpha uint16
}

// Within predicates whether the given coordinate is inside or outside of the stencil surface
func (s ImageAlphaStencil) Within(x, y int) bool {
	b := s.Image.Bounds()
	_, _, _, a := s.Image.At(b.Min.X+x, b.Min.Y+y).RGBA()
	return a >= uint32(s.Alpha)
}

// Size returns the width and height of the ImageAlphaStencil
func (s ImageAlphaStencil) Size() (int, int) {
	size := s.Image.Bounds().Size()
	return size.X, size.Y
}

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

// Calculate a new SDF from the given Stencil
func Calculate(s Stencil) *SDF {
	return calcSDF(s, findBoundaries(s))
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

func calcSDF(s Stencil, pts []point) *SDF {
	w, h := s.Size()
	sdf := New(w, h)

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			dst := point{x, y}.dstFromPts(pts)

			// use -ve sign if we are inside and +ve if outside
			if s.Within(x, y) {
				dst = -dst
			}

			sdf.Set(x, y, dst)
		}
	}

	return sdf
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

// DrawImplicitSurface renders the implicit surface defined by a Signed-Distance-Field into an image.
// Using the given field-value to define the boundary, the given color to use for surface pixels &
// the given background color to use for background pixels.
func (sdf *SDF) DrawImplicitSurface(fv float64, c color.Color, bg color.Color) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, sdf.Width, sdf.Height))

	for y := 0; y < sdf.Height; y++ {
		for x := 0; x < sdf.Width; x++ {
			dst := sdf.At(x, y)

			if dst <= fv { // on-surface
				img.Set(x, y, c)
			} else { // off-surface (background)
				img.Set(x, y, bg)
			}
		}
	}

	return img
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

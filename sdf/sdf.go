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

// SDF models a rectangular & discretized Signed Distance Field
type SDF struct {
	Field []float64
	W     int
	H     int
}

// New returns a zeroed SDF of the given size
func New(w, h int) *SDF {
	return &SDF{
		Field: make([]float64, w*h),
		W:     w,
		H:     h,
	}
}

// Width returns the width of this SDF
func (sdf *SDF) Width() int {
	return sdf.W
}

// Height returns the height of this SDF
func (sdf *SDF) Height() int {
	return sdf.H
}

// At returns the field value at the given coordinate
func (sdf *SDF) At(x, y int) float64 {
	return sdf.Field[y*sdf.W+x]
}

// Set writes the field value at the given coordinate
func (sdf *SDF) Set(x, y int, v float64) {
	sdf.Field[y*sdf.W+x] = v
}

// FromImageAlpha returns a Signed Distance Field generated from the
// alpha channel of the image after thresholding against the given alpha-threshold.
func FromImageAlpha(img image.Image, at uint16) *SDF {
	binMap := alphaToBinaryMap(img, at)
	boundaryPts := findBoundaries(binMap)
	return calcSDF(binMap, boundaryPts)
}

func alphaToBinaryMap(img image.Image, at uint16) [][]bool {
	b := img.Bounds()
	w := b.Size().X
	h := b.Size().Y

	binMap := make([][]bool, h)
	for y := range binMap {
		binMap[y] = make([]bool, w)

		for x := range binMap[y] {
			_, _, _, alpha := img.At(b.Min.X+x, b.Min.Y+y).RGBA()
			binMap[y][x] = alpha >= uint32(at)
		}
	}

	return binMap
}

func findBoundaries(binMap [][]bool) []point {
	boundaryPts := []point{}

	for y, row := range binMap {
		for x, opaque := range row {
			// a boundary must be an opaque pixel with a transparent (or off-image) neighbour
			if opaque {
				lftTransparent := x <= 0 || !binMap[y][x-1]
				rgtTransparent := x >= len(row)-1 || !binMap[y][x+1]
				topTransparent := y <= 0 || !binMap[y-1][x]
				botTransparent := y >= len(binMap)-1 || !binMap[y+1][x]

				if lftTransparent || rgtTransparent || topTransparent || botTransparent {
					boundaryPts = append(boundaryPts, point{x, y})
				}
			}
		}
	}

	return boundaryPts
}

func calcSDF(binMap [][]bool, pts []point) *SDF {
	h := len(binMap)
	w := len(binMap[0])
	sdf := New(w, h)

	for y, row := range binMap {
		for x, opaque := range row {
			dst := point{x, y}.dstFromPts(pts)

			// use -ve sign if we are inside (opaque) and +ve if outside (transparent)
			if opaque {
				dst = -dst
			}

			sdf.Set(x, y, dst)
		}
	}

	return sdf
}

// Draw returns an 8-bit grayscale representation of a Signed-Distance-Field
func (sdf *SDF) Draw() *image.Gray {
	gray := image.NewGray(image.Rect(0, 0, sdf.Width(), sdf.Height()))

	for y := 0; y < sdf.Height(); y++ {
		for x := 0; x < sdf.Width(); x++ {
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
	img := image.NewRGBA(image.Rect(0, 0, sdf.Width(), sdf.Height()))

	for y := 0; y < sdf.Height(); y++ {
		for x := 0; x < sdf.Width(); x++ {
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
	if a.W != b.W {
		return nil, errors.New("SDF a and SDF b must have matching width")
	}
	if a.H != b.H {
		return nil, errors.New("SDF a and SDF b must have matching height")
	}

	ret := New(a.W, a.H)
	for i := range a.Field {
		ret.Field[i] = a.Field[i] + (b.Field[i]-a.Field[i])*t
	}

	return ret, nil
}

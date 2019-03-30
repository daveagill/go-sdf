package sdf

import (
	"image"
	"image/color"
	"math"
)

// Stencil defines a binary surface where pixels are either inside or outside the stencil
type Stencil interface {
	// Within predicates whether the given coordinate is inside or outside of the stencil surface
	Within(x, y int) bool
	// Size returns the width and height of the Stencil
	Size() (int, int)
}

const (
	// OpaqueAlpha is an alpha-threshold so only fully-opaque pixels will be within the stencil
	OpaqueAlpha = uint16(math.MaxUint16)
	// HalfAlpha is an alpha-threshold so 50% opaque pixels will be within the stencil
	HalfAlpha = uint16(math.MaxUint16 / 2)
)

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

// ImplicitSurfaceStencil implements a Stencil where a Signed-Distance-Field is thresholded
// against the Threshold value to define an implicit surface.
type ImplicitSurfaceStencil struct {
	SDF       *SDF
	Threshold float64
}

// Within predicates whether the given coordinate is inside or outside of the stencil surface
func (s ImplicitSurfaceStencil) Within(x, y int) bool {
	return s.SDF.At(x, y) <= s.Threshold
}

// Size returns the width and height of the ImplicitSurfaceStencil
func (s ImplicitSurfaceStencil) Size() (int, int) {
	return s.SDF.Width, s.SDF.Height
}

// DrawStencil renders a Stencil into a 2-color image.
// Using color c for pixels within the stencil and color bg for pixels outside.
func DrawStencil(s Stencil, c color.Color, bg color.Color) *image.RGBA {
	w, h := s.Size()
	img := image.NewRGBA(image.Rect(0, 0, w, h))

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if s.Within(x, y) {
				img.Set(x, y, c)
			} else {
				img.Set(x, y, bg)
			}
		}
	}

	return img
}

// DrawStencilImage stencils a given source image and returns a new image where pixels within
// the stencil are taken from the source image and pixels outside default to the given bg color.
func DrawStencilImage(s Stencil, srcImg image.Image, bg color.Color) *image.RGBA {
	w, h := s.Size()
	img := image.NewRGBA(image.Rect(0, 0, w, h))

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if s.Within(x, y) {
				img.Set(x, y, srcImg.At(x, y))
			} else {
				img.Set(x, y, bg)
			}
		}
	}

	return img
}

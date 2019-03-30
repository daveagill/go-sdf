package sdf

import (
	"image"
	"image/color"
	"testing"
)

type stubStencil struct{}

func (s stubStencil) Size() (int, int)     { return 3, 5 }
func (s stubStencil) Within(x, y int) bool { return y < 2 }

func createImageAlphaStencil() *ImageAlphaStencil {
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	img.Set(0, 0, color.Transparent)
	img.Set(0, 1, color.White)
	img.Set(1, 0, color.Black)
	img.Set(1, 1, color.Transparent)

	return &ImageAlphaStencil{
		Image: img,
		Alpha: HalfAlpha,
	}
}
func TestImageAlphaStencil_Within(t *testing.T) {
	s := createImageAlphaStencil()

	if s.Within(0, 0) {
		t.Errorf("ImageAlphaStencil at (0, 0) should be outside of the stencil surface, not withn")
	}

	if !s.Within(0, 1) {
		t.Errorf("ImageAlphaStencil at (0, 1) should be within the stencil surface, not outside")
	}

	if !s.Within(1, 0) {
		t.Errorf("ImageAlphaStencil at (1, 0) should be within the stencil surface, not outside")
	}
}

func TestImageAlphaStencil_Bounds(t *testing.T) {
	s := createImageAlphaStencil()
	w, h := s.Size()

	if w != 2 || h != 2 {
		t.Errorf("ImageAlphaStencil width and weight should be (2, 2), not (%v, %v)", w, h)
	}
}

func createImplicitSurfaceStencil() *ImplicitSurfaceStencil {
	return &ImplicitSurfaceStencil{
		SDF: &SDF{
			Width:  2,
			Height: 2,
			Field:  []float64{-1, 49, 50, 100},
		},
		Threshold: 50,
	}
}

func TestImplicitSurfaceStencil_Within(t *testing.T) {
	s := createImplicitSurfaceStencil()

	tests := []struct {
		x, y        int
		expectation bool
	}{
		{0, 0, -1 <= 50},
		{0, 1, 49 <= 50},
		{1, 0, 50 <= 50},
		{1, 1, 100 <= 50},
	}

	for _, tt := range tests {
		if s.Within(tt.x, tt.y) != tt.expectation {
			if tt.expectation {
				t.Errorf("ImplicitSurfaceStencil at (%v, %v) should be within the stencil surface, not outside", tt.x, tt.y)
			} else {
				t.Errorf("ImplicitSurfaceStencil at (%v, %v) should be outside the stencil surface, not within", tt.x, tt.y)
			}
		}
	}
}

func TestImplicitSurfaceStencil_Bounds(t *testing.T) {
	s := createImplicitSurfaceStencil()
	w, h := s.Size()

	if w != 2 || h != 2 {
		t.Errorf("ImplicitSurfaceStencil width and weight should be (2, 2), not (%v, %v)", w, h)
	}
}

func TestDrawStencil(t *testing.T) {
	stencil := stubStencil{}
	img := DrawStencil(stencil, color.Black, color.White)

	stencilW, stencilH := stencil.Size()
	size := img.Bounds().Size()
	if size.X != stencilW {
		t.Errorf("Width should equal %v, not %v", stencilW, size.X)
	}

	if size.Y != stencilH {
		t.Errorf("Height should equal %v, not %v", stencilH, size.Y)
	}

	for y := 0; y < stencilH; y++ {
		for x := 0; x < stencilW; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			isBlack := r == 0 && g == 0 && b == 0 && a == 65535
			isWhite := r == 65535 && g == 65535 && b == 65535 && a == 65535

			if stencil.Within(x, y) {
				if !isBlack {
					t.Errorf("Image pixel at (%v, %v) is inside of stencil surface so should be black, not white", x, y)
				}
			} else {
				if !isWhite {
					t.Errorf("Image pixel at (%v, %v) is outside of stencil surface so should be white, not black", x, y)
				}
			}
		}
	}
}

func TestDrawStencilImage(t *testing.T) {
	// crete a dummy image with some black pixels both inside and outside the stencil area,
	// so we can be confident that the src image is really being used correctly.
	srcImg := image.NewRGBA(image.Rect(0, 0, 3, 5))
	srcImg.Set(0, 1, color.Black)
	srcImg.Set(1, 2, color.Black)
	srcImg.Set(4, 0, color.Black)

	stencil := stubStencil{}
	img := DrawStencilImage(stencil, srcImg, color.White)

	stencilW, stencilH := stencil.Size()
	size := img.Bounds().Size()
	if size.X != stencilW {
		t.Errorf("Width should equal %v, not %v", stencilW, size.X)
	}

	if size.Y != stencilH {
		t.Errorf("Height should equal %v, not %v", stencilH, size.Y)
	}

	for y := 0; y < stencilH; y++ {
		for x := 0; x < stencilW; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			srcR, srcG, srcB, srcA := srcImg.At(x, y).RGBA()

			isSrc := r == srcR && g == srcG && b == srcB && a == srcA
			isWhite := r == 65535 && g == 65535 && b == 65535 && a == 65535

			if stencil.Within(x, y) {
				if !isSrc {
					t.Errorf("Image pixel at (%v, %v) is inside of stencil surface "+
						"so should equal the src image pixel (%v, %v, %v, %v), not (%v, %v, %v, %v)",
						x, y, srcR, srcG, srcB, srcA, r, g, b, a)
				}
			} else {
				if !isWhite {
					t.Errorf("Image pixel at (%v, %v) is outside of stencil surface "+
						"so should equal white background color, not (%v, %v, %v, %v)",
						x, y, r, g, b, a)
				}
			}
		}
	}
}

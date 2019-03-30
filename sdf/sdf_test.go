package sdf

import (
	"testing"
)

func TestNew(t *testing.T) {
	sdf := New(3, 5)

	if sdf.Width != 3 {
		t.Errorf("Width should equal 3, not %v", sdf.Width)
	}

	if sdf.Height != 5 {
		t.Errorf("Height should equal 5, not %v", sdf.Height)
	}

	if len(sdf.Field) != 3*5 {
		t.Errorf("Field should have a total of 3*5=15 values, not %v", len(sdf.Field))
	}

	for i, v := range sdf.Field {
		if v != 0 {
			t.Errorf("Field value at index [%v] should be initialised to 0, not %v", i, v)
		}
	}
}

func TestIndexing(t *testing.T) {
	sdf := New(3, 5)

	sdf.Set(1, 1, 0.5)

	res := sdf.At(1, 1)
	if res != 0.5 {
		t.Errorf("Field value at (1,1) should equal 0.5, not %v", res)
	}

	res = sdf.At(2, 2)
	if res != 0 {
		t.Errorf("Field value at (2, 2) should equal 0, not %v", res)
	}
}

func TestNearestBoundaryAt(t *testing.T) {
	df := DisplacementField{
		New(2, 2),
		[]point{
			point{1, 0},
			point{10, 50},
			point{100, 100},
			point{200, 200},
		},
	}

	nx, ny := df.NearestBoundaryAt(1, 0)
	if nx != 10 || ny != 50 {
		t.Errorf("Nearest point should be (10, 50), not (%v, %v)", nx, ny)
	}

	nx, ny = df.NearestBoundaryAt(1, 1)
	if nx != 200 || ny != 200 {
		t.Errorf("Nearest point should be (200, 200), not (%v, %v)", nx, ny)
	}
}

func TestCalculate(t *testing.T) {
	stencil := stubStencil{}
	df := Calculate(stencil)

	if df.Width != 3 {
		t.Errorf("Width should equal 3, not %v", df.Width)
	}

	if df.Height != 5 {
		t.Errorf("Height should equal 5, not %v", df.Height)
	}

	for y := 0; y < df.Height; y++ {
		for x := 0; x < df.Width; x++ {
			f := df.At(x, y)
			nx, ny := df.NearestBoundaryAt(x, y)

			if stencil.Within(x, y) {
				if f > 0 {
					t.Errorf("Field value at (%v, %v) should have been <=0, not %v", x, y, f)
				}

				if nx != x || ny != y {
					t.Errorf("Nearest boundary at (%v, %v) should have equalled itself, not (%v, %v)", x, y, nx, ny)
				}
			} else {
				if f <= 0 {
					t.Errorf("Field value at (%v, %v) should have been >0, not %v", x, y, f)
				}

				if nx == x && ny == y {
					t.Errorf("Nearest boundary at (%v, %v) should NOT have equalled itself", x, y)
				}
			}

			if !stencil.Within(nx, ny) {
				t.Errorf("Nearest boundary point at (%v, %v) should have been within the stencil surface", x, y)
			}
		}
	}
}

func TestDraw(t *testing.T) {
	sdf := New(3, 5)

	sdf.Set(0, 0, 0)

	sdf.Set(1, 1, -3)
	sdf.Set(1, 2, 3)

	sdf.Set(2, 1, -130)
	sdf.Set(2, 2, 130)

	img := sdf.Draw()

	px := img.GrayAt(0, 0).Y
	if px != 127 {
		t.Errorf("Field values of 0 should be drawn at grayscale 127, not %v", px)
	}

	px = img.GrayAt(1, 1).Y
	if px != 127-3 {
		t.Errorf("Field values within the range [-127, -1] such as -3 should be drawn at grayscale 127-3=124, not %v", px)
	}

	px = img.GrayAt(1, 2).Y
	if px != 127+3 {
		t.Errorf("Field values within the range [1, 128] such as 3 should be drawn at grayscale 127+3=130, not %v", px)
	}

	px = img.GrayAt(2, 1).Y
	if px != 0 {
		t.Errorf("Negative field values beyond -127 should be clamped and drawn at grayscale 0, not %v", px)
	}

	px = img.GrayAt(2, 2).Y
	if px != 255 {
		t.Errorf("Positive field values beyond 127 should be clamped and drawn at grayscale 255, not %v", px)
	}
}

func TestLerp(t *testing.T) {
	sdf1 := &SDF{
		Width:  3,
		Height: 1,
		Field:  []float64{0, 1, 50},
	}

	sdf2 := &SDF{
		Width:  3,
		Height: 1,
		Field:  []float64{100, 2, 6},
	}

	tests := []struct {
		tween            float64
		exp0, exp1, exp2 float64
	}{
		{0, 0, 1, 50},
		{0.5, 50, 1.5, 28},
		{1, 100, 2, 6},
	}

	for _, tt := range tests {
		percent := tt.tween * 100
		lerp, err := Lerp(sdf1, sdf2, tt.tween)

		if err != nil {
			t.Errorf("Error should be been nil, not %v", err)
		}

		if lerp.Width != 3 || lerp.Height != 1 {
			t.Errorf("%v%% lerped SDF should have matching width and height of (3, 5), not (%v, %v)", percent, lerp.Width, lerp.Height)
		}

		if len(lerp.Field) != 3 {
			t.Errorf("%v%% lerped SDF should have matching field size length 3, not %v", percent, len(lerp.Field))
		}

		if lerp.Field[0] != tt.exp0 || lerp.Field[1] != tt.exp1 || lerp.Field[2] != tt.exp2 {
			t.Errorf("%v%% lerped SDF should have field values [%v %v %v], not %v", percent, tt.exp0, tt.exp1, tt.exp2, lerp.Field)
		}
	}
}

func TestLerpWithMismatchedSDFs(t *testing.T) {
	sdf1 := New(3, 4)
	sdfMismatchedY := New(3, 5)
	sdfMismatchedX := New(2, 3)

	res1, err1 := Lerp(sdf1, sdfMismatchedY, 0.5)
	res2, err2 := Lerp(sdf1, sdfMismatchedX, 0.5)

	if res1 != nil || res2 != nil {
		t.Errorf("Lerp should not return a result when SDFs are mismatched sized")
	}

	if err1 == nil || err2 == nil {
		t.Errorf("Lerp should return an erro when SDFs are mismatched sizes")
	}
}

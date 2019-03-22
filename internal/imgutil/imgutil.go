package imgutil

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"

	"github.com/daveagill/go-sdf/sdf"
)

// Load an image from file
func Load(path string) image.Image {
	reader, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	img, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	return img
}

// Save an image to file
func Save(path string, img image.Image) {
	f, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	err = png.Encode(f, img)
	if err != nil {
		log.Fatal(err)
	}
}

// SDFToImage converts a Signed-Distance-Field to an 8-bit grayscale Image
func SDFToImage(sdf *sdf.SDF) image.Image {
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

package imgutil

import (
	"bytes"
	"image"
	"image/color"
	"image/gif"
	_ "image/jpeg"
	"image/png"
	"log"
	"math"
	"os"

	"github.com/daveagill/go-sdf/sdf"
)

// Load will load an image from file
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

// SavePNG will save an image to a PNG file
func SavePNG(path string, img image.Image) {
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

// SaveGIF returns an animated GIF given a series of frames
func SaveGIF(path string, frames []image.Image, delay int) {
	outGIF := gif.GIF{
		Image: make([]*image.Paletted, len(frames)),
		Delay: make([]int, len(frames)),
	}

	// convert each frame to a palleted image within the GIF
	for i := range frames {
		buf := bytes.Buffer{}
		gif.Encode(&buf, frames[i], nil)
		palettedFrame, err := gif.Decode(&buf)
		if err != nil {
			log.Fatal(err)
		}

		outGIF.Image[i] = palettedFrame.(*image.Paletted)
		outGIF.Delay[i] = delay
	}

	// open output file
	f, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// write GIF to file
	err = gif.EncodeAll(f, &outGIF)
	if err != nil {
		log.Fatal(err)
	}
}

// FillFromBoundaryPixels returns an image where off-surface pixels are sourced from
// the nearest boundary pixels. Effectively extruding the boundary pixels out to the
// borders of the image.
func FillFromBoundaryPixels(img image.Image, df *sdf.DisplacementField) image.Image {
	outImg := image.NewRGBA(image.Rect(0, 0, df.Width, df.Height))

	for y := 0; y < df.Height; y++ {
		for x := 0; x < df.Width; x++ {
			if df.At(x, y) < 0 {
				outImg.Set(x, y, img.At(x, y))
			} else {
				boundaryX, boundaryY := df.NearestBoundaryAt(x, y)
				outImg.Set(x, y, img.At(boundaryX, boundaryY))
			}
		}
	}

	return outImg
}

// BlendedImage is an Image that tweens between two images according to its Ratio
type BlendedImage struct {
	From  image.Image
	To    image.Image
	Ratio float64
}

// ColorModel implements image.Image's ColorModel() for BlendedImage
func (img *BlendedImage) ColorModel() color.Model {
	return color.RGBA64Model
}

// Bounds implements image.Image's Bounds() for BlendedImage
func (img *BlendedImage) Bounds() image.Rectangle {
	return img.From.Bounds()
}

// At implements image.Image's At(int, int) for BlendedImage
func (img *BlendedImage) At(x, y int) color.Color {
	r1, g1, b1 := toRGB(img.From.At(x, y))
	r2, g2, b2 := toRGB(img.To.At(x, y))

	return color.RGBA{
		R: lerpCol(r1, r2, img.Ratio),
		G: lerpCol(g1, g2, img.Ratio),
		B: lerpCol(b1, b2, img.Ratio),
		A: 255,
	}
}

func lerpCol(u uint16, v uint16, t float64) uint8 {
	return uint8(uint16(float64(u)+(float64(v)-float64(u))*t) >> 8)
}

// toRGB removes the alpha component to return fully opaque non-pre-multiplied RGB values
func toRGB(c color.Color) (uint16, uint16, uint16) {
	r, g, b, a := c.RGBA()
	return uint16((float64(r) / float64(a)) * math.MaxUint16),
		uint16((float64(g) / float64(a)) * math.MaxUint16),
		uint16((float64(b) / float64(a)) * math.MaxUint16)
}

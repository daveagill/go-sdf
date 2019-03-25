package main

import (
	"flag"
	"image"
	"image/color"
	"log"

	"github.com/daveagill/go-sdf/internal/imgutil"
	"github.com/daveagill/go-sdf/sdf"
)

func main() {
	var (
		startPath  string
		endPath    string
		outPath    string
		numFrames  int
		frameDelay int
		boomerang  bool
		blackBg    bool
	)

	flag.StringVar(&startPath, "from", "", "Required. The start image (png/jpeg/gif)")
	flag.StringVar(&endPath, "to", "", "Required. The end image (png/jpeg/gif)")
	flag.StringVar(&outPath, "out", "out.gif", "The output gif")
	flag.IntVar(&numFrames, "frames", 10, "The number of frames to generate")
	flag.IntVar(&frameDelay, "framedelay", 0, "The delay between frames in 100ths of a second")
	flag.BoolVar(&boomerang, "boomerang", true, "Whether to animate back to the initial image (Doubles the number of frames)")
	flag.BoolVar(&blackBg, "blackbg", false, "Uses black as the background color instead of white")
	flag.Parse()

	if startPath == "" {
		log.Fatal("-from not specified")
	}
	if endPath == "" {
		log.Fatal("-to not specified")
	}

	startImg := imgutil.Load(startPath)
	endImg := imgutil.Load(endPath)

	if startImg.Bounds().Size() != endImg.Bounds().Size() {
		log.Fatal("Images do not have the same dimensions")
	}

	startStencil := sdf.ImageAlphaStencil{Image: startImg, Alpha: sdf.HalfAlpha}
	endStencil := sdf.ImageAlphaStencil{Image: endImg, Alpha: sdf.HalfAlpha}

	startField := sdf.Calculate(startStencil)
	endField := sdf.Calculate(endStencil)

	blendedImg := &imgutil.BlendedImage{
		From: imgutil.FillFromBoundaryPixels(startImg, startField),
		To:   imgutil.FillFromBoundaryPixels(endImg, endField),
	}

	bgCol := color.White
	if blackBg {
		bgCol = color.Black
	}

	// animate from start to end
	frames := make([]image.Image, numFrames, numFrames*2)
	for i := 0; i < numFrames; i++ {
		blendedImg.Ratio = float64(i) / float64(numFrames-1)
		blendedSDF, _ := sdf.Lerp(startField.SDF, endField.SDF, blendedImg.Ratio)
		frames[i] = blendedSDF.DrawStenciledImage(blendedImg, bgCol)
	}

	// create the reverse sequence of frames to 'boomerang' back to the start
	if boomerang {
		for i := numFrames - 1; i >= 0; i-- {
			frames = append(frames, frames[i])
		}
	}

	imgutil.SaveGIF(outPath, frames, frameDelay)
}

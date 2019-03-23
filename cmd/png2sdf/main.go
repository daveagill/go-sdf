package main

import (
	"log"
	"os"

	"github.com/daveagill/go-sdf/internal/imgutil"
	"github.com/daveagill/go-sdf/sdf"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatalf("2 arguments expected: %s input.png sdf_output.png", os.Args[0])
	}
	inpath := os.Args[1]
	outpath := os.Args[2]

	img := imgutil.Load(inpath)
	field := sdf.FromImageAlpha(img, sdf.HalfAlpha)
	grayImg := field.Draw()
	imgutil.SavePNG(outpath, grayImg)
}

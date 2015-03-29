package main

import (
	"fmt"
	"image"
	"image/png"

	"github.com/thinkofdeath/steven/resource"
)

var (
	grassBiomeColors   *image.NRGBA
	foliageBiomeColors *image.NRGBA
)

func init() {
	grassBiomeColors = loadBiomeColors("grass")
	foliageBiomeColors = loadBiomeColors("foliage")
}

func loadBiomeColors(name string) *image.NRGBA {
	f, _ := resource.Open("minecraft", fmt.Sprintf("textures/colormap/%s.png", name))
	defer f.Close()
	img, err := png.Decode(f)
	if err != nil {
		panic(err)
	}
	return img.(*image.NRGBA)
}

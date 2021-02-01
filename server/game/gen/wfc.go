package gen

import (
	"errors"
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"image/color"
	"image/png"

	"github.com/nfnt/resize"
	"github.com/shawnridgeway/wfc"

	_ "image/png"
)

var wfcTextures = loadTextures()

var (
	colorGround = color.RGBA{255, 255, 255, 255}
	colorWall   = color.RGBA{0, 0, 0, 255}
)

func loadTextures() map[string]image.Image {
	textures := map[string]image.Image{}
	files, err := ioutil.ReadDir("../data/textures")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		split := strings.Split(file.Name(), ".")
		name := split[0]
		f, err := os.Open(fmt.Sprintf("../data/textures/%s", file.Name()))
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		img, _, err := image.Decode(f)
		if err != nil {
			log.Fatal(err)
		}

		textures[name] = img
	}

	return textures
}

func runWfc(textureName string, width, height, scale int, saveGen bool) (image.Image, error) {
	texture, ok := wfcTextures[textureName]
	if !ok {
		return nil, errors.New("texture not found")
	}

	model := wfc.NewOverlappingModel(texture, 3, width, height, true, false, 4, false)
	outputImg, success := model.Generate()
	if !success {
		return nil, errors.New("image not generated")
	}

	outputImg = resize.Resize(uint(width*scale), uint(height*scale), outputImg, resize.NearestNeighbor)

	if saveGen {
		f, err := os.Create(fmt.Sprintf("../data/textures/out/%s.png", textureName))
		if err != nil {
			return nil, err
		}
		defer f.Close()

		err = png.Encode(f, outputImg)
		if err != nil {
			return nil, err
		}
	}

	return outputImg, nil
}

func colorsEqual(a, b color.Color) bool {
	ar, ag, ab, aa := a.RGBA()
	br, bg, bb, ba := b.RGBA()
	return ar == br && ag == bg && ab == bb && aa == ba
}

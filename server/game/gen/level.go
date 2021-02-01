package gen

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"math/rand"
	"os"

	"github.com/nickdavies/go-astar/astar"
)

const (
	levelWidth  = 84
	levelHeight = 84
	wfcScale    = 3
)

// Level ...
type Level struct {
	Width, Height int
	Tiles         [][]*Tile // x, y indexed
}

// BuildLevel ...
func BuildLevel() (*Level, error) {
	l := &Level{
		Width:  levelWidth,
		Height: levelHeight,
	}

	l.Tiles = make([][]*Tile, l.Width)
	for x := 0; x < l.Width; x++ {
		l.Tiles[x] = make([]*Tile, l.Height)
		for y := 0; y < l.Height; y++ {
			l.Tiles[x][y] = &Tile{X: x, Y: y}
		}
	}

	// run wfc for level
	log.Println("running wfc")
	wfcImage, err := runWfc("rooms", levelWidth/wfcScale, levelHeight/wfcScale, wfcScale, false)
	if err != nil {
		return nil, err
	}
	l.parseImage(wfcImage)

	// add a wall border because wfc always goes to the edge
	l.addBorder()

	// connect non-contiguous regions (first region is always walls)
	log.Println("connecting regions")
	regions := l.regions()
	for len(regions) > 2 {
		l.connectRegions(regions[1], regions[2]) // TODO: try randomly selecting regions instead of first 2
		regions = l.regions()
	}

	log.Println("removing superfluous walls")
	l.removeSuperfluousWalls()

	log.Println("placing doors")
	l.placeDoors()

	l.saveImage("dungeon")

	return l, nil
}

// read an image into level data (from wfc likely)
func (l *Level) parseImage(img image.Image) error {
	for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
		for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
			at := img.At(x, y)
			for tType, c := range tileTypeToColor {
				if colorsEqual(at, c) {
					l.Tiles[x][y].Type = tType
					break
				}
			}

		}
	}

	return nil
}

func (l *Level) addBorder() {
	for x := 0; x < l.Width; x++ {
		for y := 0; y < l.Height; y++ {
			if x == 0 || x == l.Width-1 || y == 0 || y == l.Height-1 {
				l.Tiles[x][y].Type = TileTypeWall
			}
		}
	}
}

func (l *Level) connectRegions(regionA, regionB []*Tile) {
	start := regionA[rand.Intn(len(regionA))]
	end := regionB[rand.Intn(len(regionB))]

	a := astar.NewAStar(l.Width, l.Height)
	p2p := astar.NewPointToPoint()
	for x := 0; x < l.Width; x++ {
		for y := 0; y < l.Height; y++ {
			t := l.Tiles[x][y]
			if t.Type == TileTypeWall {
				a.FillTile(astar.Point{Row: x, Col: y}, 3)
			}
		}
	}

	source := []astar.Point{{Row: start.X, Col: start.Y}}
	target := []astar.Point{{Row: end.X, Col: end.Y}}

	path := a.FindPath(p2p, source, target)
	for path != nil {
		if l.Tiles[path.Row][path.Col].Type == TileTypeWall {
			l.Tiles[path.Row][path.Col].Type = TileTypeHall
		}
		path = path.Parent
	}
}

func (l *Level) removeSuperfluousWalls() {
	superfluousWalls := map[*Tile]bool{}

	for x := 0; x < l.Width; x++ {
		for y := 0; y < l.Height; y++ {
			if l.superfluousWall(x, y) {
				superfluousWalls[l.Tiles[x][y]] = true
			}
		}
	}

	for t := range superfluousWalls {
		t.Type = TileTypeAir
	}
}

func (l *Level) placeDoors() {
	doors := map[*Tile]bool{}

	for x := 0; x < l.Width; x++ {
		for y := 0; y < l.Height; y++ {
			t := l.Tiles[x][y]
			if t.Type == TileTypeHall {
				neighbors := l.neighbors(x, y)
				for _, n := range neighbors {
					if n.Type == TileTypeGround {
						doors[t] = true
					}
				}
			}
		}
	}

	for t := range doors {
		t.Type = TileTypeDoor
	}
}

// for debugging
func (l *Level) saveImage(name string) error {
	f, err := os.Create(fmt.Sprintf("../data/textures/out/%s.png", name))
	if err != nil {
		return err
	}
	defer f.Close()

	outImg := image.NewRGBA(image.Rectangle{Max: image.Point{X: l.Width, Y: l.Height}})
	for x := 0; x < l.Width; x++ {
		for y := 0; y < l.Height; y++ {
			t := l.Tiles[x][y]
			outImg.Set(x, y, tileTypeToColor[t.Type])
		}
	}

	err = png.Encode(f, outImg)
	if err != nil {
		return err
	}

	return nil
}

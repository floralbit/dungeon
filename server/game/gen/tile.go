package gen

import "image/color"

// TileType ...
type TileType int

const (
	// TileTypeWall ...
	TileTypeWall = iota
	// TileTypeGround ...
	TileTypeGround
	// TileTypeHall ...
	TileTypeHall
	// TileTypeAir ...
	TileTypeAir
)

// Tile ...
type Tile struct {
	Type TileType
	X, Y int

	RegionColor int // contiguous region id
}

// for debug image saving
var tileTypeToColor = map[TileType]color.Color{
	TileTypeWall:   color.Black,
	TileTypeGround: color.White,
	TileTypeHall:   color.RGBA{128, 128, 128, 255},
	TileTypeAir:    color.RGBA{255, 255, 255, 0},
}

func (l *Level) inBounds(x, y int) bool {
	if x < 0 || x >= l.Width || y < 0 || y >= l.Height {
		return false
	}
	return true
}

func (l *Level) neighbors(x, y int) []*Tile {
	tiles := []*Tile{}

	// just up, down, left, right for now - probably a better way to do this but I'm lazy rn
	if l.inBounds(x-1, y) {
		tiles = append(tiles, l.Tiles[x-1][y])
	}
	if l.inBounds(x+1, y) {
		tiles = append(tiles, l.Tiles[x+1][y])
	}
	if l.inBounds(x, y-1) {
		tiles = append(tiles, l.Tiles[x][y-1])
	}
	if l.inBounds(x, y+1) {
		tiles = append(tiles, l.Tiles[x][y+1])
	}

	return tiles
}

func (l *Level) allNeighbors(x, y int) []*Tile {
	tiles := []*Tile{}
	for i := x - 1; i <= x+1; i++ {
		for j := y - 1; j <= y+1; j++ {
			if i == x && j == y || !l.inBounds(i, j) {
				continue
			}
			tiles = append(tiles, l.Tiles[i][j])
		}
	}
	return tiles
}

func (l *Level) superfluousWall(x, y int) bool {
	t := l.Tiles[x][y]
	if t.Type != TileTypeWall {
		return false
	}

	neighbors := l.allNeighbors(x, y)
	for _, n := range neighbors {
		if n.Type != TileTypeWall {
			return false
		}
	}

	return true
}

func (l *Level) regions() [][]*Tile {
	regionCount := l.colorRegions()
	regions := make([][]*Tile, regionCount+1)

	for x := 0; x < l.Width; x++ {
		for y := 0; y < l.Height; y++ {
			color := l.Tiles[x][y].RegionColor
			regions[color] = append(regions[color], l.Tiles[x][y])
		}
	}

	return regions
}

func (l *Level) colorRegions() int {
	visited := map[*Tile]bool{}

	currentColor := 0
	for x := 0; x < l.Width; x++ {
		for y := 0; y < l.Height; y++ {
			t := l.Tiles[x][y]
			if _, in := visited[t]; in {
				continue
			}
			if t.Type == TileTypeWall {
				continue
			}

			currentColor++
			l.colorVisit(t, currentColor, visited)
		}
	}

	return currentColor
}

func (l *Level) colorVisit(t *Tile, color int, visited map[*Tile]bool) {
	if _, in := visited[t]; in {
		return
	}
	if t.Type == TileTypeWall {
		return
	}

	visited[t] = true
	t.RegionColor = color

	neighbors := l.neighbors(t.X, t.Y)
	for _, n := range neighbors {
		l.colorVisit(n, color, visited)
	}
}

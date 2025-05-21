package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// clear clears the terminal screen using ANSI escape codes.
func clear() {
	fmt.Print("\033[H\033[2J")
}

const (
	width       = 80
	height      = 24
	frameDelay  = 80 * time.Millisecond
	perspective = 5.0
	cubeSize    = 1.5
	starCount   = 40
)

type point3D struct{ x, y, z float64 }
type star struct{ pos point3D }

type point2D struct{ x, y int }

func rotateX(p point3D, angle float64) point3D {
	cosA := math.Cos(angle)
	sinA := math.Sin(angle)
	return point3D{
		x: p.x,
		y: p.y*cosA - p.z*sinA,
		z: p.y*sinA + p.z*cosA,
	}
}

func project(p point3D) (int, int, bool) {
	z := p.z + perspective
	if z <= 0 {
		return 0, 0, false
	}
	scale := perspective / z
	x := int(p.x*scale*float64(width)/2 + float64(width)/2)
	y := int(p.y*scale*float64(height)/2 + float64(height)/2)
	if x < 0 || x >= width || y < 0 || y >= height {
		return 0, 0, false
	}
	return x, y, true
}

var cubeVertices = []point3D{
	{-cubeSize, -cubeSize, -cubeSize},
	{cubeSize, -cubeSize, -cubeSize},
	{-cubeSize, cubeSize, -cubeSize},
	{cubeSize, cubeSize, -cubeSize},
	{-cubeSize, -cubeSize, cubeSize},
	{cubeSize, -cubeSize, cubeSize},
	{-cubeSize, cubeSize, cubeSize},
	{cubeSize, cubeSize, cubeSize},
}

var edges = [][2]int{
	{0, 1}, {1, 3}, {3, 2}, {2, 0},
	{4, 5}, {5, 7}, {7, 6}, {6, 4},
	{0, 4}, {1, 5}, {2, 6}, {3, 7},
}

var stars []star

func initStars() {
	stars = make([]star, starCount)
	for i := range stars {
		stars[i].pos = point3D{
			x: (rand.Float64() - 0.5) * cubeSize * 6,
			y: (rand.Float64() - 0.5) * cubeSize * 6,
			z: (rand.Float64() - 0.5) * cubeSize * 6,
		}
	}
}

func drawLine(buf [][]rune, x1, y1, x2, y2 int, char rune) {
	dx := int(math.Abs(float64(x2 - x1)))
	dy := int(math.Abs(float64(y2 - y1)))
	sx := -1
	if x1 < x2 {
		sx = 1
	}
	sy := -1
	if y1 < y2 {
		sy = 1
	}
	err := dx - dy
	for {
		buf[y1][x1] = char
		if x1 == x2 && y1 == y2 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x1 += sx
		}
		if e2 < dx {
			err += dx
			y1 += sy
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	initStars()
	angle := 0.0
	for {
		buffer := make([][]rune, height)
		for i := range buffer {
			buffer[i] = make([]rune, width)
			for j := range buffer[i] {
				buffer[i][j] = ' '
			}
		}

		// rotate cube vertices
		rotated := make([]point3D, len(cubeVertices))
		for i, v := range cubeVertices {
			rotated[i] = rotateX(v, angle)
		}

		// draw cube edges
		for _, e := range edges {
			x1, y1, ok1 := project(rotated[e[0]])
			x2, y2, ok2 := project(rotated[e[1]])
			if ok1 && ok2 {
				drawLine(buffer, x1, y1, x2, y2, '*')
			}
		}

		// draw moving stars
		for i := range stars {
			stars[i].pos.x += 0.02
			if stars[i].pos.x > cubeSize*3 {
				stars[i].pos.x = -cubeSize * 3
			}
			p := rotateX(stars[i].pos, angle)
			x, y, ok := project(p)
			if ok {
				buffer[y][x] = '.'
			}
		}

		clear()
		for _, row := range buffer {
			fmt.Println(string(row))
		}
		time.Sleep(frameDelay)
		angle += 0.05
		if angle > 2*math.Pi {
			angle -= 2 * math.Pi
		}
	}
}

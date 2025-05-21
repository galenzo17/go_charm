package main

import (
	"fmt"
	"math"
	"strings"
	"time"
)

// clear clears the terminal screen using ANSI escape codes.
func clear() {
	fmt.Print("\033[H\033[2J")
}

const diamondSize = 8 // base size of the diamond
const frameDelay = 100 * time.Millisecond

// drawDiamond returns an ASCII diamond with vertical scaling applied.
func drawDiamond(scale float64) string {
	if scale < 0.1 {
		scale = 0.1
	}
	height := int(float64(diamondSize) * scale)
	var b strings.Builder
	for y := -height; y <= height; y++ {
		width := int(float64(diamondSize) - math.Abs(float64(y)/scale))
		if width < 0 {
			width = 0
		}
		spaces := diamondSize - width
		b.WriteString(strings.Repeat(" ", spaces))
		b.WriteString(strings.Repeat("*", width*2))
		b.WriteString("\n")
	}
	return b.String()
}

func main() {
	angle := 0.0
	for {
		scale := math.Abs(math.Cos(angle))*0.9 + 0.1 // keep scale > 0
		clear()
		fmt.Print(drawDiamond(scale))
		time.Sleep(frameDelay)
		angle += 0.15
		if angle > 2*math.Pi {
			angle -= 2 * math.Pi
		}
	}
}

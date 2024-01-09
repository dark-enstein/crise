package main

import (
	"github.com/dark-enstein/crise/internal/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image/color"
	"log"
	"sync"
)

var (
	SCREEN_W = 640
	SCREEN_H = 640
)

type P_CONTENT struct {
	Text string
}

type Game struct {
	// pFunc is to update the current game state
	pFunc func()
	// LeadText holds the structure P_CONTENT that holds the content for the current iteration of update for the current phase.
	LeadText P_CONTENT
	// utils holds extra meta data and utils of each frame and phase
	utils utils.Screen
	// handling concurrent modification of game state
	sync.Mutex
}

func (g *Game) Update() error {
	x, y := ebiten.CursorPosition()
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		log.Printf("loc: x=%d; y=%d\n", x, y)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	BGColor := color.Black
	screen.Fill(BGColor)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return SCREEN_W, SCREEN_W
}

func main() {
	ebiten.SetWindowSize(640, 640)
	ebiten.SetWindowTitle("Tetris")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}

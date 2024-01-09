package game

import (
	"context"
	"github.com/dark-enstein/crise/internal/tetra"
	utils2 "github.com/dark-enstein/crise/internal/utils"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	MAX_TETROMINO_ONSCREEN = 300
)

// TetrominoManager manages the spawn and despqwn of Tetrominoes in game.
type TetrominoManager struct {
	// activeMember is the currently active of current onScreen elements.
	activeMember *tetra.Tetromino
	// creationCounter defines the cycles of updates before a new Tetromino is spawned
	creationCounter int
	// activeN is the integer location of the currently active Tetromino in onScreen
	activeN int
	// onScreenBank stores all the currently displayed Tetromino on screen
	onScreenBank []*tetra.Tetromino
	// settings holds config settings that TetrominoManager uses to manage the tetris animation
	settings *TetroSettings
}

// TetroSettings defines the settings used by the manager to manage tetris animation
type TetroSettings struct {
	tIncrement int
}

// NewTetrominoMananger creates a new Tetromino manager. inc is the preferred increment or speed of the Tetromino on key direction directive. Right now it is measured in pixels on the screen, but later it would be changed to be a multiple of utils2.SPRITE_HEIGHT
func NewTetrominoMananger(inc int) *TetrominoManager {
	return &TetrominoManager{
		onScreenBank: make([]*tetra.Tetromino, 0, MAX_TETROMINO_ONSCREEN),
		settings:     &TetroSettings{tIncrement: inc},
	}
}

// Add adds a new Tetromino to the screen
func (t *TetrominoManager) Add() {
	//rand.Intn(len(tetra.TetraCoordinates))
	t.onScreenBank = append(t.onScreenBank, tetra.NewTetromino(tetra.TetraCoordinates[tetra.J], utils2.WALL_WIDTH, utils2.WALL_HEIGHT))
	t.activeN = len(t.onScreenBank) - 1
	t.activeMember = t.onScreenBank[t.activeN]
}

// Active returns the currently active Tetromino. This is what is currently controllable by the player.
func (t *TetrominoManager) Active() *tetra.Tetromino {
	return t.activeMember
}

// OnScreenBank returns all the Tetrominoes currently on screen.
func (t *TetrominoManager) OnScreenBank() []*tetra.Tetromino {
	return t.onScreenBank
}

// WillCollide implements collission detection for Tetris
func (t *TetrominoManager) WillCollide(t0 *tetra.Tetromino) {}

// toRealIncrement converts an increment (measured in number of sprites) to real increment (based on the dimensions of the screen)
func toRealIncrement(n int) int {
	return n * utils2.SPRITE_HEIGHT
}

func (t *TetrominoManager) Display(ctx context.Context, screen *ebiten.Image) {
	//log.Printf("calling build on N tetromino: %v\n", len(t.onScreenBank))
	for i := 0; i < len(t.onScreenBank); i++ {
		//log.Printf("calling build on tetromino: %v\n", t.onScreenBank[i])
		t.onScreenBank[i].Build(ctx)(screen)
	}
}

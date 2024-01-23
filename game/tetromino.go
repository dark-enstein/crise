package game

import (
	"context"
	"fmt"
	"github.com/dark-enstein/crise/internal/tetra"
	utils2 "github.com/dark-enstein/crise/internal/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
	"math/rand"
	"sync"
	"time"
)

const (
	MAX_TETROMINO_ONSCREEN = 300
	ACCELERATION_FACTOR    = 0.5
	TICKER_DURATION        = 1
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
	// TetrisTicker tickes per duration sending current time via the channel at every interval
	TetrisTicker *time.Ticker
	// DoneChan channel receives a completion signal on done, or error during game play
	DoneChan chan bool
	// settings holds config settings that TetrominoManager uses to manage the tetris animation
	settings *TetroSettings
	// handling syncronization and change of state
	sync.Mutex
}

// TetroSettings defines the settings used by the manager to manage tetris animation
type TetroSettings struct {
	// tIncrement holds the increment value valid, given the current input. It can always be mutated
	tIncrement float32
	// tIncrementSaved holds the increment value as defined in settings. It never changes.
	tIncrementSaved float32
}

// Accelerate increases the value of the increment by a factor to simulate acceleration
func (ts *TetroSettings) Accelerate() {
	ts.tIncrement += ACCELERATION_FACTOR
}

// Reset resets the currnetly active increment value to original as set within from settings
func (ts *TetroSettings) Reset() {
	ts.tIncrement = ts.tIncrementSaved
}

// NewTetrominoMananger creates a new Tetromino manager. inc is the preferred increment or speed of the Tetromino on key direction directive. Right now it is measured in pixels on the screen, but later it would be changed to be a multiple of utils2.SPRITE_HEIGHT
func NewTetrominoMananger(inc int) *TetrominoManager {
	log.Println("creating new tetromino")
	manager := newTetrominoMananger(inc)
	log.Println("created new tetromino, adding a tetromino")
	manager.Add()
	log.Println("added a tetromino:", manager.onScreenBank)
	return manager
}

// newTetrominoMananger creates a new Tetromino manager. inc is the preferred increment or speed of the Tetromino on key direction directive. Right now it is measured in pixels on the screen, but later it would be changed to be a multiple of utils2.SPRITE_HEIGHT
func newTetrominoMananger(inc int) *TetrominoManager {
	return &TetrominoManager{
		onScreenBank: make([]*tetra.Tetromino, 0, MAX_TETROMINO_ONSCREEN),
		settings:     &TetroSettings{tIncrement: float32(inc)},
		TetrisTicker: time.NewTicker(TICKER_DURATION * time.Second),
		DoneChan:     make(chan bool, 1),
	}
}

// Add adds a new Tetromino to the screen
func (t *TetrominoManager) Add() {
	//rand.Intn(len(tetra.TetraCoordinates))
	//tetra.Mini(rand.Intn(len(tetra.TetraCoordinates)))
	t.onScreenBank = append(t.onScreenBank, tetra.NewTetromino(tetra.TetraCoordinates[tetra.Mini(rand.Intn(len(tetra.TetraCoordinates)))], utils2.SPRITE_X0, utils2.SPRITE_Y0))
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

// ActiveTArray returns the 2D array of the currently active Tetromino
func (t *TetrominoManager) ActiveTArray() [][]int {
	return t.onScreenBank[t.activeN].Arr
}

// MoveDown moves the tetromino downward
func (t *TetrominoManager) MoveDown() {
	log.Println("acquiring move down lock")
	t.Lock()
	time.Sleep(3)
	log.Println("bank state before moving down:", t.OnScreenBank())
	for i := 0; i < len(t.ActiveTArray()); i++ {
		fX, fY := t.ActiveTArray()[i][0], t.ActiveTArray()[i][1]
		//g.sample[i][0] += INCREMENT
		t.onScreenBank[t.activeN].Arr[i][1] += t.inc()
		fmt.Printf("moved sprites from %v,%v to %v,%v\n", fX, fY, t.onScreenBank[t.activeN].Arr[i][0], t.onScreenBank[t.activeN].Arr[i][1])
	}
	log.Println("releasing move down lock")
	t.Unlock()
}

// MoveUp moves the tetromino upwards
func (t *TetrominoManager) MoveUp() {
	log.Println("acquiring move up lock")
	t.Lock()
	for i := 0; i < len(t.ActiveTArray()); i++ {
		fX, fY := t.ActiveTArray()[i][0], t.ActiveTArray()[i][1]
		//g.sample[i][0] += INCREMENT
		t.onScreenBank[t.activeN].Arr[i][1] -= t.inc()
		fmt.Printf("moved sprites from %v,%v to %v,%v\n", fX, fY, t.onScreenBank[t.activeN].Arr[i][0], t.onScreenBank[t.activeN].Arr[i][1])
	}
	log.Println("releasing move up lock")
	t.Unlock()
}

// MoveLeft moves the tetromino leftward
func (t *TetrominoManager) MoveLeft() {
	log.Println("acquiring move left lock")
	t.Lock()
	for i := 0; i < len(t.ActiveTArray()); i++ {
		fX, fY := t.ActiveTArray()[i][0], t.ActiveTArray()[i][1]
		t.onScreenBank[t.activeN].Arr[i][0] -= t.inc()
		//g.sample[i][1] += INCREMENT
		fmt.Printf("moved sprites from %v,%v to %v,%v\n", fX, fY, t.onScreenBank[t.activeN].Arr[i][0], t.onScreenBank[t.activeN].Arr[i][1])
	}
	log.Println("releasing move left lock")
	t.Unlock()
}

// MoveRight moves the tetromino rightward
func (t *TetrominoManager) MoveRight() {
	log.Println("acquiring move right lock")
	t.Lock()
	for i := 0; i < len(t.ActiveTArray()); i++ {
		fX, fY := t.ActiveTArray()[i][0], t.ActiveTArray()[i][1]
		t.onScreenBank[t.activeN].Arr[i][0] += t.inc()
		//g.sample[i][1] += INCREMENT
		fmt.Printf("moved sprites from %v,%v to %v,%v\n", fX, fY, t.onScreenBank[t.activeN].Arr[i][0], t.onScreenBank[t.activeN].Arr[i][1])
	}
	log.Println("releasing move right lock")
	t.Unlock()
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

// Accelerate accelerates the current tetromino. On key press
func (t *TetrominoManager) Accelerate() {
	t.settings.Accelerate()
}

// ResetInc resets the increment of the current tetromino. On key release
func (t *TetrominoManager) ResetInc() {
	t.settings.Reset()
}

// inc returns the current increment value cast into int
func (t *TetrominoManager) inc() int {
	return int(t.settings.tIncrement)
}

// IsAccelerated checks if the currently active tetromino is accelerated
func (t *TetrominoManager) IsAccelerated() bool {
	if t.settings.tIncrement < t.settings.tIncrementSaved || t.settings.tIncrement == t.settings.tIncrementSaved {
		return false
	}
	return true
}

// CronMove checks if the active tetromino has already moved in the current cycle
func (t *TetrominoManager) CronMove() {
	log.Println("stepping down")
	t.MoveDown()
}

// Drain empties the onScreenBank in preparation for shutdown or restart
func (t *TetrominoManager) Drain() {
	log.Println("stepping down")
	log.Println("about cleaning onScreenBank:", t.onScreenBank)
	t.onScreenBank = t.onScreenBank[:0]
	log.Println("done cleaning onScreenBank:", t.onScreenBank)
}

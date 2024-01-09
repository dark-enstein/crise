package game

import (
	"context"
	"fmt"
	"github.com/dark-enstein/crise/config"
	utils2 "github.com/dark-enstein/crise/internal/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"image/color"
	"log"
	"os"
	"path/filepath"
	"sync"
	"text/tabwriter"
	"time"
)

const (
	EXEC_DEBUG = iota + 1
	EXEC_PRODUCTION
)

var (
	DefaultFont font.Face
	COLOR_RED   = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	RUNTIME_DIR = "build/runtime"
	HELP_TEXT   = `
Here's how to play:

1. Moving Tetromino: Use the left and right arrow keys to move the falling tetromino left or right.
2. Rotating Tetromino: Press the up arrow key to rotate the tetromino.
3. Speed Up Drop: Press the down arrow key to speed up the tetromino's fall.
3.Hard Drop: Press the spacebar to instantly drop the tetromino to the bottom.
4.Pause Game: Press the 'P' key to pause and resume the game.
5. Scoring: Complete lines to score points. The more lines you clear at once, the higher your score.
6. Game Over: The game ends when the tetrominoes stack up to the top.
7. Restart Game: Press the 'R' key to restart the game after a game over.
`
	INCREMENT = 1 // effectively val * Sprite height
)

type Game struct {
	// ctx holds the main context
	ctx *context.Context
	// CFunc for graceful termination
	CFunc context.CancelCauseFunc
	// Settings holds the config for the entire game
	Settings     *config.Gameplay
	FeatureFlags map[string]config.FFlags
	// P holds the current phase of the game. See Phase, WELCOME
	P Phase
	// pFunc is to update the current game state
	pFunc func()
	// LeadText holds the structure P_CONTENT that holds the content for the current iteration of update for the current phase.
	LeadText P_CONTENT
	// utils holds extra meta data and utils of each frame and phase
	utils utils2.Screen
	// handling concurrent modification of game state
	sync.Mutex
	// helping with debugging mouse click positions
	counter int
	x, y    int
	// tetris state
	sample    [][]int
	Tetromino *TetrominoManager
	// Running mode
	ExecMode int
}

func NewGameWithContext(ctx *context.Context, cFunc context.CancelCauseFunc, startPhase Phase, settings *config.Gameplay, mode int, fflags map[string]config.FFlags) *Game {
	log.Println("creating new game")
	return &Game{
		ctx:          ctx,
		CFunc:        cFunc,
		P:            startPhase,
		Settings:     settings,
		ExecMode:     mode,
		FeatureFlags: fflags,
		sample:       [][]int{{305, 316}, {275, 346}, {305, 346}},
		Tetromino:    NewTetrominoMananger(INCREMENT),
	}
}

func NewGame(cFunc context.CancelCauseFunc, startPhase Phase, settings *config.Gameplay, mode int, fflags map[string]config.FFlags) *Game {
	log.Println("creating new game")
	return &Game{
		ctx:          nil,
		CFunc:        cFunc,
		P:            startPhase,
		Settings:     settings,
		ExecMode:     mode,
		FeatureFlags: fflags,
		sample:       [][]int{{305, 316}, {275, 346}, {305, 346}},
		Tetromino:    NewTetrominoMananger(INCREMENT),
	}
}

func prepDir(cf context.CancelCauseFunc, ff map[string]config.FFlags) {
	if _, err := os.Stat(RUNTIME_DIR); err != nil && os.IsNotExist(err) {
		err := os.Mkdir(RUNTIME_DIR, 0777)
		if err != nil {
			log.Printf("error while setting up settings dir: %s\n", err.Error())
			if len(ff) > 0 && ff[config.FLAG_GRACEFULTERM] == config.FlagEnabled {
				cf(fmt.Errorf("error while setting up settings dir: %s\n", err.Error()))
			}
		}
	} else {
		return
	}
}

func (g *Game) genWelcomeText() string {
	file, err := os.OpenFile(filepath.Join(RUNTIME_DIR, ".save"), os.O_TRUNC|os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		log.Printf("error while setting up temp file: %s\n", err.Error())
		if len(g.FeatureFlags) > 0 && g.FeatureFlags[config.FLAG_GRACEFULTERM] == config.FlagEnabled {
			g.CFunc(fmt.Errorf("error while setting up temp file: %s\n", err.Error()))
		}
	}
	w := tabwriter.NewWriter(file, 0, 0, 1, ' ', 0)
	fmt.Fprintf(w, "\nHigh Scores:\n\t\nYou:\t%v\nAll time:\t%v\n", g.Settings.You, g.Settings.Highest)
	w.Flush()
	fName := file.Name()
	file.Close()
	f, _ := os.ReadFile(fName)
	wellText := string(f)
	os.Remove(fName)
	return wellText
}

func (g *Game) Update() error {
	ok := false
	if g.LeadText, ok = PHASES[g.P]; !ok {
		log.Printf("phase %v unrecognized\n", g.P)
		os.Exit(1)
	}
	x, y := ebiten.CursorPosition()
	switch g.P {
	case INIT:
	case WELCOME:
		g.debugPos("left")
		if (x >= 120 && y >= 152) && (x <= 240 && y <= 222) {
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				log.Printf("play pressed: x=%d, y=%d\n", x, y)
				g.Lock()
				g.P = START
				g.Unlock()
			}
		}
		if (x >= 400 && y >= 152) && (x <= 520 && y <= 223) {
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				log.Printf("help pressed: x=%d, y=%d\n", x, y)
				g.Lock()
				g.P = HELP
				g.Unlock()
			}
		}
		if (x >= 400 && y >= 249) && (x <= 519 && y <= 321) {
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				log.Printf("debug pressed: x=%d, y=%d\n", x, y)
				g.Lock()
				g.P = DEBUG
				g.Unlock()
			}
		}
	case START:
		// reset func // later random gen
		fmt.Println("UPDATE_START")
		g.debugPos("left")
		if (x >= 404 && y >= 449) && (x <= 574 && y <= 533) {
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				log.Printf("quit pressed: x=%d, y=%d\n", x, y)
				g.Lock()
				g.P = QUIT
				g.Unlock()
			}
		}
		if (x >= 404 && y >= 47) && (x <= 585 && y <= 132) {
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				log.Printf("restart pressed: x=%d, y=%d\n", x, y)
				g.Lock()
				g.P = WELCOME
				g.Unlock()
			}
		}
		if len(g.Tetromino.onScreenBank) > 0 {
			currActiveArr := g.Tetromino.onScreenBank[g.Tetromino.activeN].Arr
			if inpututil.IsKeyJustReleased(ebiten.KeyArrowDown) {
				log.Printf("key arrow down pressed: x=%d, y=%d\n", x, y)
				g.Lock()
				for i := 0; i < len(currActiveArr); i++ {
					fX, fY := currActiveArr[i][0], currActiveArr[i][1]
					//g.sample[i][0] += INCREMENT
					g.Tetromino.onScreenBank[g.Tetromino.activeN].Arr[i][1] += g.Tetromino.settings.tIncrement
					fmt.Printf("moved sprites from %v,%v to %v,%v\n", fX, fY, g.Tetromino.onScreenBank[g.Tetromino.activeN].Arr[i][0], g.Tetromino.onScreenBank[g.Tetromino.activeN].Arr[i][1])
				}
				g.Unlock()
			}
			if inpututil.IsKeyJustReleased(ebiten.KeyArrowUp) {
				log.Printf("key arrow up pressed: x=%d, y=%d\n", x, y)
				g.Lock()
				for i := 0; i < len(currActiveArr); i++ {
					fX, fY := currActiveArr[i][0], currActiveArr[i][1]
					//g.sample[i][0] += INCREMENT
					g.Tetromino.onScreenBank[g.Tetromino.activeN].Arr[i][1] -= g.Tetromino.settings.tIncrement
					fmt.Printf("moved sprites from %v,%v to %v,%v\n", fX, fY, g.Tetromino.onScreenBank[g.Tetromino.activeN].Arr[i][0], g.Tetromino.onScreenBank[g.Tetromino.activeN].Arr[i][1])
				}
				g.Unlock()
			}
			if inpututil.IsKeyJustReleased(ebiten.KeyArrowRight) {
				log.Printf("key arrow right pressed: x=%d, y=%d\n", x, y)
				g.Lock()
				for i := 0; i < len(currActiveArr); i++ {
					fX, fY := currActiveArr[i][0], currActiveArr[i][1]
					g.Tetromino.onScreenBank[g.Tetromino.activeN].Arr[i][0] += g.Tetromino.settings.tIncrement
					//g.sample[i][1] += INCREMENT
					fmt.Printf("moved sprites from %v,%v to %v,%v\n", fX, fY, g.Tetromino.onScreenBank[g.Tetromino.activeN].Arr[i][0], g.Tetromino.onScreenBank[g.Tetromino.activeN].Arr[i][1])
				}
				g.Unlock()
			}
			if inpututil.IsKeyJustReleased(ebiten.KeyArrowLeft) {
				log.Printf("key arrow left pressed: x=%d, y=%d\n", x, y)
				g.Lock()
				for i := 0; i < len(currActiveArr); i++ {
					fX, fY := currActiveArr[i][0], currActiveArr[i][1]
					g.Tetromino.onScreenBank[g.Tetromino.activeN].Arr[i][0] -= g.Tetromino.settings.tIncrement
					//g.sample[i][1] += INCREMENT
					fmt.Printf("moved sprites from %v,%v to %v,%v\n", fX, fY, g.Tetromino.onScreenBank[g.Tetromino.activeN].Arr[i][0], g.Tetromino.onScreenBank[g.Tetromino.activeN].Arr[i][1])
				}
				g.Unlock()
			}
		} else {
			// at the start of the game
			g.Tetromino.Add()
		}
		if g.Tetromino.creationCounter >= 300 {
			g.Lock()
			g.counter = 0
			g.Tetromino.Add()
			g.Unlock()
		}
		g.counter++
	case HELP:
		if (x >= 404 && y >= 47) && (x <= 585 && y <= 132) {
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				log.Printf("restart pressed: x=%d, y=%d\n", x, y)
				g.Lock()
				g.P = WELCOME
				g.Unlock()
			}
		}
	case QUIT:
	case DEBUG:
		if g.ExecMode == EXEC_DEBUG {
			g.debugPos("left")
			if (x >= 404 && y >= 47) && (x <= 585 && y <= 132) {
				if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
					log.Printf("back pressed: x=%d, y=%d\n", x, y)
					g.Lock()
					g.P = WELCOME
					g.Unlock()
				}
			}
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	prepDir(g.CFunc, g.FeatureFlags)
	BGColor := color.Black
	var BorderColor color.Color = color.RGBA{R: 0, G: 0, B: 255, A: 255}
	screen.Fill(BGColor)
	u := utils2.NewScreenUtils(context.Background(), utils2.SCREEN_W, utils2.SCREEN_H, utils2.ScreenOptions{BgColor: BGColor})
	u.WithBorder(utils2.LIGHT, BorderColor).Display(screen)

	switch g.P {
	case WELCOME:
		text.Draw(screen, g.LeadText.Text, DefaultFont, 200, 80, color.White)
		u.WithContainerHollow(utils2.LIGHT, 120, 150, 120, 70, utils2.NewContainerOptions(utils2.UseWeightPack(2, 7, 10), utils2.UseText("Play", DefaultFont), utils2.UseBorderColor(BorderColor), utils2.UseFillColor(BGColor))).Display(screen)
		u.WithContainerHollow(utils2.LIGHT, 400, 150, 120, 70, utils2.NewContainerOptions(utils2.UseWeightPack(2, 7, 10), utils2.UseText("Help", DefaultFont), utils2.UseBorderColor(BorderColor), utils2.UseFillColor(BGColor))).Display(screen)
		if g.ExecMode == EXEC_DEBUG {
			u.WithContainerHollow(utils2.LIGHT, 400, 250, 120, 70, utils2.NewContainerOptions(utils2.UseWeightPack(2, 7, 10), utils2.UseText("Debug", DefaultFont), utils2.UseBorderColor(COLOR_RED), utils2.UseFillColor(BGColor))).Display(screen)
		}
		text.Draw(screen, g.genWelcomeText(), DefaultFont, 120, 320, color.White)
	case START:
		fmt.Println("DISPLAY_START")
		ctx := context.Background()
		//g.Tetromino.Add()
		u.WithContainerHollow(utils2.LIGHT, float32(utils2.WALL_X0), float32(utils2.WALL_Y0), float32(utils2.WALL_WIDTH), float32(utils2.WALL_HEIGHT), utils2.NewContainerOptions(utils2.UseWeightPack(2, 7, 10), utils2.UseText("", DefaultFont), utils2.UseBorderColor(BorderColor), utils2.UseFillColor(BGColor))).Display(screen)
		u.WithContainerHollow(utils2.LIGHT, 405, 45, 179, 85, utils2.NewContainerOptions(utils2.UseWeightPack(2, 7, 10), utils2.UseText("Restart", DefaultFont), utils2.UseBorderColor(BorderColor), utils2.UseFillColor(BGColor))).Display(screen)
		u.WithContainerHollow(utils2.LIGHT, 405, 446, 168, 85, utils2.NewContainerOptions(utils2.UseWeightPack(2, 7, 10), utils2.UseText("Quit", DefaultFont), utils2.UseBorderColor(BorderColor), utils2.UseFillColor(BGColor))).Display(screen)
		//utils2.ContainerFill(ctx, float32(g.sample[0][0]), float32(g.sample[0][1]), 30, 30, WALL_WIDTH, WALL_HEIGHT, utils2.NewContainerOptions(utils2.UseFillColor(utils2.BLUE))).Display(screen)
		//utils2.ContainerFill(ctx, float32(g.sample[1][0]), float32(g.sample[1][1]), 30, 30, WALL_WIDTH, WALL_HEIGHT, utils2.NewContainerOptions(utils2.UseFillColor(utils2.BLUE), utils2.UseBorderColor(utils2.WHITE))).Display(screen)
		//utils2.ContainerFill(ctx, float32(g.sample[2][0]), float32(g.sample[2][1]), 30, 30, WALL_WIDTH, WALL_HEIGHT, utils2.NewContainerOptions(utils2.UseFillColor(utils2.BLUE), utils2.UseBorderColor(utils2.WHITE))).Display(screen)
		//utils2.MultiContainerFill(ctx, 30, 30, WALL_WIDTH, WALL_HEIGHT, [][]int{{305, 316}, {275, 346}, {305, 346}}, utils2.NewContainerOptions(utils2.UseWeightPack(2, 7, 10), utils2.UseBorderColor(utils2.WHITE), utils2.UseFillColor(utils2.BLUE)))(screen)
		//u.WithWall(float32(SPRITE_W), float32(SPRITE_H), float32(WALL_X0), float32(WALL_Y0), color.RGBA{R: 255, G: 0, B: 0, A: 255}).Display(screen)
		spawnWall(ctx)(screen)
		if len(g.FeatureFlags) > 0 && g.FeatureFlags[config.FLAG_PLAYGAME] == config.FlagEnabled {
			// super buggy, hehe. when I find time, I'll continue.
			g.Tetromino.Display(ctx, screen)
		}
	case HELP:
		u.WithContainerHollow(utils2.LIGHT, 120, 120, 300, 300, utils2.NewContainerOptions(utils2.UseWeightPack(2, 7, 10), utils2.UseText(HELP_TEXT, DefaultFont), utils2.UseBorderColor(color.Transparent), utils2.UseFillColor(BGColor))).Display(screen)
		u.WithContainerHollow(utils2.LIGHT, 405, 45, 179, 85, utils2.NewContainerOptions(utils2.UseWeightPack(2, 7, 10), utils2.UseText("Restart", DefaultFont), utils2.UseBorderColor(BorderColor), utils2.UseFillColor(BGColor))).Display(screen)
	case QUIT:
		u.WithContainerHollow(utils2.LIGHT, 171, 240, 300, 175, utils2.NewContainerOptions(utils2.UseWeightPack(2, 7, 10), utils2.UseText("Quitting...", DefaultFont), utils2.UseBorderColor(color.Transparent), utils2.UseFillColor(BGColor))).Display(screen)
		time.Sleep(5)
		os.Exit(0)
	case DEBUG:
		u.WithBorder(utils2.LIGHT, COLOR_RED).Display(screen)
		u.WithContainerHollow(utils2.LIGHT, 405, 45, 179, 85, utils2.NewContainerOptions(utils2.UseWeightPack(2, 7, 10), utils2.UseText("Back", DefaultFont), utils2.UseBorderColor(BorderColor), utils2.UseFillColor(BGColor))).Display(screen)
		text.Draw(screen, "Sample", DefaultFont, 120, 300, color.White)
		//u.WithContainerHollow(utils2.LIGHT, 40, 45, 250, 300, utils2.NewContainerOptions(utils2.UseWeightPack(2, 7, 10), utils2.UseText("Back", DefaultFont), utils2.UseBorderColor(BorderColor), utils2.UseFillColor(BGColor))).Display(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return utils2.SCREEN_W, utils2.SCREEN_W
}

func (g *Game) debugPos(s string) {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		log.Println("debug pos")
	}
	if g.counter == 2 {
		log.Println("counter resetting to zero")
		g.counter = 0
	}
	x, y := ebiten.CursorPosition()
	switch s {
	case "left":
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			log.Println("inner loop left")
			g.Lock()
			g.counter++
			g.Unlock()
			log.Println("counter at", g.counter)
			if g.counter == 2 {
				log.Println("counter equals 2, computing diff")
				if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
					log.Printf("Left: mouse pressed at: fx=%d, fy=%d, x=%d, y=%d, width=%d, height=%d", g.x, g.y, x, y, x-g.x, y-g.y)
				}
			}
			if g.counter == 1 {
				log.Println("counter equals 1, recording state")
				if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
					g.x, g.y = x, y
				}
			}
		}
	case "right":
		log.Println("inner loop right")
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
			log.Printf("Right: mouse pressed at: x=%d, y=%d\n", x, y)
		}
	}
}

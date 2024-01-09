package tetra

import (
	"context"
	"github.com/dark-enstein/crise/internal/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
	"strconv"
	"strings"
)

type Mini int

const (
	I = iota + 1
	O
	T
	S
	Z
	J
	L
	WALL
)

var (
	SPRITE_PIXELTOWHOLE_RATIO = 30
)

// TetraCoordinates 5x5
var TetraCoordinates = map[Mini]string{
	I:    "0,0;1,0;2,0;3,0",
	O:    "0,0;1,0;0,1;1,1",
	T:    "0,0;1,0;2,0;1,1",
	S:    "1,0;2,0;0,1;1,1",
	Z:    "0,0;1,0;1,1;2,1",
	J:    "1,0;1,1;1,2;0,2",
	L:    "0,0;0,1;0,2;1,2",
	WALL: "0,0;1,0;2,0;3,0;4,0;5,0;6,0;7,0;8,0;9,0;10,0;11,0;12,0;13,0;14,0;15,0;16,0;17,0;18,0;19,0;20,0;21,0;22,0;23,0;24,0;25,0;26,0;27,0;28,0;29,0;29,1;29,2;29,3;29,4;29,5;29,6;29,7;29,8;29,9;29,10;29,11;29,12;29,13;29,14;29,15;29,16;29,17;29,18;29,19;29,20;29,21;29,22;29,23;29,24;29,25;29,26;29,27;29,28;29,29;29,30;29,31;29,32;29,33;29,34;29,35;29,36;29,37;29,38;29,39;29,40;29,41;29,42;29,43;29,44;29,45;29,46;29,47;29,48;29,49;28,49;27,49;26,49;25,49;24,49;23,49;22,49;21,49;20,49;19,49;18,49;17,49;16,49;15,49;14,49;13,49;12,49;11,49;10,49;9,49;8,49;7,49;6,49;5,49;4,49;3,49;2,49;1,49;0,49;0,48;0,47;0,46;0,45;0,44;0,43;0,42;0,41;0,40;0,39;0,38;0,37;0,36;0,35;0,34;0,33;0,32;0,31;0,30;0,29;0,28;0,27;0,26;0,25;0,24;0,23;0,22;0,21;0,20;0,19;0,18;0,17;0,16;0,15;0,14;0,13;0,12;0,11;0,10;0,9;0,8;0,7;0,6;0,5;0,4;0,3;0,2;0,1",
}

type Tetromino struct {
	Arr                       [][]int
	wholeWidth, wholeHeight   int
	spriteWidth, spriteHeight float32
}

// NewTetromino converts a sequence of object coordinates delimited by colon (;) into the corresponding slice of byte slices
func NewTetromino(coordinateMap string, wholeWidth, wholeHeight int) *Tetromino {
	var t = &Tetromino{wholeWidth: wholeWidth, wholeHeight: wholeHeight}
	splices := strings.Split(coordinateMap, ";")
	var buf = [][]int{}
	for i := 0; i < len(splices); i++ {
		buf = append(buf, parseCoord(splices[i]))
	}
	t.Arr = buf[:len(splices)]
	log.Println(buf)
	return t
}

// parseCoord parses a coordinate in notation (x, y), to its corresponding 2 byte array
func parseCoord(s string) []int {
	splices := strings.Split(s, ",")
	if len(splices) != 2 {
		log.Printf("coordinates %s parsed is in the wrong format\n", s)
		return nil
	}
	first, _ := strconv.Atoi(splices[0])
	second, _ := strconv.Atoi(splices[1])
	return []int{first, second}
}

// Coordinates converts an object from Tetranomics notation, to its coordinates form
func (t *Tetromino) Coordinates() string {
	coord := ""
	for i := 0; i < len(t.Arr); i++ {
		for j := 0; j < len(t.Arr[i]); j++ {
			pair := t.Arr[i]
			if j == 0 {
				coord += strconv.Itoa(pair[j]) + ","
			} else {
				coord += strconv.Itoa(pair[j])
			}
		}
		if i < len(t.Arr)-1 {
			coord += ";"
		}
	}
	return coord
}

// Build sets up a Tetromino ready for display. By calling the Display function against the return value, the Tetromino is rendered on screen
func (t *Tetromino) Build(ctx context.Context) func(screen *ebiten.Image) {
	t.spriteWidth = float32(t.wholeWidth / SPRITE_PIXELTOWHOLE_RATIO)
	t.spriteWidth = float32(t.wholeHeight / SPRITE_PIXELTOWHOLE_RATIO)
	return MultiContainerFill(ctx, t.spriteWidth, t.spriteHeight, t.wholeWidth, t.wholeHeight, t.Arr, utils.NewContainerOptions(utils.UseFillColor(utils.BLUE), utils.UseBorderColor(utils.WHITE)))
}

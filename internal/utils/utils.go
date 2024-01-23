package utils

import (
	"context"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"log"
	"os"
)

var (
	GREEN        color.Color = &color.RGBA{R: 0, G: 128, B: 70, A: 255}
	RED          color.Color = &color.RGBA{R: 255, G: 0, B: 0, A: 255}
	BLUE         color.Color = &color.RGBA{R: 0, G: 0, B: 255, A: 255}
	WHITE                    = &color.White
	BLACK                    = &color.Black
	BORDER_WIDTH float32     = 2
)

type BORDER_STRENGTH int
type TEXT_ALIGNMENT int

const (
	ALIGN_CENTER TEXT_ALIGNMENT = iota
	ALIGN_LEFT
	ALIGN_RIGHT
)

const (
	LIGHT BORDER_STRENGTH = iota
	MEDIUM
	BOLD
)

type Screen struct {
	ctx              context.Context
	WHOLE_W, WHOLE_H int
	BG_COLOR         color.Color
	Border           BorderUtils
}

// Utils defined a methodset that the calling Screen util must implement
type Utils interface {
	// Display renders the object constructed to the destination
	Display(screen *ebiten.Image)
}

// BorderUtils defines the structure of a Border
type BorderUtils ContainerUtils

func (b BorderUtils) Display(screen *ebiten.Image) {
	b.display(screen)
}

func NewScreenUtils(ctx context.Context, layoutW, layoutH int, opts ScreenOptions) *Screen {
	return &Screen{ctx: ctx, WHOLE_H: layoutH, WHOLE_W: layoutW, BG_COLOR: opts.BgColor}
}

type ScreenOptions struct {
	BgColor color.Color
}

func (s *Screen) WithBorder(b BORDER_STRENGTH, borderColor color.Color, opts ...BorderOptions) *BorderUtils {
	_, cancel := context.WithCancel(s.ctx)
	defer cancel()
	if s.WHOLE_W == 0 {
		log.Println("screen width not passed in")
		cancel()
	} else if s.WHOLE_H == 0 {
		log.Println("screen height not passed in")
		cancel()
	}
	var borderStartX, borderStartY float32 = 20, 20
	sb := BorderUtils(*s.WithContainerHollow(b, 20, 20, float32(s.WHOLE_H)-(2*borderStartY), float32(s.WHOLE_W)-(2*borderStartX), NewContainerOptions(UseMarginBorder(borderColor))))
	return &sb
}

// ContainerUtils holds the core prestate of a render ready to display. The state of the render is displayed by calling its Display() function, and passing in the ebiten.Image as an argument.
// It satisfies the Utils interface
type ContainerUtils struct {
	ctx                             context.Context
	BORDER_W, BORDER_H              int
	bgColor, borderColor, fillColor color.Color
	OFFSET                          float32
	display                         func(screen *ebiten.Image)
}

func (b ContainerUtils) Display(screen *ebiten.Image) {
	b.display(screen)
}

func (s *Screen) WithContainerHollow(b BORDER_STRENGTH, x, y, width, height float32, opts *ContainerOptions) *ContainerUtils {
	_, cancel := context.WithCancel(s.ctx)
	defer cancel()
	if s.WHOLE_W == 0 {
		log.Println("screen width not passed in")
		cancel()
	} else if s.WHOLE_H == 0 {
		log.Println("screen height not passed in")
		cancel()
	}
	var offset float32

	BorderColor := opts.BorderColor
	switch b {
	case LIGHT:
		offset = opts.OffsetPack.Light
	case MEDIUM:
		offset = opts.OffsetPack.Medium
	case BOLD:
		offset = opts.OffsetPack.Bold
	default:
		log.Println("border strength undefined. doing nothing.")
		return nil
	}
	con := ContainerUtils{OFFSET: offset, bgColor: s.BG_COLOR, borderColor: BorderColor, display: func(screen *ebiten.Image) {
		var length float32 = 0
		if opts.Text != TEXT_UNSET {
			length = float32(len(opts.Text))
			var boxHeightAdjustment float32 = 22
			var charWidth float32 = 15
			var charHeight float32 = 38
			var boxWPadding = (width - (length * charWidth)) / 2
			var boxHPadding = (height - charHeight) / 2
			text.Draw(screen, opts.Text, opts.FontFace, int(x+boxWPadding), int(y+boxHeightAdjustment+boxHPadding), color.White)
		}
		vector.StrokeRect(screen, x, y, width, height, offset, BorderColor, false)
	}}
	return &con
}

func (s *Screen) WithContainerFill(x, y, width, height float32, opts *ContainerOptions) *ContainerUtils {
	return ContainerFill(s.ctx, x, y, width, height, s.WHOLE_W, s.WHOLE_H, opts)
}

// ContainerFill draws a bounded box. x, y is the point on the window where the bounded box begins (top-left) vertex. width and height are the dimensions of thw box to be rendered, and wholeH and wholeW are the dimensions for the window itself.
// opts an ContainerOptions holds the customization options for the bounded box to be rendered
func ContainerFill(ctx context.Context, x, y, width, height float32, wholeW, wholeH int, opts *ContainerOptions) *ContainerUtils {
	_, cancel := context.WithCancel(ctx)
	defer cancel()
	if wholeW == 0 {
		log.Println("screen width not passed in")
		cancel()
	} else if wholeH == 0 {
		log.Println("screen height not passed in")
		cancel()
	}

	FillColor := opts.FillColor
	con := ContainerUtils{fillColor: FillColor, display: func(screen *ebiten.Image) {
		var length float32 = 0
		if opts.Text != TEXT_UNSET {
			length = float32(len(opts.Text))
			var boxHeightAdjustment float32 = 22
			var charWidth float32 = 15
			var charHeight float32 = 38
			var boxWPadding = (width - (length * charWidth)) / 2
			var boxHPadding = (height - charHeight) / 2
			text.Draw(screen, opts.Text, opts.FontFace, int(x+boxWPadding), int(y+boxHeightAdjustment+boxHPadding), color.White)
		}
		vector.DrawFilledRect(screen, x, y, width+10, height+10, FillColor, false)
		if opts.BorderColor != nil {
			vector.StrokeRect(screen, x, y, width, height, BORDER_WIDTH, opts.BorderColor, false)
		}
	}}
	return &con
}

type WALL_FACE int

const (
	WALL_LEFT = iota
	WALL_RIGHT
	WALL_UP
	WALL_DOWN
)

func (s *Screen) WithWall(spriteWidth, spriteHeight, wallX, wallY float32, col color.Color) *ContainerUtils {
	var normX = func(i float32) float32 {
		return i + wallX
	}
	var normY = func(i float32) float32 {
		return i + wallY
	}
	wallConfig := []struct {
		face           WALL_FACE
		member         float32
		startX, startY float32
	}{
		{
			WALL_LEFT,
			30,
			0,
			0,
		}, {
			WALL_RIGHT,
			30,
			30,
			0,
		}, {
			WALL_UP,
			29,
			1,
			0,
		}, {
			WALL_DOWN,
			29,
			1,
			29,
		},
	}

	realStartXLeft, realStartYLeft, x, y := normX(wallConfig[0].startX*spriteWidth), normY(wallConfig[0].startY*spriteHeight), wallConfig[0].startX, wallConfig[0].startY
	memLeft := wallConfig[0].member
	//startX, startY := realStartX+(mem*spriteWidth), realStartY+(mem*spriteWidth)
	if int(realStartXLeft) > s.WHOLE_W || int(realStartYLeft) > s.WHOLE_H {
		log.Println("wall start metrics is overbound")
		os.Exit(1)
	}
	log.Printf("start: (%.1f,%.1f); realStart: (%.1f,%.1f); members: %.1f\n", x, y, realStartXLeft, realStartYLeft, memLeft)
	realStartXRight, realStartYRight, x, y := normX(wallConfig[1].startX*spriteWidth), normY(wallConfig[1].startY*spriteHeight), wallConfig[1].startX, wallConfig[1].startY
	memRight := wallConfig[1].member
	//startX, startY := realStartX+(mem*spriteWidth), realStartY+(mem*spriteWidth)
	if int(realStartXRight) > s.WHOLE_W || int(realStartYRight) > s.WHOLE_H {
		log.Println("wall start metrics is overbound")
		os.Exit(1)
	}
	log.Printf("start: (%.1f,%.1f); realStart: (%.1f,%.1f); members: %.1f\n", x, y, realStartXRight, realStartYRight, memRight)
	realStartXUp, realStartYUp, x, y := normX(wallConfig[2].startX*spriteWidth), normY(wallConfig[2].startY*spriteHeight), wallConfig[2].startX, wallConfig[2].startY
	memUp := wallConfig[2].member
	//startX, startY := realStartX+(mem*spriteWidth), realStartY+(mem*spriteWidth)
	if int(realStartXUp) > s.WHOLE_W || int(realStartYUp) > s.WHOLE_H {
		log.Println("wall start metrics is overbound")
		os.Exit(1)
	}
	log.Printf("start: (%.1f,%.1f); realStart: (%.1f,%.1f); members: %.1f\n", x, y, realStartXUp, realStartYUp, memLeft)
	realStartXDown, realStartYDown, x, y := normX(wallConfig[3].startX*spriteWidth), normY(wallConfig[3].startY*spriteHeight), wallConfig[3].startX, wallConfig[3].startY
	memDown := wallConfig[3].member
	//startX, startY := realStartX+(mem*spriteWidth), realStartY+(mem*spriteWidth)
	if int(realStartXDown) > s.WHOLE_W || int(realStartYDown) > s.WHOLE_H {
		log.Println("wall start metrics is overbound")
		os.Exit(1)
	}
	log.Printf("start: (%.1f,%.1f); realStart: (%.1f,%.1f); members: %.1f\n", x, y, realStartXDown, realStartYDown, memRight)
	return &ContainerUtils{display: func(screen *ebiten.Image) {
		for iLeft := 0; iLeft <= int(memLeft); iLeft++ {
			s.WithContainerFill(realStartXLeft, realStartYLeft, spriteWidth, spriteHeight, NewContainerOptions(UseFillColor(col))).Display(screen)
			//realStartX += spriteWidth
			realStartYLeft += spriteHeight
		}
		for iRight := 0; iRight <= int(memRight); iRight++ {
			s.WithContainerFill(realStartXRight, realStartYRight, spriteWidth, spriteHeight, NewContainerOptions(UseFillColor(col))).Display(screen)
			//startX += spriteWidth
			realStartYRight += spriteHeight
		}
		for iUp := 0; iUp <= int(memUp); iUp++ {
			s.WithContainerFill(realStartXUp, realStartYUp, spriteWidth, spriteHeight, NewContainerOptions(UseFillColor(col))).Display(screen)
			realStartXUp += spriteWidth
			//realStartXUp += spriteHeight
		}
		for iDown := 0; iDown <= int(memDown); iDown++ {
			s.WithContainerFill(realStartXDown, realStartYDown, spriteWidth, spriteHeight, NewContainerOptions(UseFillColor(col))).Display(screen)
			realStartXDown += spriteWidth
			//realStartYRight += spriteHeight
		}
	}}
	return nil
}

package utils

import (
	"context"
	"github.com/dark-enstein/crise/config"
)

type PlayArea struct {
	TopLeftX, TopLeftY int
	Width, Height      int
	Settings           *config.Gameplay
}

func NewPlayArea(tLX, tLY, w, h int, cfg *config.Gameplay) *PlayArea {
	return &PlayArea{
		TopLeftX: tLX,
		TopLeftY: tLY,
		Width:    w,
		Height:   h,
		Settings: cfg,
	}
}

func (p *PlayArea) PlaceIn(ctx context.Context, x, y int) *ContainerUtils {
	realX, realY := (x*SPRITE_WIDTH)+p.Width, (y*SPRITE_HEIGHT)+p.Height
	return ContainerFill(ctx, float32(realX), float32(realY), SPRITE_WIDTH, SPRITE_HEIGHT, p.Width, p.Height, NewContainerOptions(UseWeightPack(2, 7, 10), UseBorderColor(WHITE), UseFillColor(BLUE)))
}

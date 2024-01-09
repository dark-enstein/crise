package game

import (
	"context"
	"github.com/dark-enstein/crise/internal/tetra"
	"github.com/dark-enstein/crise/internal/utils"
	"github.com/hajimehoshi/ebiten/v2"
)

func spawnWall(ctx context.Context) func(screen *ebiten.Image) {
	return func(screen *ebiten.Image) {
		tetra.NewTetromino(tetra.TetraCoordinates[tetra.WALL], utils.WALL_X0, utils.WALL_Y0).Build(ctx)(screen)
	}
}

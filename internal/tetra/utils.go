package tetra

import (
	"context"
	"github.com/dark-enstein/crise/internal/utils"
	"github.com/hajimehoshi/ebiten/v2"
)

// // MultiContainerFill used to create multicontainer objects
func MultiContainerFill(ctx context.Context, spriteWidth, spriteHeight float32, wholeWidth, wholeHeight int, arr [][]int, opts *utils.ContainerOptions) func(screen *ebiten.Image) {
	var pa = utils.NewPlayArea(utils.WALL_X0, utils.WALL_Y0, wholeWidth, wholeHeight, nil)
	//fmt.Println("Entered Multi container fill")
	bulkDisplay := []*utils.ContainerUtils{}
	for i := 0; i < len(arr); i++ {
		//fmt.Println("appending container fill")
		//bulkDisplay = append(bulkDisplay, ContainerFill(ctx, float32(arr[i][0]), float32(arr[i][1]), spriteWidth, spriteHeight, wholeWidth, wholeHeight, opts))
		bulkDisplay = append(bulkDisplay, pa.PlaceIn(ctx, arr[i][0], arr[i][1]))
	}

	return func(screen *ebiten.Image) {
		//fmt.Println("about running each display")
		for i := 0; i < len(bulkDisplay); i++ {
			//fmt.Println("running each display")
			bulkDisplay[i].Display(screen)
		}
	}
}

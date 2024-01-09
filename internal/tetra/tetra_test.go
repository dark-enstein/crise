package tetra

import (
	"reflect"
	"testing"
)

func TestNewTetranomics(t *testing.T) {
	var testTetra = []struct {
		coord string
		array [][]byte
	}{
		{
			"1,0;2,0;3,0;4,0;5,0",
			[][]byte{{1, 0}, {2, 0}, {3, 0}, {4, 0}, {5, 0}},
		},
	}
	for i := 0; i < len(testTetra); i++ {
		tetra := NewTetromino(testTetra[i].coord, 640, 640)
		if !reflect.DeepEqual(tetra.Arr, testTetra[i].array) {
			t.Errorf("Got %v, 1st: %v\n", tetra.Arr, tetra.Arr[0])
		}
	}
}

func TestTetrominoes_Coordinates(t *testing.T) {
	var testTetra = []struct {
		coord string
		array [][]byte
	}{
		{
			"1,0;2,0;3,0;4,0;5,0",
			[][]byte{{1, 0}, {2, 0}, {3, 0}, {4, 0}, {5, 0}},
		},
	}
	for i := 0; i < len(testTetra); i++ {
		coordinates := NewTetromino(testTetra[i].coord, 640, 640).Coordinates()
		if !(coordinates == testTetra[i].coord) {
			t.Errorf("Expected %s, Got %s\n", testTetra[i].coord, coordinates)
		}
	}
}

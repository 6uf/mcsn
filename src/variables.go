package src

import (
	"image"
	"image/color"

	"github.com/6uf/apiGO"
)

type SkinUrls struct {
	Url     interface{} `json:"url"`
	Varient interface{} `json:"variant"`
}

type Pixel struct {
	Point image.Point
	Color color.Color
}

var (
	Bearers   apiGO.MCbearers
	Acc       apiGO.Config
	Proxys    apiGO.Proxys
	images    []image.Image
	thirdRow  [][]int = [][]int{{64, 16, 72, 24}, {56, 16, 64, 24}, {48, 16, 56, 24}, {40, 16, 48, 24}, {32, 16, 40, 24}, {24, 16, 32, 24}, {16, 16, 24, 24}, {8, 16, 16, 24}, {0, 16, 8, 24}}
	secondRow [][]int = [][]int{{64, 8, 72, 16}, {56, 8, 64, 16}, {48, 8, 56, 16}, {40, 8, 48, 16}, {32, 8, 40, 16}, {24, 8, 32, 16}, {16, 8, 24, 16}, {8, 8, 16, 16}, {0, 8, 8, 16}}
	firstRow  [][]int = [][]int{{64, 0, 72, 8}, {56, 0, 64, 8}, {48, 0, 56, 8}, {40, 0, 48, 8}, {32, 0, 40, 8}, {24, 0, 32, 8}, {16, 0, 24, 8}, {8, 0, 16, 8}, {0, 0, 8, 8}}
)

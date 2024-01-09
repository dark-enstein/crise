package utils

import (
	"golang.org/x/image/font"
	"image/color"
)

var (
	TEXT_UNSET = "UNSET"
)

type ContainerOptions struct {
	BorderColor color.Color
	FillColor   color.Color
	Text        string
	Alignment   *TEXT_ALIGNMENT
	FontFace    font.Face
	OffsetPack  *OffsetPack
}

type OffsetPack struct {
	Light, Medium, Bold float32
}

type ContainerOption func(*ContainerOptions)

func NewContainerOptions(cops ...ContainerOption) *ContainerOptions {
	var defBorCol color.Color = color.White
	defaultOpts := &ContainerOptions{
		OffsetPack: &OffsetPack{
			Light:  5,
			Medium: 7,
			Bold:   10,
		},
		Text:        TEXT_UNSET,
		BorderColor: defBorCol,
	}

	for _, cop := range cops {
		cop(defaultOpts)
	}

	return defaultOpts
}

func UseWeightPack(light, medium, bold float32) ContainerOption {
	return func(c *ContainerOptions) {
		c.OffsetPack = &OffsetPack{
			Light:  light,
			Medium: medium,
			Bold:   bold,
		}
	}
}

func UseText(t string, ft font.Face) ContainerOption {
	return UseTextWithAlignment(t, ALIGN_CENTER, ft)
}

func UseTextWithAlignment(t string, align TEXT_ALIGNMENT, ft font.Face) ContainerOption {
	return func(c *ContainerOptions) {
		c.Text = t
		c.Alignment = &align
		c.FontFace = ft
	}
}

func UseBorderColor(col color.Color) ContainerOption {
	return func(c *ContainerOptions) {
		c.BorderColor = col
	}
}

func UseFillColor(col color.Color) ContainerOption {
	return func(c *ContainerOptions) {
		c.FillColor = col
	}
}

func UseMarginBorder(borderColor color.Color) ContainerOption {
	return func(c *ContainerOptions) {
		c.OffsetPack = &OffsetPack{
			Light:  2,
			Medium: 7,
			Bold:   10,
		}
		c.BorderColor = borderColor
	}
}

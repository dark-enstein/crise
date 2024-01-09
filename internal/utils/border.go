package utils

import "image/color"

type BorderOptions ContainerOptions

type BorderOption func(options *BorderOptions)

func NewBorderOptions(cops ...BorderOption) *BorderOptions {
	var defBorCol color.Color = color.White
	defaultOpts := &BorderOptions{
		OffsetPack: &OffsetPack{
			Light:  5,
			Medium: 7,
			Bold:   10,
		},
		Text:        "UNSET",
		BorderColor: defBorCol,
	}

	for _, cop := range cops {
		cop(defaultOpts)
	}

	return defaultOpts
}

//func UseWeightPack(light, medium, bold float32) BorderOption {
//	return func(c *BorderOptions) {
//		c.OffsetPack = &OffsetPack{
//			Light:  light,
//			Medium: medium,
//			Bold:   bold,
//		}
//	}
//}
//
//func UseText(t string) BorderOption {
//	return UseTextWithAlignment(t, ALIGN_CENTER)
//}
//
//func UseTextWithAlignment(t string, align TEXT_ALIGNMENT) BorderOption {
//	return func(c *BorderOptions) {
//		c.Text = t
//		c.Alignment = &align
//	}
//}
//
//func UseBorderColor(col color.Color) BorderOption {
//	return func(c *BorderOptions) {
//		c.BorderColor = &col
//	}
//}
//
//func UseFillColor(col color.Color) BorderOption {
//	return func(c *BorderOptions) {
//		c.FillColor = &col
//	}
//}

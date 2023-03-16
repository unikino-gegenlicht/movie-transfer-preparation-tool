package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"image/color"
	"movie-transfer-preparation-tool/resources"
)

type CustomTheme struct{}

// change the default colors of texts and other widgets in the dark mode
func (c CustomTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if variant == theme.VariantDark {
		// return the color from the cd (#ffdd00)
		return color.NRGBA{
			R: 211,
			G: 255,
			B: 0,
			A: 255,
		}
	} else if variant == theme.VariantLight {
		return color.Black
	}
	return theme.DefaultTheme().Color(name, variant)
}

// change the built in fonts to the fonts used directly in the project
func (c CustomTheme) Font(style fyne.TextStyle) fyne.Resource {
	if style.Monospace {
		if style.Bold {
			return resources.NotoSansMonoBold
		}
		return resources.NotoSansMonoRegular
	}
	if style.Bold {
		if style.Italic {
			return resources.OpenSansBoldItalic
		}
		return resources.OpenSansBold
	}
	if style.Italic {
		return resources.OpenSansItalic
	}
	return resources.OpenSans
}

// keep the icon lookup as is
func (c CustomTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

// keep the size lookup as is
func (c CustomTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

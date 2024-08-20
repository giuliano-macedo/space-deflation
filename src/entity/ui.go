package entity

import (
	"image/color"

	"github.com/abelroes/gmtk2024/src/constants"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Ui struct {
	ShowRestartText bool
	font            *text.GoTextFaceSource
}

func NewUi(font *text.GoTextFaceSource) *Ui {
	return &Ui{
		font: font,
	}
}

func (ui *Ui) Draw(screen *ebiten.Image) {
	if ui.ShowRestartText {
		op := &text.DrawOptions{
			LayoutOptions: text.LayoutOptions{PrimaryAlign: text.AlignCenter},
		}
		op.GeoM.Translate(constants.Width/2, (constants.Height/2)-50)

		op.ColorScale.ScaleWithColor(color.White)
		text.Draw(screen, "You died! press Enter or R to try again", &text.GoTextFace{
			Source: ui.font,
			Size:   15,
		}, op)
	}
}

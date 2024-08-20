package entity

import (
	"math/rand/v2"

	"github.com/abelroes/gmtk2024/src/constants"
	"github.com/hajimehoshi/ebiten/v2"
)

type Background struct {
	bgWidth        float64
	bgHeight       float64
	backgroundList []*ebiten.Image
	currentBG      *ebiten.Image
}

func NewBackground(backgroundList []*ebiten.Image) *Background {
	bg := &Background{
		bgWidth:        1280,
		bgHeight:       960,
		backgroundList: backgroundList,
	}
	bg.ChangeBG()
	return bg
}

func (bg *Background) ChangeBG() {
	bg.currentBG = bg.backgroundList[rand.IntN(len(bg.backgroundList))]
}

func (bg *Background) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.ColorScale.Scale(0.5, 0.5, 0.5, 0.8)
	op.GeoM.Scale(constants.Width/bg.bgWidth, constants.Height/bg.bgHeight)

	op.Filter = ebiten.FilterLinear
	screen.DrawImage(bg.currentBG, op)
}

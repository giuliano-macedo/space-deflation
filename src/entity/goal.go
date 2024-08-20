package entity

import (
	"image/color"

	"github.com/abelroes/gmtk2024/src/collision"
	"github.com/abelroes/gmtk2024/src/settings"
	"github.com/abelroes/gmtk2024/src/vector"
	"github.com/hajimehoshi/ebiten/v2"
	ebivector "github.com/hajimehoshi/ebiten/v2/vector"
)

type Goal struct {
	Pos           vector.Vector2
	Radius        float64
	Collider      collision.CollisionRect
	img           *ebiten.Image
	angle         float64
	debugSettings *settings.SettingsDebug
}

const (
	rotationSpeed = .01
)

func NewGoal(img *ebiten.Image, debugSettings *settings.SettingsDebug, radius float64) Goal {
	return Goal{
		Radius:        radius,
		img:           img,
		debugSettings: debugSettings,
	}
}

func (goal *Goal) SetPos(pos vector.Vector2) {
	goal.angle = 0
	goal.Pos = pos
	colliderFactor := .75
	r := goal.Radius * colliderFactor
	goal.Collider = collision.CollisionRect{
		Pos: vector.New(goal.Pos.X-r, goal.Pos.Y-r),
		W:   2 * r,
		H:   2 * r,
	}
}

func (goal *Goal) Update() {
	goal.angle += rotationSpeed
}

func (goal *Goal) Draw(screen *ebiten.Image) {

	bounds := goal.img.Bounds()
	w, h := float64(bounds.Dx()), float64(bounds.Dy())

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-w/2, -h/2)
	op.GeoM.Scale((goal.Radius*2)/w, (goal.Radius*2)/h)
	op.GeoM.Rotate(goal.angle)
	op.GeoM.Translate(goal.Pos.X, goal.Pos.Y)
	op.Filter = ebiten.FilterLinear
	screen.DrawImage(goal.img, op)

	if goal.debugSettings.GoalHitbox {
		ebivector.DrawFilledRect(screen, float32(goal.Collider.Pos.X), float32(goal.Collider.Pos.Y), float32(goal.Collider.W), float32(goal.Collider.H), color.RGBA{R: 255}, true)
	}
}

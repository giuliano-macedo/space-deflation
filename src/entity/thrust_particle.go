package entity

import (
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/abelroes/gmtk2024/src/vector"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	particleThrustTtl = time.Duration(1 * time.Second)
)

var (
	Dot = ebiten.NewImage(1, 1)
)

func init() {
	Dot.Fill(color.RGBA{
		R: 255,
		G: 255,
		B: 0,
		A: 255,
	})
}

type ThrustParticle struct {
	pos   vector.Vector2
	vel   vector.Vector2
	start time.Time
}

func (particle *ThrustParticle) Update() {
	particle.pos.Add(particle.vel)
}

func (particle *ThrustParticle) Draw(screen *ebiten.Image, now time.Time) {
	op := &ebiten.DrawImageOptions{}
	timeRemainingRatio := float32(now.Sub(particle.start)) / float32(particleThrustTtl)
	gradient := (1 - timeRemainingRatio)
	colorGradient := gradient * gradient * gradient
	redGradient := gradient
	op.ColorScale.Scale(redGradient, colorGradient, colorGradient, 1)

	op.GeoM.Translate(particle.pos.X, particle.pos.Y)
	screen.DrawImage(Dot, op)
}

func SpawnParticle(w, h, rot, x, y float64) ThrustParticle {
	s, c := math.Sincos(rot) // NOTE: this is computed in player.Draw, maybe pass to this method instead of recomputing it

	os, oc := math.Sincos(rot - (math.Pi / 2))
	gaussR := rand.NormFloat64()
	sigma := math.Sqrt(h)
	mean := .5
	r := gaussR*sigma + mean

	opVec := vector.Vector2{X: r * oc, Y: r * -os}

	velr := 1 - (math.Abs(mean - gaussR))
	vel := 1 + .5*velr

	particlePos := vector.Vector2{X: opVec.X + x + (-c * (w / 2)), Y: opVec.Y + y - (-s * (w / 2))}
	return ThrustParticle{
		pos:   particlePos,
		vel:   vector.Vector2{X: -c * vel, Y: s * vel},
		start: time.Now(),
	}
}

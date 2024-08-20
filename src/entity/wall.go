package entity

import (
	"time"

	"github.com/abelroes/gmtk2024/src/collision"
	"github.com/abelroes/gmtk2024/src/vector"
	"github.com/hajimehoshi/ebiten/v2"
)

type Wall struct {
	W, H     float64
	Pos      vector.Vector2
	img      *ebiten.Image
	Collisor collision.CollisionRect
	Movement *WallMovement
}

type WallMovementState byte

const (
	MovementGoing WallMovementState = iota
	MovementGoingPause
	MovementReturning
	MovementReturningPause
)

type WallMovement struct {
	Direction  vector.Vector2
	Speed      float64
	Cooldown   time.Duration
	state      WallMovementState
	pauseStart time.Time
}

func NewWall(img *ebiten.Image, x, y, width, height float64, movement *WallMovement) Wall {
	pos := vector.Vector2{X: x, Y: y}
	return Wall{
		Pos: pos,
		W:   width, H: height,
		img: img,
		Collisor: collision.CollisionRect{
			Pos: pos,
			W:   width, H: height,
		},
		Movement: movement,
	}
}

func (wall *Wall) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	bounds := wall.img.Bounds()
	op.GeoM.Scale(wall.W/float64(bounds.Dx()), wall.H/float64(bounds.Dy()))
	op.GeoM.Translate(wall.Collisor.Pos.X, wall.Collisor.Pos.Y)
	op.Filter = ebiten.FilterLinear
	screen.DrawImage(wall.img, op)
}

func (wall *Wall) Update() {
	if wall.Movement == nil {
		return
	}

	dir := wall.Movement.Direction.Normalize()

	switch wall.Movement.state {
	case MovementGoing:
		wall.Collisor.Pos.Add(dir.MulScalar(wall.Movement.Speed))

		finalPos := wall.Pos.AddOut(wall.Movement.Direction)
		if wall.Collisor.Pos.Distance(finalPos) <= 5 {
			wall.Movement.state = MovementGoingPause
			wall.Movement.pauseStart = time.Now()
		}

	case MovementGoingPause:
		if time.Since(wall.Movement.pauseStart) >= wall.Movement.Cooldown {
			wall.Movement.state = MovementReturning
		}
	case MovementReturning:
		wall.Collisor.Pos.Add(dir.MulScalar(-wall.Movement.Speed))

		startPos := wall.Pos
		if wall.Collisor.Pos.Distance(startPos) <= 5 {
			wall.Movement.pauseStart = time.Now()
			wall.Movement.state = MovementReturningPause
		}
	case MovementReturningPause:
		if time.Since(wall.Movement.pauseStart) >= wall.Movement.Cooldown {
			wall.Movement.state = MovementGoing
		}
	}
}

package entity

import (
	_ "embed"
	"image/color"
	"math"
	"time"

	"github.com/abelroes/gmtk2024/src/collision"
	"github.com/abelroes/gmtk2024/src/constants"
	"github.com/abelroes/gmtk2024/src/settings"
	"github.com/abelroes/gmtk2024/src/vector"
	"github.com/hajimehoshi/ebiten/v2"
	ebivector "github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	initialManeuverability = 0.1
	airFriction            = 0.1
	scaleFactor            = 0.005
	initialScale           = .5
	constantPropulsion     = 0.0
	minimumScale           = 0.03
	oobMargin              = 50
)

type PlayerEvent int

const (
	PlayerDiedByCollision = iota
	PlayerDiedByShrinking
	PlayerDiedByOutOfBounds
)

type Player struct {
	Dead            bool
	Scale           float64
	Rot             float64
	Propulsion      float64
	Pos             vector.Vector2
	Vel             vector.Vector2
	Acl             vector.Vector2
	Maneuverability float64
	MinimumScale    float64
	Collisor        collision.CollisionPolygon
	basePolygon     []vector.Vector2
	eventHandler    func(PlayerEvent)
	bounds          collision.CollisionRect
	debugSettings   *settings.SettingsDebug

	//NOTE: this is allocating more memory than needed
	thrustParticles []ThrustParticle

	img *ebiten.Image
}

func NewPlayer(img *ebiten.Image, debugSettings *settings.SettingsDebug, basePolygon []vector.Vector2, eventHandler func(PlayerEvent)) Player {
	return Player{
		Scale:           initialScale,
		img:             img,
		Maneuverability: initialManeuverability,
		thrustParticles: make([]ThrustParticle, 0, 1000),
		Collisor:        collision.CollisionPolygon{Vertices: make([]vector.Vector2, len(basePolygon))},
		basePolygon:     basePolygon,
		MinimumScale:    minimumScale,
		eventHandler:    eventHandler,
		bounds: collision.CollisionRect{
			Pos: vector.New(-oobMargin, -oobMargin),
			W:   constants.Width + (oobMargin * 2),
			H:   constants.Height + (oobMargin * 2),
		},
		debugSettings: debugSettings,
	}
}

func (player *Player) Reset() {
	player.Dead = false
	player.Rot = 0
	player.Vel.Set(0, 0)
	player.Acl.Set(0, 0)
	player.Scale = initialScale
}

func (player *Player) updateCollider(sin, cos float64) {
	for i := 0; i < len(player.Collisor.Vertices); i++ {
		var (
			base = player.basePolygon[i]
			// vector rotation
			x = cos*base.X - sin*base.Y
			y = sin*base.X + cos*base.Y
		)

		player.Collisor.Vertices[i] = vector.Vector2{
			X: (player.Pos.X + x*player.Scale),
			Y: (player.Pos.Y - y*player.Scale),
		}
	}
}

func (player *Player) DieByCollision() {
	player.die(PlayerDiedByCollision)
}

func (player *Player) die(event PlayerEvent) {
	player.Dead = true
	player.eventHandler(event)
}

func (player *Player) GetDimensions() (float64, float64) {
	bounds := player.img.Bounds()
	w, h := float64(bounds.Dx()), float64(bounds.Dy())
	w *= player.Scale
	h *= player.Scale

	return w, h
}

func (player *Player) GetTopLeftPos() vector.Vector2 {
	w, h := player.GetDimensions()
	return vector.Vector2{
		X: player.Pos.X - w/2,
		Y: player.Pos.Y - h/2,
	}
}

func (player *Player) Update() {
	player.updateThrustParticles()

	if player.Dead {
		return
	}

	if player.checkOob() {
		return
	}

	w, h := player.GetDimensions()

	if player.Scale <= player.MinimumScale {
		player.die(PlayerDiedByShrinking)
		return
	}

	sin, cos := math.Sincos(player.Rot)

	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		player.Rot += player.Maneuverability
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		player.Rot -= player.Maneuverability
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeySpace) {
		player.Propulsion = 1.0

		vol := 5
		for i := 0; i < vol; i++ {
			player.thrustParticles = append(player.thrustParticles, SpawnParticle(w, h, player.Rot, player.Pos.X, player.Pos.Y))
		}
	}

	if player.Propulsion > 0.0 {
		player.Acl.Set(player.Propulsion*cos, -player.Propulsion*sin)
		player.Scale = max(0, player.Scale-(player.Propulsion*scaleFactor))
	}

	player.Vel.Add(player.Acl)
	player.Pos.Add(player.Vel)
	player.Acl.Set(0, 0)
	player.Propulsion = constantPropulsion

	player.Vel.Sub(player.Vel.MulScalar(airFriction))

	player.updateCollider(sin, cos)
}

func (player *Player) checkOob() bool {
	var (
		pos    = player.Pos
		bounds = player.bounds

		horizontal = pos.X < bounds.Pos.X || pos.X >= bounds.Pos.X+bounds.W
		vertical   = pos.Y < bounds.Pos.Y || pos.Y >= bounds.Pos.Y+bounds.H
		isOob      = horizontal || vertical
	)
	if isOob {
		player.die(PlayerDiedByOutOfBounds)
		return true
	}
	return false
}

func (player *Player) updateThrustParticles() {
	//NOTE: Probably refact to Particles or smthing like that
	now := time.Now()

	lastDead := -1
	for i := len(player.thrustParticles) - 1; i >= 0; i-- {
		particle := player.thrustParticles[i]
		isDead := now.Sub(particle.start) > particleThrustTtl
		if isDead {
			lastDead = i
			break
		}
	}

	if lastDead != -1 {
		player.thrustParticles = player.thrustParticles[lastDead:]
	}

	for i := range player.thrustParticles {
		player.thrustParticles[i].Update()
	}
}

func (player *Player) Draw(screen *ebiten.Image, camera Camera) {
	now := time.Now()
	for _, particle := range player.thrustParticles {
		particle.Draw(screen, now)
	}

	if player.Dead {
		return
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(player.Scale, player.Scale)
	w, h := player.GetDimensions()

	// center position based on the image * scale
	op.GeoM.Translate(-w/2, -h/2)
	// rotate
	op.GeoM.Rotate(-player.Rot)
	// shift position based on camera position
	op.GeoM.Translate(camera.X, -camera.Y)
	// position image based on player position
	op.GeoM.Translate(player.Pos.X, player.Pos.Y)
	op.Filter = ebiten.FilterLinear
	screen.DrawImage(player.img, op)

	if player.debugSettings.PlayerHitbox {
		player.drawPolygonDebugHitBox(screen)
	}
}

func (player *Player) drawPolygonDebugHitBox(screen *ebiten.Image) {
	polygon := player.Collisor

	next := 0
	verticesQtd := len(polygon.Vertices)

	color := color.RGBA{R: 255}
	for current := 0; current < verticesQtd; current++ {
		next = current + 1
		if next == verticesQtd {
			next = 0
		}

		currentVec := polygon.Vertices[current]
		nextVec := polygon.Vertices[next]

		ebivector.StrokeLine(screen, float32(currentVec.X), float32(currentVec.Y), float32(nextVec.X), float32(nextVec.Y), 2, color, true)
		ebivector.DrawFilledCircle(screen, float32(currentVec.X), float32(currentVec.Y), 2, color, true)
	}
}

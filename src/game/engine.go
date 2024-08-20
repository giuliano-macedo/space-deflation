package game

import (
	"fmt"
	"image/color"

	"github.com/abelroes/gmtk2024/assets"
	"github.com/abelroes/gmtk2024/assets/levels"
	"github.com/abelroes/gmtk2024/src/audio"
	"github.com/abelroes/gmtk2024/src/collision"
	"github.com/abelroes/gmtk2024/src/entity"
	"github.com/abelroes/gmtk2024/src/settings"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Engine struct {
	player       entity.Player
	enemies      []entity.Wall
	goal         entity.Goal
	camera       entity.Camera
	ui           *entity.Ui
	audioManager *audio.Manager
	asset        *assets.Asset
	background   *entity.Background
	settings     *settings.Settings
	onGameWin    func()

	currentLevelIndex int
}

func NewEngine(asset *assets.Asset, audioManager *audio.Manager, settings *settings.Settings) *Engine {
	gameEngine := &Engine{
		enemies:      []entity.Wall{},
		audioManager: audioManager,
		asset:        asset,
		goal:         entity.NewGoal(asset.GetImage(assets.GoalImgIndex), &settings.Debug, 40),
		ui:           entity.NewUi(asset.Font),
		background:   entity.NewBackground(asset.Backgrounds),

		currentLevelIndex: settings.Debug.InitialLevel,
		settings:          settings,
	}
	gameEngine.player = entity.NewPlayer(asset.GetImage(assets.PlayerImgIndex), &settings.Debug, asset.PlayerPolygon, gameEngine.handlePlayerEvents)

	gameEngine.resetLevel()
	return gameEngine
}

func (e *Engine) SetOnGameWin(onGameWin func()) {
	e.onGameWin = onGameWin
}

func (g *Engine) win() {
	g.audioManager.PlaySoundFx(audio.WinFx)

	if g.currentLevelIndex == len(g.asset.Levels)-1 {
		g.onGameWin()
		g.currentLevelIndex = 0
		g.resetLevel()
		return
	}
	g.goToNextLevel()
}

func (e *Engine) goToNextLevel() {
	e.background.ChangeBG()
	e.currentLevelIndex++
	e.resetLevel()
}

func (g *Engine) drawPlayer(screen *ebiten.Image) {
	g.player.Draw(screen, g.camera)
}

func (g *Engine) drawEnemies(screen *ebiten.Image) {
	for _, enemy := range g.enemies {
		enemy.Draw(screen)
	}
}

func (g *Engine) handlePlayerEvents(event entity.PlayerEvent) {
	switch event {

	case entity.PlayerDiedByShrinking:
		g.audioManager.PlaySoundFx(audio.PopFx)
		g.ui.ShowRestartText = true

	case entity.PlayerDiedByCollision, entity.PlayerDiedByOutOfBounds:
		g.audioManager.PlaySoundFx(audio.ExplosionFx)
		g.ui.ShowRestartText = true
	}
}

func (g *Engine) drawBg(screen *ebiten.Image) {
	g.background.Draw(screen)
}

func (g *Engine) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 0xff})
	g.drawBg(screen)
	g.drawPlayer(screen)
	g.drawEnemies(screen)
	g.goal.Draw(screen)
	g.ui.Draw(screen)

	if g.settings.Debug.Fps {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("Fps: %.2f Tps: %.2f\ncurrentLevel: %d", ebiten.ActualFPS(), ebiten.ActualTPS(), g.currentLevelIndex+1))
	}
}

func (g *Engine) resetLevel() {
	g.setLevel(g.asset.Levels[g.currentLevelIndex])
}

func (g *Engine) setLevel(level levels.Level) {
	g.player.Reset()
	g.player.Pos = level.PlayerStartPos
	g.goal.SetPos(level.GoalPos)

	walls := make([]entity.Wall, 0, len(level.Walls))

	for _, wallInfo := range level.Walls {
		var movement *entity.WallMovement
		if wallInfo.Movement != nil {
			movement = &entity.WallMovement{
				Direction: wallInfo.Movement.Direction,
				Speed:     wallInfo.Movement.Speed,
				Cooldown:  wallInfo.Movement.Cooldown,
			}
		}

		wall := entity.NewWall(g.asset.GetImage(assets.EnemyImgIndex), wallInfo.Pos.X, wallInfo.Pos.Y, wallInfo.W, wallInfo.H, movement)
		walls = append(walls, wall)
	}

	g.enemies = walls
}

func (g *Engine) Update() error {
	if g.player.Dead && (inpututil.IsKeyJustPressed(ebiten.KeyR) || inpututil.IsKeyJustPressed(ebiten.KeyEnter)) {
		g.ui.ShowRestartText = false
		g.resetLevel()
	}

	//Note: could gain perfomance by only updating moving walls
	for i := range g.enemies {
		g.enemies[i].Update()
	}

	g.audioManager.PlaySoundTrackInLoop()

	g.player.Update()
	g.goal.Update()
	g.collisionDetection()

	return nil
}

func (g *Engine) collisionDetection() {
	if !g.player.Dead {
		for _, wall := range g.enemies {
			if collision.HasCollidedRectPolygon(wall.Collisor, g.player.Collisor) {
				g.player.DieByCollision()
			}
		}

		if collision.HasCollidedRectPolygon(g.goal.Collider, g.player.Collisor) {
			g.win()
		}
	}
}

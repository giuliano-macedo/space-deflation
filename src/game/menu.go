package game

import (
	"image/color"
	"time"

	"github.com/abelroes/gmtk2024/assets"
	"github.com/abelroes/gmtk2024/src/audio"
	"github.com/abelroes/gmtk2024/src/constants"
	"github.com/abelroes/gmtk2024/src/entity"
	"github.com/abelroes/gmtk2024/src/settings"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
	gameName        = "Space Deflation"
	creditsRolStart = 8 * time.Second
	creditRollSpeed = .25
	maxCreditRoll   = -400.0
)

type MenuState int

const (
	OnMenuState MenuState = iota
	PlayingState
	CreditsState
)

type Menu struct {
	gameEngine             *Engine
	menuTransitionHappened bool
	audioManager           *audio.Manager
	font                   *text.GoTextFaceSource
	background             *entity.Background
	settings               *settings.Settings
	state                  MenuState
	credits                string
	creditsStartedAt       time.Time
	creditsY               float64
}

func NewMenu(assets *assets.Asset, audioManager *audio.Manager, gameEngine *Engine, settings *settings.Settings, credits string) *Menu {
	initialState := OnMenuState
	if settings.Debug.SkipMenu {
		initialState = PlayingState
	}

	menu := &Menu{
		state:        initialState,
		audioManager: audioManager,
		gameEngine:   gameEngine,
		background:   entity.NewBackground(assets.Backgrounds),
		font:         assets.Font,
		settings:     settings,
		credits:      credits,
	}
	return menu
}

func (m *Menu) Update() error {
	switch m.state {
	case PlayingState:
		m.gameEngine.Update()
	case OnMenuState:
		m.audioManager.PlaySoundTrackInLoop()
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeyKPEnter) {
			m.state = PlayingState
		}
	case CreditsState:
		if m.creditsStartedAt == (time.Time{}) {
			m.creditsStartedAt = time.Now()
		}

		if !m.audioManager.IsPlayingWinTrack() {
			m.audioManager.PlaySoundTrackInLoop()
		}

		if time.Since(m.creditsStartedAt) > creditsRolStart && m.creditsY > maxCreditRoll {
			m.creditsY -= creditRollSpeed
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			m.state = OnMenuState
			m.creditsStartedAt = time.Time{}
			m.creditsY = 0
		}
	}

	return nil
}

func (m *Menu) OnGameWin() {
	m.state = CreditsState
	m.audioManager.StopSoundTrack()
	m.audioManager.PlayWinTrack()
	m.background.ChangeBG()
}

func (m *Menu) Draw(screen *ebiten.Image) {
	switch m.state {
	case PlayingState:
		if !m.menuTransitionHappened {
			m.menuTransitionHappened = true
		}

		m.gameEngine.Draw(screen)
	case OnMenuState:
		m.background.Draw(screen)
		w := constants.Width / 2.0
		h := constants.Height / 2.0

		m.DrawText(screen, gameName, 1, w, h-150)
		m.DrawText(screen, "Use UP/W for propulsion and LEFT/A or RIGHT/D to steer", 7, w, h-50)
		m.DrawText(screen, "Press ENTER to start", 7, w, h)

	case CreditsState:
		m.background.Draw(screen)
		w := constants.Width / 2.0
		m.DrawText(screen, "Congratulations!", 1, w, m.creditsY)
		m.DrawText(screen, "thanks for playing", 2, w, m.creditsY+50)
		m.DrawText(screen, "Credits", 3, w, m.creditsY+120)
		m.DrawText(screen, m.credits+"\nPress enter to go back to main menu", 6, w, m.creditsY+150)
	}
}

func (m *Menu) DrawText(screen *ebiten.Image, textStr string, size int, x, y float64) {
	var scale float64 = 1.0
	switch size {
	case 2:
		scale = 0.88
	case 3:
		scale = 0.7
	case 4:
		scale = 0.58
	case 5:
		scale = 0.52
	case 6:
		scale = 0.47
	case 7:
		scale = 0.25
	}
	fontSize := 50.0

	op := &text.DrawOptions{
		LayoutOptions: text.LayoutOptions{PrimaryAlign: text.AlignCenter, LineSpacing: 42.0},
	}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(x, y)

	op.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, textStr, &text.GoTextFace{
		Source: m.font,
		Size:   fontSize,
	}, op)
}

func (g *Menu) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return constants.Width, constants.Height
}

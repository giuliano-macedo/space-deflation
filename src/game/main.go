package game

import (
	"encoding/json"
	"os"
	"runtime"

	"github.com/abelroes/gmtk2024/assets"
	"github.com/abelroes/gmtk2024/src/audio"
	"github.com/abelroes/gmtk2024/src/settings"
	"github.com/hajimehoshi/ebiten/v2"
)

func loadSettings() (*settings.Settings, error) {
	if runtime.GOOS == "js" {
		return settings.DefaultSettings, nil
	}

	f, err := os.Open("settings.json")
	if os.IsNotExist(err) {
		return settings.DefaultSettings, nil
	}

	if err != nil {
		return nil, err
	}

	settings := &settings.Settings{}
	err = json.NewDecoder(f).Decode(settings)
	return settings, err
}

func Main() error {
	settings, err := loadSettings()
	if err != nil {
		return err
	}

	assets, err := assets.New()
	if err != nil {
		return err
	}

	audioManager, err := audio.NewManager(assets.Sounds, settings.Volume)
	if err != nil {
		return err
	}

	gameEngine := NewEngine(assets, audioManager, settings)
	menu := NewMenu(assets, audioManager, gameEngine, settings, assets.Credits)
	gameEngine.SetOnGameWin(menu.OnGameWin)

	ebiten.SetWindowSize(settings.Screen.Width, settings.Screen.Height)
	ebiten.SetWindowTitle(gameName)

	return ebiten.RunGame(menu)
}

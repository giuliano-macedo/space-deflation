package settings

import (
	"github.com/abelroes/gmtk2024/src/constants"
)

type SettingsVolume struct {
	SoundTrack float64
	SoundFx    float64
}

type SettingsDebug struct {
	Fps          bool
	PlayerHitbox bool
	GoalHitbox   bool
	SkipSplash   bool
	SkipMenu     bool
	InitialLevel int
}

type SettingsScreen struct {
	Width, Height int
}

type Settings struct {
	Volume SettingsVolume
	Screen SettingsScreen
	Debug  SettingsDebug
}

var DefaultSettings = &Settings{
	Volume: SettingsVolume{
		SoundTrack: .65,
		SoundFx:    1,
	},
	Screen: SettingsScreen{
		Width:  constants.Width * 1.5,
		Height: constants.Height * 1.5,
	},
}

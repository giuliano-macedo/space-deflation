package audio

import (
	"math/rand"
	"time"

	"github.com/abelroes/gmtk2024/src/settings"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

const audioSampleRate = 44000

type SoundIndex int

const (
	SoundTrackIndex SoundIndex = iota
	WinTrackIndex
	Explosion1Index
	Explosion2Index
	Explosion3Index
	PlopIndex
	PopIndex
	Rocket1Index
	Rocket2Index
	SwooshIndex
)

type SoundFx int

const (
	ExplosionFx SoundFx = iota
	PopFx
	RocketFx
	WinFx
)

type Manager struct {
	ctx          *audio.Context
	audioPlayers []*audio.Player
	rng          *rand.Rand
}

func (manager *Manager) createAudioPlayer(audio *mp3.Stream) (*audio.Player, error) {
	audioPlayer, err := manager.ctx.NewPlayer(audio)
	if err != nil {
		return nil, err
	}

	audioPlayer.SetVolume(1)
	return audioPlayer, nil
}

func NewManager(streams []*mp3.Stream, volume settings.SettingsVolume) (*Manager, error) {
	manager := &Manager{
		ctx:          audio.NewContext(audioSampleRate),
		audioPlayers: make([]*audio.Player, len(streams)),
		rng:          rand.New(rand.NewSource(time.Now().UnixMilli())),
	}

	for key, value := range streams {
		audioPlayer, err := manager.createAudioPlayer(value)
		if err != nil {
			return nil, err
		}

		manager.audioPlayers[key] = audioPlayer
	}

	manager.SetVolumes(volume)
	return manager, nil
}

func (manager *Manager) PlaySoundFx(soundFx SoundFx) {
	switch soundFx {
	case ExplosionFx:
		switch manager.rng.Int() % 3 {
		case 0:
			manager.play(Explosion1Index, 1650)
		case 1:
			manager.play(Explosion2Index, 1600)
		case 2:
			manager.play(Explosion3Index, 1250)
		}
	case WinFx:
		// manager.play(SwooshIndex, 831)
		manager.play(PlopIndex, 693)
	case PopFx:
		manager.play(PopIndex, 1800)
	case RocketFx:
	}
}

func (manager *Manager) PlaySoundTrackInLoop() {
	if !manager.audioPlayers[SoundTrackIndex].IsPlaying() {
		manager.play(SoundTrackIndex, 3000)
	}
}

func (manager *Manager) StopSoundTrack() {
	player := manager.audioPlayers[SoundTrackIndex]
	player.SetPosition(0)
	player.Pause()
}

func (manager *Manager) PlayWinTrack() {
	manager.play(WinTrackIndex, 3700)
}

func (manager *Manager) IsPlayingWinTrack() bool {
	return manager.audioPlayers[WinTrackIndex].IsPlaying()
}

func (manager *Manager) SetVolumes(volume settings.SettingsVolume) {
	for key, player := range manager.audioPlayers {
		isSoundTrack := SoundIndex(key) == SoundTrackIndex || SoundIndex(key) == WinTrackIndex
		if isSoundTrack {
			player.SetVolume(volume.SoundTrack)
		} else {
			player.SetVolume(volume.SoundFx)
		}
	}
}

/*
Offset in milliseconds
*/
func (manager *Manager) play(key SoundIndex, offset int) {
	player := manager.audioPlayers[key]
	player.SetPosition(time.Duration(offset) * time.Millisecond)
	player.Play()
}

package assets

import (
	_ "embed"
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"image"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/abelroes/gmtk2024/assets/levels"
	"github.com/abelroes/gmtk2024/src/audio"
	"github.com/abelroes/gmtk2024/src/vector"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type ImageIndex int

const (
	PlayerImgIndex ImageIndex = iota
	GoalImgIndex
	EnemyImgIndex
)

type Asset struct {
	PlayerPolygon []vector.Vector2

	Images      []*ebiten.Image
	Backgrounds []*ebiten.Image

	Sounds []*mp3.Stream
	Font   *text.GoTextFaceSource
	Levels []levels.Level

	Credits string
}

var imagesFilenames = []string{
	PlayerImgIndex: "images/rocket.png",
	EnemyImgIndex:  "images/pipe.png",
	GoalImgIndex:   "images/blackhole.png",
}

var soundsFileNames = []string{
	audio.SoundTrackIndex: "music/EscapeFromDeathStar.mp3",
	audio.WinTrackIndex:   "music/WinSong.mp3",
	audio.Explosion1Index: "sounds/explosion1.mp3",
	audio.Explosion2Index: "sounds/explosion2.mp3",
	audio.Explosion3Index: "sounds/explosion3.mp3",
	audio.PopIndex:        "sounds/pop.mp3",
	audio.Rocket1Index:    "sounds/rocket1.mp3",
	audio.Rocket2Index:    "sounds/rocket2.mp3",
	audio.SwooshIndex:     "sounds/win.mp3",
	audio.PlopIndex:       "sounds/plop.mp3",
}

func New() (*Asset, error) {
	imgs, err := readImgs()
	if err != nil {
		return nil, err
	}

	playerPolygon, err := readPlayerPolygon(imgs[PlayerImgIndex])
	if err != nil {
		return nil, err
	}

	soundMap, err := readSounds()
	if err != nil {
		return nil, err
	}

	fontFile, err := embededFs.Open("fonts/PixelifySans-Regular.ttf")
	if err != nil {
		return nil, err
	}

	assetFont, err := text.NewGoTextFaceSource(fontFile)
	if err != nil {
		return nil, err
	}

	backgrounds, err := readBackgrounds()
	if err != nil {
		return nil, err
	}

	credits, err := embededFs.ReadFile("credits.txt")
	if err != nil {
		return nil, err
	}

	lvls, err := readLevels()
	if err != nil {
		return nil, err
	}
	_ = lvls

	return &Asset{
		PlayerPolygon: playerPolygon,
		Sounds:        soundMap,
		Font:          assetFont,
		Levels:        lvls,
		// Levels:      levels.Levels,
		Backgrounds: backgrounds,
		Images:      imgs,
		Credits:     string(credits),
	}, nil
}

func (asset *Asset) GetImage(index ImageIndex) *ebiten.Image {
	return asset.Images[index]
}

func readSounds() ([]*mp3.Stream, error) {
	m := make([]*mp3.Stream, len(soundsFileNames))
	for index, filename := range soundsFileNames {
		f, err := embededFs.Open(filename)
		if err != nil {
			return nil, err
		}

		stream, err := mp3.DecodeWithoutResampling(f)
		if err != nil {
			return nil, err
		}

		m[index] = stream
	}

	return m, nil
}

func readImg(fileName string) (*ebiten.Image, error) {
	f, err := embededFs.Open(fileName)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	return ebiten.NewImageFromImage(img), nil
}

func readImgs() ([]*ebiten.Image, error) {
	imgs := make([]*ebiten.Image, len(imagesFilenames))
	for index, filename := range imagesFilenames {
		img, err := readImg(filename)
		if err != nil {
			return nil, err
		}

		imgs[index] = img
	}

	return imgs, nil
}

func readLevels() ([]levels.Level, error) {
	f, err := embededFs.Open("levels/levels.tmx")
	if err != nil {
		return nil, err
	}

	var tmxMap levels.Map
	xml.NewDecoder(f).Decode(&tmxMap)

	lvls := make([]levels.Level, len(tmxMap.ObjectGroups))

	for i, group := range tmxMap.ObjectGroups {
		lvl := &lvls[i]

		playerObj := group.FindObjectByName("player")
		if playerObj == nil {
			return nil, fmt.Errorf("player not found in %s", group.Name)
		}
		goalObj := group.FindObjectByName("goal")
		if goalObj == nil {
			return nil, fmt.Errorf("goal not found in %s", group.Name)
		}

		lvl.PlayerStartPos = playerObj.CenterPos()

		lvl.GoalPos = goalObj.CenterPos()

		for _, obj := range group.Objects {
			if obj.Name == "pipe" {

				movement, err := getWallMovementFromObj(obj)
				if err != nil {
					return nil, err
				}

				wall := levels.WallInfo{
					W:        obj.Width,
					H:        obj.Height,
					Pos:      obj.TopLeftPos(),
					Movement: movement,
				}
				lvl.Walls = append(lvl.Walls, wall)
			}
		}
	}

	return lvls, nil
}

func getWallMovementFromObj(obj levels.Object) (*levels.WallMovementInfo, error) {
	if obj.Props == nil {
		return nil, nil
	}

	directionStr, _ := obj.Props.GetPropString("direction")
	if directionStr == "" {
		return nil, nil
	}
	xStr, yStr, found := strings.Cut(directionStr, ",")
	if !found {
		return nil, fmt.Errorf("failed parsing direction")
	}

	x, err := strconv.ParseFloat(xStr, 64)
	if err != nil {
		return nil, err
	}

	y, err := strconv.ParseFloat(yStr, 64)
	if err != nil {
		return nil, err
	}

	speed, err := obj.Props.GetPropFloat("speed")
	if err != nil {
		return nil, err
	}

	cooldown, err := obj.Props.GetPropFloat("cooldown")
	if err != nil {
		return nil, err
	}

	return &levels.WallMovementInfo{
		Direction: vector.New(x, y),
		Speed:     speed,
		Cooldown:  time.Duration(cooldown),
	}, nil
}

func readBackgrounds() ([]*ebiten.Image, error) {
	backgrounds := []*ebiten.Image{}

	dirName := "images"
	entries, err := embededFs.ReadDir(dirName)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		fname := entry.Name()
		if strings.HasPrefix(fname, "background") {
			img, err := readImg(path.Join(dirName, fname))
			if err != nil {
				return nil, err
			}
			backgrounds = append(backgrounds, img)
		}
	}

	if len(backgrounds) == 0 {
		return nil, fmt.Errorf("no background images found in %s", dirName)
	}

	return backgrounds, nil
}

func readPlayerPolygon(playerImg image.Image) ([]vector.Vector2, error) {
	f, err := embededFs.Open("hitboxes/player_poly.csv")
	if err != nil {
		return nil, err
	}

	records, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, err
	}

	bounds := playerImg.Bounds()
	imgW, imgH := float64(bounds.Max.X), float64(bounds.Max.Y)

	polygons := make([]vector.Vector2, 0, len(records)-1)
	for _, record := range records[1:] {
		x, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			return nil, err
		}

		y, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			return nil, err
		}

		polygons = append(polygons, vector.New(x-imgW*.5, y-imgH*.5))
	}

	return polygons, nil
}

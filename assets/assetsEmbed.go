package assets

import "embed"

//go:embed fonts/PixelifySans-Regular.ttf
//go:embed hitboxes/player_poly.csv
//go:embed sounds/*.mp3
//go:embed music/*.mp3
//go:embed images/*.png
//go:embed levels/levels.tmx
//go:embed credits.txt
var embededFs embed.FS

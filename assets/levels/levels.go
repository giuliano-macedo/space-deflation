package levels

import (
	"time"

	"github.com/abelroes/gmtk2024/src/vector"
)

type WallMovementInfo struct {
	Direction vector.Vector2
	Speed     float64
	Cooldown  time.Duration
}

type WallInfo struct {
	W, H     float64
	Pos      vector.Vector2
	Movement *WallMovementInfo
}

type Level struct {
	PlayerStartPos vector.Vector2
	GoalPos        vector.Vector2
	Walls          []WallInfo
}

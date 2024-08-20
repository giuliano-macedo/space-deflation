package collision

import "github.com/abelroes/gmtk2024/src/vector"

type CollisionRect struct {
	Pos vector.Vector2
	W   float64
	H   float64
}

type CollisionPolygon struct {
	Vertices []vector.Vector2
}

func lineLine(x1, y1, x2, y2, x3, y3, x4, y4 float64) bool {
	// calculate the direction of the lines
	uA := ((x4-x3)*(y1-y3) - (y4-y3)*(x1-x3)) / ((y4-y3)*(x2-x1) - (x4-x3)*(y2-y1))
	uB := ((x2-x1)*(y1-y3) - (y2-y1)*(x1-x3)) / ((y4-y3)*(x2-x1) - (x4-x3)*(y2-y1))

	// if uA and uB are between 0-1, lines are colliding
	if uA >= 0 && uA <= 1 && uB >= 0 && uB <= 1 {
		return true
	}
	return false
}

func lineRect(x1, y1, x2, y2, rx, ry, rw, rh float64) bool {
	// check if the line has hit any of the rectangle's sides
	// uses the Line/Line function below
	left := lineLine(x1, y1, x2, y2, rx, ry, rx, ry+rh)
	right := lineLine(x1, y1, x2, y2, rx+rw, ry, rx+rw, ry+rh)
	top := lineLine(x1, y1, x2, y2, rx, ry, rx+rw, ry)
	bottom := lineLine(x1, y1, x2, y2, rx, ry+rh, rx+rw, ry+rh)

	// if ANY of the above are true,
	// the line has hit the rectangle
	if left || right || top || bottom {
		return true
	}
	return false
}

// Source: https://www.jeffreythompson.org/collision-detection/poly-rect.php
func HasCollidedRectPolygon(rect CollisionRect, polygon CollisionPolygon) bool {
	next := 0
	verticesQtd := len(polygon.Vertices)

	for current := 0; current < verticesQtd; current++ {
		next = current + 1
		if next == verticesQtd {
			next = 0
		}

		currentVec := polygon.Vertices[current]
		nextVec := polygon.Vertices[next]

		collision := lineRect(currentVec.X, currentVec.Y, nextVec.X, nextVec.Y, rect.Pos.X, rect.Pos.Y, rect.W, rect.H)

		if collision {
			return true
		}

	}

	return false
}

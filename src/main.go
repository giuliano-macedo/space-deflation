package main

import "github.com/abelroes/gmtk2024/src/game"

func main() {
	err := game.Main()
	if err != nil {
		panic(err)
	}
}

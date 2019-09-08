package main

import (
	"fmt"

	"github.com/sindriing/godoku/game"
)

func main() {
	game := game.InitializeBoard("presets/board1.txt")
	fmt.Println("Initial Sudoku board")
	game.PrintBoard()
	fmt.Println()
	game.FirstSweep()

	progress := game.Solve()
	fmt.Println("Solved Sudoku board")
	progress.PrintBoard()

}

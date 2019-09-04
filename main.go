package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
)

func main() {
	fmt.Println("Hello World")
	game := initializeBoard()

}

type cell struct {
	blocked  [9]bool //Which numbers can't be placed
	exploded bool    //Have we used this value for removing possibilities
	value    int
}

type sudoku struct {
	board  [9][9]cell
	filled int //How many total squares filled
}

func (s *sudoku) explode(x int, y int, val int) {
	s.board[x][y].exploded = true

	s.board[x][y].value = val
	s.explodeInner(x, y, val)
	s.explodeVertical(x, val)
	s.explodeHorizontal(y, val)
}

//Remve possibilities in the small 3x3 chunks
func (s *sudoku) explodeInner(x int, y int, val int) {
	x = x - (x % 3)
	y = y - (y % 3)

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if s.board[i][j].blocked[val] == false {
				s.board[i][j].blocked[val] = true
			}
		}
	}
}

func (s *sudoku) explodeVertical(x int, val int) {
	for i := 0; i < 9; i++ {
		if s.board[x][i].blocked[val] == false {
			s.board[x][i].blocked[val] = true
		}
	}
}

func (s *sudoku) explodeHorizontal(y int, val int) {
	for i := 0; i < 9; i++ {
		if s.board[i][y].blocked[val] == false {
			s.board[i][y].blocked[val] = true
		}
	}
}

func initializeBoard(path string) sudoku {
	Sudoku := sudoku{readBoard(path), 0}
	return board
}

func (s sudoku)printBoard() {
	var vals [9]interface{}
	for i:=0; i<9; i++ {
		for j:=0; j<9; j++{
			vals = Append(s.board[i][j].value)
		}
		if i%3 == 0 {
			fmt.Println("-------------")
		}
		fmt.Printf{"|%d%d%d|%d%d%d|%d%d%d|", vals...}
	}
	fmt.Println("-------------")
}

func readBoard(path string) [9][9]cell {
	var board [9][9]cell
	var noBlock [9]bool
	dat, err := ioutil.ReadFile(path)
	check(err)

	//Split the input file on commas and newlines
	rule := regexp.MustCompile("\n,")
	list := rule.Split(string(dat), -1)

	c := 0
	var value int
	for i, c := 0, 0; i < 9; i++ {
		for j := 0; j < 9; j, c = j+1, c+1 {
			value, err = strconv.Atoi(list[c])
			board[i][j] = cell{noBlock, false, value}
		}
	}
	return board
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

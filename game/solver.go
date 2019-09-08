package game

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var allBlock = [9]bool{true, true, true, true, true, true, true, true, true}

type cell struct {
	blocked [9]bool //Which numbers can't be placed
	value   int
}

type sudoku struct {
	board  [9][9]cell
	filled int //How many total squares filled
}

//Attempts to solve the board, cannot make guesses
func (s sudoku) Solve() sudoku {
	for s.filled < 81 {
		change := false
		for i := 0; i < 9; i++ {
			for j := 0; j < 9; j++ {
				if s.isFree(i, j) {
					// Check if there is only one available value
					if val := s.checkOneVal(i, j); val != 0 {
						s.explode(i, j, val)
						change = true
					}
					// Check row
					if val := s.checkRow(i, j); val != 0 {
						s.explode(i, j, val)
						change = true
					}

					// Check column
					if val := s.checkCol(i, j); val != 0 {
						s.explode(i, j, val)
						change = true
					}

					// Check 3x3 chunk
					if val := s.checkChunk(i, j); val != 0 {
						s.explode(i, j, val)
						change = true
					}
				}
			}
		}
		if !change {
			fmt.Println("Current rule set provides no possible moves")
			break
		}
	}
	return s
}

//Check if this cell has only one possible value left
func (s *sudoku) checkOneVal(x int, y int) int {
	var last, sum int
	for i, blocked := range s.board[x][y].blocked {
		if !blocked {
			sum++
			last = i
		}
		if sum > 1 {
			return 0
		}
	}
	return last + 1
}

//Check if I'm the only one in my row who can receive a specific value
func (s *sudoku) checkRow(x int, y int) int {
	var success bool
	for v, block := range s.board[x][y].blocked {
		if !block {
			success = true
			for i := 0; i < 9; i++ {
				if !s.board[x][i].blocked[v] && i != y {
					//Another cell can have this value
					success = false
					break
				}
			}
			if success {
				return v + 1
			}
		}
	}
	return 0
}

//Check if I'm the only one in my column who can receive a specific value
func (s *sudoku) checkCol(x int, y int) int {
	var success bool
	for v, block := range s.board[x][y].blocked {
		if !block {
			success = true
			for i := 0; i < 9; i++ {
				if !s.board[i][y].blocked[v] && i != x {
					//Another cell can have this value
					success = false
					break
				}
			}
			if success {
				return v + 1
			}
		}
	}
	return 0
}

//Initial solve run adds no new numbers but explodes all preset values
func (s *sudoku) FirstSweep() {
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if curr := s.board[i][j].value; curr != 0 {
				s.explode(i, j, curr)
			}
		}
	}
}

//Check if I'm the only one in my 3x3 chunk who can receive a specific value
func (s *sudoku) checkChunk(x int, y int) int {
	xn := x - (x % 3)
	yn := y - (y % 3)
	var success bool
	for v, block := range s.board[x][y].blocked {
		if !block {
			success = true
			for i := xn; i < xn+3; i++ {
				for j := yn; j < yn+3; j++ {
					if !s.board[i][j].blocked[v] && (i != x || j != y) {
						//Another cell can have this value
						success = false
						break
					}
				}
			}
			if success {
				return v + 1
			}
		}
	}
	return 0
}

// "Explode" means to spread the influence of a value to restrict other cells
// from taking that value
func (s *sudoku) explode(x int, y int, val int) {
	s.board[x][y].blocked = allBlock
	s.filled++

	s.board[x][y].value = val
	val-- //indexing is 0-8 while sudoku is 1-9
	s.explodeChunk(x, y, val)
	s.explodeVertical(x, val)
	s.explodeHorizontal(y, val)
}

//Remve possibilities in the small 3x3 chunks
func (s *sudoku) explodeChunk(x int, y int, val int) {
	x = x - (x % 3)
	y = y - (y % 3)

	for i := x; i < x+3; i++ {
		for j := y; j < y+3; j++ {
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

func (s *sudoku) isFree(x int, y int) bool {
	return s.board[x][y].value == 0
}

func (s sudoku) PrintBoard() {
	for i := 0; i < 9; i++ {
		if i%3 == 0 {
			fmt.Println("-------------------------")
		}
		for j := 0; j < 9; j++ {
			if j%3 == 0 {
				fmt.Printf("| ")
			}
			fmt.Printf("%d ", s.board[i][j].value)
		}
		fmt.Printf("|\n")
	}
	fmt.Println("-------------------------")
}

// Used for debugging
func (s sudoku) printBlock() {
	for i := 0; i < 9; i++ {
		if i%3 == 0 {
			fmt.Println("")
		}
		for j := 0; j < 9; j++ {
			if j%3 == 0 {
				fmt.Printf("| ")
			}
			var form [9]int
			for index, value := range s.board[i][j].blocked {
				if !value {
					form[index] = index + 1
				}
			}
			for _, value := range form {
				fmt.Print(value)
			}
			fmt.Print(" ")
		}
		fmt.Print("|\n")
	}
	fmt.Println()
}

//InitializeBoard creates the initial game board
func InitializeBoard(path string) sudoku {
	Sudoku := sudoku{readBoard(path), 0}
	return Sudoku
}

func readBoard(path string) [9][9]cell {
	var board [9][9]cell
	var noBlock [9]bool
	dat, err := ioutil.ReadFile(path)
	check(err)

	//Split the input file on commas and newlines
	f := func(c rune) bool {
		return c == '\n' || c == ','
	}
	list := strings.FieldsFunc(string(dat), f)

	var value int
	for i, c := 0, 0; i < 9; i++ {
		for j := 0; j < 9; j, c = j+1, c+1 {
			value, err = strconv.Atoi(list[c])
			board[i][j] = cell{noBlock, value}
		}
	}
	return board
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

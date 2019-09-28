package solver

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var allBlock = [9]bool{true, true, true, true, true, true, true, true, true}

// Cell represents a single cell/square on the sudoku board
type Cell struct {
	blocked [9]bool //Which numbers can't be placed
	Value   int8    //Value represents the digit Value of a square
	Unsure  bool    //Did we guess this value or are we cetrain of its correctness
}

// Sudoku struct to represent the game board
type Sudoku struct {
	board  [9][9]Cell
	filled int //How many total squares filled
}

// Begin will read in a board and get things ready and start solving.
// if the gui argument is true then a channel will be initialized which
// feeds the gui with board states to display
func (s *Sudoku) Begin(path string, feeder chan<- [9][9]Cell) {
	s.board = readBoard(path)
	s.FirstSweep(feeder)
}

// if the cell doesn't have a value and all values are blocked then something is wrong
func sanity(c Cell) bool {
	if c.Value == 0 {
		for _, block := range c.blocked {
			if !block {
				return true
			}
		}
		return false
	}
	return true
}

func (s *Sudoku) checkCell(i int, j int) (change bool) {
	// Check if there is only one available Value
	if val := s.checkOneVal(i, j); val != 0 {
		s.explode(i, j, val, false)
		change = true
	}
	// Check row
	if val := s.checkRow(i, j); val != 0 {
		s.explode(i, j, val, false)
		change = true
	}
	// Check column
	if val := s.checkCol(i, j); val != 0 {
		s.explode(i, j, val, false)
		change = true
	}
	// Check 3x3 chunk
	if val := s.checkChunk(i, j); val != 0 {
		s.explode(i, j, val, false)
		change = true
	}
	return
}

//Solve attempts to solve the board, cannot make guesses
func (s Sudoku) Solve(feeder chan<- [9][9]Cell) bool {
	for s.filled < 81 {
		change := false
		for i := 0; i < 9; i++ {
			for j := 0; j < 9; j++ {
				// fmt.Println(i, j, s.board[0][1].blocked[1])
				if !sanity(s.board[i][j]) {
					fmt.Println("Shit hit the fan!")
					return false
				} else if s.isFree(i, j) {
					if feeder != nil && s.checkCell(i, j) {
						feeder <- s.board
						change = true
					}
				}
			}
		}
		if !change {
			//Todo implementing guessing
			fmt.Println("Making a guess")
			s.guess(feeder)

			// if s.Solve()
		}
	}
	return true
}

func (s *Sudoku) guess(feeder chan<- [9][9]Cell) {
	backup := makeBackup(*s)

	//Attempt to make high probability guesses by prioritizing cells with few options
	for standards := 2; standards < 9; standards++ {
		for i := 0; i < 9; i++ {
			for j := 0; j < 9; j++ {
				if s.isFree(i, j) && countOptions(s.board[i][j]) == standards {
					for nr, block := range s.board[i][j].blocked {
						if block == false {
							//make a guess and attempt to solve the board
							//if we can't we can at least rule out the guessed number
							backup.explode(i, j, int8(nr+1), true)
							if backup.Solve(feeder) {
								s.explode(i, j, int8(nr+1), false)
							} else {
								s.board[i][j].blocked[nr] = true
								s.checkCell(i, j)
							}
							return
						}
					}
				}
			}
		}
	}
}

func countOptions(c Cell) (options int) {
	for _, block := range c.blocked {
		if !block {
			options++
		}
	}
	return
}

func makeBackup(s Sudoku) Sudoku {
	bac := new(Sudoku)
	bac.filled = s.filled
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			bac.board[i][j] = s.board[i][j]
		}
	}
	return *bac
}

//Check if this Cell has only one possible Value left
func (s *Sudoku) checkOneVal(x int, y int) int8 {
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
	return int8(last + 1)
}

//Check if I'm the only one in my row who can receive a specific Value
func (s *Sudoku) checkRow(x int, y int) int8 {
	var success bool
	for v, block := range s.board[x][y].blocked {
		if !block {
			success = true
			for i := 0; i < 9; i++ {
				if !s.board[x][i].blocked[v] && i != y {
					//Another Cell can have this Value
					success = false
					break
				}
			}
			if success {
				return int8(v + 1)
			}
		}
	}
	return 0
}

//Check if I'm the only one in my column who can receive a specific Value
func (s *Sudoku) checkCol(x int, y int) int8 {
	var success bool
	for v, block := range s.board[x][y].blocked {
		if !block {
			success = true
			for i := 0; i < 9; i++ {
				if !s.board[i][y].blocked[v] && i != x {
					//Another Cell can have this Value
					success = false
					break
				}
			}
			if success {
				return int8(v + 1)
			}
		}
	}
	return 0
}

// FirstSweep is run before solving. It sets some required internal variables
// based on the initial Values in the board
func (s *Sudoku) FirstSweep(feeder chan<- [9][9]Cell) {
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if curr := s.board[i][j].Value; curr != 0 {
				s.explode(i, j, curr, false)
			}
		}
	}
	if feeder != nil {
		feeder <- s.board
	}
}

//Check if I'm the only one in my 3x3 chunk who can receive a specific Value
func (s *Sudoku) checkChunk(x int, y int) int8 {
	xn := x - (x % 3)
	yn := y - (y % 3)
	var success bool
	for v, block := range s.board[x][y].blocked {
		if !block {
			success = true
			for i := xn; i < xn+3; i++ {
				for j := yn; j < yn+3; j++ {
					if !s.board[i][j].blocked[v] && (i != x || j != y) {
						//Another Cell can have this Value
						success = false
						break
					}
				}
			}
			if success {
				return int8(v + 1)
			}
		}
	}
	return 0
}

// "Explode" means to spread the influence of a Value to restrict other Cells
// from taking that Value
func (s *Sudoku) explode(x int, y int, val int8, unsure bool) {
	s.board[x][y].blocked = allBlock
	s.board[x][y].Unsure = unsure
	s.filled++
	s.board[x][y].Value = val

	val-- //indexing is 0-8 while Sudoku is 1-9
	s.explodeChunk(x, y, val)
	s.explodeVertical(x, val)
	s.explodeHorizontal(y, val)

	// Send a board state to the GUI to draw
	// if s.Feeder != nil {
	// 	s.Feeder <- s.board
	// }
}

//Remve possibilities in the small 3x3 chunks
func (s *Sudoku) explodeChunk(x int, y int, val int8) {
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

func (s *Sudoku) explodeVertical(x int, val int8) {
	for i := 0; i < 9; i++ {
		if s.board[x][i].blocked[val] == false {
			s.board[x][i].blocked[val] = true
		}
	}
}

func (s *Sudoku) explodeHorizontal(y int, val int8) {
	for i := 0; i < 9; i++ {
		if s.board[i][y].blocked[val] == false {
			s.board[i][y].blocked[val] = true
		}
	}
}

func (s *Sudoku) isFree(x int, y int) bool {
	return s.board[x][y].Value == 0
}

// PrintBoard prints the sudoku board to the console
func (s Sudoku) PrintBoard() {
	for i := 0; i < 9; i++ {
		if i%3 == 0 {
			fmt.Println("-------------------------")
		}
		for j := 0; j < 9; j++ {
			if j%3 == 0 {
				fmt.Printf("| ")
			}
			fmt.Printf("%d ", s.board[i][j].Value)
		}
		fmt.Printf("|\n")
	}
	fmt.Println("-------------------------")
}

// Used for debugging
func (s Sudoku) printBlock() {
	for i := 0; i < 9; i++ {
		if i%3 == 0 {
			fmt.Println("")
		}
		for j := 0; j < 9; j++ {
			if j%3 == 0 {
				fmt.Printf("| ")
			}
			var form [9]int
			for index, Value := range s.board[i][j].blocked {
				if !Value {
					form[index] = index + 1
				}
			}
			for _, Value := range form {
				fmt.Print(Value)
			}
			fmt.Print(" ")
		}
		fmt.Print("|\n")
	}
	fmt.Println()
}

func readBoard(path string) [9][9]Cell {
	var board [9][9]Cell
	var noBlock [9]bool
	dat, err := ioutil.ReadFile(path)
	check(err)

	//Split the input file on commas and newlines
	f := func(c rune) bool {
		return c == '\n' || c == ','
	}
	list := strings.FieldsFunc(string(dat), f)

	var Value int
	for i, c := 0, 0; i < 9; i++ {
		for j := 0; j < 9; j, c = j+1, c+1 {
			Value, err = strconv.Atoi(list[c])
			board[i][j] = Cell{noBlock, int8(Value), false}
		}
	}
	return board
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

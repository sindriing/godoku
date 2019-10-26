package solver

import (
	"testing"
)

// TestConnection tests the connection to sudoku.com
func TestConnection(t *testing.T) {
	ans, err := GetLevelOnline("hard")
	if err != nil {
		panic(err)
	}
	t.Log("Hello There Obi-Wan-Kenoby!")
	t.Log(ans)
}

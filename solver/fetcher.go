package solver

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

type sudokudata struct {
	// Desc [][5]interface{} `json:"desc"`
	Desc [5]interface{} `json:"desc"`
}

// GetLevelOnline sends a request to an api of sudoku.com for a board of a specified difficulty
func GetLevelOnline(diff string) (string, error) {
	diff = strings.ToLower(diff)
	if diff != "expert" && diff != "hard" && diff != "medium" && diff != "easy" {
		differr := errors.New("Incorrect difficulty! please select easy, medium, hard or expert")
		return "", differr
	}
	resp, err := http.Get("https://sudoku.com/api/getLevel/" + diff)
	if err != nil {
		return "", err
	}

	var b sudokudata
	body, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &b)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	rawboard := b.Desc[0].(string)
	return rawboard, err
}

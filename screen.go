package main

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"time"

	_ "image/png"

	"github.com/sindriing/godoku/solver"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Sindri's Sudoku Solver",
		Bounds: pixel.R(0, 0, 410, 410),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	spritesheet, err := loadPicture("resources/digitsPalette.png")
	if err != nil {
		panic(err)
	}

	batch := pixel.NewBatch(&pixel.TrianglesData{}, spritesheet)

	var digitFrames []pixel.Rect
	for x := spritesheet.Bounds().Min.X; x < spritesheet.Bounds().Max.X; x += 27 {
		digitFrames = append(digitFrames, pixel.R(x, 0, x+27, 40))
	}

	var (
		frames = 0
		second = time.Tick(time.Second)
	)
	var board [9][9]solver.Cell

	feeder := make(chan [9][9]solver.Cell)

	//Initialize the board

	go func(path string, feeder chan [9][9]solver.Cell) {
		Solver := new(solver.Sudoku)
		Solver.Begin(path, feeder)

		//Solve the board
		Solver.Solve(feeder)

	}("presets/"+os.Args[1], feeder)

	//Drawing the gridlines
	var thickness float64
	var x, y float64
	line := imdraw.New(nil)
	line.Color = color.RGBA{30, 30, 30, 255}
	for i := 0.0; i < 10; i++ {
		x = 3 + 45*i
		if int(i)%3 == 0 {
			thickness = 7
		} else {
			thickness = 3
		}
		line.Push(pixel.V(x, 0), pixel.V(x, 420))
		line.Line(thickness)
		line.Push(pixel.V(0, x), pixel.V(420, x))
		line.Line(thickness)
	}

	for !win.Closed() {
		if win.JustPressed(pixelgl.MouseButtonLeft) || win.Pressed(pixelgl.KeyRight) {
			select {
			case board = <-feeder:
			default:
			}
		}

		win.Clear(color.RGBA{215, 215, 215, 255})
		batch.Clear()

		//Drawing the Digits
		for i := 0.0; i < 9; i++ {
			for j := 0.0; j < 9; j++ {
				x = 25 + 45*i
				y = 25 + 45*j
				digit := pixel.NewSprite(spritesheet, digitFrames[board[int(i)][int(j)].Value])
				digit.Draw(batch, pixel.IM.Moved(pixel.V(x, y)))
				if board[int(i)][int(j)].Unsure == true {
					highlight := imdraw.New(nil)
					highlight.Color = color.RGBA{230, 110, 110, 150}
					highlight.Push(pixel.V(x, y))
					highlight.Circle(20, 0)
					highlight.Draw(win)
				}
			}
		}
		batch.Draw(win)
		line.Draw(win)
		win.Update()

		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0
		default:
		}
	}
}

func main() {
	pixelgl.Run(run)
}

package main

import (
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const (
	screenWidth  = 400.0
	screenHeight = 300.0
)

type Coord struct {
	X, Y float64
}

const (
	modeTitle = iota
	modeGame
	modeGameOver
)

var (
	err error

	road  *ebiten.Image
	wRoad float64
	hRoad float64

	finish  *ebiten.Image
	wFinish float64
	hFinish float64
	yFinish int

	car  *ebiten.Image
	wCar float64
	hCar float64
	xCar float64
	yCar float64

	y1Car float64 = -500

	ytotal = 50
	yoff   = 50

	mode int
)

// 88 to 248
var cars []*Coord = []*Coord{
	&Coord{Y: -500, X: 90},
	&Coord{Y: -550, X: 150},

	&Coord{Y: -750, X: 100},
	&Coord{Y: -750, X: 200},

	&Coord{Y: -1000, X: 150},
	&Coord{Y: -1000, X: 248},

	&Coord{Y: -1200, X: 130},
	&Coord{Y: -1200, X: 240},

	&Coord{Y: -1500, X: 100},
	&Coord{Y: -1500, X: 150},
	&Coord{Y: -1500, X: 180},
}

func init() {
	{
		road, _, err = ebitenutil.NewImageFromFile("images/background/road.png", ebiten.FilterDefault)
		if err != nil {
			log.Fatal(err)
		}
		w, h := road.Size()
		wRoad = float64(w)
		hRoad = float64(h)
	}
	{
		car, _, err = ebitenutil.NewImageFromFile("images/sprites/topcar.png", ebiten.FilterDefault)
		if err != nil {
			log.Fatal(err)
		}
		w, h := car.Size()
		wCar = float64(w)
		hCar = float64(h)
		xCar = (screenWidth / 2) - (wCar / 2)
		yCar = screenHeight - hCar - 5
	}
	{
		finish, _, err = ebitenutil.NewImageFromFile("images/background/finish.png", ebiten.FilterDefault)
		if err != nil {
			log.Fatal(err)
		}
		w, h := finish.Size()
		wFinish = float64(w)
		hFinish = float64(h)
		yFinish = 2000
	}
	mode = modeGame
}

func update(screen *ebiten.Image) error {
	switch mode {
	case modeGame:
		ytotal += 5
		yoff = ytotal % screenHeight
		if ebiten.IsKeyPressed(ebiten.KeyLeft) {
			xCar -= 5
			if xCar < 88 {
				xCar = 88
			}
		}
		if ebiten.IsKeyPressed(ebiten.KeyRight) {
			xCar += 5
			if xCar > 248 {
				xCar = 248
			}
		}

		y1Car += 5
		for _, c := range cars {
			c.Y += 5
		}

		if ytotal > yFinish {
			mode = modeGameOver
		}

		if ebiten.IsDrawingSkipped() {
			return nil
		}

		// Top
		opRoad := &ebiten.DrawImageOptions{}
		opRoad.GeoM.Scale(screenWidth/wRoad, screenHeight/hRoad)
		opRoad.GeoM.Translate(0, float64(yoff))
		screen.DrawImage(road, opRoad)

		// Bottom
		opRoadBtm := &ebiten.DrawImageOptions{}
		opRoadBtm.GeoM.Scale(screenWidth/wRoad, screenHeight/hRoad)
		opRoadBtm.GeoM.Translate(0, float64(yoff-screenHeight))
		screen.DrawImage(road, opRoadBtm)

		// Car
		{
			opCar := &ebiten.DrawImageOptions{}
			opCar.GeoM.Translate(xCar, yCar)
			screen.DrawImage(car, opCar)
		}

		for _, c := range cars {
			opCar := &ebiten.DrawImageOptions{}
			opCar.GeoM.Translate(c.X, c.Y)
			screen.DrawImage(car, opCar)
		}

		ebitenutil.DebugPrint(screen, fmt.Sprintf("%02f, ytotal: %d", xCar, ytotal))
	case modeGameOver:
		// Top
		opRoad := &ebiten.DrawImageOptions{}
		opRoad.GeoM.Scale(screenWidth/wRoad, screenHeight/hRoad)
		opRoad.GeoM.Translate(0, float64(yoff))
		screen.DrawImage(road, opRoad)

		// Bottom
		opRoadBtm := &ebiten.DrawImageOptions{}
		opRoadBtm.GeoM.Scale(screenWidth/wRoad, screenHeight/hRoad)
		opRoadBtm.GeoM.Translate(0, float64(yoff-screenHeight))
		screen.DrawImage(road, opRoadBtm)

		// Car
		{
			opCar := &ebiten.DrawImageOptions{}
			opCar.GeoM.Translate(xCar, yCar)
			screen.DrawImage(car, opCar)
		}

		ebitenutil.DebugPrint(screen, fmt.Sprintf("%02f, ytotal: %d", xCar, ytotal))

		cmd := exec.Command("./eject.exe")
		err = cmd.Run()
		os.Exit(0)
	}

	return nil
}

func main() {
	if err := ebiten.Run(update, screenWidth, screenHeight, 2, "Infinite Scroll (Ebiten Demo)"); err != nil {
		log.Fatal(err)
	}
}

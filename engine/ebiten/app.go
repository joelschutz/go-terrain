package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/png"
	_ "image/png"
	"log"
	"math/rand"
	"os"
	"strconv"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/joelschutz/go-terrain/pkg"
)

const (
	screenWidth  = 512
	screenHeight = 512
	genSize      = 10000
	maxGenCount  = 10
)

var modes = [3]string{
	"Normals", "HeightMap", "Position",
}

type Game struct {
	frames   []image.Image
	hm       pkg.HeightMap
	nm       pkg.NormalMap
	rain     pkg.RainSim
	mode     int
	age, gen uint
}

func (g *Game) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		g.rain = append(g.rain, pkg.NewParticle(mgl32.Vec2{float32(x), float32(y)}, mgl32.Vec2{}, 1, 1, 1))
	} else if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		g.mode--
		if g.mode < 0 {
			g.mode = 1
		}
	}

	g.rain.Update(&g.hm, false)
	if len(g.rain) < 1 {
		for i := 0; i < genSize; i++ {
			g.rain = append(g.rain, pkg.NewParticle(mgl32.Vec2{rand.Float32() * screenHeight, rand.Float32() * screenWidth}, mgl32.Vec2{}, 1, 1, 1))
		}
		save(fmt.Sprintf("images/gen%v.png", g.gen), g.hm)
		g.nm = pkg.NewNormalMap(g.hm, 30)
		g.gen++
		if g.gen >= maxGenCount {
			return fmt.Errorf("Stop")
		}
	}
	g.age++
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Select Background
	var bg *ebiten.Image
	switch g.mode {
	case 0:
		// Perlin Noise
		bg = ebiten.NewImageFromImage(g.hm)
	case 1:
		// Normal
		bg = ebiten.NewImageFromImage(g.nm)
	}

	// Draw Backgroung
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(bg, op)

	// Draw Particles
	for i, p := range g.rain {
		screen.Set(int(p.X()), int(p.Y()), color.RGBA{255, 0, 0, 255})
		ebitenutil.DebugPrintAt(screen, strconv.Itoa(i), int(p.X()), int(p.Y()))
	}

	// Draw HUD
	ebitenutil.DebugPrint(screen, modes[g.mode])
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Generation: %v", g.gen), 0, 10)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Particles: %v", len(g.rain)), 0, 20)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%.2f", ebiten.ActualTPS()), 0, 30)

	g.frames = append(g.frames, screen.SubImage(screen.Bounds()))
	return
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func save(filePath string, img image.Image) {
	imgFile, err := os.Create(filePath)
	defer imgFile.Close()
	if err != nil {
		log.Println("Cannot create file:", err)
	}
	png.Encode(imgFile, img)
}

func drawGif(frames []image.Image, filename string) {
	g := gif.GIF{}
	for _, v := range frames {
		fr := image.NewPaletted(frames[0].Bounds(), color.Palette{
			color.Black, color.White,
		})
		draw.Draw(fr, v.Bounds(), v, image.Point{}, draw.Src)
		g.Image = append(g.Image, fr)
		g.Delay = append(g.Delay, 0)
	}
	f, _ := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()
	gif.EncodeAll(f, &g)
}

func EbitenApp() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("My Game")

	g := &Game{}
	g.hm = pkg.NewPerlinMap(512, 30)
	save("hm.png", g.hm)
	g.nm = pkg.NewNormalMap(g.hm, 30)
	for i := 0; i < genSize; i++ {
		g.rain = append(g.rain, pkg.NewParticle(mgl32.Vec2{rand.Float32() * screenHeight, rand.Float32() * screenWidth}, mgl32.Vec2{}, 1, 1, 1))
	}
	save("nm.png", g.nm)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"github.com/gen2brain/raylib-go/physics"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Player struct {
	Body     *physics.Body
	Position rl.Vector2
	Size     rl.Vector2
	Velocity rl.Vector2
}

type Game struct {
	GameOver      bool
	Player        Player
	Camera        rl.Camera2D
	Platform      *physics.Body
	Obstacles     []rl.Rectangle
	PlatformCount int
	Generator     rl.Rectangle
	Score         int
	Restart       bool
}

func NewGame() *Game {
	g := &Game{}
	g.Init()
	return g
}

func (g *Game) Init() {
	physics.Init()
	g.Player.Position = rl.NewVector2(screenWidth/2-500, screenHeight/4*3)
	g.Player.Size = rl.NewVector2(80, 80)

	g.Player.Body = physics.NewBodyRectangle(g.Player.Position, g.Player.Size.X, g.Player.Size.Y, 1)
	g.Player.Body.FreezeOrient = true

	g.Platform = physics.NewBodyRectangle(rl.NewVector2(screenWidth-screenWidth/2+150, screenHeight), screenWidth*3, 150, 1)
	g.Platform.Enabled = false

	g.Generator = rl.NewRectangle(-1260, screenHeight-100, 1, 200)

	g.Camera.Offset = rl.NewVector2(screenWidth/2+50, screenHeight/2)
	g.Camera.Target = rl.NewVector2(screenWidth/2, screenHeight/2)
	g.Camera.Rotation = 0.0
	g.Camera.Zoom = 1

	g.Score = 0
	g.PlatformCount = 1
	g.Restart = false
}

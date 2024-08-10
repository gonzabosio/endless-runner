package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gen2brain/raylib-go/physics"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth   = 1280
	screenHeight  = 720
	velocity      = 0.5
	platfDisplace = 12
)

var (
	canJump       = true
	genCollision  bool
	generatingMap bool
	platf2        *physics.Body
)

func main() {
	game := NewGame()
	game.GameOver = true

	rl.InitWindow(screenWidth, screenHeight, "Go-Nowhere")
	defer rl.CloseWindow()

	rl.SetTargetFPS(75)

	for !rl.WindowShouldClose() {
		if game.GameOver {
			rl.BeginDrawing()
			rl.ClearBackground(rl.Black)
			rl.BeginMode2D(game.Camera)
			rl.DrawText("GO NO-WHERE", screenWidth/2-rl.MeasureText("GO NO-WHERE", 40)+80, screenHeight/2-100, 40, rl.Blue)
			rl.DrawText("Press ENTER to start", screenWidth/2-rl.MeasureText("Press ENTER to start", 24)+75, screenHeight/2, 24, rl.DarkBlue)
			if game.Score > 0 {
				rl.DrawText(fmt.Sprintf("Last Score: %v", game.Score), screenWidth/2-rl.MeasureText("Score: 0000", 24), screenHeight/2+100, 18, rl.Blue)
			}
			if rl.IsKeyPressed(rl.KeyEnter) {
				if game.Restart {
					time.Sleep(1 * time.Second)
					game = NewGame()
					game.Draw()
				}
				game.GameOver = false
			}
			rl.EndDrawing()
		} else {
			game.Draw()
		}
	}
	physics.Close()
}

func (g *Game) Draw() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)
	rl.BeginMode2D(g.Camera)
	physics.Update()

	g.Platform.Position.X -= platfDisplace

	if rl.IsKeyPressed(rl.KeyUp) && canJump {
		g.Player.Body.Velocity.Y = -velocity * 5
		go func() {
			canJump = false
			time.Sleep(500 * time.Millisecond)
			canJump = true
		}()
	}
	if rl.IsKeyDown(rl.KeyDown) && canJump {
		g.Player.Body.Position.Y += 40
	}

	// Draw physics
	for i, body := range physics.GetBodies() {
		vertexCount := physics.GetShapeVerticesCount(i)
		for j := 0; j < vertexCount; j++ {
			vertexA := body.GetShapeVertex(j)
			jj := 0
			if j+1 < vertexCount {
				jj = j + 1
			}
			vertexB := body.GetShapeVertex(jj)
			rl.DrawLineV(vertexA, vertexB, rl.Blue)
		}
	}

	//Map generation
	rl.DrawRectangleRec(g.Generator, rl.Blank)
	switch g.PlatformCount {
	case 1:
		genCollision = rl.CheckCollisionRecs(rl.NewRectangle(g.Platform.Position.X*2, g.Platform.Position.Y, screenWidth, 150), g.Generator)
	case 2:
		genCollision = rl.CheckCollisionRecs(rl.NewRectangle(platf2.Position.X*2, platf2.Position.Y, screenWidth, 150), g.Generator)
	}
	if genCollision && !generatingMap {
		generatingMap = true
		if g.PlatformCount == 1 {
			g.PlatformCount++
			go func() {
				platf2 = physics.NewBodyRectangle(rl.NewVector2(screenWidth*2+screenWidth/2+150, screenHeight), screenWidth*3, 150, 1)
				platf2.Enabled = false
				time.Sleep(1 * time.Second)
				g.Obstacles = []rl.Rectangle{}
				g.CreateObstacles()
				time.Sleep(2 * time.Second)
				g.Platform.Destroy()
				generatingMap = false
			}()
		} else {
			g.PlatformCount--
			go func() {
				g.Platform = physics.NewBodyRectangle(rl.NewVector2(screenWidth*2+screenWidth/2+150, screenHeight), screenWidth*3, 150, 1)
				g.Platform.Enabled = false
				time.Sleep(1 * time.Second)
				g.Obstacles = []rl.Rectangle{}
				g.CreateObstacles()
				time.Sleep(2 * time.Second)
				platf2.Destroy()
				generatingMap = false
			}()
		}
	}
	if platf2 != nil {
		platf2.Position.X -= platfDisplace
	}

	//Obstacles movement
	if len(g.Obstacles) == 6 {
		for _, rec := range g.Obstacles {
			rl.DrawRectangleRec(rec, rl.Red)
		}
		g.Obstacles[0].X -= platfDisplace
		g.Obstacles[1].X -= platfDisplace
		g.Obstacles[2].X -= platfDisplace
		g.Obstacles[3].X -= platfDisplace
		g.Obstacles[4].X -= platfDisplace
		g.Obstacles[5].X -= platfDisplace

	}

	if g.Player.Body.Position.Y > 800 {
		physics.Close()
		g.Restart = true
		g.GameOver = true
	}
	for _, rec := range g.Obstacles {
		if rl.CheckCollisionRecs(rl.NewRectangle(g.Player.Body.Position.X-40, g.Player.Body.Position.Y-40, 80, 40), rec) {
			physics.Close()
			g.Restart = true
			g.GameOver = true
		}
	}
	g.Score += 1
	rl.DrawText(fmt.Sprintf("SCORE: %v", g.Score), -45, 5, 20, rl.Green)
	rl.DrawText("Arrow Up to Jump", -45, 28, 14, rl.Gray)
	rl.DrawText("Arrow Down to Crouch", -45, 46, 14, rl.Gray)
	rl.DrawFPS(screenWidth-130, 0)
	rl.EndDrawing()
}

func (g *Game) CreateObstacles() {
	var spacer int32 = rl.GetRandomValue(270, 370)

	for i := 0; i < 6; i++ {
		up := rand.Int31n(2)
		x := 400 + int32(spacer)
		if up == 1 {
			up = rand.Int31n(2)
			if up == 1 {
				g.Obstacles = append(g.Obstacles, rl.NewRectangle(float32(x+spacer), screenHeight-250, 110, 60))
			} else {
				g.Obstacles = append(g.Obstacles, rl.NewRectangle(float32(x+spacer), screenHeight-190, 110, 60))
			}
		} else {
			up = rand.Int31n(2)
			if up == 1 {
				g.Obstacles = append(g.Obstacles, rl.NewRectangle(float32(x+spacer), screenHeight-205, 80, 130))
			} else {
				g.Obstacles = append(g.Obstacles, rl.NewRectangle(float32(x+spacer), screenHeight-185, 60, 110))
			}
		}
		spacer += rl.GetRandomValue(200, 300)
	}
}

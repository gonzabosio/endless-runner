package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gen2brain/raylib-go/physics"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const screenWidth = 1280
const screenHeight = 720

const velocity = 0.5
const platfDisplace = 10

var canJump = true

var genCollision bool
var generatingMap bool
var platfCounter = 1
var platf2 *physics.Body

var obstacles []rl.Rectangle
var score int64 = 0

func main() {
	rl.InitWindow(screenWidth, screenHeight, "Go-Nowhere")
	defer rl.CloseWindow()

	physics.Init()
	player := physics.NewBodyRectangle(rl.NewVector2(screenWidth/12, screenHeight/4*3), 80, 80, 1)
	player.FreezeOrient = true

	platf := physics.NewBodyRectangle(rl.NewVector2(screenWidth-screenWidth/2+150, screenHeight), screenWidth*3, 150, 1)
	platf.Enabled = false

	generator := rl.NewRectangle(-1260, screenHeight-100, 1, 200)

	rl.SetTargetFPS(75)
	camera := rl.Camera2D{}
	camera.Offset = rl.Vector2{X: screenWidth/2 + 50, Y: screenHeight / 2}
	camera.Rotation = 0.0
	camera.Zoom = 1

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		rl.BeginMode2D(camera)
		camera.Target = rl.NewVector2(screenWidth/2, screenHeight/2)

		physics.Update()
		platf.Position.X -= platfDisplace
		if rl.IsKeyPressed(rl.KeyUp) && canJump {
			player.Velocity.Y = -velocity * 5
			go func() {
				canJump = false
				time.Sleep(500 * time.Millisecond)
				canJump = true
			}()
		}
		if rl.IsKeyDown(rl.KeyDown) && canJump {
			player.Position.Y += 40
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
		rl.DrawRectangleRec(generator, rl.White)
		switch platfCounter {
		case 1:
			genCollision = rl.CheckCollisionRecs(rl.NewRectangle(platf.Position.X*2, platf.Position.Y, screenWidth, 150), generator)
		case 2:
			genCollision = rl.CheckCollisionRecs(rl.NewRectangle(platf2.Position.X*2, platf2.Position.Y, screenWidth, 150), generator)
		}
		if genCollision && !generatingMap {
			generatingMap = true
			if platfCounter == 1 {
				platfCounter++
				go func() {
					platf2 = physics.NewBodyRectangle(rl.NewVector2(screenWidth*2+screenWidth/2+150, screenHeight), screenWidth*3, 150, 1)
					platf2.Enabled = false
					time.Sleep(2 * time.Second)
					obstacles = []rl.Rectangle{}
					createObstacles()
					time.Sleep(2 * time.Second)
					platf.Destroy()
					generatingMap = false
				}()
			} else {
				platfCounter--
				go func() {
					platf = physics.NewBodyRectangle(rl.NewVector2(screenWidth*2+screenWidth/2+150, screenHeight), screenWidth*3, 150, 1)
					platf.Enabled = false
					time.Sleep(2 * time.Second)
					obstacles = []rl.Rectangle{}
					createObstacles()
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
		if len(obstacles) != 0 {
			for _, rec := range obstacles {
				rl.DrawRectangleRec(rec, rl.Red)
			}
			obstacles[0].X -= platfDisplace
			obstacles[1].X -= platfDisplace
			obstacles[2].X -= platfDisplace
			obstacles[3].X -= platfDisplace
			obstacles[4].X -= platfDisplace
		}
		score += 1
		rl.DrawText(fmt.Sprintf("SCORE: %v", score), -45, 5, 20, rl.Green)

		for _, rec := range obstacles {
			if rl.CheckCollisionRecs(rl.NewRectangle(player.Position.X-40, player.Position.Y-40, 80, 40), rec) {
				score = 0
			}
		}
		// rl.DrawFPS(-50, 0)
		rl.EndDrawing()
	}
	physics.Close()
}

func createObstacles() {
	var spacer int32 = rl.GetRandomValue(200, 250)

	for i := 0; i < 5; i++ {
		up := rand.Int31n(2)
		x := 400 + int32(spacer)
		if up == 1 {
			obstacles = append(obstacles, rl.NewRectangle(float32(x+spacer), screenHeight-190, 90, 50))
		} else {
			obstacles = append(obstacles, rl.NewRectangle(float32(x+spacer), screenHeight-175, 50, 100))
		}
		spacer += rl.GetRandomValue(180, 300)
	}
}

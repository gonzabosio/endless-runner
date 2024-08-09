package main

import (
	"fmt"
	"time"

	"github.com/gen2brain/raylib-go/physics"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const screenWidth = 1280
const screenHeight = 720

const velocity = 0.5
const displacePlatf = 8

var canJump = true

var collision bool
var generatingMap bool
var platfCounter = 1

func main() {
	rl.InitWindow(screenWidth, screenHeight, "Go-Nowhere")
	defer rl.CloseWindow()

	physics.Init()
	player := physics.NewBodyRectangle(rl.NewVector2(screenWidth/12, screenHeight/4*3), 40, 40, 1)

	platf := physics.NewBodyRectangle(rl.NewVector2(screenWidth-screenWidth/2, screenHeight), screenWidth*2, 150, 1)
	platf.Enabled = false

	var platf2 *physics.Body
	generator := rl.NewRectangle(-800, screenHeight-100, 1, 200)

	rl.SetTargetFPS(75)
	camera := rl.Camera2D{}
	camera.Offset = rl.Vector2{X: screenWidth / 2, Y: screenHeight / 2}
	camera.Rotation = 0.0
	camera.Zoom = 0.24

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		rl.BeginMode2D(camera)
		camera.Target = rl.NewVector2(screenWidth/2, screenHeight/2)

		physics.Update()

		platf.Position.X -= displacePlatf
		if rl.IsKeyPressed(rl.KeyUp) && canJump {
			player.Velocity.Y = -velocity * 5
			go func() {
				canJump = false
				time.Sleep(500 * time.Millisecond)
				canJump = true
			}()
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
			collision = rl.CheckCollisionRecs(rl.NewRectangle(platf.Position.X*2, platf.Position.Y, screenWidth, 150), generator)
		case 2:
			collision = rl.CheckCollisionRecs(rl.NewRectangle(platf2.Position.X*2, platf2.Position.Y, screenWidth, 150), generator)
		}
		if collision && !generatingMap {
			generatingMap = true
			fmt.Println("COLLISION in Platform: ", platfCounter)
			if platfCounter == 1 {
				platfCounter++
				go func() {
					platf2 = physics.NewBodyRectangle(rl.NewVector2(screenWidth-screenWidth/2, screenHeight), screenWidth*2, 150, 1)
					platf2.Enabled = false
					platf2.Position.X -= displacePlatf
					time.Sleep(2 * time.Second)
					platf.Destroy()
					generatingMap = false
				}()
			} else {
				platfCounter--
				go func() {
					platf = physics.NewBodyRectangle(rl.NewVector2(screenWidth-screenWidth/2, screenHeight), screenWidth*2, 150, 1)
					platf.Enabled = false
					platf.Position.X -= displacePlatf
					time.Sleep(2 * time.Second)
					platf2.Destroy()
					generatingMap = false
				}()
			}

		}
		if platf2 != nil {
			platf2.Position.X -= displacePlatf
		}
		rl.DrawFPS(0, 0)
		rl.EndDrawing()
	}
	physics.Close()
}

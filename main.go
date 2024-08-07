package main

import (
	"time"

	"github.com/gen2brain/raylib-go/physics"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const screenWidth = 1280
const screenHeight = 720

const velocity = 0.5

var canJump = true

func main() {
	rl.InitWindow(screenWidth, screenHeight, "Go-Nowhere")
	defer rl.CloseWindow()

	physics.Init()
	body := physics.NewBodyRectangle(rl.NewVector2(screenWidth/12, screenHeight/4*3), 40, 40, 1)
	body.FreezeOrient = true

	floor := physics.NewBodyRectangle(rl.NewVector2(screenWidth/2, screenHeight), screenWidth, 150, 1)
	floor.Enabled = false

	rl.SetTargetFPS(75)
	camera := rl.Camera2D{}
	camera.Offset = rl.Vector2{X: screenWidth / 2, Y: screenHeight / 2}
	camera.Rotation = 0.0
	camera.Zoom = 1.0

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		rl.BeginMode2D(camera)
		camera.Target = rl.NewVector2(body.Position.X, screenHeight/2)

		physics.Update()

		body.Velocity.X = +velocity
		if rl.IsKeyPressed(rl.KeyUp) && canJump {
			body.Velocity.Y = -velocity * 5

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
		rl.DrawFPS(int32(body.Position.X-screenWidth/2), 5)
		rl.EndDrawing()
	}
	physics.Close()
}

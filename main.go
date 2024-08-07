package main

import (
	"github.com/gen2brain/raylib-go/physics"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const screenWidth = 1280
const screenHeight = 720

const velocity = 0.5

func main() {
	rl.InitWindow(screenWidth, screenHeight, "Go-Nowhere")
	defer rl.CloseWindow()

	physics.Init()
	body := physics.NewBodyRectangle(rl.NewVector2(screenWidth/12, screenHeight/4*3), 40, 40, 1)
	body.FreezeOrient = true

	floor := physics.NewBodyRectangle(rl.NewVector2(screenWidth/2, screenHeight), screenWidth, 150, 1)
	floor.Enabled = false

	rl.SetTargetFPS(75)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		physics.Update()

		body.Velocity.X = +velocity
		if rl.IsKeyPressed(rl.KeyUp) {
			body.Velocity.Y = -velocity * 5
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

		rl.DrawFPS(0, 0)

		rl.EndDrawing()
	}
	physics.Close()
}

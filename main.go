package main

import rl "github.com/gen2brain/raylib-go/raylib"

func main() {
	rl.InitWindow(1280, 720, "Go-Nowhere")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)
		rl.EndDrawing()
	}
}

package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = int32(800)
	screenHeight = int32(450)
	gridSize     = int32(15)
	second       = int32(6000) //To keep track of scale

	timeScale = 1000  //To slow down the simulation
	g         = -9.89 //Standard gravity
)

var (
	rX = float32(24.0)
	rY = float32(35.0)
	rZ = float32(0.0)

	time = float32(0.0)
)

func renderStats() {
	rl.DrawText("T: "+fmt.Sprintf("%.2f", time)+" s", 10, 40, 20, rl.Gray)
	rl.DrawFPS(10, 10)

	if sim_started {
		rl.DrawText("Simulation Started", 10, 70, 20, rl.NewColor(255, 255, 255, 50))
	}

}

func main() {
	fmt.Println("AIM Rocket Started")

	//Window Set-up
	rl.InitWindow(screenWidth, screenHeight, "AIM: Autonomous Impact Management")
	rl.SetTargetFPS(60)

	//Camera setup
	c := rl.Camera3D{}
	c.Position = rl.NewVector3(0.0, 50.0, 35.0)
	c.Target = rl.NewVector3(0.0, 0.0, 0.0)
	c.Up = rl.NewVector3(0.0, 1.0, 0.0)
	c.Fovy = 45.0
	rl.SetCameraMode(c, rl.CameraFree)

	//Rocket Setup
	rocket := Rocket{
		color:        rl.NewColor(128, 128, 200, 255),
		texture:      rl.LoadTexture("black.png"),
		model:        rl.LoadModel("./rocketshape.obj"),
		position:     rl.NewVector3(rX, rY, rZ),
		scale:        rl.NewVector3(1.0, 1.0, 1.0),
		rotationAxis: rl.NewVector3(1.0, 1.0, 1.0),

		//Based on SpaceX's Falcon 9
		fuel: int32(1),
		Ve:   int32(2),
		Mo:   int32(2),
		Mf:   int32(1),
	}
	rocket.setup()

	for !rl.WindowShouldClose() {

		// Update camera
		rl.UpdateCamera(&c)

		rl.BeginDrawing()
		rl.BeginMode3D(c)
		rl.ClearBackground(rl.NewColor(28, 28, 28, 255))

		//Grid
		rl.DrawGrid(gridSize, 1.0)

		rocket.update() //Recalculate stuff
		rocket.draw()   //Redraw the rocket

		//Mount Pad
		rl.DrawCylinderWires(rl.NewVector3(0.0, 0.0, 0.0), 2.0, 2.0, 0.01, 20, rl.NewColor(255, 255, 255, 105))
		rl.DrawCylinder(rl.NewVector3(0.0, 0.0, 0.0), 2.0, 2.0, 0.01, 20, rl.NewColor(255, 255, 255, 50))

		rl.DrawCylinderWires(rl.NewVector3(0.0, 10.0, 0.0), 2.0, 2.0, 0.01, 20, rl.NewColor(255, 255, 255, 105))

		//Bounding thrust zone
		// rl.DrawLine3D(rl.NewVector3(-2.0, 10.0, 2.0), rl.NewVector3(rX, rY, rZ), rl.NewColor(155, 155, 255, 55))
		// rl.DrawLine3D(rl.NewVector3(2.0, 10.0, -2.0), rl.NewVector3(rX, rY, rZ), rl.NewColor(155, 155, 255, 55))
		// rl.DrawLine3D(rl.NewVector3(2.0, 10.0, 2.0), rl.NewVector3(rX, rY, rZ), rl.NewColor(155, 155, 255, 55))
		// rl.DrawLine3D(rl.NewVector3(-2.0, 10.0, -2.0), rl.NewVector3(rX, rY, rZ), rl.NewColor(155, 155, 255, 55))

		rl.EndMode3D()

		/* 2D MODE */

		//Drawing stats
		renderStats()
		rocket.drawInfo()
		//Time management (lol)
		if sim_started {
			time += 0.001
		}

		rl.EndDrawing()

	}
}

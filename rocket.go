package main

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

//Rocket ...
type Rocket struct {

	//How it looks
	color   rl.Color
	model   rl.Model
	texture rl.Texture2D

	//How it is
	position     rl.Vector3
	rotationAxis rl.Vector3
	scale        rl.Vector3

	// defaultTransform rl.Matrix

	//Velocity
	vx float32
	vy float32
	vz float32

	//Angles
	yaw   float32
	pitch float32
	roll  float32

	//What it has
	fuel int32 //Liters
	Ve   int32 //Espace Velocity (m/s^2)
	Mo   int32 //Wet mass (kg)
	Mf   int32 //Dry Mass (kg)

	u float32 //Change in velocity

	thrust float32
}

var (
	sim_started        = false
	zero_gx_at_landing = false
)

func (r *Rocket) updateYaw(angle float32) {
	r.yaw += angle
	if r.yaw > 360 || r.yaw < -360 {
		r.yaw = 0.0
	}
	yawUpdate := rl.MatrixRotateX(rl.Deg2rad * angle)
	r.model.Transform = rl.MatrixMultiply(r.model.Transform, yawUpdate)
}

func (r *Rocket) updatePitch(angle float32) {
	r.pitch += angle
	if r.pitch > 360 || r.pitch < -360 {
		r.pitch = 0.0
	}
	pitchUpdate := rl.MatrixRotateZ(rl.Deg2rad * angle)
	r.model.Transform = rl.MatrixMultiply(r.model.Transform, pitchUpdate)
}

func (r *Rocket) setup() {
	sim_started = false
	// r.defaultTransform = r.model.Transform
	r.position.X = rX
	r.position.Y = rY
	r.position.Z = rZ

	r.vx = 0.0
	r.vy = 0.0
	r.vz = 0.0

	r.pitch = 0.0
	r.yaw = 0.0

	lg := (math.Log(float64(r.Mo) / float64(r.Mf)))
	r.u = (float32(r.Ve) * float32(lg))
	// r.thrust = 20000 * float32(r.Mo)
	// fmt.Println(r.thrust)
	r.thrust = 30
	time = 0

	r.updateYaw(0)
	r.updatePitch(0)
}

func (r *Rocket) thrustRocket() {
	r.vy += float32(r.thrust * float32(math.Sin(float64((r.pitch+90)*math.Pi)/180)))
	r.vx += float32(r.thrust * float32(math.Sin(float64((r.pitch)*math.Pi)/180)))
	r.vz += float32(r.thrust * float32(math.Cos(float64((r.yaw+90)*math.Pi)/180)))
	rl.DrawCube(rl.NewVector3(r.position.X, r.position.Y-1, r.position.Z), 0.25, 2, 0.25, rl.Orange)

}

func (r *Rocket) manageKeys() {

	if rl.IsKeyPressed(rl.KeyR) {
		r.setup()
	}

	if rl.IsKeyPressed(rl.KeyS) {
		sim_started = true
	}

	//Thrust
	if rl.IsKeyDown((rl.KeyT)) {
		// r.vy += float32(r.u)
		// lg := (math.Log(float64(r.Mo) / float64(r.Mf)))
		// r.u = (float32(r.Ve) * float32(lg))
		// r.thrust = r.u * float32(r.Mo)
		// r.vy += r.thrust
		// fmt.Println(r.u)
		// fmt.Println(r.thrust)
		r.thrustRocket()
	}

	//Yaw control
	if rl.IsKeyDown(rl.KeyY) {
		r.updateYaw(10)

	} else if rl.IsKeyDown(rl.KeyH) {
		r.updateYaw(-10)
	}

	//Pitch control
	if rl.IsKeyDown(rl.KeyP) {
		r.updatePitch(10)
	} else if rl.IsKeyDown(rl.KeyL) {
		r.updatePitch(-10)
	}

}

func (r *Rocket) flightControl() {

	/* Inside the pad: Make it land*/
	rocket_in_pad := (r.position.X > -10.0 && r.position.X < 10.0 && r.position.Z > -10.0 && r.position.Z < 10.0)
	if rocket_in_pad {

		//Landing Rocket vertically
		if zero_gx_at_landing {
			if -r.vy > r.position.Y {

				if r.pitch < 0 {
					r.updatePitch(0.5)
				}

				if r.pitch > 0 {
					r.updatePitch(-0.5)
				}

				if int32(r.pitch) != 0 {
					r.thrustRocket()
				} else {
					if -r.vy < r.position.Y {
						r.thrustRocket()
					}
				}
				return
			}

			return
		}

		if int32(r.vx) == 0 {
			zero_gx_at_landing = true
			return
		}

		if r.vx != 0 {
			// a := float32(math.Abs(float64(r.pitch)))
			if r.vx < -2 {
				r.updatePitch(0.1)
			} else if r.vx > 2 {
				r.updatePitch(-0.1)
			}
		}

		fmt.Println("ROCKET VERTICAL")
		//If rocket is vertical, land it
		//If velocity is negative & less than height
		if r.vy < -r.position.Y*5 {
			r.thrustRocket()
			return
		}

	} else {
		//Head to 0,0
		if r.vy < -r.position.Y*5 || r.position.Y < 10 {

			if r.position.X > 0 {
				if r.vx > -5 {
					r.updatePitch(-0.1 * (r.position.X / 10))
				}
			} else {
				if r.vx < 5 {
					r.updatePitch(0.1 * (r.position.X / 10))
				}
			}
			r.thrustRocket()
			return
		}

	}

}

func (r *Rocket) update() {
	//Check if user pressed a key
	r.manageKeys()

	//Flight the rocket!!

	//Calculate forces
	if sim_started {
		//F_abajo = m*g
		r.vy += (g * float32(r.Mo))
		r.position.X = r.position.X + (r.vx / timeScale)
		r.position.Y = r.position.Y + (r.vy / timeScale)
		r.position.Z = r.position.Z + (r.vz / timeScale)
		r.flightControl()
	}

	if r.position.Y < 0 {

		r.position.Y = 0
		// r.vy = 0
		sim_started = false
	}

	fmt.Println(r.vx, r.vy, r.vz, "<", r.position.X, ",", r.position.Y, ",", r.position.Z, ">")

	if r.position.Y == 0 {
		fmt.Println("------ FINAL STATS ---------")
		fmt.Println(r.vx, r.vy, r.vz)
	}

}

func (r *Rocket) drawInfo() {
	rl.DrawText("Yaw: "+fmt.Sprintf("%g", r.yaw)+"°", screenWidth-160, 20, 20, rl.DarkGray)
	rl.DrawText("Pitch: "+fmt.Sprintf("%g", r.pitch+90)+"°", screenWidth-160, 40, 20, rl.DarkGray)
	rl.DrawText("Roll: "+fmt.Sprintf("%g", r.roll)+"°", screenWidth-160, 60, 20, rl.DarkGray)
}

func (r *Rocket) draw() {
	rl.DrawModelEx(r.model, r.position, r.rotationAxis, 0.0, r.scale, r.color) // Draw 3d model with texture

	//Draw yellow sphere
	rl.DrawSphere(rl.NewVector3(r.position.X, r.position.Y, r.position.Z), 0.30, rl.NewColor(255, 255, 0, 255))

	//Draw Vector Directions
	if r.vx < 0 {
		rl.DrawCube(rl.NewVector3(r.position.X-1, r.position.Y, r.position.Z), 1.5, 0.25, 0.25, rl.Red)
	} else {
		rl.DrawCube(rl.NewVector3(r.position.X+1, r.position.Y, r.position.Z), 1.5, 0.25, 0.25, rl.Red)
	}

	if r.vy < 0 {
		rl.DrawCube(rl.NewVector3(r.position.X, r.position.Y-1, r.position.Z), 0.25, 1.5, 0.25, rl.Blue)
	} else {
		rl.DrawCube(rl.NewVector3(r.position.X, r.position.Y+1, r.position.Z), 0.25, 1.5, 0.25, rl.Blue)

	}

	if r.vz < 0 {
		rl.DrawCube(rl.NewVector3(r.position.X, r.position.Y, r.position.Z-1), 0.25, 0.25, 1.5, rl.Green)
	} else {
		rl.DrawCube(rl.NewVector3(r.position.X, r.position.Y, r.position.Z+1), 0.25, 0.25, 1.5, rl.Green)
	}
}

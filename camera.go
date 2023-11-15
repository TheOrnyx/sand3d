package main

import (
	glm "github.com/go-gl/mathgl/mgl32"
	"github.com/chewxy/math32"
)

const (
	FORWARD = iota
	BACKWARD
	LEFT
	RIGHT
)

const (
	INIT_YAW         = -90
	INIT_PITCH       = 0.0
	INIT_SPEED       = 0.1
	INIT_SENSITIVITY = 0.1
	INIT_ZOOM        = 45.0
)

type Camera struct {
	Position         glm.Vec3
	Front            glm.Vec3
	Up               glm.Vec3
	Right            glm.Vec3
	WorldUp          glm.Vec3
	Yaw              float32
	Pitch            float32
	MovementSpeed    float32
	MouseSensitivity float32
	Zoom             float32
}

func MakeCamera(position, up glm.Vec3, yaw, pitch float32) *Camera {
	newCam := &Camera{
		Front: glm.Vec3{0.0, 0.0, -1.0},
		MovementSpeed: INIT_SPEED,
		MouseSensitivity: INIT_SENSITIVITY,
		Zoom: INIT_ZOOM,
		
		Position: position,
		WorldUp: up,
		Yaw: yaw,
		Pitch: pitch,
	}
	newCam.updateCameraVectors()
	return newCam
}

func MakeCameraWithScalars(posX, posY, posZ, upX, upY, upZ, yaw, pitch float32) *Camera {
	return MakeCamera(glm.Vec3{posX, posY, posZ}, glm.Vec3{upX, upY, upZ}, yaw, pitch)
}

func (c *Camera) GetViewMatrix() glm.Mat4 {
	return glm.LookAtV(c.Position, c.Position.Add(c.Front), c.Up)
}

func (c *Camera) ProcessKeyPress(moveDir int, deltaTime float32) {
	velocity := c.MovementSpeed + deltaTime
	switch moveDir {
	case FORWARD:
		c.Position = c.Position.Add(c.Front.Mul(velocity))
	case BACKWARD:
		c.Position = c.Position.Sub(c.Front.Mul(velocity))
	case LEFT:
		c.Position = c.Position.Sub(c.Right.Mul(velocity))
	case RIGHT:
		c.Position = c.Position.Add(c.Right.Mul(velocity))
	}
}

// ProcessMouseMovement handle the mouse movement
func (c *Camera) ProcessMouseMovement(xOffset, yOffset float32, constrainPitch bool) {
	xOffset *= c.MouseSensitivity
	yOffset *= c.MouseSensitivity

	c.Yaw += xOffset
	c.Pitch += yOffset

	if constrainPitch {
		c.Pitch = glm.Clamp(c.Pitch, -89.0, 89.0)
	}
	c.updateCameraVectors()
}

// ProcessMouseScroll handle mouse scroll
func (c *Camera) ProcessMouseScroll(yOffset float32)  {
	c.Zoom -= yOffset
	c.Zoom = glm.Clamp(c.Zoom, 1.0, 45.0)
}

// updateCameraVectors calculates the front amera vectors euler angles
func (c *Camera) updateCameraVectors() {
	yawRad, pitchRad := glm.DegToRad(c.Yaw), glm.DegToRad(c.Pitch)
	frontX := math32.Cos(yawRad) * math32.Cos(pitchRad)
	frontY := math32.Sin(pitchRad)
	frontZ := math32.Sin(yawRad) * math32.Cos(pitchRad)

	c.Front = glm.Vec3{frontX, frontY, frontZ}.Normalize()
	c.Right = c.Front.Cross(c.WorldUp).Normalize()
	c.Up = c.Right.Cross(c.Front).Normalize()
}



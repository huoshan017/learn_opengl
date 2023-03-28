package common

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

type CameraMovement int32

const (
	Forward  CameraMovement = iota
	Backward CameraMovement = 1
	Left     CameraMovement = 2
	Right    CameraMovement = 3
)

const (
	YAW         = -90.0
	PITCH       = 0.0
	SPEED       = 2.5
	SENSITIVITY = 0.1
	ZOOM        = 45.0
)

type Camera struct {
	// camera attributes
	position mgl32.Vec3
	front    mgl32.Vec3
	up       mgl32.Vec3
	right    mgl32.Vec3
	worldUp  mgl32.Vec3

	// euler angles
	yaw   float64
	pitch float64

	// camera options
	movementSpeed    float64
	mouseSensitivity float64
	zoom             float64
}

func NewCamera(position, up mgl32.Vec3, yaw, pitch float64) *Camera {
	c := &Camera{
		position: position,
		worldUp:  up,
		yaw:      yaw,
		pitch:    pitch,
	}
	c.front = mgl32.Vec3{0.0, 0.0, -1.0}
	c.movementSpeed = SPEED
	c.mouseSensitivity = SENSITIVITY
	c.zoom = ZOOM
	c.updateCameraVectors()
	return c
}

func NewCameraDefault() *Camera {
	return NewCamera(mgl32.Vec3{0.0, 0.0, 0.0}, mgl32.Vec3{0.0, 1.0, 0.0}, YAW, PITCH)
}

func NewCameraDefaultExceptPosition(position mgl32.Vec3) *Camera {
	return NewCamera(position, mgl32.Vec3{0.0, 1.0, 0.0}, YAW, PITCH)
}

func NewCamera2(posX, posY, posZ float64, upX, upY, upZ float64, yaw, pitch float64) *Camera {
	c := &Camera{}
	c.front = mgl32.Vec3{0.0, 0.0, -1.0}
	c.movementSpeed = SPEED
	c.mouseSensitivity = SENSITIVITY
	c.zoom = ZOOM
	c.position = mgl32.Vec3{float32(posX), float32(posY), float32(posZ)}
	c.worldUp = mgl32.Vec3{float32(upX), float32(upY), float32(upZ)}
	c.yaw = yaw
	c.pitch = pitch
	return c
}

func (c Camera) Position() mgl32.Vec3 {
	return c.position
}

func (c Camera) Front() mgl32.Vec3 {
	return c.front
}

func (c Camera) Zoom() float64 {
	return c.zoom
}

func (c *Camera) YawAdd(degree float64) {
	c.yaw += degree
}

// returns the view matrix calculated using Euler Angles and the LookAt Matrix
func (c Camera) GetViewMatrix() mgl32.Mat4 {
	return mgl32.LookAtV(c.position, c.position.Add(c.front), c.up)
}

// processes input received from any keyboard-like input systems. Accepts input parameters in the form of camera defined ENUM (to abstract it from windowing systems)
func (c *Camera) ProcessKeyboard(direction CameraMovement, deltaTime float64) {
	velocity := float32(c.movementSpeed * deltaTime)
	if direction == Forward {
		c.position = c.position.Add(c.front.Mul(velocity))
	} else if direction == Backward {
		c.position = c.position.Sub(c.front.Mul(velocity))
	} else if direction == Left {
		c.position = c.position.Sub(c.right.Mul(velocity))
	} else if direction == Right {
		c.position = c.position.Add(c.right.Mul(velocity))
	}
}

// processes input received from a mouse input system. Expects the offset value in both the x and y direction.
func (c *Camera) ProcessMouseMovement(xoffset, yoffset float64, constrainPitch bool) {
	xoffset *= c.mouseSensitivity
	yoffset *= c.mouseSensitivity

	c.yaw += xoffset
	c.pitch += yoffset

	// make sure that when pitch is out of bounds, screen doesn't get flipped
	if constrainPitch {
		if c.pitch > 89.0 {
			c.pitch = 89.0
		}
		if c.pitch < -89.0 {
			c.pitch = -89.0
		}
	}

	// update front, right and up vectors using the updated Euler angles
	c.updateCameraVectors()
}

// processes input receives from a mouse scroll-wheel event. Only requires input on the vertical wheel-axis.
func (c *Camera) ProcessMouseScroll(yoffset float64) {
	c.zoom -= yoffset
	if c.zoom < 1.0 {
		c.zoom = 1.0
	}
	if c.zoom > 45.0 {
		c.zoom = 45.0
	}
}

// calculates the front vector from the Camera's (updated) Euler Angles
func (c *Camera) updateCameraVectors() {
	var front = mgl32.Vec3{
		float32(math.Cos(math.Pi*c.yaw/180) * math.Cos(math.Pi*c.pitch/180)),
		float32(math.Sin(math.Pi * c.pitch / 180)),
		float32(math.Sin(math.Pi*c.yaw/180) * math.Cos(math.Pi*c.pitch/180)),
	}
	c.front = front.Normalize()
	// also re-calculate the Right and Up vectors
	c.right = c.front.Cross(c.worldUp).Normalize() // normalize the vectors, because their length gets closer to 0 the more you look up or down which results in slower movement.
	c.up = c.right.Cross(c.front).Normalize()
}

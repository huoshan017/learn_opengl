package main

import (
	"learn_opengl/assimp"
	"learn_opengl/common"
	"learn_opengl/gl"
	"log"
	"math"
	"math/rand"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	SRC_WIDTH  = 800
	SRC_HEIGHT = 600
)

var (
	// camera
	camera     *common.Camera = common.NewCameraDefaultExceptPosition(mgl32.Vec3{0.0, 0.0, 3.0})
	lastX      float64        = SRC_WIDTH / 2.0
	lastY      float64        = SRC_HEIGHT / 2.0
	firstMouse bool           = true
	// timing
	deltaTime, lastFrame float64
)

func main() {
	// glfw: initialize and configure
	// ------------------------------
	glfw.Init()
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)

	// glfw: window creation
	// ---------------------
	window, err := glfw.CreateWindow(SRC_WIDTH, SRC_HEIGHT, "LearnOpenGL", nil, nil)
	if err != nil {
		log.Fatalf("Failed to create GLFW window")
	}

	window.MakeContextCurrent()
	gl.Init()
	window.SetFramebufferSizeCallback(func(_ *glfw.Window, width, height int) {
		gl.Viewport(0, 0, int32(width), int32(height))
	})
	window.SetCursorPosCallback(mouseCallback)
	window.SetScrollCallback(scrollCallback)

	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)

	// configure global opengl state
	// -----------------------------
	gl.Enable(gl.DEPTH_TEST)

	// build and compile shaders
	// -------------------------
	shader := gl.NewShader("10.2.instancing.vs", "10.2.instancing.fs")

	// load models
	rock := assimp.NewModelDefault("../resources/objects/rock/rock.obj")
	planet := assimp.NewModelDefault("../resources/objects/planet/planet.obj")

	// generate a large list of semi-random model transformation matrices
	// ------------------------------------------------------------------
	var amount uint32 = 1000
	modelMatrices := make([]mgl32.Mat4, amount)
	rand.Seed(int64(glfw.GetTime()))
	var radius float32 = 50.0
	var offset float32 = 2.5
	for i := uint32(0); i < amount; i++ {
		// 1. translation: displace along circle with 'radius' in range [-offset, offset]
		var angle float32 = float32(i) / float32(amount) * 360.0
		var displacement float32 = float32(rand.Int31()%(int32)(2*offset*100.0))/100.0 - offset
		var x float32 = float32(math.Sin(float64(angle)))*radius + displacement
		displacement = float32(rand.Int31()%(int32)(2*offset*100))/100.0 - offset
		var y float32 = displacement * 0.4 // keep height of asteroid field smaller compared to width of x and z
		displacement = float32(rand.Int31()%(int32)(2*offset*100))/100.0 - offset
		var z float32 = float32(math.Cos(float64(angle)))*radius + displacement
		model := mgl32.Translate3D(x, y, z)

		// 2. scale: Scale between 0.05 and 0.25
		var scale float32 = float32(rand.Int31()%20)/100.0 + 0.05
		model = model.Mul4(mgl32.Scale3D(scale, scale, scale))

		// 3. rotation: add random rotation around a (semi)randomly picked rotation axis vector
		var rotAngle float32 = float32(rand.Int31() % 360)
		model = model.Mul4(mgl32.HomogRotate3D(rotAngle, mgl32.Vec3{0.4, 0.6, 0.8}))

		// 4. now add to list of matrices
		modelMatrices[i] = model
	}

	// render loop
	// -----------
	for !window.ShouldClose() {
		// per-frame time logic
		// --------------------
		currentFrame := glfw.GetTime()
		deltaTime = currentFrame - lastFrame
		lastFrame = currentFrame

		// input
		// -----
		processInput(window)

		// render
		// ------
		gl.ClearColor(0.1, 0.1, 0.1, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// configure transformation matrices
		projection := mgl32.Perspective(common.Degree2Radian(45.0), SRC_WIDTH/SRC_HEIGHT, 1.0, 100.0)
		view := camera.GetViewMatrix()
		shader.Use()
		shader.SetMat4("projection\x00", &projection)
		shader.SetMat4("view\x00", &view)

		// draw planet
		model := mgl32.Translate3D(0.0, -3.0, 0.0)
		model = model.Mul4(mgl32.Scale3D(4.0, 4.0, 4.0))
		shader.SetMat4("model\x00", &model)
		planet.Draw(&shader)

		// draw meteorites
		for i := uint32(0); i < amount; i++ {
			shader.SetMat4("model\x00", (&modelMatrices[i]))
			rock.Draw(&shader)
		}

		// glfw: swap buffers and poll IO events (keys pressed/released, mouse moved etc.)
		window.SwapBuffers()
		glfw.WaitEventsTimeout(0.01)
	}

	glfw.Terminate()
}

func processInput(window *glfw.Window) {
	if window.GetKey(glfw.KeyEscape) == glfw.Press {
		window.SetShouldClose(true)
	}

	if window.GetKey(glfw.KeyW) == glfw.Press {
		camera.ProcessKeyboard(common.Forward, deltaTime)
	}
	if window.GetKey(glfw.KeyS) == glfw.Press {
		camera.ProcessKeyboard(common.Backward, deltaTime)
	}
	if window.GetKey(glfw.KeyA) == glfw.Press {
		camera.ProcessKeyboard(common.Left, deltaTime)
	}
	if window.GetKey(glfw.KeyD) == glfw.Press {
		camera.ProcessKeyboard(common.Right, deltaTime)
	}
}

// glfw: whenever the mouse moves, this callback is called
// -------------------------------------------------------
func mouseCallback(_ *glfw.Window, xposIn, yposIn float64) {
	xpos := xposIn
	ypos := yposIn

	if firstMouse {
		lastX = xpos
		lastY = ypos
		firstMouse = false
	}

	xoffset := xpos - lastX
	yoffset := lastY - ypos // reversed since y-coordinates go from bottom to top

	lastX = xpos
	lastY = ypos

	camera.ProcessMouseMovement(xoffset, yoffset, true)
}

// glfw: whenever the mouse scroll wheel scrolls, this callback is called
// ----------------------------------------------------------------------
func scrollCallback(_ *glfw.Window, _, yoffset float64) {
	camera.ProcessMouseScroll(yoffset)
}

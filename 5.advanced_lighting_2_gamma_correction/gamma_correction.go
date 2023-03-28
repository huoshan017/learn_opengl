package main

import (
	"learn_opengl/common"
	"learn_opengl/gl"
	"log"
	"unsafe"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/huoshan017/go-stbi"
)

const (
	SRC_WIDTH  = 800
	SRC_HEIGHT = 600
)

var (
	gammaEnabled    = false
	gammaKeyPressed = false
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
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	// build and compile shaders
	// -------------------------
	shader := gl.NewShader("2.gamma_correction.vs", "2.gamma_correction.fs")

	// set up vertex data (and buffer(s)) and configure vertex attributes
	// ------------------------------------------------------------------
	planeVertices := []float32{
		// positions            // normals         // texcoords
		10.0, -0.5, 10.0, 0.0, 1.0, 0.0, 10.0, 0.0,
		-10.0, -0.5, 10.0, 0.0, 1.0, 0.0, 0.0, 0.0,
		-10.0, -0.5, -10.0, 0.0, 1.0, 0.0, 0.0, 10.0,

		10.0, -0.5, 10.0, 0.0, 1.0, 0.0, 10.0, 0.0,
		-10.0, -0.5, -10.0, 0.0, 1.0, 0.0, 0.0, 10.0,
		10.0, -0.5, -10.0, 0.0, 1.0, 0.0, 10.0, 10.0,
	}

	// plane vao
	var planeVao, planeVbo uint32
	gl.GenVertexArrays(1, &planeVao)
	gl.GenBuffers(1, &planeVbo)
	gl.BindVertexArray(planeVao)
	gl.BindBuffer(gl.ARRAY_BUFFER, planeVbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(planeVertices)*4, unsafe.Pointer(&planeVertices[0]), gl.STATIC_DRAW)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 8*4, 0)
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 8*4, 3*4)
	gl.EnableVertexAttribArray(2)
	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, 8*4, 6*4)
	gl.BindVertexArray(0)

	// load texture
	// ------------
	floorTexture := loadTexture("../resources/textures/wood.png")
	floorTextureGammaCorrected := loadTexture("../resources/textures/wood.png")

	// shader configuration
	// --------------------
	shader.Use()
	shader.SetInt32("floorTexture\x00", 0)

	// lighting info
	// -------------
	lightPositions := []mgl32.Vec3{
		{-3.0, 0.0, 0.0},
		{-1.0, 0.0, 0.0},
		{1.0, 0.0, 0.0},
		{3.0, 0.0, 0.0},
	}
	lightColors := []mgl32.Vec3{
		{0.25},
		{0.50},
		{0.75},
		{1.00},
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

		// draw objects
		shader.Use()
		projection := mgl32.Perspective(common.Degree2Radian(float32(camera.Zoom())), SRC_WIDTH/SRC_HEIGHT, 0.1, 100.0)
		view := camera.GetViewMatrix()
		shader.SetMat4("projection\x00", &projection)
		shader.SetMat4("view\x00", &view)
		// set light uniforms
		gl.Uniform3fv(gl.GetUniformLocation(shader.Id(), "lightPositions\x00"), 4, &lightPositions[0][0])
		gl.Uniform3fv(gl.GetUniformLocation(shader.Id(), "lightColors\x00"), 4, &lightColors[0][0])
		cameraPos := camera.Position()
		shader.SetVec3("viewPos\x00", &cameraPos)
		shader.SetBool("gamma\x00", gammaEnabled)
		// floor
		gl.BindVertexArray(planeVao)
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, func() uint32 {
			if gammaEnabled {
				return floorTextureGammaCorrected
			} else {
				return floorTexture
			}
		}())
		gl.DrawArrays(gl.TRIANGLES, 0, 6)

		// glfw: swap buffers and poll IO events (keys pressed/released, mouse moved etc.)
		window.SwapBuffers()
		glfw.WaitEventsTimeout(0.01)
	}

	// optional: de-allcate all resources once they've outlived their purpose:
	// -----------------------------------------------------------------------
	gl.DeleteVertexArrays(1, &planeVao)
	gl.DeleteBuffers(1, &planeVbo)

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

	if window.GetKey(glfw.KeyB) == glfw.Press && !gammaKeyPressed {
		gammaEnabled = !gammaEnabled
		gammaKeyPressed = true
	}
	if window.GetKey(glfw.KeyB) == glfw.Release {
		gammaKeyPressed = false
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

// utility function for loading a 2D texture from file
func loadTexture(path string) uint32 {
	var textureId uint32
	gl.GenTextures(1, &textureId)

	var nChannels int32
	image, err := stbi.Load(path, &nChannels, 0)
	if err != nil {
		log.Fatalf("Texture failed to load, err: %v", err)
	}

	var format int32
	if nChannels == 0 {
		format = gl.RED
	} else if nChannels == 3 {
		format = gl.RGB
	} else if nChannels == 4 {
		format = gl.RGBA
	}

	gl.BindTexture(gl.TEXTURE_2D, textureId)
	width := image.Rect.Dx()
	height := image.Rect.Dy()
	gl.TexImage2D(gl.TEXTURE_2D, 0, format, int32(width), int32(height), 0, uint32(format), gl.UNSIGNED_BYTE, unsafe.Pointer(&image.Pix[0]))
	gl.GenerateMipmap(gl.TEXTURE_2D)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	return textureId
}

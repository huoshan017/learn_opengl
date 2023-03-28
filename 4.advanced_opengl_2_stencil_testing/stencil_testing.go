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
	// camera
	camera     *common.Camera = common.NewCameraDefaultExceptPosition(mgl32.Vec3{0.0, 0.0, 3.0})
	lastX      float64        = SRC_WIDTH / 2.0
	lastY      float64        = SRC_HEIGHT / 2.0
	firstMouse bool           = true
	// timing
	deltaTime, lastFrame float64
	// lighting
	lightPos = mgl32.Vec3{1.2, 1.0, 2.0}
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
	gl.DepthFunc(gl.LESS)
	gl.Enable(gl.STENCIL_TEST)
	gl.StencilFunc(gl.NOTEQUAL, 1, 0xff)
	gl.StencilOp(gl.KEEP, gl.KEEP, gl.REPLACE)

	// build and compile our shader zprogram
	// -------------------------------------
	var shader = gl.NewShader("2.stencil_testing.vs", "2.stencil_testing.fs")
	var shaderSingleColor = gl.NewShader("2.stencil_testing.vs", "2.stencil_single_color.fs")

	// set up vertex data (and buffer(s)) and configure vertex attributes
	// ------------------------------------------------------------------
	cubeVertices := []float32{
		// positions      // normals      // texture coords
		-0.5, -0.5, -0.5, 0.0, 0.0,
		0.5, -0.5, -0.5, 1.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		-0.5, 0.5, -0.5, 0.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 0.0,

		-0.5, -0.5, 0.5, 0.0, 0.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 1.0,
		0.5, 0.5, 0.5, 1.0, 1.0,
		-0.5, 0.5, 0.5, 0.0, 1.0,
		-0.5, -0.5, 0.5, 0.0, 0.0,

		-0.5, 0.5, 0.5, 1.0, 0.0,
		-0.5, 0.5, -0.5, 1.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 1.0,
		-0.5, -0.5, 0.5, 0.0, 0.0,
		-0.5, 0.5, 0.5, 1.0, 0.0,

		0.5, 0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, -0.5, 0.5, 0.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 0.0,

		-0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, -0.5, -0.5, 1.0, 1.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		-0.5, -0.5, 0.5, 0.0, 0.0,
		-0.5, -0.5, -0.5, 0.0, 1.0,

		-0.5, 0.5, -0.5, 0.0, 1.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, 0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 0.0,
		-0.5, 0.5, 0.5, 0.0, 0.0,
		-0.5, 0.5, -0.5, 0.0, 1.0,
	}
	planeVertices := []float32{
		5.0, -0.5, 5.0, 2.0, 0.0,
		-5.0, -0.5, 5.0, 0.0, 0.0,
		-5.0, -0.5, -5.0, 0.0, 2.0,
		5.0, -0.5, 5.0, 2.0, 0.0,
		-5.0, -0.5, -5.0, 0.0, 2.0,
		5.0, -0.5, -5.0, 2.0, 2.0,
	}

	// first, configure the cube's VAO (and VBO)
	var cubeVao, cubeVbo uint32
	gl.GenVertexArrays(1, &cubeVao)
	gl.GenBuffers(1, &cubeVbo)
	gl.BindVertexArray(cubeVao)
	gl.BindBuffer(gl.ARRAY_BUFFER, cubeVbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(cubeVertices)*4, unsafe.Pointer(&cubeVertices[0]), gl.STATIC_DRAW)
	// position attribute
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, 0)
	gl.EnableVertexAttribArray(0)
	// textures coords attribute
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, 3*4)
	gl.EnableVertexAttribArray(1)
	gl.BindVertexArray(0)

	var planeVao, planeVbo uint32
	gl.GenVertexArrays(1, &planeVao)
	gl.GenBuffers(1, &planeVbo)
	gl.BindVertexArray(planeVao)

	gl.BindBuffer(gl.ARRAY_BUFFER, planeVbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(planeVertices)*4, unsafe.Pointer(&planeVertices[0]), gl.STATIC_DRAW)
	// note that we update the lamp's position attribute's stride to reflect the updated buffer data
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, 0)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, 3*4)
	gl.EnableVertexAttribArray(1)
	gl.BindVertexArray(0)

	// load textures (we now use a utility function to keep the code more organized)
	// -----------------------------------------------------------------------------
	cubeTexture := loadTexture("../resources/textures/marble.jpg")
	floorTexture := loadTexture("../resources/textures/metal.png")

	// shader configuration
	// --------------------
	shader.Use()
	shader.SetInt32("texture1\x00", 0)

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
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT | gl.STENCIL_BUFFER_BIT)

		// be sure to activate shader when setting uniforms/drawing objects
		shaderSingleColor.Use()
		model := mgl32.Ident4()
		view := camera.GetViewMatrix()
		projection := mgl32.Perspective(common.Degree2Radian(float32(camera.Zoom())), SRC_WIDTH/SRC_HEIGHT, 0.1, 100)
		shaderSingleColor.SetMat4("view\x00", &view)
		shaderSingleColor.SetMat4("projection\x00", &projection)

		shader.Use()
		shader.SetMat4("view\x00", &view)
		shader.SetMat4("projection\x00", &projection)

		// draw floor as normal, but don't write the floor to the stencil buffer, we only care about the containers. We set its mask to 0x00 to not write to the stencil buffer.
		gl.StencilMask(0x00)
		// floor
		gl.BindVertexArray(planeVao)
		gl.BindTexture(gl.TEXTURE_2D, floorTexture)
		shader.SetMat4("model\x00", &model)
		gl.DrawArrays(gl.TRIANGLES, 0, 6)
		gl.BindVertexArray(0)

		// 1st. render pass, draw objects as normal, writing to the stencil buffer
		// -----------------------------------------------------------------------
		gl.StencilFunc(gl.ALWAYS, 1, 0xff)
		gl.StencilMask(0xff)
		// cubes
		gl.BindVertexArray(cubeVao)
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, cubeTexture)
		model = mgl32.Translate3D(-1.0, 0.0, -1.0)
		shader.SetMat4("model\x00", &model)
		gl.DrawArrays(gl.TRIANGLES, 0, 36)
		model = mgl32.Ident4()
		model = model.Mul4(mgl32.Translate3D(2.0, 0.0, 0.0))
		shader.SetMat4("model\x00", &model)
		gl.DrawArrays(gl.TRIANGLES, 0, 36)

		// 2nd. render pass, now draw slightly scaled versions of the objects, this time disabling stencil writing.
		// Because the stencil buffer is now filled with several 1s. The parts of the buffer that are 1 are not drawn, thus only drawing
		// the objects' size differences, making it look like borders.
		// -----------------------------------------------------------------------------------------------------------------------------
		gl.StencilFunc(gl.NOTEQUAL, 1, 0xff)
		gl.StencilMask(0x00)
		gl.Disable(gl.DEPTH_TEST)
		shaderSingleColor.Use()
		var scale float32 = 1.1
		// cubes
		gl.BindVertexArray(cubeVao)
		gl.BindTexture(gl.TEXTURE0, cubeTexture)
		model = mgl32.Translate3D(-1.0, 0.0, -1.0)
		model = model.Mul4(mgl32.Scale3D(scale, scale, scale))
		shaderSingleColor.SetMat4("model\x00", &model)
		gl.DrawArrays(gl.TRIANGLES, 0, 36)
		// floor
		model = mgl32.Translate3D(2.0, 0.0, 0.0)
		model = model.Mul4(mgl32.Scale3D(scale, scale, scale))
		shaderSingleColor.SetMat4("model\x00", &model)
		gl.DrawArrays(gl.TRIANGLES, 0, 36)
		gl.BindVertexArray(0)

		gl.StencilMask(0xff)
		gl.StencilFunc(gl.ALWAYS, 0, 0xff)
		gl.Enable(gl.DEPTH_TEST)

		// glfw: swap buffers and poll IO events (keys pressed/released, mouse moved etc.)
		// -------------------------------------------------------------------------------
		window.SwapBuffers()
		//glfw.PollEvents()
		glfw.WaitEventsTimeout(0.01)
	}

	// optional: de-allocate all resources once they've outlived their purpose:
	// ------------------------------------------------------------------------
	gl.DeleteVertexArrays(1, &cubeVao)
	gl.DeleteVertexArrays(1, &planeVao)
	gl.DeleteBuffers(1, &cubeVbo)
	gl.DeleteBuffers(1, &planeVbo)

	// glfw: terminate, clearing all previously allocated GLFW resources
	// -----------------------------------------------------------------
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
	if nChannels == 1 {
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

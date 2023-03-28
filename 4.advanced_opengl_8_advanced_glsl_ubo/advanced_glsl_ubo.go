package main

import (
	"learn_opengl/common"
	"learn_opengl/gl"
	"log"
	"unsafe"

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

	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)

	// configure global opengl state
	// -----------------------------
	gl.Enable(gl.DEPTH_TEST)

	// build and compile our shader zprogram
	// -------------------------------------
	var shaderRed = gl.NewShader("8.advanced_glsl.vs", "8.red.fs")
	var shaderGreen = gl.NewShader("8.advanced_glsl.vs", "8.green.fs")
	var shaderBlue = gl.NewShader("8.advanced_glsl.vs", "8.blue.fs")
	var shaderYellow = gl.NewShader("8.advanced_glsl.vs", "8.yellow.fs")

	// set up vertex data (and buffer(s)) and configure vertex attributes
	// ------------------------------------------------------------------
	cubeVertices := []float32{
		// positions
		-0.5, -0.5, -0.5,
		0.5, -0.5, -0.5,
		0.5, 0.5, -0.5,
		0.5, 0.5, -0.5,
		-0.5, 0.5, -0.5,
		-0.5, -0.5, -0.5,

		-0.5, -0.5, 0.5,
		0.5, -0.5, 0.5,
		0.5, 0.5, 0.5,
		0.5, 0.5, 0.5,
		-0.5, 0.5, 0.5,
		-0.5, -0.5, 0.5,

		-0.5, 0.5, 0.5,
		-0.5, 0.5, -0.5,
		-0.5, -0.5, -0.5,
		-0.5, -0.5, -0.5,
		-0.5, -0.5, 0.5,
		-0.5, 0.5, 0.5,

		0.5, 0.5, 0.5,
		0.5, 0.5, -0.5,
		0.5, -0.5, -0.5,
		0.5, -0.5, -0.5,
		0.5, -0.5, 0.5,
		0.5, 0.5, 0.5,

		-0.5, -0.5, -0.5,
		0.5, -0.5, -0.5,
		0.5, -0.5, 0.5,
		0.5, -0.5, 0.5,
		-0.5, -0.5, 0.5,
		-0.5, -0.5, -0.5,

		-0.5, 0.5, -0.5,
		0.5, 0.5, -0.5,
		0.5, 0.5, 0.5,
		0.5, 0.5, 0.5,
		-0.5, 0.5, 0.5,
		-0.5, 0.5, -0.5,
	}
	// cube
	var cubeVao, cubeVbo uint32
	gl.GenVertexArrays(1, &cubeVao)
	gl.GenBuffers(1, &cubeVbo)
	gl.BindVertexArray(cubeVao)
	gl.BindBuffer(gl.ARRAY_BUFFER, cubeVbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(cubeVertices)*4, unsafe.Pointer(&cubeVertices[0]), gl.STATIC_DRAW)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, 0)
	gl.EnableVertexAttribArray(0)

	// configure a uniform buffer object
	// ---------------------------------
	// first. We get the relevant block indices
	uniformBlockIndexRed := gl.GetUniformBlockIndex(shaderRed.Id(), &([]byte("Matrices"))[0])
	uniformBlockIndexGreen := gl.GetUniformBlockIndex(shaderGreen.Id(), &([]byte("Matrices"))[0])
	uniformBlockIndexBlue := gl.GetUniformBlockIndex(shaderBlue.Id(), &([]byte("Matricess"))[0])
	uniformBlockIndexYellow := gl.GetUniformBlockIndex(shaderYellow.Id(), &([]byte("Matrices"))[0])
	// then we link each shader's uniform block to this uniform binding point
	gl.UniformBlockBinding(shaderRed.Id(), uniformBlockIndexRed, 0)
	gl.UniformBlockBinding(shaderGreen.Id(), uniformBlockIndexGreen, 0)
	gl.UniformBlockBinding(shaderBlue.Id(), uniformBlockIndexBlue, 0)
	gl.UniformBlockBinding(shaderYellow.Id(), uniformBlockIndexYellow, 0)
	// Now actually create the buffer
	var uboMatrices uint32
	gl.GenBuffers(1, &uboMatrices)
	gl.BindBuffer(gl.UNIFORM_BUFFER, uboMatrices)
	gl.BufferData(gl.UNIFORM_BUFFER, int(2*unsafe.Sizeof(mgl32.Ident4())), nil, gl.STATIC_DRAW)
	gl.BindBuffer(gl.UNIFORM_BUFFER, 0)
	// define the range of the buffer that links to a uniform binding point
	gl.BindBufferRange(gl.UNIFORM_BUFFER, 0, uboMatrices, 0, int(2*unsafe.Sizeof(mgl32.Ident4())))

	// store the projection matrix (we only do this once now) (note: we're not using zoom anymore by changing the FOV)
	projection := mgl32.Perspective(45.0, SRC_WIDTH/SRC_HEIGHT, 0.1, 100.0)
	gl.BindBuffer(gl.UNIFORM_BUFFER, uboMatrices)
	gl.BufferSubData(gl.UNIFORM_BUFFER, 0, int(unsafe.Sizeof(mgl32.Ident4())), unsafe.Pointer(&projection))
	gl.BindBuffer(gl.UNIFORM_BUFFER, 0)

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

		// set the view and projection matrix in the uniform block - we only have to do this once per loop iteration.
		view := camera.GetViewMatrix()
		gl.BindBuffer(gl.UNIFORM_BUFFER, uboMatrices)
		gl.BufferSubData(gl.UNIFORM_BUFFER, int(unsafe.Sizeof(mgl32.Mat4{})), int(unsafe.Sizeof(mgl32.Mat4{})), unsafe.Pointer(&view))
		gl.BindBuffer(gl.UNIFORM_BUFFER, 0)

		// draw 4 cubes
		// RED
		gl.BindVertexArray(cubeVao)
		shaderRed.Use()
		model := mgl32.Translate3D(-0.75, 0.75, 0.0) // move top-left
		shaderRed.SetMat4("model\x00", &model)
		gl.DrawArrays(gl.TRIANGLES, 0, 36)
		// GREEN
		shaderGreen.Use()
		model = mgl32.Translate3D(0.75, 0.75, 0.0) // move top-right
		shaderGreen.SetMat4("model\x00", &model)
		gl.DrawArrays(gl.TRIANGLES, 0, 36)
		// YELLOW
		shaderYellow.Use()
		model = mgl32.Translate3D(-0.75, -0.75, 0.0) // move bottom-left
		shaderYellow.SetMat4("model\x00", &model)
		gl.DrawArrays(gl.TRIANGLES, 0, 36)
		// BLUE
		shaderBlue.Use()
		model = mgl32.Translate3D(0.75, -0.75, 0.0) // move bottom-right
		shaderBlue.SetMat4("model\x00", &model)
		gl.DrawArrays(gl.TRIANGLES, 0, 36)

		// glfw: swap buffers and poll IO events (keys pressed/released, mouse moved etc.)
		// -------------------------------------------------------------------------------
		window.SwapBuffers()
		//glfw.PollEvents()
		glfw.WaitEventsTimeout(0.01)
	}

	// optional: de-allocate all resources once they've outlived their purpose:
	// ------------------------------------------------------------------------
	gl.DeleteVertexArrays(1, &cubeVao)
	gl.DeleteBuffers(1, &cubeVbo)

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

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
	window.SetScrollCallback(scrollCallback)

	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)

	// configure global opengl state
	// -----------------------------
	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.MULTISAMPLE) // enabled by default on some drivers, but not all so always enable to make sure

	// build and compile shaders
	// -------------------------
	shader := gl.NewShader("11.2.anti_aliasing.vs", "11.2.anti_aliasing.fs")
	screenShader := gl.NewShader("11.2.aa_post.vs", "11.2.aa_post.fs")

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
	quadVertices := []float32{
		-1.0, 1.0, 0.0, 1.0,
		-1.0, -1.0, 0.0, 0.0,
		1.0, -1.0, 1.0, 0.0,

		-1.0, 1.0, 0.0, 1.0,
		1.0, -1.0, 0.0, 1.0,
		1.0, 1.0, 1.0, 1.0,
	}

	// setup cube vao
	var cubeVao, cubeVbo uint32
	gl.GenVertexArrays(1, &cubeVao)
	gl.GenBuffers(1, &cubeVbo)
	gl.BindVertexArray(cubeVao)
	gl.BindBuffer(gl.ARRAY_BUFFER, cubeVbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(cubeVertices)*4, unsafe.Pointer(&cubeVertices[0]), gl.STATIC_DRAW)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, 0)
	// set screen vao
	var quadVao, quadVbo uint32
	gl.GenVertexArrays(1, &quadVao)
	gl.GenBuffers(1, &quadVbo)
	gl.BindVertexArray(quadVao)
	gl.BindBuffer(gl.ARRAY_BUFFER, quadVbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(quadVertices)*4, unsafe.Pointer(&quadVertices[0]), gl.STATIC_DRAW)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 4*4, 0)
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 4*4, 2*4)

	// configure MSAA framebuffer
	// --------------------------
	var framebuffer uint32
	gl.GenFramebuffers(1, &framebuffer)
	gl.BindFramebuffer(gl.FRAMEBUFFER, framebuffer)
	// create a multisampled color attachment texture
	var textureColorBufferMultiSampled uint32
	gl.GenTextures(1, &textureColorBufferMultiSampled)
	gl.BindTexture(gl.TEXTURE_2D_MULTISAMPLE, textureColorBufferMultiSampled)
	gl.TexImage2DMultisample(gl.TEXTURE_2D_MULTISAMPLE, 4, gl.RGB, SRC_WIDTH, SRC_HEIGHT, true)
	gl.BindTexture(gl.TEXTURE_2D_MULTISAMPLE, 0)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D_MULTISAMPLE, textureColorBufferMultiSampled, 0)
	// create a (also multisampled) renderbuffer object for depth and stencil attachments
	var rbo uint32
	gl.GenRenderbuffers(1, &rbo)
	gl.BindRenderbuffer(gl.RENDERBUFFER, rbo)
	gl.RenderbufferStorageMultisample(gl.RENDERBUFFER, 4, gl.DEPTH24_STENCIL8, SRC_WIDTH, SRC_HEIGHT)
	gl.BindRenderbuffer(gl.RENDERBUFFER, 0)
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_STENCIL_ATTACHMENT, gl.RENDERBUFFER, rbo)

	if gl.CheckFramebufferStatus(gl.FRAMEBUFFER) != gl.FRAMEBUFFER_COMPLETE {
		log.Fatalf("ERROR::FRAMEBUFFER:: Framebuffer is not complete!")
	}
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	// configure second post-processing framebuffer
	var intermediateFbo uint32
	gl.GenFramebuffers(1, &intermediateFbo)
	gl.BindBuffer(gl.FRAMEBUFFER, intermediateFbo)
	// create a color attachment texture
	var screenTexture uint32
	gl.GenTextures(1, &screenTexture)
	gl.BindTexture(gl.TEXTURE_2D, screenTexture)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB, SRC_WIDTH, SRC_HEIGHT, 0, gl.RGB, gl.UNSIGNED_BYTE, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, screenTexture, 0) // we only need a color buffer

	if gl.CheckFramebufferStatus(gl.FRAMEBUFFER) != gl.FRAMEBUFFER_COMPLETE {
		log.Fatalf("ERROR::FRAMEBUFFER:: Intermediate framebuffer is not complete!")
	}
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	// shader configuration
	// --------------------
	screenShader.Use()
	screenShader.SetInt32("screenTexture\x00", 0)

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

		// 1. draw screen as normal in multisampled buffers
		gl.BindFramebuffer(gl.FRAMEBUFFER, framebuffer)
		gl.ClearColor(0.1, 0.1, 0.1, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.Enable(gl.DEPTH_TEST)

		// configure transformation matrices
		shader.Use()
		projection := mgl32.Perspective(common.Degree2Radian(45.0), SRC_WIDTH/SRC_HEIGHT, 1.0, 100.0)
		view := camera.GetViewMatrix()
		model := mgl32.Ident4()
		shader.SetMat4("projection\x00", &projection)
		shader.SetMat4("view\x00", &view)
		shader.SetMat4("model\x00", &model)

		gl.BindVertexArray(cubeVao)
		gl.DrawArrays(gl.TRIANGLES, 0, 36)

		// 2. now blit multisampled buffer(s) to normal colorbuffer of intermediate fbo. Image is stored in screenTexture
		gl.BindFramebuffer(gl.READ_FRAMEBUFFER, framebuffer)
		gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, intermediateFbo)
		gl.BlitFramebuffer(0, 0, SRC_WIDTH, SRC_HEIGHT, 0, 0, SRC_WIDTH, SRC_HEIGHT, gl.COLOR_BUFFER_BIT, gl.NEAREST)

		// 3. now render quad with scene's visuals as its texture image
		gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
		gl.ClearColor(1.0, 1.0, 1.0, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.Disable(gl.DEPTH_TEST)

		// draw screen quad
		screenShader.Use()
		gl.BindVertexArray(quadVao)
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, screenTexture) // use the now resolved color attachment as the quad's texture
		gl.DrawArrays(gl.TRIANGLES, 0, 6)

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

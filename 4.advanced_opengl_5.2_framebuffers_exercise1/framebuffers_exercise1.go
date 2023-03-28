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

	// build and compile our shader zprogram
	// -------------------------------------
	var shader = gl.NewShader("5.2.framebuffers.vs", "5.2.framebuffers.fs")
	var screenShader = gl.NewShader("5.2.framebuffers_screen.vs", "5.2.framebuffers_screen.fs")

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
		// positions    // texture coords
		5.0, -0.5, 5.0, 2.0, 0.0,
		-5.0, -0.5, 5.0, 0.0, 0.0,
		-5.0, -0.5, -5.0, 0.0, 2.0,
		5.0, -0.5, 5.0, 2.0, 0.0,
		-5.0, -0.5, -5.0, 0.0, 2.0,
		5.0, -0.5, -5.0, 2.0, 2.0,
	}
	quadVertices := []float32{ // vertex attributes for a quad that fills the entire screen in Normalized Device Coordiantes.
		-0.3, 1.0, 0.0, 1.0,
		-0.3, 0.7, 0.0, 0.0,
		0.3, 0.7, 1.0, 0.0,

		-0.3, 1.0, 0.0, 1.0,
		0.3, 0.7, 1.0, 0.0,
		0.3, 1.0, 1.0, 1.0,
	}

	// cube
	var cubeVao, cubeVbo uint32
	gl.GenVertexArrays(1, &cubeVao)
	gl.GenBuffers(1, &cubeVbo)
	gl.BindVertexArray(cubeVao)
	gl.BindBuffer(gl.ARRAY_BUFFER, cubeVbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(cubeVertices)*4, unsafe.Pointer(&cubeVertices[0]), gl.STATIC_DRAW)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, 0)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, 3*4)
	gl.EnableVertexAttribArray(1)
	gl.BindVertexArray(0)

	// plane
	var planeVao, planeVbo uint32
	gl.GenVertexArrays(1, &planeVao)
	gl.GenBuffers(1, &planeVbo)
	gl.BindVertexArray(planeVao)
	gl.BindBuffer(gl.ARRAY_BUFFER, planeVbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(planeVertices)*4, unsafe.Pointer(&planeVertices[0]), gl.STATIC_DRAW)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, 0)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, 3*4)
	gl.EnableVertexAttribArray(1)
	gl.BindVertexArray(0)

	// screen quad
	var quadVao, quadVbo uint32
	gl.GenVertexArrays(1, &quadVao)
	gl.GenBuffers(1, &quadVbo)
	gl.BindVertexArray(quadVao)
	gl.BindBuffer(gl.ARRAY_BUFFER, quadVbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(quadVertices)*4, unsafe.Pointer(&quadVertices[0]), gl.STATIC_DRAW)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 4*4, 0)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 4*4, 2*4)
	gl.EnableVertexAttribArray(1)
	gl.BindVertexArray(0)

	// load textures
	// -------------
	cubeTexture := loadTexture("../resources/textures/container.jpg")
	floorTexture := loadTexture("../resources/textures/metal.png")

	// shader configuration
	// --------------------
	shader.Use()
	shader.SetInt32("texture1\x00", 0)

	screenShader.Use()
	screenShader.SetInt32("screenTexture\x00", 0)

	// framebuffer configuration
	// -------------------------
	var framebuffer uint32
	gl.GenFramebuffers(1, &framebuffer)
	gl.BindFramebuffer(gl.FRAMEBUFFER, framebuffer)
	// create a color attachment texture
	var textureColorbuffer uint32
	gl.GenTextures(1, &textureColorbuffer)
	gl.BindTexture(gl.TEXTURE_2D, textureColorbuffer)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB, SRC_WIDTH, SRC_HEIGHT, 0, gl.RGB, gl.UNSIGNED_BYTE, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, textureColorbuffer, 0)
	// create a renderbuffer object for depth and stencil attachment (we won't be sampling these)
	var rbo uint32
	gl.GenRenderbuffers(1, &rbo)
	gl.BindRenderbuffer(gl.RENDERBUFFER, rbo)
	gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH24_STENCIL8, SRC_WIDTH, SRC_HEIGHT)           // use a single renderbuffer object for both a depth AND stencil buffer.
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_STENCIL_ATTACHMENT, gl.RENDERBUFFER, rbo) // now actually attach it
	// now that we actually created the framebuffer and added all attachments we want to check if it is actually complete now
	if gl.CheckFramebufferStatus(gl.FRAMEBUFFER) != gl.FRAMEBUFFER_COMPLETE {
		log.Fatalf("ERROR::FRAMEBUFFER:: Framebuffer is not complete!")
	}
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

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

		// first render pass: mirror texture.
		// bind to framebuffer and draw to color texture as we normally
		// would, but with the view camera reversed.
		// bind to framebuffer and draw scene as we normally would to color texture
		// ------------------------------------------------------------------------
		gl.BindFramebuffer(gl.FRAMEBUFFER, framebuffer)
		gl.Enable(gl.DEPTH_TEST) // enable depth testing (is disabled for rendering screen-space quad)

		// make sure we clear the framebuffer's content
		gl.ClearColor(0.1, 0.1, 0.1, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		shader.Use()
		model := mgl32.Ident4()
		camera.YawAdd(180.0)                     // rotate the camera's yaw 180 degrees around
		camera.ProcessMouseMovement(0, 0, false) // call this to make sure it updates its camera vectors, note that we disable pitch constrains for this specific case (otherwise we can't reverse camera's pitch values)
		view := camera.GetViewMatrix()
		camera.YawAdd(-180.0) // reset it back to its original orientation
		camera.ProcessMouseMovement(0, 0, true)
		projection := mgl32.Perspective(common.Degree2Radian(float32(camera.Zoom())), SRC_WIDTH/SRC_HEIGHT, 0.1, 100)
		shader.SetMat4("view\x00", &view)
		shader.SetMat4("projection\x00", &projection)

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
		// floor
		gl.BindVertexArray(planeVao)
		gl.BindTexture(gl.TEXTURE_2D, floorTexture)
		model = mgl32.Ident4()
		shader.SetMat4("model\x00", &model)
		gl.DrawArrays(gl.TRIANGLES, 0, 6)
		gl.BindVertexArray(0)

		// second render pass: draw as normal
		// ----------------------------------
		gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
		// clear all relevant buffers
		gl.ClearColor(0.1, 0.1, 0.1, 1.0) // set clear color to white (not really necessary actually, since we won't be able to see the behind the quad anyways)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		view = camera.GetViewMatrix()
		shader.SetMat4("view\x00", &view)

		// cubes
		gl.BindVertexArray(cubeVao)
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, cubeTexture)
		model = model.Mul4(mgl32.Translate3D(-1.0, 0.0, -1.0))
		shader.SetMat4("model\x00", &model)
		gl.DrawArrays(gl.TRIANGLES, 0, 36)
		model = mgl32.Ident4()
		model = model.Mul4(mgl32.Translate3D(2.0, 0.0, 0.0))
		shader.SetMat4("model\x00", &model)
		gl.DrawArrays(gl.TRIANGLES, 0, 36)
		// floor
		gl.BindVertexArray(planeVao)
		gl.BindTexture(gl.TEXTURE_2D, floorTexture)
		model = mgl32.Ident4()
		shader.SetMat4("model\x00", &model)
		gl.DrawArrays(gl.TRIANGLES, 0, 6)
		gl.BindVertexArray(0)

		// now draw the mirror quad with screen texture
		// --------------------------------------------
		gl.Disable(gl.DEPTH_TEST) // disable depth test so screen-space quad isn't discarded due to depth test.

		screenShader.Use()
		gl.BindVertexArray(quadVao)
		gl.BindTexture(gl.TEXTURE_2D, textureColorbuffer) // use the color attachment texture as the texture of the quad plane
		gl.DrawArrays(gl.TRIANGLES, 0, 6)

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
	gl.DeleteVertexArrays(1, &quadVao)
	gl.DeleteBuffers(1, &cubeVbo)
	gl.DeleteBuffers(1, &planeVbo)
	gl.DeleteBuffers(1, &quadVbo)

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

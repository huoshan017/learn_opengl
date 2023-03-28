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
	// meshes
	planeVao uint32
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
	shader := gl.NewShader("3.1.2.shadow_mapping.vs", "3.1.2.shadow_mapping.fs")
	simpleDepthShader := gl.NewShader("3.1.2.shadow_mapping_depth.vs", "3.1.2.shadow_mapping_depth.fs")
	debugDepthQuad := gl.NewShader("3.1.2.debug_quad.vs", "3.1.2.debug_quad_depth.fs")

	// set up vertex data (and buffer(s)) and configure vertex attributes
	// ------------------------------------------------------------------
	planeVertices := []float32{
		// positions            // normals         // texcoords
		25.0, -0.5, 25.0, 0.0, 1.0, 0.0, 25.0, 0.0,
		-25.0, -0.5, 25.0, 0.0, 1.0, 0.0, 0.0, 0.0,
		-25.0, -0.5, -25.0, 0.0, 1.0, 0.0, 0.0, 25.0,

		25.0, -0.5, 25.0, 0.0, 1.0, 0.0, 25.0, 0.0,
		-25.0, -0.5, -25.0, 0.0, 1.0, 0.0, 0.0, 25.0,
		25.0, -0.5, -25.0, 0.0, 1.0, 0.0, 25.0, 25.0,
	}

	// plane vao
	var planeVbo uint32
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
	woodTexture := loadTexture("../resources/textures/wood.png")

	// configure depth map FBO
	// -----------------------
	const (
		SHADOW_WIDTH  int32 = 1024
		SHADOW_HEIGHT int32 = 1024
	)
	var depthMapFbo uint32
	gl.GenFramebuffers(1, &depthMapFbo)
	// create depth texture
	var depthMap uint32
	gl.GenTextures(1, &depthMap)
	gl.BindTexture(gl.TEXTURE_2D, depthMap)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.DEPTH_COMPONENT, SHADOW_WIDTH, SHADOW_HEIGHT, 0, gl.DEPTH_COMPONENT, gl.FLOAT, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	// attach depth texture as FBO's depth buffer
	gl.BindFramebuffer(gl.FRAMEBUFFER, depthMapFbo)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.TEXTURE_2D, depthMap, 0)
	gl.DrawBuffer(gl.NONE)
	gl.ReadBuffer(gl.NONE)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	// shader configuration
	// --------------------
	shader.Use()
	shader.SetInt32("diffuseTexture\x00", 0)
	shader.SetInt32("shadowMap\x00", 1)
	debugDepthQuad.Use()
	debugDepthQuad.SetInt32("depthMap\x00", 0)

	// lighting info
	// -------------
	lightPos := mgl32.Vec3{-2.0, 4.0, -1.0}

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

		// 1. render depth of scene to texture (from light's perspective)
		// --------------------------------------------------------------
		var (
			lightProjection, lightView, lightSpaceMatrix mgl32.Mat4
			near_plane, far_plane                        float32 = 1.0, 7.5
		)
		lightProjection = mgl32.Ortho(-10.0, 10.0, -10.0, 10.0, near_plane, far_plane)
		lightView = mgl32.LookAtV(lightPos, mgl32.Vec3{0.0}, mgl32.Vec3{0.0, 1.0, 0.0})
		lightSpaceMatrix = lightView.Mul4(lightProjection)
		// render scene from light's point of view
		simpleDepthShader.Use()
		simpleDepthShader.SetMat4("lightSpaceMatrix\x00", &lightSpaceMatrix)

		gl.Viewport(0, 0, SRC_WIDTH, SRC_HEIGHT)
		gl.BindFramebuffer(gl.FRAMEBUFFER, depthMapFbo)
		gl.Clear(gl.DEPTH_BUFFER_BIT)
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, woodTexture)
		renderScene(&simpleDepthShader)
		gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

		// reset viewport
		gl.Viewport(0, 0, SRC_WIDTH, SRC_HEIGHT)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// 2. render scene as normal using the generated depth/shadow map
		// --------------------------------------------------------------
		shader.Use()
		projection := mgl32.Perspective(common.Degree2Radian(float32(camera.Zoom())), SRC_WIDTH/SRC_HEIGHT, 0.1, 100.0)
		view := camera.GetViewMatrix()
		shader.SetMat4("projection\x00", &projection)
		shader.SetMat4("view\x00", &view)
		// set light uniforms
		cameraPos := camera.Position()
		shader.SetVec3("viewPos\x00", &cameraPos)
		shader.SetVec3("lightPos\x00", &lightPos)
		shader.SetMat4("lightSpaceMatrix\x00", &lightSpaceMatrix)
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, woodTexture)
		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_2D, depthMap)
		renderScene(&shader)

		// render depth map to quad for visual debugging
		// ---------------------------------------------
		debugDepthQuad.Use()
		debugDepthQuad.SetFloat32("near_plane\x00", near_plane)
		debugDepthQuad.SetFloat32("far_plane\x00", far_plane)
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, depthMap)
		//renderQuad()

		// glfw: swap buffers and poll IO events (keys pressed/released, mouse moved etc.)
		// -------------------------------------------------------------------------------
		window.SwapBuffers()
		glfw.WaitEventsTimeout(0.01)
	}

	// optional: de-allcate all resources once they've outlived their purpose:
	// -----------------------------------------------------------------------
	gl.DeleteVertexArrays(1, &planeVao)
	gl.DeleteBuffers(1, &planeVbo)

	glfw.Terminate()
}

// render the 3D scene
// -------------------
func renderScene(shader *gl.Shader) {
	// floor
	model := mgl32.Ident4()
	shader.SetMat4("model\x00", &model)
	gl.BindVertexArray(planeVao)
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
	// cubes
	model = mgl32.Ident4()
	model = model.Mul4(mgl32.Translate3D(0.0, 1.5, 0.0))
	model = model.Mul4(mgl32.Scale3D(0.5, 0.5, 0.5))
	shader.SetMat4("model\x00", &model)
	renderCube()
	model = mgl32.Ident4()
	model = model.Mul4(mgl32.Translate3D(2.0, 0.0, 1.0))
	model = model.Mul4(mgl32.Scale3D(0.5, 0.5, 0.5))
	shader.SetMat4("model\x00", &model)
	renderCube()
	model = mgl32.Ident4()
	model = model.Mul4(mgl32.Translate3D(-1.0, 0.0, 2.0))
	model = model.Mul4(mgl32.HomogRotate3D(common.Degree2Radian(60.0), mgl32.Vec3{1.0, 0.0, 1.0}.Normalize()))
	model = model.Mul4(mgl32.Scale3D(0.25, 0.25, 0.25))
	shader.SetMat4("model\x00", &model)
	renderCube()
}

// renderCube() renders a 1x1 3D cube in NDC.
// ------------------------------------------
var (
	cubeVao, cubeVbo uint32
)

func renderCube() {
	// initialize (if necessary)
	if cubeVao == 0 {
		vertices := []float32{
			// back face
			-1.0, -1.0, -1.0, 0.0, 0.0, -1.0, 0.0, 0.0, // bottom-left
			1.0, 1.0, -1.0, 0.0, 0.0, -1.0, 1.0, 1.0, // top-right
			1.0, -1.0, -1.0, 0.0, 0.0, -1.0, 1.0, 0.0, // bottom-right
			1.0, 1.0, -1.0, 0.0, 0.0, -1.0, 1.0, 1.0, // top-right
			-1.0, -1.0, -1.0, 0.0, 0.0, -1.0, 0.0, 0.0, // bottom-left
			-1.0, 1.0, -1.0, 0.0, 0.0, -1.0, 0.0, 1.0, // top-left
			// front face
			-1.0, -1.0, 1.0, 0.0, 0.0, 1.0, 0.0, 0.0, // bottom-left
			1.0, -1.0, 1.0, 0.0, 0.0, 1.0, 1.0, 0.0, // bottom-right
			1.0, 1.0, 1.0, 0.0, 0.0, 1.0, 1.0, 1.0, // top-right
			1.0, 1.0, 1.0, 0.0, 0.0, 1.0, 1.0, 1.0, // top-right
			-1.0, 1.0, 1.0, 0.0, 0.0, 1.0, 0.0, 1.0, // top-left
			-1.0, -1.0, 1.0, 0.0, 0.0, 1.0, 0.0, 0.0, // bottom-left
			// left face
			-1.0, 1.0, 1.0, -1.0, 0.0, 0.0, 1.0, 0.0, // top-right
			-1.0, 1.0, -1.0, -1.0, 0.0, 0.0, 1.0, 1.0, // top-left
			-1.0, -1.0, -1.0, -1.0, 0.0, 0.0, 0.0, 1.0, // bottom-left
			-1.0, -1.0, -1.0, -1.0, 0.0, 0.0, 0.0, 1.0, // bottom-left
			-1.0, -1.0, 1.0, -1.0, 0.0, 0.0, 0.0, 0.0, // bottom-right
			-1.0, 1.0, 1.0, -1.0, 0.0, 0.0, 1.0, 0.0, // top-right
			// right face
			1.0, 1.0, 1.0, 1.0, 0.0, 0.0, 1.0, 0.0, // top-left
			1.0, -1.0, -1.0, 1.0, 0.0, 0.0, 0.0, 1.0, // bottom-right
			1.0, 1.0, -1.0, 1.0, 0.0, 0.0, 1.0, 1.0, // top-right
			1.0, -1.0, -1.0, 1.0, 0.0, 0.0, 0.0, 1.0, // bottom-right
			1.0, 1.0, 1.0, 1.0, 0.0, 0.0, 1.0, 0.0, // top-left
			1.0, -1.0, 1.0, 1.0, 0.0, 0.0, 0.0, 0.0, // bottom-left
			// bottom face
			-1.0, -1.0, -1.0, 0.0, -1.0, 0.0, 0.0, 1.0, // top-right
			1.0, -1.0, -1.0, 0.0, -1.0, 0.0, 1.0, 1.0, // top-left
			1.0, -1.0, 1.0, 0.0, -1.0, 0.0, 1.0, 0.0, // bottom-left
			1.0, -1.0, 1.0, 0.0, -1.0, 0.0, 1.0, 0.0, // bottom-left
			-1.0, -1.0, 1.0, 0.0, -1.0, 0.0, 0.0, 0.0, // bottom-right
			-1.0, -1.0, -1.0, 0.0, -1.0, 0.0, 0.0, 1.0, // top-right
			// top face
			-1.0, 1.0, -1.0, 0.0, 1.0, 0.0, 0.0, 1.0, // top-left
			1.0, 1.0, 1.0, 0.0, 1.0, 0.0, 1.0, 0.0, // bottom-right
			1.0, 1.0, -1.0, 0.0, 1.0, 0.0, 1.0, 1.0, // top-right
			1.0, 1.0, 1.0, 0.0, 1.0, 0.0, 1.0, 0.0, // bottom-right
			-1.0, 1.0, -1.0, 0.0, 1.0, 0.0, 0.0, 1.0, // top-left
			-1.0, 1.0, 1.0, 0.0, 1.0, 0.0, 0.0, 0.0, // bottom-left
		}
		gl.GenVertexArrays(1, &cubeVao)
		gl.GenBuffers(1, &cubeVbo)
		// fill buffer
		gl.BindBuffer(gl.ARRAY_BUFFER, cubeVbo)
		gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, unsafe.Pointer(&vertices[0]), gl.STATIC_DRAW)
		// link vertex attributes
		gl.BindVertexArray(cubeVao)
		gl.EnableVertexAttribArray(0)
		gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 8*4, 0)
		gl.EnableVertexAttribArray(1)
		gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 8*4, 3*4)
		gl.EnableVertexAttribArray(2)
		gl.VertexAttribPointer(2, 2, gl.FLOAT, false, 8*4, 6*4)
		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
		gl.BindVertexArray(0)
	}
	// render Cube
	gl.BindVertexArray(cubeVao)
	gl.DrawArrays(gl.TRIANGLES, 0, 36)
	gl.BindVertexArray(0)
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

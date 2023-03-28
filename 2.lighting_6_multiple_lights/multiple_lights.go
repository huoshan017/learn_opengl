package main

import (
	"fmt"
	"learn_opengl/common"
	"learn_opengl/gl"
	"log"
	"math"
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

	// build and compile our shader zprogram
	// -------------------------------------
	var lightingShader = gl.NewShader("6.multiple_lights.vs", "6.multiple_lights.fs")
	var lightCubeShader = gl.NewShader("6.light_cube.vs", "6.light_cube.fs")

	// set up vertex data (and buffer(s)) and configure vertex attributes
	// ------------------------------------------------------------------
	vertices := []float32{
		// positions      // normals      // texture coords
		-0.5, -0.5, -0.5, 0.0, 0.0, -1.0, 0.0, 0.0,
		0.5, -0.5, -0.5, 0.0, 0.0, -1.0, 1.0, 0.0,
		0.5, 0.5, -0.5, 0.0, 0.0, -1.0, 1.0, 1.0,
		0.5, 0.5, -0.5, 0.0, 0.0, -1.0, 1.0, 1.0,
		-0.5, 0.5, -0.5, 0.0, 0.0, -1.0, 0.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 0.0, -1.0, 0.0, 0.0,

		-0.5, -0.5, 0.5, 0.0, 0.0, 1.0, 0.0, 0.0,
		0.5, -0.5, 0.5, 0.0, 0.0, 1.0, 1.0, 0.0,
		0.5, 0.5, 0.5, 0.0, 0.0, 1.0, 1.0, 1.0,
		0.5, 0.5, 0.5, 0.0, 0.0, 1.0, 1.0, 1.0,
		-0.5, 0.5, 0.5, 0.0, 0.0, 1.0, 0.0, 1.0,
		-0.5, -0.5, 0.5, 0.0, 0.0, 1.0, 0.0, 0.0,

		-0.5, 0.5, 0.5, -1.0, 0.0, 0.0, 1.0, 0.0,
		-0.5, 0.5, -0.5, -1.0, 0.0, 0.0, 1.0, 1.0,
		-0.5, -0.5, -0.5, -1.0, 0.0, 0.0, 0.0, 1.0,
		-0.5, -0.5, -0.5, -1.0, 0.0, 0.0, 0.0, 1.0,
		-0.5, -0.5, 0.5, -1.0, 0.0, 0.0, 0.0, 0.0,
		-0.5, 0.5, 0.5, -1.0, 0.0, 0.0, 1.0, 0.0,

		0.5, 0.5, 0.5, 1.0, 0.0, 0.0, 1.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 0.0, 0.0, 1.0, 1.0,
		0.5, -0.5, -0.5, 1.0, 0.0, 0.0, 0.0, 1.0,
		0.5, -0.5, -0.5, 1.0, 0.0, 0.0, 0.0, 1.0,
		0.5, -0.5, 0.5, 1.0, 0.0, 0.0, 0.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 0.0, 0.0, 1.0, 0.0,

		-0.5, -0.5, -0.5, 0.0, -1.0, 0.0, 0.0, 1.0,
		0.5, -0.5, -0.5, 0.0, -1.0, 0.0, 1.0, 1.0,
		0.5, -0.5, 0.5, 0.0, -1.0, 0.0, 1.0, 0.0,
		0.5, -0.5, 0.5, 0.0, -1.0, 0.0, 1.0, 0.0,
		-0.5, -0.5, 0.5, 0.0, -1.0, 0.0, 0.0, 0.0,
		-0.5, -0.5, -0.5, 0.0, -1.0, 0.0, 0.0, 1.0,

		-0.5, 0.5, -0.5, 0.0, 1.0, 0.0, 0.0, 1.0,
		0.5, 0.5, -0.5, 0.0, 1.0, 0.0, 1.0, 1.0,
		0.5, 0.5, 0.5, 0.0, 1.0, 0.0, 1.0, 0.0,
		0.5, 0.5, 0.5, 0.0, 1.0, 0.0, 1.0, 0.0,
		-0.5, 0.5, 0.5, 0.0, 1.0, 0.0, 0.0, 0.0,
		-0.5, 0.5, -0.5, 0.0, 1.0, 0.0, 0.0, 1.0,
	}
	// positions all containers
	cubePositions := []mgl32.Vec3{
		{0.0, 0.0, 0.0},
		{2.0, 5.0, -15.0},
		{-1.5, -2.2, -2.5},
		{-3.8, -2.0, -12.3},
		{2.4, -0.4, -3.5},
		{-1.7, 3.0, -7.5},
		{1.3, -2.0, -2.5},
		{1.5, 2.0, -2.5},
		{1.5, 0.2, -1.5},
		{-1.3, 1.0, -1.5},
	}
	// positions of the point lights
	pointLightPositions := []mgl32.Vec3{
		{0.7, 0.2, 2.0},
		{2.3, -3.3, -4.0},
		{-4.0, 2.0, -12.0},
		{0.0, 0.0, -3.0},
	}

	// first, configure the cube's VAO (and VBO)
	var vbo, cubeVao uint32
	gl.GenVertexArrays(1, &cubeVao)
	gl.GenBuffers(1, &vbo)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, unsafe.Pointer(&vertices[0]), gl.STATIC_DRAW)

	gl.BindVertexArray(cubeVao)
	// position attribute
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 8*4, 0)
	gl.EnableVertexAttribArray(0)
	// normal attribute
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 8*4, 3*4)
	gl.EnableVertexAttribArray(1)
	// textures coords attribute
	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, 8*4, 6*4)
	gl.EnableVertexAttribArray(2)

	// second, configure the light's VAO (VBO stays the same; the vertices are the same for the light object which is also a 3D cube)
	var lightCubeVAO uint32
	gl.GenVertexArrays(1, &lightCubeVAO)
	gl.BindVertexArray(lightCubeVAO)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	// note that we update the lamp's position attribute's stride to reflect the updated buffer data
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 8*4, 0)
	gl.EnableVertexAttribArray(0)

	// load textures (we now use a utility function to keep the code more organized)
	// -----------------------------------------------------------------------------
	diffuseMap := loadTexture("../resources/textures/container2.png")
	specularMap := loadTexture("../resources/textures/container2_specular.png")

	// shader configuration
	// --------------------
	lightingShader.Use()
	lightingShader.SetInt32("material.diffuse\x00", 0)
	lightingShader.SetInt32("material.specular\x00", 1)

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

		// be sure to activate shader when setting uniforms/drawing objects
		lightingShader.Use()
		cameraPos := camera.Position()
		lightingShader.SetVec3("viewPos\x00", &cameraPos)
		lightingShader.SetFloat32("material.shininess\x00", 32.0)

		/*
			Here we set all the uniforms for the 5/6 types of lights we have. we have to set them manually and index
			the proper PointLight struct in the array to set each uniform variable. This can be done more code-friendly
			by defining light types as classes and set their values in there, or by using a more efficient uniform approach
			by using 'Uniform buffer objects', but that is something we'll discuss in the 'Advanced GLSL' tutorial.
		*/
		// directional light
		lightingShader.SetVec3("dirLight.direction\x00", &mgl32.Vec3{-0.2, -1.0, -0.3})
		lightingShader.SetVec3("dirLight.ambient\x00", &mgl32.Vec3{0.05, 0.05, 0.05})
		lightingShader.SetVec3("dirLight.diffuse\x00", &mgl32.Vec3{0.4, 0.4, 0.4})
		lightingShader.SetVec3("dirLight.specular\x00", &mgl32.Vec3{0.5, 0.5, 0.5})
		// point lights
		for i := 0; i < 4; i++ {
			ls := fmt.Sprintf("pointLights[%v]", i)
			lightingShader.SetVec3(ls+".position\x00", &(pointLightPositions[i]))
			lightingShader.SetVec3(ls+".ambient\x00", &mgl32.Vec3{0.05, 0.05, 0.05})
			lightingShader.SetVec3(ls+".diffuse\x00", &mgl32.Vec3{0.8, 0.8, 0.8})
			lightingShader.SetVec3(ls+".specular\x00", &mgl32.Vec3{1.0, 1.0, 1.0})
			lightingShader.SetFloat32(ls+".constant\x00", 1.0)
			lightingShader.SetFloat32(ls+".linear\x00", 0.09)
			lightingShader.SetFloat32(ls+".quadratic\x00", 0.032)
		}
		// spot light
		lightingShader.SetVec3("spotLight.position\x00", &cameraPos)
		cameraFront := camera.Front()
		lightingShader.SetVec3("spotLight.direction\x00", &cameraFront)
		lightingShader.SetVec3("spotLight.ambient\x00", &mgl32.Vec3{0.0, 0.0, 0.0})
		lightingShader.SetVec3("spotLight.diffuse\x00", &mgl32.Vec3{1.0, 1.0, 1.0})
		lightingShader.SetVec3("spotLight.specular\x00", &mgl32.Vec3{1.0, 1.0, 1.0})
		lightingShader.SetFloat32("spotLight.constant\x00", 1.0)
		lightingShader.SetFloat32("spotLight.linear\x00", 0.09)
		lightingShader.SetFloat32("spotLight.quadratic\x00", 0.032)
		lightingShader.SetFloat32("spotLight.cutOff\x00", float32(math.Cos(float64(common.Degree2Radian(12.5)))))
		lightingShader.SetFloat32("spotLight.outerCutOff\x00", float32(math.Cos(float64(common.Degree2Radian(15.0)))))

		// view/projection transformations
		projection := mgl32.Perspective(common.Degree2Radian(float32(camera.Zoom())), SRC_WIDTH/SRC_HEIGHT, 0.1, 100.0)
		view := camera.GetViewMatrix()
		lightingShader.SetMat4("projection\x00", &projection)
		lightingShader.SetMat4("view\x00", &view)

		// world transformation
		model := mgl32.Ident4()
		lightingShader.SetMat4("model\x00", &model)

		// bind diffuse map
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, diffuseMap)
		// bind specular map
		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_2D, specularMap)

		// render container
		gl.BindVertexArray(cubeVao)
		for i := 0; i < 10; i++ {
			cp := &cubePositions[i]
			// calculate the model matrix for each object and pass it to shader before drawing
			model = mgl32.Translate3D(cp.X(), cp.Y(), cp.Z())
			angle := float32(20.0 * i)
			model = model.Mul4(mgl32.HomogRotate3D(common.Degree2Radian(angle), mgl32.Vec3{1.0, 0.3, 0.5}))
			lightCubeShader.SetMat4("model\x00", &model)

			gl.DrawArrays(gl.TRIANGLES, 0, 36)
		}

		// also draw the lamp object(s)
		lightCubeShader.Use()
		lightCubeShader.SetMat4("projection\x00", &projection)
		lightCubeShader.SetMat4("view\x00", &view)

		// we now draw as many light bulbs as we have point lights.
		gl.BindVertexArray(lightCubeVAO)
		for i := 0; i < 4; i++ {
			pos := &pointLightPositions[i]
			model = mgl32.Translate3D(pos.X(), pos.Y(), pos.Z())
			model = model.Mul4(mgl32.Scale3D(0.2, 0.2, 0.2))
			lightCubeShader.SetMat4("model\x00", &model)
			gl.DrawArrays(gl.TRIANGLES, 0, 36)
		}

		// glfw: swap buffers and poll IO events (keys pressed/released, mouse moved etc.)
		// -------------------------------------------------------------------------------
		window.SwapBuffers()
		//glfw.PollEvents()
		glfw.WaitEventsTimeout(0.01)
	}

	// optional: de-allocate all resources once they've outlived their purpose:
	// ------------------------------------------------------------------------
	gl.DeleteVertexArrays(1, &cubeVao)
	gl.DeleteVertexArrays(1, &lightCubeVAO)
	gl.DeleteBuffers(1, &vbo)

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

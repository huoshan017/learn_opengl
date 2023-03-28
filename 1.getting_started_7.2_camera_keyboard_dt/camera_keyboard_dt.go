package main

import (
	"learn_opengl/gl"
	"log"
	"math"
	"unsafe"

	"github.com/huoshan017/go-stbi"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	SRC_WIDTH  = 800
	SRC_HEIGHT = 600
)

var (
	cameraPos            = mgl32.Vec3{0.0, 0.0, 3.0}
	cameraFront          = mgl32.Vec3{0.0, 0.0, -1.0}
	cameraUp             = mgl32.Vec3{0.0, 1.0, 0.0}
	deltaTime, lastFrame float32
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

	// configure global opengl state
	// -----------------------------
	gl.Enable(gl.DEPTH_TEST)

	// build and compile our shader zprogram
	// -------------------------------------
	var shader = gl.NewShader("7.2.camera.vs", "7.2.camera.fs")

	// set up vertex data (and buffer(s)) and configure vertex attributes
	// ------------------------------------------------------------------
	vertices := []float32{
		// position    // texture coords
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
	// world space positions of our cubes
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

	var vbo, vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)

	gl.BindVertexArray(vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, unsafe.Pointer(&vertices[0]), gl.STATIC_DRAW)

	// position attribute
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, 0)
	gl.EnableVertexAttribArray(0)

	// color attribute
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, 3*4)
	gl.EnableVertexAttribArray(1)

	// load and create texture
	// -----------------------
	var texture1, texture2 uint32
	// texture1
	gl.GenTextures(1, &texture1)
	gl.BindTexture(gl.TEXTURE_2D, texture1)
	// set the texture wrapping parameters
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	// set texture filtering parameters
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	//
	stbi.SetFlipVerticallyOnLoad(true)
	// load image
	var nChannels int32
	image, err := stbi.Load("../resources/textures/container.jpg", &nChannels, 0)
	if err != nil {
		log.Fatalf("Failed to load texture, err: %v", err)
	}
	width := image.Rect.Dx()
	height := image.Rect.Dy()
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB, int32(width), int32(height), 0, gl.RGB, gl.UNSIGNED_BYTE, unsafe.Pointer(&image.Pix[0]))
	gl.GenerateMipmap(gl.TEXTURE_2D)
	// texture2
	gl.GenTextures(1, &texture2)
	gl.BindTexture(gl.TEXTURE_2D, texture2)
	// set the texture wrapping parameters
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	// set the texture filtering parameters
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	// load image, create texture and generate mipmaps
	image, err = stbi.Load("../resources/textures/awesomeface.png", &nChannels, 0)
	if err != nil {
		log.Fatalf("Failed to load texture, err: %v", err)
	}
	width = image.Rect.Dx()
	height = image.Rect.Dy()
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB, int32(width), int32(height), 0, gl.RGBA, gl.UNSIGNED_BYTE, unsafe.Pointer(&image.Pix[0]))
	gl.GenerateMipmap(gl.TEXTURE_2D)

	// tell opengl for each sampler to which texture unit it belongs to (only has to done once)
	// ----------------------------------------------------------------------------------------
	shader.Use()
	shader.SetInt32("texture1\x00", 0)
	shader.SetInt32("texture2\x00", 1)

	// pass projection matrix to shader (as projection matrix rarely changes there's no need to do this per frame)
	projection := mgl32.Perspective(45*math.Pi/180, SRC_WIDTH/SRC_HEIGHT, 0.1, 100.0)
	shader.SetMat4("projection\x00", &projection)

	// render loop
	// -----------
	for !window.ShouldClose() {
		// per-frame time logic
		// --------------------
		currentFrame := glfw.GetTime()
		deltaTime = float32(currentFrame) - lastFrame
		lastFrame = float32(currentFrame)

		// input
		// -----
		processInput(window)

		// render
		// ------
		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// bind texture
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture1)
		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_2D, texture2)

		// active shader
		shader.Use()

		// create transformations
		view := mgl32.LookAtV(cameraPos, cameraPos.Add(cameraFront), cameraUp)
		// note: currently we set the projection matrix each frame, but since the projection matrix rarely changes it's often best practice to set it outside the main loop only once.
		shader.SetMat4("view\x00", &view)

		// render container
		gl.BindVertexArray(vao)
		for i := 0; i < 10; i++ {
			model := mgl32.Ident4()
			cp := &cubePositions[i]
			model = model.Mul4(mgl32.Translate3D(cp.X(), cp.Y(), cp.Z()))
			angle := float32(20.0 * i)
			model = model.Mul4(mgl32.HomogRotate3D(angle*math.Pi/180, mgl32.Vec3{1.0, 0.3, 0.5}))
			shader.SetMat4("model\x00", &model)

			gl.DrawArrays(gl.TRIANGLES, 0, 36)
		}

		// glfw: swap buffers and poll IO events (keys pressed/released, mouse moved etc.)
		// -------------------------------------------------------------------------------
		window.SwapBuffers()
		//glfw.PollEvents()
		glfw.WaitEventsTimeout(0.001)
	}

	// optional: de-allocate all resources once they've outlived their purpose:
	// ------------------------------------------------------------------------
	gl.DeleteVertexArrays(1, &vao)
	gl.DeleteBuffers(1, &vbo)

	// glfw: terminate, clearing all previously allocated GLFW resources
	// -----------------------------------------------------------------
	glfw.Terminate()
}

func processInput(window *glfw.Window) {
	if window.GetKey(glfw.KeyEscape) == glfw.Press {
		window.SetShouldClose(true)
	}

	cameraSpeed := 2.5 * deltaTime
	if window.GetKey(glfw.KeyW) == glfw.Press {
		cameraPos = cameraPos.Add(cameraFront.Mul(cameraSpeed))
	}
	if window.GetKey(glfw.KeyS) == glfw.Press {
		cameraPos = cameraPos.Sub(cameraFront.Mul(cameraSpeed))
	}
	if window.GetKey(glfw.KeyA) == glfw.Press {
		cameraPos = cameraPos.Sub(cameraFront.Cross(cameraUp).Normalize().Mul(cameraSpeed))
	}
	if window.GetKey(glfw.KeyD) == glfw.Press {
		cameraPos = cameraPos.Add(cameraFront.Cross(cameraUp).Normalize().Mul(cameraSpeed))
	}
}

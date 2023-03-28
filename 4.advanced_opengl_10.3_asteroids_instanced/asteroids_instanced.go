package main

import (
	"learn_opengl/assimp"
	"learn_opengl/common"
	"learn_opengl/gl"
	"log"
	"math"
	"math/rand"
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

	// build and compile shaders
	// -------------------------
	asteroidsShader := gl.NewShader("10.3.asteroids.vs", "10.3.asteroids.fs")
	planetShader := gl.NewShader("10.3.planet.vs", "10.3.planet.fs")

	// load models
	rock := assimp.NewModelDefault("../resources/objects/rock/rock.obj")
	planet := assimp.NewModelDefault("../resources/objects/planet/planet.obj")

	// generate a large list of semi-random model transformation matrices
	// ------------------------------------------------------------------
	var amount int32 = 100000
	modelMatrices := make([]mgl32.Mat4, amount)
	rand.Seed(int64(glfw.GetTime()))
	var radius float32 = 150.0
	var offset float32 = 25.0
	for i := int32(0); i < amount; i++ {
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

	// configure instanced array
	// -------------------------
	var buffer uint32
	gl.GenBuffers(1, &buffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, buffer)
	gl.BufferData(gl.ARRAY_BUFFER, int(amount)*int(unsafe.Sizeof(mgl32.Mat4{})), unsafe.Pointer(&modelMatrices[0]), gl.STATIC_DRAW)

	// set transformation matrices as an instance vertex attribute (with divisor 1)
	// note: we're cheating a little by taking the, now publicly declared, VAO of the model's mesh(es) and adding new vertexAttribPointers
	// normally you'd want to do this in a more organized fashion, but for learning purposes this will do.
	// -----------------------------------------------------------------------------------------------------------------------------------
	for i := 0; i < len(rock.Meshes()); i++ {
		vao := rock.Meshes()[i].Vao()
		gl.BindVertexArray(vao)
		// set attribute pointers for matrix (4 times vec4)
		gl.EnableVertexAttribArray(3)
		gl.VertexAttribPointer(3, 4, gl.FLOAT, false, int32(unsafe.Sizeof(mgl32.Mat4{})), 0)
		gl.EnableVertexAttribArray(4)
		gl.VertexAttribPointer(4, 4, gl.FLOAT, false, int32(unsafe.Sizeof(mgl32.Mat4{})), int(unsafe.Sizeof(mgl32.Vec4{})))
		gl.EnableVertexAttribArray(5)
		gl.VertexAttribPointer(5, 4, gl.FLOAT, false, int32(unsafe.Sizeof(mgl32.Mat4{})), 2*int(unsafe.Sizeof(mgl32.Vec4{})))
		gl.EnableVertexAttribArray(6)
		gl.VertexAttribPointer(6, 4, gl.FLOAT, false, int32(unsafe.Sizeof(mgl32.Mat4{})), 3*int(unsafe.Sizeof(mgl32.Vec4{})))

		gl.VertexAttribDivisor(3, 1)
		gl.VertexAttribDivisor(4, 1)
		gl.VertexAttribDivisor(5, 1)
		gl.VertexAttribDivisor(6, 1)

		gl.BindVertexArray(0)
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
		asteroidsShader.Use()
		asteroidsShader.SetMat4("projection\x00", &projection)
		asteroidsShader.SetMat4("view\x00", &view)
		planetShader.Use()
		planetShader.SetMat4("projection\x00", &projection)
		planetShader.SetMat4("view\x00", &view)

		// draw planet
		model := mgl32.Translate3D(0.0, -3.0, 0.0)
		model = model.Mul4(mgl32.Scale3D(4.0, 4.0, 4.0))
		planetShader.SetMat4("model\x00", &model)
		planet.Draw(&planetShader)

		// draw meteorites
		asteroidsShader.Use()
		asteroidsShader.SetInt32("texture_diffuse1\x00", 0)
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, rock.TextureLoaded()[0].Id()) // note: we also made the textures_loaded vector public (instead of private) from the model class.
		for i := 0; i < len(rock.Meshes()); i++ {
			gl.BindVertexArray(rock.Meshes()[i].Vao())
			gl.DrawElementsInstanced(gl.TRIANGLES, int32(len(rock.Meshes()[i].Indices())), gl.UNSIGNED_INT, nil, amount)
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

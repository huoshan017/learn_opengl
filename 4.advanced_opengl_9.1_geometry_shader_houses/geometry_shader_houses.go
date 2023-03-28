package main

import (
	"learn_opengl/gl"
	"log"
	"unsafe"

	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	SRC_WIDTH  = 800
	SRC_HEIGHT = 600
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

	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)

	// configure global opengl state
	// -----------------------------
	gl.Enable(gl.DEPTH_TEST)

	// build and compile shaders
	// -------------------------
	shader := gl.NewShader2("9.1.geometry_shader.vs", "9.1.geometry_shader.fs", "9.1.geometry_shader.gs")

	// set up vertex data (and buffer(s)) and configure vertex attributes
	// ------------------------------------------------------------------
	points := []float32{
		-0.5, 0.5, 1.0, 0.0, 0.0, // top-left
		0.5, 0.5, 0.0, 1.0, 0.0, // top-right
		0.5, -0.5, 0.0, 0.0, 1.0, // bottom-right
		-0.5, -0.5, 1.0, 1.0, 0.0, // bottom-left
	}
	var vao, vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(points)*4, unsafe.Pointer(&points[0]), gl.STATIC_DRAW)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 5*4, 0)
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 5*4, 2*4)
	gl.BindVertexArray(0)

	// render loop
	// -----------
	for !window.ShouldClose() {
		// render
		// ------
		gl.ClearColor(0.1, 0.1, 0.1, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// draw points
		shader.Use()
		gl.BindVertexArray(vao)
		gl.DrawArrays(gl.POINTS, 0, 4)

		// glfw: swap buffers and poll IO events (keys pressed/released, mouse moved etc.)
		window.SwapBuffers()
		glfw.WaitEventsTimeout(0.01)
	}

	// optional: de-allocate all resources once they've outlived their purpose:
	// ------------------------------------------------------------------------
	gl.DeleteVertexArrays(1, &vao)
	gl.DeleteBuffers(1, &vbo)

	glfw.Terminate()
}

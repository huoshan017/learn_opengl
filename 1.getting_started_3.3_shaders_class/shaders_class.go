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

	window, err := glfw.CreateWindow(SRC_WIDTH, SRC_HEIGHT, "LearnOpenGL", nil, nil)
	if err != nil {
		log.Fatalf("Failed to create GLFW window")
	}

	window.MakeContextCurrent()
	gl.Init()
	window.SetFramebufferSizeCallback(func(_ *glfw.Window, width, height int) {
		gl.Viewport(0, 0, int32(width), int32(height))
	})

	// build and compile our shader program
	// ------------------------------------
	ourShader := gl.NewShader("3.3.shader.vs", "3.3.shader.fs")

	// set up vertex data (and buffer(s)) and configure vertex attributes
	// ------------------------------------------------------------------
	vertices := []float32{
		// positions    // colors
		0.5, -0.5, 0.0, 1.0, 0.0, 0.0, // bottom left
		-0.5, -0.5, 0.0, 0.0, 1.0, 0.0, // bottom right
		0.0, 0.5, 0.0, 0.0, 0.0, 1.0, // top
	}

	var vao, vbo uint32
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)

	// bind the Vertex Array Object first, then bind and set vertex buffer(s), and then configure vertex attributes.
	gl.BindVertexArray(vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, unsafe.Pointer(&vertices[0]), gl.STATIC_DRAW)

	// position attribute
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 6*4, 0)
	gl.EnableVertexAttribArray(0)
	// color attribute
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 6*4, 3*4)
	gl.EnableVertexAttribArray(1)

	// render loop
	// -----------
	for !window.ShouldClose() {
		// input
		// -----
		if window.GetKey(glfw.KeyEscape) == glfw.Press {
			window.SetShouldClose(true)
		}

		// render
		// ------
		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		// render the triangle
		ourShader.Use()
		gl.BindVertexArray(vao)
		gl.DrawArrays(gl.TRIANGLES, 0, 3)

		// glfw: swap buffers and poll IO events (keys pressed/released, mouse moved etc.)
		window.SwapBuffers()
		glfw.PollEvents()
	}

	// optional: de-allocate all resources once they've outlived their purpose:
	// ------------------------------------------------------------------------
	gl.DeleteVertexArrays(1, &vao)
	gl.DeleteBuffers(1, &vbo)

	// glfw: terminate, clearing all previously allocated GLFW resources
	glfw.Terminate()
}

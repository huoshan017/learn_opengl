package main

import (
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

func main() {
	// glfw: initialize and configure
	// ------------------------------
	glfw.Init()
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)

	window, err := glfw.CreateWindow(SRC_WIDTH, SRC_HEIGHT, "LearnOpenGL", nil, nil)
	if err != nil {
		log.Fatalf("Failed to create GLFW window, err %v", err)
	}
	window.MakeContextCurrent()
	gl.Init()

	// configure global opengl state
	// -----------------------------
	gl.Enable(gl.DEPTH_TEST)

	// build and compile shaders
	// -------------------------
	shader := gl.NewShader("10.1.instancing.vs", "10.1.instancing.fs")

	// generate a list of 100 quad location/translation-vectors
	// --------------------------------------------------------
	var translations [100]mgl32.Vec2
	var index int
	var offset float32 = 0.1
	for y := -10; y < 10; y += 2 {
		for x := -10; x < 10; x += 2 {
			translations[index] = mgl32.Vec2{float32(x/10.0) + offset, float32(y/10.0) + offset}
			index += 1
		}
	}

	// store instance data in an array buffer
	// --------------------------------------
	var instanceVbo uint32
	gl.GenBuffers(1, &instanceVbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, instanceVbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(translations)*int(unsafe.Sizeof(mgl32.Vec2{})), unsafe.Pointer(&translations[0]), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	// set up vertex data (and buffer(s)) and configure vertex attributes
	// ------------------------------------------------------------------
	var quadVertices = []float32{
		// positions // colors
		-0.05, 0.05, 1.0, 0.0, 0.0,
		0.05, -0.05, 0.0, 1.0, 0.0,
		-0.05, -0.05, 0.0, 0.0, 1.0,

		-0.05, 0.05, 1.0, 0.0, 0.0,
		0.05, -0.05, 0.0, 1.0, 0.0,
		0.05, 0.05, 0.0, 1.0, 1.0,
	}

	var quadVao, quadVbo uint32
	gl.GenVertexArrays(1, &quadVao)
	gl.GenBuffers(1, &quadVbo)
	gl.BindVertexArray(quadVao)
	gl.BindBuffer(gl.ARRAY_BUFFER, quadVbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(quadVertices)*4, unsafe.Pointer(&quadVertices[0]), gl.STATIC_DRAW)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 5*4, 0)
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 5*4, 2*4)
	// also set instance data
	gl.EnableVertexAttribArray(2)
	gl.BindBuffer(gl.ARRAY_BUFFER, instanceVbo) // this attribute comes from a different vertex buffer
	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, 2*4, 0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.VertexAttribDivisor(2, 1) // tell OpenGL this is an instanced vertex attribute

	// render loop
	// -----------
	for !window.ShouldClose() {
		// render
		// ------
		gl.ClearColor(0.1, 0.1, 0.1, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// draw 100 instanced quads
		shader.Use()
		gl.BindVertexArray(quadVao)
		gl.DrawArraysInstanced(gl.TRIANGLES, 0, 6, 100) // 100 triangles of 6 vertices each
		gl.BindVertexArray(0)

		// glfw: swap buffers and poll IO events (keys pressed/released, mouse moved etc.)
		// -------------------------------------------------------------------------------
		window.SwapBuffers()
		glfw.WaitEventsTimeout(0.01)
	}

	// optional: de-allocate all resources once they've outlived their purpose:
	// ------------------------------------------------------------------------
	gl.DeleteVertexArrays(1, &quadVao)
	gl.DeleteBuffers(1, &quadVbo)

	glfw.Terminate()
}

package main

import (
	"learn_opengl/gl"
	"log"
	"unsafe"

	"github.com/huoshan017/go-stbi"

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
	window.SetFramebufferSizeCallback(func(_ *glfw.Window, width, height int) {
		gl.Viewport(0, 0, int32(width), int32(height))
	})

	// build and compile our shader zprogram
	// -------------------------------------
	var shader = gl.NewShader("4.1.texture.vs", "4.1.texture.fs")

	// set up vertex data (and buffer(s)) and configure vertex attributes
	// ------------------------------------------------------------------
	vertices := []float32{
		// position
		0.5, 0.5, 0.0, 1.0, 0.0, 0.0, 1.0, 1.0, // top right
		0.5, -0.5, 0.0, 0.0, 1.0, 0.0, 1.0, 0.0, // bottom right
		-0.5, -0.5, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, // bottom left
		-0.5, 0.5, 0.0, 1.0, 1.0, 0.0, 0.0, 1.0, // top left
	}
	indices := []uint32{
		0, 1, 3, // first triangle
		1, 2, 3, // second triangle
	}

	var vbo, vao, ebo uint32
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)
	gl.GenBuffers(1, &ebo)

	gl.BindVertexArray(vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, unsafe.Pointer(&vertices[0]), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, unsafe.Pointer(&indices[0]), gl.STATIC_DRAW)

	// position attribute
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 8*4, 0)
	gl.EnableVertexAttribArray(0)

	// color attribute
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 8*4, 3*4)
	gl.EnableVertexAttribArray(1)

	// texture coord attribute
	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, 8*4, 6*4)
	gl.EnableVertexAttribArray(2)

	// load and create texture
	// -----------------------
	var texture uint32
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	// set the texture wrapping parameters
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	// set texture filtering parameters
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

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

		// bind texture
		gl.BindTexture(gl.TEXTURE_2D, texture)

		// render container
		shader.Use()
		gl.BindVertexArray(vao)
		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, 0)

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
	gl.DeleteBuffers(1, &ebo)

	// glfw: terminate, clearing all previously allocated GLFW resources
	// -----------------------------------------------------------------
	glfw.Terminate()
}

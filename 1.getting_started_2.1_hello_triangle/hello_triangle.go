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

var (
	vertexShaderSource = `
#version 330 core
layout (location = 0) in vec3 aPos;
void main()
{
    gl_Position = vec4(aPos.x, aPos.y, aPos.z, 1.0);
}
` + "\x00"
	fragmentShaderSource = `
#version 330 core
out vec4 FragColor;
void main()
{
    FragColor = vec4(1.0f, 0.5f, 0.2f, 1.0f);
}
` + "\x00"
)

func main() {
	glfw.Init()
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)

	window, err := glfw.CreateWindow(SRC_WIDTH, SRC_HEIGHT, "LearnOpenGL", nil, nil)
	if err != nil {
		log.Fatalf("glfw create window err: %v", err)
	}

	window.MakeContextCurrent()
	gl.Init()
	window.SetFramebufferSizeCallback(func(_ *glfw.Window, width int, height int) {
		gl.Viewport(0, 0, int32(width), int32(height))
	})

	// build and compile our shader program
	// ------------------------------------------
	// vertex shader
	vertexShader := gl.CreateShader(gl.VERTEX_SHADER)
	gl.ShaderSource(vertexShader, vertexShaderSource)
	gl.CompileShader(vertexShader)
	// check for shader compile error
	var success int32
	gl.GetShaderiv(vertexShader, gl.COMPILE_STATUS, &success)
	if success == 0 {
		log.Fatalf("ERROR::SHADER::VERTEX::COMPILATION_FAILED: %v", gl.GetShaderInfoLog(vertexShader))
	}

	// fragment shader
	fragmentShader := gl.CreateShader(gl.FRAGMENT_SHADER)
	gl.ShaderSource(fragmentShader, fragmentShaderSource)
	gl.CompileShader(fragmentShader)
	// check for shader compile error
	gl.GetShaderiv(fragmentShader, gl.COMPILE_STATUS, &success)
	if success == 0 {
		log.Fatalf("ERROR::SHADER::FRAGMENT::COMPILATION_FAILED: %v\n", gl.GetShaderInfoLog(fragmentShader))
	}

	// link shaders
	shaderProgram := gl.CreateProgram()
	gl.AttachShader(shaderProgram, vertexShader)
	gl.AttachShader(shaderProgram, fragmentShader)
	gl.LinkProgram(shaderProgram)
	// check for link errors
	gl.GetProgramiv(shaderProgram, gl.LINK_STATUS, &success)
	if success == 0 {
		log.Fatalf("ERROR::SHADER::PROGRAM::LINKING_FAILED: %v\n", gl.GetProgramInfoLog(shaderProgram))
	}
	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	// set up vertex data (and buffer(s)) and configure vertex attributes
	// ------------------------------------------------------------------
	var vertices = []float32{
		-0.5, -0.5, 0.0, // left
		0.5, -0.5, 0.0, // right
		0.0, 0.5, 0.0, // top
	}

	var vbo, vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)
	// bind the Vertex Array Object first, then bind and set vertex buffer(s), and then configure vertex attribute(s)
	gl.BindVertexArray(vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, unsafe.Pointer(&vertices[0]), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, 0)
	gl.EnableVertexAttribArray(0)

	// note that this is allowed, the call to gl.VertexAttribPointer registered VBO as the vertex attribute's bound vertex buffer object so afterwards we can safely unbind
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	// You can unbind the VAO afterwards so other VAO calls won't accidentally modify this VAO, but this rarely happens. Modifying other
	// VAOs requires a call to gl.BindVertexArray anyways so we generally don't unbind VAOs (nor VBOs) when it's not directly necessary
	gl.BindVertexArray(0)

	for !window.ShouldClose() {
		// input
		if window.GetKey(glfw.KeyEscape) == glfw.Press {
			window.SetShouldClose(true)
		}

		// render
		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		// draw our first triangle
		gl.UseProgram(shaderProgram)
		gl.BindVertexArray(vao)
		gl.DrawArrays(gl.TRIANGLES, 0, 3)

		window.SwapBuffers()
		glfw.PollEvents()
	}

	gl.DeleteVertexArrays(1, &vao)
	gl.DeleteBuffers(1, &vbo)
	gl.DeleteProgram(shaderProgram)

	glfw.Terminate()
}

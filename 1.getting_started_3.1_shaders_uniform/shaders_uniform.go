package main

import (
	"learn_opengl/gl"
	"log"
	"math"
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
void main() {
	gl_Position = vec4(aPos, 1.0);
}
` + "\x00"
	fragmentShaderSource = `
#version 330 core
out vec4 FragColor;
uniform vec4 ourColor;
void main() {
	FragColor = ourColor;
}
` + "\x00"
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
	// vertex shader
	vertexShader := gl.CreateShader(gl.VERTEX_SHADER)
	gl.ShaderSource(vertexShader, vertexShaderSource)
	gl.CompileShader(vertexShader)
	// check for shader compile error
	var success int32
	gl.GetShaderiv(vertexShader, gl.COMPILE_STATUS, &success)
	if success == 0 {
		log.Fatalf("ERROR::SHADER::VERTEX::COMPILATION_FAILED\n%v", gl.GetShaderInfoLog(vertexShader))
	}

	// fragment shader
	fragmentShader := gl.CreateShader(gl.FRAGMENT_SHADER)
	gl.ShaderSource(fragmentShader, fragmentShaderSource)
	gl.CompileShader(fragmentShader)
	// check for shader compile error
	gl.GetShaderiv(fragmentShader, gl.COMPILE_STATUS, &success)
	if success == 0 {
		log.Fatalf("ERROR::SHADER::FRAGMENT::COMPILATION_FAILED\n%v", gl.GetShaderInfoLog(fragmentShader))
	}

	// link shaders
	shaderProgram := gl.CreateProgram()
	gl.AttachShader(shaderProgram, vertexShader)
	gl.AttachShader(shaderProgram, fragmentShader)
	gl.LinkProgram(shaderProgram)
	// check for link errors
	gl.GetProgramiv(shaderProgram, gl.LINK_STATUS, &success)
	if success == 0 {
		log.Fatalf("ERROR::SHADER::PROGRAM::LINKING_FAILED\n%v", gl.GetProgramInfoLog(shaderProgram))
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	// set up vertex data (and buffer(s)) and configure vertex attributes
	// ------------------------------------------------------------------
	vertices := []float32{
		0.5, -0.5, 0.0, // bottom right
		-0.5, -0.5, 0.0, // bottom left
		0.0, 0.5, 0.0, // top
	}

	var vbo, vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)
	gl.BindVertexArray(vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, unsafe.Pointer(&vertices[0]), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, 0)
	gl.EnableVertexAttribArray(0)

	gl.BindVertexArray(vao)

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

		//be sure to activate the shader before any calls to glUniform
		gl.UseProgram(shaderProgram)

		// update shader uniform
		timeValue := glfw.GetTime()
		greenValue := math.Sin(timeValue)/2.0 + 0.5
		vertexColorLocation := gl.GetUniformLocation(shaderProgram, "ourColor\x00")
		gl.Uniform4f(vertexColorLocation, 0.0, float32(greenValue), 0.0, 1.0)

		// render the triangle
		gl.DrawArrays(gl.TRIANGLES, 0, 3)

		// glfw: swap buffers and poll IO events (keys pressed/released, mouse moved etc.)
		// -------------------------------------------------------------------------------
		window.SwapBuffers()
		glfw.PollEvents()
	}

	// optional: de-allocate all resources once they've outlived their purpose:
	// ------------------------------------------------------------------------
	gl.DeleteVertexArrays(1, &vao)
	gl.DeleteBuffers(1, &vbo)
	gl.DeleteProgram(shaderProgram)

	// glfw: terminate, clearing all previously allocated GLFW resources
	// -----------------------------------------------------------------
	glfw.Terminate()
}

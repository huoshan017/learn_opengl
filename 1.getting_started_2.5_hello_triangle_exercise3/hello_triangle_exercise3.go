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
void main() {
	gl_Position = vec4(aPos.x, aPos.y, aPos.z, 1.0);
}
` + "\x00"
	fragmentShaderSourceOrange = `
#version 330 core
out vec4 FragColor;
void main() {
	FragColor = vec4(1.0f, 0.5f, 0.2f, 1.0f);
}
` + "\x00"
	fragmentShaderSourceYellow = `
#version 330 core
out vec4 FragColor;
void main() {
	FragColor = vec4(1.0f, 1.0f, 0.0f, 1.0f);
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
	// we skipped compile log checks this time for readability (if you do encounter issues, add the compile-checks! see previous code samples)
	vertexShader := gl.CreateShader(gl.VERTEX_SHADER)
	fragmentShaderOrange := gl.CreateShader(gl.FRAGMENT_SHADER)
	fragmentShaderYellow := gl.CreateShader(gl.FRAGMENT_SHADER)
	shaderProgramOrange := gl.CreateProgram()
	shaderProgramYellow := gl.CreateProgram()
	gl.ShaderSource(vertexShader, vertexShaderSource)
	gl.CompileShader(vertexShader)
	// check for shader compile error
	var success int32
	gl.GetShaderiv(vertexShader, gl.COMPILE_STATUS, &success)
	if success == 0 {
		log.Fatalf("ERROR::SHADER::VERTEX::COMPILATION_FAILED: %v", gl.GetShaderInfoLog(vertexShader))
	}
	gl.ShaderSource(fragmentShaderOrange, fragmentShaderSourceOrange)
	gl.CompileShader(fragmentShaderOrange)
	// check for shader compile error
	gl.GetShaderiv(fragmentShaderOrange, gl.COMPILE_STATUS, &success)
	if success == 0 {
		log.Fatalf("ERROR::SHADER::FRAGMENTORANGE::COMPILATION_FAILED: %v\n", gl.GetShaderInfoLog(fragmentShaderOrange))
	}
	gl.ShaderSource(fragmentShaderYellow, fragmentShaderSourceYellow)
	gl.CompileShader(fragmentShaderYellow)
	// check for shader compile error
	gl.GetShaderiv(fragmentShaderYellow, gl.COMPILE_STATUS, &success)
	if success == 0 {
		log.Fatalf("ERROR::SHADER::FRAGMENTYELLOW::COMPILATION_FAILED: %v\n", gl.GetShaderInfoLog(fragmentShaderYellow))
	}
	// link the first program object
	gl.AttachShader(shaderProgramOrange, vertexShader)
	gl.AttachShader(shaderProgramOrange, fragmentShaderOrange)
	gl.LinkProgram(shaderProgramOrange)
	// then link the second program object using a different fragment shader (but same vertex shader)
	// this is perfectly allowed since the inputs and outputs of both the vertex and fragment shaders are equally matched.
	gl.AttachShader(shaderProgramYellow, vertexShader)
	gl.AttachShader(shaderProgramYellow, fragmentShaderYellow)
	gl.LinkProgram(shaderProgramYellow)

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShaderOrange)
	gl.DeleteShader(fragmentShaderYellow)

	// set up vertex data (and buffer(s)) and configure vertex attributes
	// ------------------------------------------------------------------
	var firstTriangle = []float32{
		-0.9, -0.5, 0.0, // left
		0.0, -0.5, 0.0, // right
		-0.45, 0.5, 0.0, // top
	}
	var secondTriangle = []float32{
		0.0, -0.5, 0.0, // left
		0.9, -0.5, 0.0, // right
		0.45, 0.5, 0.0, // top
	}

	var vaos, vbos [2]uint32
	gl.GenVertexArrays(2, &vaos[0]) // we can also generate multiple vaos or buffers at the same time
	gl.GenBuffers(2, &vbos[0])
	// first triangle setup
	// --------------------
	gl.BindVertexArray(vaos[0])
	gl.BindBuffer(gl.ARRAY_BUFFER, vbos[0])
	gl.BufferData(gl.ARRAY_BUFFER, len(firstTriangle)*4, unsafe.Pointer(&firstTriangle[0]), gl.STATIC_DRAW)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, 0) // Vertex attributes stay the same
	gl.EnableVertexAttribArray(0)
	// second triangle setup
	// ---------------------
	gl.BindVertexArray(vaos[1])             // note that we bind to a different vao now
	gl.BindBuffer(gl.ARRAY_BUFFER, vbos[1]) // and a different vbo
	gl.BufferData(gl.ARRAY_BUFFER, len(secondTriangle)*4, unsafe.Pointer(&secondTriangle[0]), gl.STATIC_DRAW)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, 0) // because the vertex data is tightly packed we can also specify 0 as the vertex attribute's stride to let OpenGL figure it out
	gl.EnableVertexAttribArray(0)

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

		// now when we draw the triangle we first use the vertex and orange fragment shader from the first program
		gl.UseProgram(shaderProgramOrange)
		// draw the first triangle using the data from our first vao
		gl.BindVertexArray(vaos[0])
		gl.DrawArrays(gl.TRIANGLES, 0, 3) // this call should output an orange triangle
		// then we draw the second triangle using the data from second vao
		// when we draw the second triangle we want to use a different shader program so we switch to the shader program with our yellow fragment shader.
		gl.UseProgram(shaderProgramYellow)
		gl.BindVertexArray(vaos[1])
		gl.DrawArrays(gl.TRIANGLES, 0, 3) // this call should output a yellow triangle

		// glfw: swap buffers and poll IO events (keys pressed/released, mouse moved etc.)
		// -------------------------------------------------------------------------------
		window.SwapBuffers()
		glfw.PollEvents()
	}

	// optional: de-allocate all resources once they've outlived their purpose:
	// ------------------------------------------------------------------------
	gl.DeleteVertexArrays(2, &vaos[0])
	gl.DeleteBuffers(2, &vbos[0])
	gl.DeleteProgram(shaderProgramOrange)
	gl.DeleteProgram(shaderProgramYellow)

	// glfw: terminate, clearing all previously allocated GLFW resources
	// -----------------------------------------------------------------
	glfw.Terminate()
}

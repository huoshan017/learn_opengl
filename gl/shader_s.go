package gl

import (
	"io/ioutil"
	"log"

	"github.com/go-gl/mathgl/mgl32"
)

type Shader struct {
	id uint32
}

func vsAndFs(vertexPath, fragmentPath string) (uint32, uint32) {
	vertexShaderSource, err := ioutil.ReadFile(vertexPath)
	if err != nil {
		log.Fatalf("Failed to read vertex shader file %v", vertexPath)
	}
	fragmentShaderSource, err := ioutil.ReadFile(fragmentPath)
	if err != nil {
		log.Fatalf("Failed to read fragment shader file %v", fragmentPath)
	}
	// vertex shader
	vertexShader := CreateShader(VERTEX_SHADER)
	ShaderSource(vertexShader, string(vertexShaderSource)+"\x00")
	CompileShader(vertexShader)
	checkCompileErrors(vertexShader, "VERTEX")
	// fragment shader
	fragmentShader := CreateShader(FRAGMENT_SHADER)
	ShaderSource(fragmentShader, string(fragmentShaderSource)+"\x00")
	CompileShader(fragmentShader)
	checkCompileErrors(fragmentShader, "FRAGMENT")
	return vertexShader, fragmentShader
}

func NewShader(vertexPath, fragmentPath string) Shader {
	vertexShader, fragmentShader := vsAndFs(vertexPath, fragmentPath)
	// shader program
	id := CreateProgram()
	AttachShader(id, vertexShader)
	AttachShader(id, fragmentShader)
	LinkProgram(id)
	checkCompileErrors(id, "PROGRAM")
	// delete the shaders as they're linked into our program now and no longer necessary
	DeleteShader(vertexShader)
	DeleteShader(fragmentShader)
	return Shader{id: id}
}

func NewShader2(vertexPath, fragmentPath, geometryPath string) Shader {
	vertexShader, fragmentShader := vsAndFs(vertexPath, fragmentPath)
	geometryShaderSource, err := ioutil.ReadFile(geometryPath)
	if err != nil {
		log.Fatalf("Failed to read geometry shader file %v", geometryPath)
	}
	geometryShader := CreateShader(GEOMETRY_SHADER)
	ShaderSource(geometryShader, string(geometryShaderSource)+"\x00")
	CompileShader(geometryShader)
	checkCompileErrors(geometryShader, "GEOMETRY")
	// shader program
	id := CreateProgram()
	AttachShader(id, vertexShader)
	AttachShader(id, fragmentShader)
	AttachShader(id, geometryShader)
	LinkProgram(id)
	checkCompileErrors(id, "PROGRAM")
	// delete the shaders as they're linked into our program now and no longer necessary
	DeleteShader(vertexShader)
	DeleteShader(fragmentShader)
	DeleteShader(geometryShader)
	return Shader{id: id}
}

func (s *Shader) Id() uint32 {
	return s.id
}

func (s *Shader) Use() {
	UseProgram(s.id)
}

func (s *Shader) SetBool(name string, value bool) {
	v := func() int32 {
		if value {
			return 1
		} else {
			return 0
		}
	}()
	Uniform1i(GetUniformLocation(s.id, name), v)
}

func (s *Shader) SetInt32(name string, value int32) {
	Uniform1i(GetUniformLocation(s.id, name), value)
}

func (s *Shader) SetFloat32(name string, value float32) {
	Uniform1f(GetUniformLocation(s.id, name), value)
}

func (s *Shader) SetVec2(name string, value *mgl32.Vec2) {
	Uniform2fv(GetUniformLocation(s.id, name), 1, &((*value)[0]))
}

func (s *Shader) SetVec2WithXY(name string, x, y float32) {
	Uniform2f(GetUniformLocation(s.id, name), x, y)
}

func (s *Shader) SetVec3(name string, value *mgl32.Vec3) {
	Uniform3fv(GetUniformLocation(s.id, name), 1, &((*value)[0]))
}

func (s *Shader) SetVec3WithXYZ(name string, x, y, z float32) {
	Uniform3f(GetUniformLocation(s.id, name), x, y, z)
}

func (s *Shader) SetVec4(name string, value *mgl32.Vec4) {
	Uniform4fv(GetUniformLocation(s.id, name), 1, &((*value)[0]))
}

func (s *Shader) SetVec4WithXYZW(name string, x, y, z, w float32) {
	Uniform4f(GetUniformLocation(s.id, name), x, y, z, w)
}

func (s *Shader) SetMat2(name string, value *mgl32.Mat2) {
	UniformMatrix2fv(GetUniformLocation(s.id, name), 1, false, &((*value)[0]))
}

func (s *Shader) SetMat3(name string, value *mgl32.Mat3) {
	UniformMatrix3fv(GetUniformLocation(s.id, name), 1, false, &((*value)[0]))
}

func (s *Shader) SetMat4(name string, value *mgl32.Mat4) {
	UniformMatrix4fv(GetUniformLocation(s.id, name), 1, false, &((*value)[0]))
}

func checkCompileErrors(shader uint32, typ string) {
	var success int32
	if typ == "VERTEX" || typ == "FRAGMENT" || typ == "GEOMETRY" {
		GetShaderiv(shader, COMPILE_STATUS, &success)
		if success == FALSE {
			log.Fatalf("ERROR::SHADER_COMPILATION_ERROR of type: %v\n%v", typ, GetShaderInfoLog(shader))
		}
	} else if typ == "PROGRAM" {
		GetProgramiv(shader, LINK_STATUS, &success)
		if success == FALSE {
			log.Fatalf("ERROR::PROGRAM_LINKING_ERROR of type: %v\n%v", typ, GetProgramInfoLog(shader))
		}
	} else {
		log.Fatalf("ERROR unknown type: %v", typ)
	}
}

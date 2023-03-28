package gl

import (
	"strconv"
	"unsafe"

	"github.com/go-gl/mathgl/mgl32"
)

const (
	MAX_BONE_INFLUENCE = 4
)

type Vertex struct {
	// position
	Position mgl32.Vec3
	// normal
	Normal mgl32.Vec3
	// texCoords
	TexCoords mgl32.Vec2
	// tangent
	Tangent mgl32.Vec3
	// bitangent
	Bitangent mgl32.Vec3
	// bone indexes which will influence this vertex
	BoneIds [MAX_BONE_INFLUENCE]int32
	// weights from each bone
	Weights [MAX_BONE_INFLUENCE]float32
}

var (
	_dummyVertex = Vertex{}
)

type Texture struct {
	id   uint32
	typ  string
	path string
}

func NewTexture(id uint32, typ, path string) Texture {
	return Texture{id: id, typ: typ, path: path}
}

func (t Texture) Id() uint32 {
	return t.id
}

func (t Texture) Type() string {
	return t.typ
}

func (t Texture) Path() string {
	return t.path
}

type Mesh struct {
	vertices []Vertex
	indices  []uint32
	textures []Texture
	vao      uint32
	vbo      uint32
	ebo      uint32
}

func NewMesh(vertices []Vertex, indices []uint32, textures []Texture) Mesh {
	mesh := Mesh{
		vertices: vertices,
		indices:  indices,
		textures: textures,
	}
	mesh.setupMesh()
	return mesh
}

func (m Mesh) Vao() uint32 {
	return m.vao
}

func (m Mesh) Indices() []uint32 {
	return m.indices
}

func (m *Mesh) Draw(shader *Shader) {
	// bind appropriate textures
	var (
		diffuseNr  = 1
		specularNr = 1
		normalNr   = 1
		heightNr   = 1
	)
	for i := int32(0); i < int32(len(m.textures)); i++ {
		ActiveTexture(TEXTURE0 + uint32(i)) // active proper texture unit before binding
		// retrieve texture number (the N in diffuse_textureN)
		var number string
		var name = m.textures[i].typ
		if name == "texture_diffuse" {
			number = strconv.Itoa(diffuseNr)
			diffuseNr++
		} else if name == "texture_specular" {
			number = strconv.Itoa(specularNr)
			specularNr++
		} else if name == "texture_normal" {
			number = strconv.Itoa(normalNr)
			normalNr++
		} else if name == "texture_height" {
			number = strconv.Itoa(heightNr)
			heightNr++
		}

		// now set the sampler to the correct texture unit
		Uniform1i(GetUniformLocation(shader.Id(), name+number+"\x00"), i)
		// and finally bind the texture
		BindTexture(TEXTURE_2D, m.textures[i].id)

		//log.Printf("mesh %p bind texture %v", m, m.textures[i].id)
	}

	// draw mesh
	BindVertexArray(m.vao)
	DrawElements(TRIANGLES, int32(len(m.indices)), UNSIGNED_INT, 0)
	BindVertexArray(0)

	//log.Printf("mesh %p drawn, indices len %v", m, len(m.indices))

	// always good practice to set everything back to defaults once configured
	ActiveTexture(TEXTURE0)
}

func (m *Mesh) setupMesh() {
	// create buffers/arrays
	GenVertexArrays(1, &m.vao)
	GenBuffers(1, &m.vbo)
	GenBuffers(1, &m.ebo)

	BindVertexArray(m.vao)
	// load data into vertex buffers
	BindBuffer(ARRAY_BUFFER, m.vbo)
	// A great thing about structs is that their memory layout is sequential for all its items.
	// The effect is that we can simply pass a pointer to the struct and it translates perfectly to a mgl32.Vec3/2 array which
	// again translates to 3/2 floats which translates to a byte array.
	BufferData(ARRAY_BUFFER, len(m.vertices)*int(unsafe.Sizeof(_dummyVertex)), unsafe.Pointer(&m.vertices[0]), STATIC_DRAW)

	BindBuffer(ELEMENT_ARRAY_BUFFER, m.ebo)
	BufferData(ELEMENT_ARRAY_BUFFER, len(m.indices)*4, unsafe.Pointer(&m.indices[0]), STATIC_DRAW)

	// set the vertex attribute pointers
	// vertex Positions
	EnableVertexAttribArray(0)
	VertexAttribPointer(0, 3, FLOAT, false, int32(unsafe.Sizeof(_dummyVertex)), 0)
	// vertex normals
	EnableVertexAttribArray(1)
	VertexAttribPointer(1, 3, FLOAT, false, int32(unsafe.Sizeof(_dummyVertex)), int(unsafe.Offsetof(_dummyVertex.Normal)))
	// vertex texture coords
	EnableVertexAttribArray(2)
	VertexAttribPointer(2, 2, FLOAT, false, int32(unsafe.Sizeof(_dummyVertex)), int(unsafe.Offsetof(_dummyVertex.TexCoords)))
	// vertex tangent
	EnableVertexAttribArray(3)
	VertexAttribPointer(3, 3, FLOAT, false, int32(unsafe.Sizeof(_dummyVertex)), int(unsafe.Offsetof(_dummyVertex.Tangent)))
	// vertex bitangent
	EnableVertexAttribArray(4)
	VertexAttribPointer(4, 3, FLOAT, false, int32(unsafe.Sizeof(_dummyVertex)), int(unsafe.Offsetof(_dummyVertex.Bitangent)))
	// ids
	EnableVertexAttribArray(5)
	VertexAttribIPointer(5, 4, INT, int32(unsafe.Sizeof(_dummyVertex)), int(unsafe.Offsetof(_dummyVertex.BoneIds)))
	// weights
	EnableVertexAttribArray(6)
	VertexAttribPointer(6, 4, FLOAT, false, int32(unsafe.Sizeof(_dummyVertex)), int(unsafe.Offsetof(_dummyVertex.Weights)))

	BindVertexArray(0)
}

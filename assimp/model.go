package assimp

import (
	"learn_opengl/gl"
	"log"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/huoshan017/assimp"
)

type Model struct {
	textureLoaded   []gl.Texture
	meshes          []gl.Mesh
	directory       string
	gammaCorrection bool
}

func NewModel(path string, gamma bool) *Model {
	model := &Model{
		gammaCorrection: gamma,
	}
	model.load(path)
	return model
}

func NewModelDefault(path string) *Model {
	return NewModel(path, false)
}

func (m *Model) TextureLoaded() []gl.Texture {
	return m.textureLoaded
}

func (m *Model) Meshes() []gl.Mesh {
	return m.meshes
}

func (m *Model) Draw(shader *gl.Shader) {
	for i := 0; i < len(m.meshes); i++ {
		m.meshes[i].Draw(shader)
	}
}

// loads a model with supported ASSIMP extensions from file and stores the resulting meshes in the meshes vector.
func (m *Model) load(path string) {
	// read file via assimp
	scene := assimp.ImportFile(path, uint(assimp.Process_Triangulate|assimp.Process_GenSmoothNormals|assimp.Process_FlipUVs|assimp.Process_CalcTangentSpace))

	// check for errors
	if scene == nil || scene.Flags()&assimp.SceneFlags_Incomplete > 0 || scene.RootNode() == nil {
		log.Fatalf("ERROR::ASSIMP:: %v", assimp.GetErrorString())
	}

	// retrieve the directory path of the filepath
	var found bool
	pathBytes := []byte(path)
	for i := len(pathBytes) - 1; i >= 0; i-- {
		if pathBytes[i] == '/' {
			m.directory = string(pathBytes[:i])
			found = true
			break
		}
	}
	if !found {
		log.Fatalf("ERROR::STRINGS:: not found sperate")
	}

	// process ASSIMP's root node recursively
	m.processNode(scene.RootNode(), scene)
}

// processes a node in a recursive fashion. Processes each individual mesh located at the node and repeats this process on its children nodes (if any).
func (m *Model) processNode(node *assimp.Node, scene *assimp.Scene) {
	meshes := scene.Meshes()
	// process each mesh located at the current node
	for i := 0; i < node.NumMeshes(); i++ {
		// the node object only contains indices to index the actual objects in the scene.
		// the scene contains all the data, node is just to keep stuff organized (like relations between nodes).
		nMeshes := node.Meshes()
		mesh := meshes[nMeshes[i]]
		glmesh := m.processMesh(mesh, scene)
		m.meshes = append(m.meshes, glmesh)
	}
	// after we've processed all of the meshes (if any) we then recursively process each of the children nodes
	nodeChildren := node.Children()
	for i := 0; i < node.NumChildren(); i++ {
		m.processNode(nodeChildren[i], scene)
	}

	log.Printf("num of meshes %v, num of nodeChildren %v", node.NumMeshes(), node.NumChildren())
}

func (m *Model) processMesh(mesh *assimp.Mesh, scene *assimp.Scene) gl.Mesh {
	// data to fill
	var (
		vertices []gl.Vertex
		indices  []uint32
		textures []gl.Texture
	)

	meshVertices := mesh.Vertices()
	meshNormals := mesh.Normals()
	meshTangents := mesh.Tangents()
	meshBitangents := mesh.Bitangents()
	// walk through each of the mesh's vertices
	for i := 0; i < mesh.NumVertices(); i++ {
		var vertex gl.Vertex
		// positions
		aiVertex := &meshVertices[i]
		var vector = mgl32.Vec3{aiVertex.X(), aiVertex.Y(), aiVertex.Z()}
		vertex.Position = vector
		if meshNormals != nil {
			value := &meshNormals[i]
			vector = mgl32.Vec3{value.X(), value.Y(), value.Z()}
			vertex.Normal = vector
		}
		// texture coordinates
		c := mesh.TextureCoords(0)
		if len(c) > 0 { // does the mesh contain texture coordinates?
			// a vertex can contain up to 8 different texture coordinates. We thus make the assumption that we won't
			// use models where a vertex can have multiple texture coordinates so we always take the first set (0).
			var vec = mgl32.Vec2{c[i].X(), c[i].Y()}
			vertex.TexCoords = vec
			// tangent
			tangent := &meshTangents[i]
			vector = mgl32.Vec3{tangent.X(), tangent.Y(), tangent.Z()}
			vertex.Tangent = vector
			// bitangent
			bitangent := &meshBitangents[i]
			vector = mgl32.Vec3{bitangent.X(), bitangent.Y(), bitangent.Z()}
			vertex.Bitangent = vector
		} else {
			vertex.TexCoords = mgl32.Vec2{0.0, 0.0}
		}

		vertices = append(vertices, vertex)
	}

	meshNumFaces := mesh.NumFaces()
	meshFaces := mesh.Faces()
	// now wak through each of the mesh's faces (a face is a mesh its triangle) and retrieve the corresponding vertex indices.
	for i := 0; i < meshNumFaces; i++ {
		face := &meshFaces[i]
		faceNumIndices := face.NumIndices()
		// retrieve all indices of the face and store them in the indices vector
		for j := 0; j < int(faceNumIndices); j++ {
			indices = append(indices, face.Index(j))
		}
	}

	// process materials
	material := scene.Materials()[mesh.MaterialIndex()]
	// we assume a convention for sampler names in the shaders. Each diffuse texture should be named
	// as 'texture_diffuseN' where N is a sequential number ranging from 1 to MAX_SAMPLER_NUMBER.
	// Same applies to other texture as the following list summarizes:
	// diffuse: texture_diffuseN
	// specular: texture_specularN
	// normal: texture_normalN

	// 1. diffuse maps
	diffuseMaps := m.loadMaterialTextures(material, assimp.TextureType_Diffuse, "texture_diffuse")
	textures = append(textures, diffuseMaps...)
	// 2. specular maps
	specularMaps := m.loadMaterialTextures(material, assimp.TextureType_Specular, "texture_specular")
	textures = append(textures, specularMaps...)
	// 3. normal maps
	normalMaps := m.loadMaterialTextures(material, assimp.TextureType_Height, "texture_normal")
	textures = append(textures, normalMaps...)
	// 4. height maps
	heightMaps := m.loadMaterialTextures(material, assimp.TextureType_Ambient, "texture_height")
	textures = append(textures, heightMaps...)

	log.Printf("len(vertices)=%v len(indices)=%v len(textures)=%v", len(vertices), len(indices), len(textures))
	return gl.NewMesh(vertices, indices, textures)
}

func (m *Model) loadMaterialTextures(mat *assimp.Material, typ assimp.TextureType, typeName string) []gl.Texture {
	var textures []gl.Texture
	matTextureCount := mat.GetMaterialTextureCount(typ)
	for i := 0; i < matTextureCount; i++ {
		var str string
		str, _, _, _, _, _, _, _ = mat.GetMaterialTexture(typ, i)
		if str == "" {
			log.Fatalf("material %p get empty string with type %v and index %v", mat, typ, i)
		}
		var skip bool
		for j := 0; j < len(m.textureLoaded); j++ {
			if str == m.textureLoaded[j].Path() {
				textures = append(textures, m.textureLoaded[j])
				skip = true
				break
			}
		}
		if !skip {
			// if texture hasn't been loaded already, load it
			log.Printf("before new texture, str %v", str)
			textureId := gl.TextureFromFile(str, m.directory, false)
			texture := gl.NewTexture(textureId, typeName, str)
			textures = append(textures, texture)
			m.textureLoaded = append(m.textureLoaded, texture) // store it as texture loaded for entire model, to ensure we won't unnecessary load duplicate textures.
			log.Printf("new texture %v, str %v", texture.Id(), str)
		}
	}
	log.Printf("m.textureLoaded len %v,  textures len %v", len(m.textureLoaded), len(textures))
	return textures
}

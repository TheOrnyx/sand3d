package main

import "github.com/go-gl/gl/v4.6-core/gl"

type GraphicsResources struct {
	VAO uint32
	VBO uint32
	EBO uint32
	Vertices []float32
}

// CreateResources creates a GraphicsResources struct instance to hold important stuff
func CreateResources(vertices []float32) *GraphicsResources {
	n := new(GraphicsResources)
	n.Vertices = vertices
	n.MakeObjects()
	
	return n
}

// MakeObjects create the objects such as the VBO, VAO etc
func (g *GraphicsResources) MakeObjects()  {
	gl.GenVertexArrays(1, &g.VAO)
	gl.GenBuffers(1, &g.VBO)

	gl.BindVertexArray(g.VAO)

	gl.BindBuffer(gl.ARRAY_BUFFER, g.VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(g.Vertices)*4, gl.Ptr(g.Vertices), gl.STATIC_DRAW)

	//position stuff
	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 5*4, 0)
	gl.EnableVertexAttribArray(0)

	//texture coord stuff
	gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false, 5*4, 3*4)
	gl.EnableVertexAttribArray(1)
}

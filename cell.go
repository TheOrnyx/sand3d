package main

import (
	"github.com/go-gl/gl/v4.6-core/gl"
	glm "github.com/go-gl/mathgl/mgl32"
)

const ( //cell types
	AIR = iota
	DIRT
	WALL
	WATER
)

type Cell struct {
	PosX, PosY int32 //shouldn't matter since using a grid
	Type int //the cell type, should be zero'd at AIR
	// put some other stuff here
}

// Draw draws the cell with the specified translation
func (c *Cell) Draw(posX, posY, posZ float32, shader *shader)  {
	if c.Type != AIR {
		switch c.Type {
		case DIRT: 
			shader.SetBool("Water", false)
			model := glm.Ident4()
			model = model.Mul4(glm.Translate3D(posX, posY, posZ))
			model = model.Mul4(glm.Scale3D(CELL_SIZE_SCALAR, CELL_SIZE_SCALAR, CELL_SIZE_SCALAR))
			shader.SetMat4("model", &model)
			gl.DrawArrays(gl.TRIANGLES, 0, 36)

		case WATER:
			shader.SetBool("Water", true)
			model := glm.Ident4()
			model = model.Mul4(glm.Translate3D(posX, posY, posZ))
			model = model.Mul4(glm.Scale3D(CELL_SIZE_SCALAR, CELL_SIZE_SCALAR, CELL_SIZE_SCALAR))
			shader.SetMat4("model", &model)
			gl.DrawArrays(gl.TRIANGLES, 0, 36)
		}
	}
}

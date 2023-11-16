package main

import (
	"math/rand"
)

type World struct {
	Cells                [][][]Cell
	Visited              [][][]bool
	Width, Height, Depth int
}

func MakeWorld(width, height, depth int) *World {
	newWorld := new(World) //I hate my job
	newWorld.ResetCellGrid(width, height, depth)
	newWorld.Width, newWorld.Height, newWorld.Depth = width, height, depth
	return newWorld
}

// ResetCellGrid reset the cell grid and the visited grid
func (w *World) ResetCellGrid(width, height, depth int) {
	var cubeGrid = make([][][]Cell, width)
	var visitedGrid = make([][][]bool, width)
	for i := range cubeGrid {
		cubeGrid[i] = make([][]Cell, height)
		visitedGrid[i] = make([][]bool, height)
		for j := range cubeGrid[i] {
			cubeGrid[i][j] = make([]Cell, depth)
			visitedGrid[i][j] = make([]bool, depth)
			for k := range cubeGrid[i][j] {
				cubeGrid[i][j][k] = Cell{Type: AIR}
			}
		}
	}
	w.Cells = cubeGrid
	w.Visited = visitedGrid
}

// ResetVisitedGrid resets just the visited grid
func (w *World) ResetVisitedGrid(width, height, depth int) {
	var visitedGrid = make([][][]bool, width)
	for i := range visitedGrid {
		visitedGrid[i] = make([][]bool, height)
		for j := range visitedGrid[i] {
			visitedGrid[i][j] = make([]bool, depth)
		}
	}
	w.Visited = visitedGrid
}

// Draw draw the world
func (w *World) Draw(shader *shader)  {
	var startX, startY, startZ float32
	startX = -0.5 + 0.5*CELL_SIZE_SCALAR
	startY = -0.5 + 0.5*CELL_SIZE_SCALAR
	startZ = -0.5 + 0.5*CELL_SIZE_SCALAR
	
	for x := 0; x < w.Width; x++ {
		for y := 0; y < w.Height; y++ {
			for z := 0; z < w.Depth; z++ {
				posX := startX + float32(x) * CELL_SIZE_SCALAR
				posY := startY + float32(y) * CELL_SIZE_SCALAR
				posZ := startZ + float32(z) * CELL_SIZE_SCALAR
				w.Cells[x][y][z].Draw(posX, posY, posZ, shader)
			}
		}
	}
}

// ------------------------------ Adding Things ------------------------------
// AddCell Adds a cell of cellType to point at x,y,z
func (w *World) AddCell(x, y, z, cellType int)  {
	if w.IndexInRange(x, y, z) {
		w.Cells[x][y][z] = Cell{Type: cellType}
	}
}

// ------------------------------ Stuff for Updating ------------------------------

// Update updates the world
func (w *World) Update() {
	//maybe add stuff for like only updating sections so I can maybe goroutine it
	w.ResetVisitedGrid(w.Width, w.Height, w.Depth)
	for x := 0; x < w.Width; x++ {
		for y := 0; y < w.Height; y++ {
			for z := 0; z < w.Depth; z++ {
				if !w.Visited[x][y][z] {
					w.MoveCell(x, y, z)
				}
			}
		}
	}
}

// MoveCell attempts to move the cell
func (w *World) MoveCell(x, y, z int) {
	//probs do a switch for the move direction but for now I'm just doing down
	if w.IndexInRange(x, y-1, z) && w.Cells[x][y-1][z].Type == AIR {
		w.SwapCells(x, y, z, x, y-1, z)
		w.Visited[x][y-1][z] = true
	} else {
		move1 := w.IndexInRange(x-1, y-1, z) && w.Cells[x-1][y-1][z].Type == AIR
		move2 := w.IndexInRange(x+1, y-1, z) && w.Cells[x+1][y-1][z].Type == AIR
		move3 := w.IndexInRange(x, y-1, z-1) && w.Cells[x][y-1][z-1].Type == AIR
		move4 := w.IndexInRange(x, y-1, z+1) && w.Cells[x][y-1][z+1].Type == AIR
		
		move5 := w.IndexInRange(x-1, y-1, z+1) && w.Cells[x-1][y-1][z+1].Type == AIR
		move6 := w.IndexInRange(x+1, y-1, z-1) && w.Cells[x+1][y-1][z-1].Type == AIR
		move7 := w.IndexInRange(x-1, y-1, z-1) && w.Cells[x-1][y-1][z-1].Type == AIR
		move8 := w.IndexInRange(x+1, y-1, z+1) && w.Cells[x+1][y-1][z+1].Type == AIR

		
		optionRan := false
		if !move1 && !move2 && !move3 && !move4 && !move5 && !move6 && !move7 && !move8 {
			return
		}
		
		for !optionRan {
			choice := rand.Int31n(8)
			switch {
			case move1 && choice == 0:
				w.SwapCells(x, y, z, x-1, y-1, z)
				optionRan = true
			case move2 && choice == 1:
				w.SwapCells(x, y, z, x+1, y-1, z)
				optionRan = true
			case move3 && choice == 2:
				w.SwapCells(x, y, z, x, y-1, z-1)
				optionRan = true
			case move4 && choice == 3:
				w.SwapCells(x, y, z, x, y-1, z+1)
				optionRan = true
			case move5 && choice == 4:
				w.SwapCells(x, y, z, x-1, y-1, z+1)
				optionRan = true
			case move6 && choice == 5:
				w.SwapCells(x, y, z, x+1, y-1, z-1)
				optionRan = true
			case move7 && choice == 6:
				w.SwapCells(x, y, z, x-1, y-1, z-1)
				optionRan = true
			case move8 && choice == 7:
				w.SwapCells(x, y, z, x+1, y-1, z+1)
				optionRan = true
			}
		}
	}
}

// SwapCells swaps two cells with each other
func (w *World) SwapCells(x1, y1, z1, x2, y2, z2 int)  {
	cell2 := w.Cells[x2][y2][z2]
	w.Cells[x2][y2][z2] = w.Cells[x1][y1][z1]
	w.Cells[x1][y1][z1] = cell2
	w.Visited[x1][y1][z1] = true
	w.Visited[x2][y2][z2] = true
}

// IndexInRange check if proposed movement is in range of the world
// Remember to change the var before passing it through since you check the changed position
func (w *World) IndexInRange(x, y, z int) bool {
	xInRange := x >= 0 && x < w.Width
	yInRange := y >= 0 && y < w.Height //might need to make this -1
	zInRange := z >= 0 && z < w.Depth

	return xInRange && yInRange && zInRange
}

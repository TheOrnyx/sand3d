package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/go-gl/gl/v4.6-core/gl"
	glm "github.com/go-gl/mathgl/mgl32"
	"github.com/veandco/go-sdl2/sdl"
)

const WIN_WIDTH, WIN_HEIGHT = 1000, 1000
const FRAME_RATE = 60
const frameDelay = 1000/FRAME_RATE
const WORLD_SIZE = 60 //the amount of cells in each direction (so the amount of cubes should be WORLD_SIZE^3)
const CELL_SIZE_SCALAR = 1.0 / WORLD_SIZE //scalar to use for the size of the cubes

var vertices = []float32{ //the cube vertices, possible move
	-0.5, -0.5, -0.5, 0.0, 0.0,
	0.5, -0.5, -0.5, 1.0, 0.0,
	0.5, 0.5, -0.5, 1.0, 1.0,
	0.5, 0.5, -0.5, 1.0, 1.0,
	-0.5, 0.5, -0.5, 0.0, 1.0,
	-0.5, -0.5, -0.5, 0.0, 0.0,

	-0.5, -0.5, 0.5, 0.0, 0.0,
	0.5, -0.5, 0.5, 1.0, 0.0,
	0.5, 0.5, 0.5, 1.0, 1.0,
	0.5, 0.5, 0.5, 1.0, 1.0,
	-0.5, 0.5, 0.5, 0.0, 1.0,
	-0.5, -0.5, 0.5, 0.0, 0.0,

	-0.5, 0.5, 0.5, 1.0, 0.0,
	-0.5, 0.5, -0.5, 1.0, 1.0,
	-0.5, -0.5, -0.5, 0.0, 1.0,
	-0.5, -0.5, -0.5, 0.0, 1.0,
	-0.5, -0.5, 0.5, 0.0, 0.0,
	-0.5, 0.5, 0.5, 1.0, 0.0,

	0.5, 0.5, 0.5, 1.0, 0.0,
	0.5, 0.5, -0.5, 1.0, 1.0,
	0.5, -0.5, -0.5, 0.0, 1.0,
	0.5, -0.5, -0.5, 0.0, 1.0,
	0.5, -0.5, 0.5, 0.0, 0.0,
	0.5, 0.5, 0.5, 1.0, 0.0,

	-0.5, -0.5, -0.5, 0.0, 1.0,
	0.5, -0.5, -0.5, 1.0, 1.0,
	0.5, -0.5, 0.5, 1.0, 0.0,
	0.5, -0.5, 0.5, 1.0, 0.0,
	-0.5, -0.5, 0.5, 0.0, 0.0,
	-0.5, -0.5, -0.5, 0.0, 1.0,

	-0.5, 0.5, -0.5, 0.0, 1.0,
	0.5, 0.5, -0.5, 1.0, 1.0,
	0.5, 0.5, 0.5, 1.0, 0.0,
	0.5, 0.5, 0.5, 1.0, 0.0,
	-0.5, 0.5, 0.5, 0.0, 0.0,
	-0.5, 0.5, -0.5, 0.0, 1.0,
}

var drawBoundingBox = true

var graphics *GraphicsResources
var camera *Camera = MakeCamera(glm.Vec3{0, 0, 3}, glm.Vec3{0, 1, 0}, INIT_YAW, INIT_PITCH)
var deltaTime, lastFrame float32
var lastMouseX, lastMouseY int32 = WIN_WIDTH / 2, WIN_HEIGHT / 2
var drawType int = DIRT
var selectionY float32 = WORLD_SIZE-1 //the plane at which you make selections from

func main() {
	fmt.Println("begin")
	runtime.LockOSThread()

	// ------------------------------ Window Setup ------------------------------

	window, err := sdl.CreateWindow("the zinger", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, WIN_WIDTH, WIN_HEIGHT, sdl.WINDOW_OPENGL|sdl.WINDOW_ALLOW_HIGHDPI)
	if err != nil {
		log.Fatal(err)
	}
	defer window.Destroy()

	context, err := window.GLCreateContext()
	if err != nil {
		log.Fatal("could not create context: ", err)
	}
	defer sdl.GLDeleteContext(context)

	if err := gl.Init(); err != nil {
		log.Fatal("could not initialize OpenGL: ", err)
	}
	gl.Enable(gl.DEPTH_TEST)
	gl.Viewport(0, 0, WIN_WIDTH, WIN_HEIGHT)
	sdl.SetRelativeMouseMode(true)

	// ------------------------------ Other setups ------------------------------

	graphics = CreateResources(vertices)

	//Shader setup
	vertSource, err := os.ReadFile("./data/world.vs")
	if err != nil {
		log.Fatal(err)
	}
	fragSource, err := os.ReadFile("./data/world.fs")
	if err != nil {
		log.Fatal(err)
	}
	worldShader, err := NewShader(vertSource, fragSource, "world shader")
	if err != nil {
		log.Fatal(err)
	}

	texture, err := NewTexture("./data/dirt.png", worldShader)
	if err != nil {
		log.Fatal(err)
	}

	world := MakeWorld(WORLD_SIZE, WORLD_SIZE, WORLD_SIZE)
	// world.Cells[0][20][0].Type = DIRT

	// ------------------------------ Main Loop ------------------------------
	for !handleEvents() {
		startTime := time.Now()
		currentFrame := float32(sdl.GetTicks64()) / 1000
		deltaTime = currentFrame - lastFrame
		lastFrame = currentFrame

		//clear the screen
		gl.ClearColor(0, 0, 0, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		//update the world
		world.Update()
		world.Cells[10][59][10].Type = DIRT
		world.Cells[15][59][10].Type = DIRT
		world.Cells[10][59][15].Type = DIRT
		world.Cells[15][59][15].Type = DIRT
		world.Cells[30][40][10].Type = WATER
		world.Cells[35][54][15].Type = WATER
		world.Cells[40][27][15].Type = WATER

		//bind textures
		texture.Bind(0)

		//camera stuff
		proj := glm.Perspective(glm.DegToRad(camera.Zoom), WIN_WIDTH/WIN_HEIGHT, 0.1, 100.0)
		worldShader.SetMat4("projection", &proj)
		view := camera.GetViewMatrix()
		worldShader.SetMat4("view", &view)

		//draw the outer cube
		if drawBoundingBox {
			worldShader.SetBool("white", true)
			gl.BindVertexArray(graphics.VAO)
			model := glm.Ident4()
			worldShader.SetMat4("model", &model)
			gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
			gl.DrawArrays(gl.TRIANGLES, 0, 36)
		}

		world.GetCameraCell(camera)
		
		//draw the world
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
		worldShader.SetBool("white", false)
		world.Draw(worldShader)

		//display and then delay
		window.GLSwap()
		
		elapsedTime := time.Since(startTime)
		if elapsedTime < time.Duration(frameDelay)*time.Millisecond {
			time.Sleep(time.Duration(frameDelay)*time.Millisecond - elapsedTime)
		}
	}
}

func makeOneNumArray(length int, num float32) []float32 {
	arr := make([]float32, length)
	for i := range arr {
		arr[i] = num
	}
	return arr
}

package main

import (
	"fmt"
	"github.com/go-gl/gl/v4.6-core/gl"
	glm "github.com/go-gl/mathgl/mgl32"
	"github.com/veandco/go-sdl2/sdl"
	"image"
	"image/png"
	"log"
	"os"
	"runtime"
	"time"
)

const WIN_WIDTH, WIN_HEIGHT = 800, 800
const FRAME_RATE = 60
const RADIUS = 10.0
const SENSITIVITY = 0.1

var camera *Camera = MakeCamera(glm.Vec3{0, 0, 3}, glm.Vec3{0, 1, 0}, INIT_YAW, INIT_PITCH)
var deltaTime, lastFrame float32
var lastMouseX, lastMouseY int32 = WIN_WIDTH / 2, WIN_HEIGHT / 2
var yaw, pitch float32 = -90, 0

var triVertices = []float32{
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

var indices = []uint32{
	0, 1, 3, // first triangle
	1, 2, 3, // second triangle
}

var cubePositions = []glm.Vec3{
	glm.Vec3{0.0, 0.0, 0.0},
	glm.Vec3{2.0, 5.0, -15.0},
	glm.Vec3{-1.5, -2.2, -2.5},
	glm.Vec3{-3.8, -2.0, -12.3},
	glm.Vec3{2.4, -0.4, -3.5},
	glm.Vec3{-1.7, 3.0, -7.5},
	glm.Vec3{1.3, -2.0, -2.5},
	glm.Vec3{1.5, 2.0, -2.5},
	glm.Vec3{1.5, 0.2, -1.5},
	glm.Vec3{-1.3, 1.0, -1.5},
}

func main() {
	fmt.Println("begin")
	runtime.LockOSThread()

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
	gl.Viewport(0, 0, WIN_WIDTH, WIN_HEIGHT)
	progShader := NewShader("./vertShader.vs", "./fragShader.fs")
	vao, vbo, ebo := makeVao()

	texture1, err := loadTexture("./sofa-cat.png")
	if err != nil {
		log.Fatal(err)
	}

	texture2, err := loadTexture("./transhoward.png")
	if err != nil {
		log.Fatal(err)
	}
	sdl.SetRelativeMouseMode(true)
	

	gl.Enable(gl.DEPTH_TEST)
	progShader.use()
	progShader.SetInt("texture1\x00", 0)
	progShader.SetInt("texture2\x00", 1)

	proj := glm.Perspective(glm.DegToRad(45), WIN_WIDTH/WIN_HEIGHT, 0.1, 100.0)
	progShader.SetMat4("projection\x00", &proj)

	for !handleEvents() {
		startTime := time.Now()
		currentFrame := float32(sdl.GetTicks64()) / 1000
		deltaTime = currentFrame - lastFrame
		lastFrame = currentFrame

		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture1)
		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_2D, texture2)

		progShader.use()

		view := camera.GetViewMatrix()
		progShader.SetMat4("view\x00", &view)

		gl.BindVertexArray(vao)
		for i := range cubePositions {
			cubePos := cubePositions[i]
			model := glm.Ident4()
			model = model.Mul4(glm.Translate3D(cubePos.X(), cubePos.Y(), cubePos.Z()))
			var angle float32 = float32(20.0 * i)
			model = model.Mul4(glm.HomogRotate3D(glm.DegToRad(angle), glm.Vec3{1.0, 0.3, 0.5}))
			progShader.SetMat4("model\x00", &model)
			gl.DrawArrays(gl.TRIANGLES, 0, 36)
		}
		window.GLSwap()

		elapsedTime := time.Since(startTime)
		sleepTime := time.Second/time.Duration(FRAME_RATE) - elapsedTime
		if sleepTime > 0 {
			time.Sleep(sleepTime)
		}
	}

	gl.DeleteVertexArrays(1, &vao)
	gl.DeleteBuffers(1, &vbo)
	gl.DeleteBuffers(1, &ebo)
	gl.DeleteProgram(progShader.ID)
}

func handleEvents() bool {
	handleMouse()
	handleKeys(sdl.GetKeyboardState())
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			return true

		case *sdl.MouseWheelEvent:
			camera.ProcessMouseScroll(t.PreciseY)
		}
	}
	return false
}

func handleKeys(keys []uint8) {
	if keys[sdl.SCANCODE_W] == 1 {
		camera.ProcessKeyPress(FORWARD, deltaTime)
	}
	if keys[sdl.SCANCODE_S] == 1 {
		camera.ProcessKeyPress(BACKWARD, deltaTime)
	}
	if keys[sdl.SCANCODE_A] == 1 {
		camera.ProcessKeyPress(LEFT, deltaTime)
	}
	if keys[sdl.SCANCODE_D] == 1 {
		camera.ProcessKeyPress(RIGHT, deltaTime)
	}
}

func handleMouse() {
	mouseX, mouseY, _ := sdl.GetMouseState()
	xOffset, yOffset := mouseX - lastMouseX, lastMouseY - mouseY
	lastMouseX, lastMouseY = mouseX, mouseY
	camera.ProcessMouseMovement(float32(xOffset), float32(yOffset), true)
}

func loadTexture(texPath string) (uint32, error) {
	imgFile, err := os.Open(texPath)
	if err != nil {
		return 0, fmt.Errorf("failed to open file: %v", err)
	}
	defer imgFile.Close()

	img, err := png.Decode(imgFile)
	if err != nil {
		return 0, fmt.Errorf("failed to decode image: %v", err)
	}

	rgba, ok := img.(*image.RGBA)
	if !ok {
		return 0, fmt.Errorf("image is not in RGBA format")
	}

	flipImage(rgba)

	var texture uint32
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(rgba.Rect.Size().X), int32(rgba.Rect.Size().Y), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba.Pix))
	gl.GenerateMipmap(gl.TEXTURE_2D)
	return texture, nil
}

func flipImage(img *image.RGBA) {
	width, height := img.Bounds().Size().X, img.Bounds().Size().Y
	for y := 0; y < height/2; y++ {
		for x := 0; x < width; x++ {
			tmp := img.RGBAAt(x, y)
			img.SetRGBA(x, y, img.RGBAAt(x, height-1-y))
			img.SetRGBA(x, height-1-y, tmp)
		}
	}
}

func makeVao() (uint32, uint32, uint32) {
	var vbo, vao, ebo uint32

	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)

	gl.BindVertexArray(vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(triVertices)*4, gl.Ptr(triVertices), gl.STATIC_DRAW)

	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 5*4, 0)
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false, 5*4, 3*4)
	gl.EnableVertexAttribArray(1)

	return vao, vbo, ebo
}

func makeOneNumArray(length int, num float32) []float32 {
	arr := make([]float32, length)
	for i := range arr {
		arr[i] = num
	}
	return arr
}

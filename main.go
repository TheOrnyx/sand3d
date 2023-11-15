package main

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/chewxy/math32"
	"github.com/go-gl/gl/v4.6-core/gl"
	glm "github.com/go-gl/mathgl/mgl32"
	"github.com/veandco/go-sdl2/sdl"
)

const WIN_WIDTH, WIN_HEIGHT = 800, 800
const FRAME_RATE = 60
const RADIUS = 10.0
const SENSITIVITY = 0.1

var camera *Camera = MakeCamera(glm.Vec3{0, 0, 3}, glm.Vec3{0, 1, 0}, INIT_YAW, INIT_PITCH)
var deltaTime, lastFrame float32
var lastMouseX, lastMouseY int32 = WIN_WIDTH / 2, WIN_HEIGHT / 2
var yaw, pitch float32 = -90, 0

var lightPos glm.Vec3 = glm.Vec3{1.2, 1.0, 2.0}

var triVertices = []float32{
	-0.5, -0.5, -0.5, 0.0, 0.0, -1.0, 0.0, 0.0,
	0.5, -0.5, -0.5, 0.0, 0.0, -1.0, 1.0, 0.0,
	0.5, 0.5, -0.5, 0.0, 0.0, -1.0, 1.0, 1.0,
	0.5, 0.5, -0.5, 0.0, 0.0, -1.0, 1.0, 1.0,
	-0.5, 0.5, -0.5, 0.0, 0.0, -1.0, 0.0, 1.0,
	-0.5, -0.5, -0.5, 0.0, 0.0, -1.0, 0.0, 0.0,

	-0.5, -0.5, 0.5, 0.0, 0.0, 1.0, 0.0, 0.0,
	0.5, -0.5, 0.5, 0.0, 0.0, 1.0, 1.0, 0.0,
	0.5, 0.5, 0.5, 0.0, 0.0, 1.0, 1.0, 1.0,
	0.5, 0.5, 0.5, 0.0, 0.0, 1.0, 1.0, 1.0,
	-0.5, 0.5, 0.5, 0.0, 0.0, 1.0, 0.0, 1.0,
	-0.5, -0.5, 0.5, 0.0, 0.0, 1.0, 0.0, 0.0,

	-0.5, 0.5, 0.5, -1.0, 0.0, 0.0, 1.0, 0.0,
	-0.5, 0.5, -0.5, -1.0, 0.0, 0.0, 1.0, 1.0,
	-0.5, -0.5, -0.5, -1.0, 0.0, 0.0, 0.0, 1.0,
	-0.5, -0.5, -0.5, -1.0, 0.0, 0.0, 0.0, 1.0,
	-0.5, -0.5, 0.5, -1.0, 0.0, 0.0, 0.0, 0.0,
	-0.5, 0.5, 0.5, -1.0, 0.0, 0.0, 1.0, 0.0,

	0.5, 0.5, 0.5, 1.0, 0.0, 0.0, 1.0, 0.0,
	0.5, 0.5, -0.5, 1.0, 0.0, 0.0, 1.0, 1.0,
	0.5, -0.5, -0.5, 1.0, 0.0, 0.0, 0.0, 1.0,
	0.5, -0.5, -0.5, 1.0, 0.0, 0.0, 0.0, 1.0,
	0.5, -0.5, 0.5, 1.0, 0.0, 0.0, 0.0, 0.0,
	0.5, 0.5, 0.5, 1.0, 0.0, 0.0, 1.0, 0.0,

	-0.5, -0.5, -0.5, 0.0, -1.0, 0.0, 0.0, 1.0,
	0.5, -0.5, -0.5, 0.0, -1.0, 0.0, 1.0, 1.0,
	0.5, -0.5, 0.5, 0.0, -1.0, 0.0, 1.0, 0.0,
	0.5, -0.5, 0.5, 0.0, -1.0, 0.0, 1.0, 0.0,
	-0.5, -0.5, 0.5, 0.0, -1.0, 0.0, 0.0, 0.0,
	-0.5, -0.5, -0.5, 0.0, -1.0, 0.0, 0.0, 1.0,

	-0.5, 0.5, -0.5, 0.0, 1.0, 0.0, 0.0, 1.0,
	0.5, 0.5, -0.5, 0.0, 1.0, 0.0, 1.0, 1.0,
	0.5, 0.5, 0.5, 0.0, 1.0, 0.0, 1.0, 0.0,
	0.5, 0.5, 0.5, 0.0, 1.0, 0.0, 1.0, 0.0,
	-0.5, 0.5, 0.5, 0.0, 1.0, 0.0, 0.0, 0.0,
	-0.5, 0.5, -0.5, 0.0, 1.0, 0.0, 0.0, 1.0,
}

var cubePositions = []glm.Vec3{
	glm.Vec3{ 0.0,  0.0,  0.0},
	glm.Vec3{ 2.0,  5.0, -15.0},
	glm.Vec3{-1.5, -2.2, -2.5},
	glm.Vec3{-3.8, -2.0, -12.3},
	glm.Vec3{ 2.4, -0.4, -3.5},
	glm.Vec3{-1.7,  3.0, -7.5},
	glm.Vec3{ 1.3, -2.0, -2.5},
	glm.Vec3{ 1.5,  2.0, -2.5},
	glm.Vec3{ 1.5,  0.2, -1.5},
	glm.Vec3{-1.3,  1.0, -1.5},
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
	gl.Enable(gl.DEPTH_TEST)
	gl.Viewport(0, 0, WIN_WIDTH, WIN_HEIGHT)

	lightingShader := NewShader("./lighting.vs", "./lighting.fs")
	lightCubeShader := NewShader("./light_cube.vs", "./light_cube.fs")

	cubeVAO, VBO, lightCubeVAO := makeVao()

	diffuseMap, err := loadTexture("./container2.png")
	if err != nil {
		log.Fatal(err)
	}

	specularMap, err := loadTexture("./container2_specular.png")
	if err != nil {
		log.Fatal(err)
	}

	emissionMap, err := loadTexture("./sofa-emission.png")
	if err != nil {
		log.Fatal(err)
	}
	
	lightingShader.use()
	lightingShader.SetInt("material.diffuse\x00", 0)
	lightingShader.SetInt("material.specular\x00", 1)
	lightingShader.SetInt("material.emission\x00", 2)

	sdl.SetRelativeMouseMode(true)

	for !handleEvents() {
		startTime := time.Now()
		currentFrame := float32(sdl.GetTicks64()) / 1000
		deltaTime = currentFrame - lastFrame
		lastFrame = currentFrame

		gl.ClearColor(0, 0, 0, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// lightPos = glm.Vec3{1+math32.Sin(currentFrame)*2, math32.Sin(currentFrame/2), lightPos.Z()}

		lightingShader.use()
		lightingShader.SetVec3("light.position\x00", &camera.Position)
		lightingShader.SetVec3("light.direction\x00", &camera.Front)
		lightingShader.SetFloat("light.cutOff\x00", math32.Cos(glm.DegToRad(12.5)))
		lightingShader.SetFloat("light.outerCutOff\x00", math32.Cos(glm.DegToRad(17.5)))
		lightingShader.SetVec3("viewPos\x00", &camera.Position)

		lightingShader.SetVec3f("light.ambient\x00", 0.2, 0.2, 0.2)
		lightingShader.SetVec3f("light.diffuse\x00", 0.5, 0.5, 0.5)
		lightingShader.SetVec3f("light.specular\x00", 1.0, 1.0, 1.0)

		lightingShader.SetFloat("light.constant\x00", 1.0)
		lightingShader.SetFloat("light.linear\x00", 0.09)
		lightingShader.SetFloat("light.quadratic\x00", 0.032)

		lightingShader.SetFloat("material.shininess\x00", 32.0)

		proj := glm.Perspective(glm.DegToRad(camera.Zoom), WIN_WIDTH/WIN_HEIGHT, 0.1, 100.0)
		view := camera.GetViewMatrix()
		lightingShader.SetMat4("projection\x00", &proj)
		lightingShader.SetMat4("view\x00", &view)

		model := glm.Ident4()
		lightingShader.SetMat4("model\x00", &model)

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, diffuseMap)
		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_2D, specularMap)
		gl.ActiveTexture(gl.TEXTURE2)
		gl.BindTexture(gl.TEXTURE_2D, emissionMap)
		
		gl.BindVertexArray(cubeVAO)
		gl.DrawArrays(gl.TRIANGLES, 0, 36)

		gl.BindVertexArray(cubeVAO)
		for i := 0; i < 10; i++ {
			model = glm.Ident4()
			model = model.Mul4(glm.Translate3D(cubePositions[i].Elem()))
			angle := float32(20 * i)
			model = model.Mul4(glm.HomogRotate3D(glm.DegToRad(angle), glm.Vec3{1, 0.3, 0.5}))
			lightingShader.SetMat4("model\x00", &model)

			gl.DrawArrays(gl.TRIANGLES, 0, 36)
		}

		// lamp stuff
		lightCubeShader.use()
		lightCubeShader.SetMat4("projection\x00", &proj)
		lightCubeShader.SetMat4("view\x00", &view)
		model = glm.Ident4()
		model = model.Mul4(glm.Translate3D(lightPos.X(), lightPos.Y(), lightPos.Z()))
		model = model.Mul4(glm.Scale3D(0.2, 0.2, 0.2))
		lightCubeShader.SetMat4("model\x00", &model)

		gl.BindVertexArray(lightCubeVAO)
		gl.DrawArrays(gl.TRIANGLES, 0, 36)

		window.GLSwap()

		elapsedTime := time.Since(startTime)
		sleepTime := time.Second/time.Duration(FRAME_RATE) - elapsedTime
		if sleepTime > 0 {
			time.Sleep(sleepTime)
		}
	}
	gl.DeleteVertexArrays(1, &cubeVAO)
	gl.DeleteVertexArrays(1, &lightCubeVAO)
	gl.DeleteBuffers(1, &VBO)
}

func handleEvents() bool {
	handleKeys(sdl.GetKeyboardState())
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			return true

		case *sdl.MouseWheelEvent:
			camera.ProcessMouseScroll(t.PreciseY)

		case *sdl.MouseMotionEvent:
			handleMouse(t)
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

func handleMouse(t *sdl.MouseMotionEvent) {
	mouseX, mouseY := lastMouseX+t.XRel, lastMouseY+t.YRel
	xOffset, yOffset := mouseX-lastMouseX, lastMouseY-mouseY
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
	var vbo, cubeVAO, lightCubeVAO uint32

	gl.GenVertexArrays(1, &cubeVAO)
	gl.GenBuffers(1, &vbo)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(triVertices)*4, gl.Ptr(triVertices), gl.STATIC_DRAW)

	gl.BindVertexArray(cubeVAO)

	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 8*4, 0)
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointerWithOffset(1, 3, gl.FLOAT, false, 8*4, 3*4)
	gl.EnableVertexAttribArray(1)

	gl.VertexAttribPointerWithOffset(2, 2, gl.FLOAT, false, 8*4, 6*4)
	gl.EnableVertexAttribArray(2)

	gl.GenVertexArrays(1, &lightCubeVAO)
	gl.BindVertexArray(lightCubeVAO)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)

	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 8*4, 0)
	gl.EnableVertexAttribArray(0)

	return cubeVAO, vbo, lightCubeVAO
}

func makeOneNumArray(length int, num float32) []float32 {
	arr := make([]float32, length)
	for i := range arr {
		arr[i] = num
	}
	return arr
}

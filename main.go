package main

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"runtime"
	"time"
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/veandco/go-sdl2/sdl"
	glm "github.com/go-gl/mathgl/mgl32"
)

const WIN_WIDTH, WIN_HEIGHT = 800, 800
const FRAME_RATE = 60

var triVertices = []float32{
	// positions          // colors           // texture coords
     0.5,  0.5, 0.0,   1.0, 0.0, 0.0,   1.0, 1.0,   // top right
     0.5, -0.5, 0.0,   0.0, 1.0, 0.0,   1.0, 0.0,   // bottom right
    -0.5, -0.5, 0.0,   0.0, 0.0, 1.0,   0.0, 0.0,   // bottom left
    -0.5,  0.5, 0.0,   1.0, 1.0, 0.0,   0.0, 1.0,    // top left 
}

var indices = []uint32 {
	0, 1, 3, // first triangle
	1, 2, 3,  // second triangle
}

func main() {
	fmt.Println("begin")
	runtime.LockOSThread()
	
	window, err := sdl.CreateWindow("the zinger", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, WIN_WIDTH, WIN_HEIGHT, sdl.WINDOW_OPENGL)
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
	if err != nil{
		log.Fatal(err)
	}

	texture2, err := loadTexture("./sofa-cat.png")
	if err != nil{
		log.Fatal(err)
	}

	progShader.use()
	gl.Uniform1i(gl.GetUniformLocation(progShader.ID, gl.Str("texture1\x00")), 0)
	progShader.SetInt("texture2\x00", 1)
	for !handleEvents() {
		startTime := time.Now()
		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture1)
		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_2D, texture2)
		
		progShader.use()

		trans := glm.Ident4()
		trans = trans.Mul4(glm.Translate3D(0.5, -0.5, 0.0))
		trans = glm.HomogRotate3DZ(glm.DegToRad(float32(sdl.GetTicks64()/10)))
		transformLoc := gl.GetUniformLocation(progShader.ID, gl.Str("transform\x00"))
		gl.UniformMatrix4fv(transformLoc, 1, false, &trans[0])
		
		gl.BindVertexArray(vao)
		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)
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
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch event.(type) {
		case *sdl.QuitEvent:
			return true
		}
	}
	return false
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

func makeVao() (uint32, uint32, uint32){
	var vbo, vao, ebo uint32
	
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)
	gl.GenBuffers(1, &ebo)
	
	gl.BindVertexArray(vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(triVertices) * int(unsafe.Sizeof(triVertices[0])), gl.Ptr(triVertices), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*int(unsafe.Sizeof(indices[0])), gl.Ptr(indices), gl.STATIC_DRAW)
	
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, int32(8*unsafe.Sizeof(float32(0))), nil)
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, int32(8*unsafe.Sizeof(float32(0))), unsafe.Pointer(uintptr(3*4)))
	gl.EnableVertexAttribArray(1)

	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, int32(8*unsafe.Sizeof(float32(0))), unsafe.Pointer(uintptr(6*4)))
	gl.EnableVertexAttribArray(2)

	// gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	// gl.BindVertexArray(0)
	
	return vao, vbo, ebo
}

func makeOneNumArray(length int, num float32) []float32 {
	arr := make([]float32, length)
	for i := range arr {
		arr[i] = num
	}
	return arr
}

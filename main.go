package main

import (
	"fmt"
	"log"
	"runtime"
	"time"
	"unsafe"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/veandco/go-sdl2/sdl"
)

const WIN_WIDTH, WIN_HEIGHT = 800, 800
const FRAME_RATE = 60

var triVertices = []float32{
	// positions         // colors
	0.5, -0.5, 0.0,  1.0, 0.0, 0.0,   // bottom right
    -0.5, -0.5, 0.0,  0.0, 1.0, 0.0,   // bottom left
	0.0,  0.5, 0.0,  0.0, 0.0, 1.0,    // top 
}

const vertShaderSrc = `
	#version 330 core
	layout (location = 0) in vec3 aPos;   // the position variable has attribute position 0
	layout (location = 1) in vec3 aColor; // the color variable has attribute position 1
	  
	out vec3 ourColor; // output a color to the fragment shader
	
	void main()
	{
	    gl_Position = vec4(aPos, 1.0);
	    ourColor = aColor; // set ourColor to the input color we got from the vertex data
	}
`

const fragShaderSrc = `
	#version 330 core
	out vec4 FragColor;  
	in vec3 ourColor;
	  
	void main()
	{
	    FragColor = vec4(ourColor, 1.0);
	}
`

func main() {
	fmt.Println("begin")
	//progStart := time.Now()
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
	vao, vbo := makeVao()
	
	for !handleEvents() {
		startTime := time.Now()
		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		progShader.use()

		gl.BindVertexArray(vao)
		gl.DrawArrays(gl.TRIANGLES, 0, 3)
		// gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)
		window.GLSwap()
		
		elapsedTime := time.Since(startTime)
		sleepTime := time.Second/time.Duration(FRAME_RATE) - elapsedTime
		if sleepTime > 0 {
			time.Sleep(sleepTime)
		}
	}

	gl.DeleteVertexArrays(1, &vao)
	gl.DeleteBuffers(1, &vbo)
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

func makeVao() (uint32, uint32){
	var vbo, vao uint32
	
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)
	gl.BindVertexArray(vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(triVertices) * int(unsafe.Sizeof(triVertices[0])), gl.Ptr(triVertices), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, int32(6*unsafe.Sizeof(float32(0))), nil)
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, int32(6*unsafe.Sizeof(float32(0))), unsafe.Pointer(uintptr(3*4)))
	gl.EnableVertexAttribArray(1)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	gl.BindVertexArray(0)
	
	return vao, vbo
}



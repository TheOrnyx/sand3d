package main

import "github.com/veandco/go-sdl2/sdl"

func handleEvents() bool {
	handleKeys(sdl.GetKeyboardState())
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			return true

		case *sdl.MouseWheelEvent:
			camera.ProcessMouseScroll(t.PreciseY)

		case *sdl.MouseMotionEvent:
			handleMouseMovement(t)

		case *sdl.MouseButtonEvent:
			//handle mouse draw stuff here
			
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
	if keys[sdl.SCANCODE_E] == 1 {
		camera.ProcessKeyPress(UP, deltaTime)
	}
	if keys[sdl.SCANCODE_Q] == 1 {
		camera.ProcessKeyPress(DOWN, deltaTime)
	}
	if keys[sdl.SCANCODE_F] == 1 {
		drawBoundingBox = !drawBoundingBox
	}
}



func handleMouseMovement(t *sdl.MouseMotionEvent) {
	mouseX, mouseY := lastMouseX+t.XRel, lastMouseY+t.YRel
	xOffset, yOffset := mouseX-lastMouseX, lastMouseY-mouseY
	lastMouseX, lastMouseY = mouseX, mouseY
	camera.ProcessMouseMovement(float32(xOffset), float32(yOffset), true)
}

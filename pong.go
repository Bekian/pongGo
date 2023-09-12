package main

import (
	"fmt"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

// setting window variables
const winWidth, winHeight int = 800, 600

// setting composite structs that will be used somewhere else
type color struct {
	r, g, b byte
}

type pos struct {
	x, y float32
}

// ball struct
type ball struct {
	pos    // basically inherits pos struct
	radius float32
	xv     float32 // x and y velocities
	yv     float32
	color  color
}

// drawing the ball from left to right, top to bottom
func (ball *ball) draw(pixels []byte) {
	// yagni
	for y := -ball.radius; y < ball.radius; y++ {
		for x := -ball.radius; x < ball.radius; x++ {
			if x*x+y*y < ball.radius*ball.radius { //avoiding sqroot bc its expensive
				setPixel(int(ball.x+x), int(ball.y+y), ball.color, pixels)
			}
		}
	}
}

// updating the ball
// computes new position, and collision logic
// TODO: fix bug where ball is behind paddle and gets stuck
// i think the solution is within the paddle collision logic
func (ball *ball) update(leftPaddle, rightPaddle *paddle) {
	// update ball position
	ball.x += ball.xv
	ball.y += ball.yv
	// if ball hits top or bottom boundry then invert y velocity to bounce
	if ball.y-ball.radius < 0 || int(ball.y+ball.radius) > winHeight {
		ball.yv *= -1
	}
	// if ball hits either left and right walls then reset position
	if ball.x < 0 || int(ball.x) > winWidth {
		ball.x = float32(winWidth) / 2
		ball.y = float32(winHeight) / 2
	}

	// if the balls position is inside the left paddle
	if ball.x-ball.radius < leftPaddle.x+leftPaddle.w/2 {
		if ball.y > leftPaddle.y-leftPaddle.h/2 && ball.y < leftPaddle.y+leftPaddle.h/2 {
			ball.xv *= -1
		}
	}
	// same as above but for right paddle
	if ball.x+ball.radius > rightPaddle.x-rightPaddle.w/2 {
		if ball.y > rightPaddle.y-rightPaddle.h/2 && ball.y < rightPaddle.y+rightPaddle.h/2 {
			ball.xv *= -1
		}
	}
}

type paddle struct {
	pos
	w     float32
	h     float32
	color color
}

// drawing the paddle
func (paddle *paddle) draw(pixels []byte) {
	startX := int(paddle.x - paddle.w/2)
	startY := int(paddle.y - paddle.h/2)

	for y := 0; y < int(paddle.h); y++ {
		for x := 0; x < int(paddle.w); x++ {
			setPixel(startX+x, startY+y, paddle.color, pixels)
		}
	}
}

// updating the paddle, TODO: implement bounds for moving off the screen
func (paddle *paddle) update(keyState []uint8) {
	if keyState[sdl.SCANCODE_UP] != 0 {
		paddle.y -= 5
	} else if keyState[sdl.SCANCODE_DOWN] != 0 {
		paddle.y += 5
	}
}

// utility function to clear the whole screen of pixels before drawing
func clear(pixels []byte) {
	for i := range pixels {
		pixels[i] = 0
	}
}

// assigning pixels to the byte array using a position, color, and an array of pixels
func setPixel(x, y int, c color, pixels []byte) {
	// get the index of the chosen pixel
	index := (y*winWidth + x) * 4
	// make sure index is in range
	if index < len(pixels)-4 && index >= 0 {
		pixels[index] = c.r
		pixels[index+1] = c.g
		pixels[index+2] = c.b
	}
}

// main gameloop
func main() {
	// initializing sdl
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		fmt.Println("Error creating window:", err)
		return
	}
	defer sdl.Quit()
	// creating a window using the bounds and undefined position
	window, err := sdl.CreateWindow("Pong", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int32(winWidth), int32(winHeight), sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Println("Error creating window:", err)
		return
	}
	defer window.Destroy()
	// creating a renderer to put stuff into the window
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Println("Error creating renderer:", err)
		return
	}
	defer renderer.Destroy()
	// creating a texture that we can put pixels into
	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(winWidth), int32(winHeight))
	if err != nil {
		fmt.Println("Error creating renderer:", err)
		return
	}
	defer tex.Destroy()
	// a byte array that will store our pixels, essentially the screen in variable form before its drawn
	pixels := make([]byte, int(winWidth*winHeight)*4)
	// initializing the game entities
	player1 := paddle{pos{100, 300}, 20, 100, color{255, 255, 255}}
	player2 := paddle{pos{700, 300}, 20, 100, color{255, 255, 255}}
	ball := ball{pos{float32(winWidth) / 2, float32(winHeight) / 2}, 20, 2, 2, color{255, 255, 255}}

	keyState := sdl.GetKeyboardState()

	// Main loop
	running := true
	for running {

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
			}
		}

		// Update object positions
		clear(pixels)

		//updates
		player1.update(keyState)
		player2.update(keyState)
		ball.update(&player1, &player2)

		// Draw the objects
		player1.draw(pixels)
		player2.draw(pixels)
		ball.draw(pixels)

		// Update the texture and renderer
		tex.Update(nil, unsafe.Pointer(&pixels[0]), winWidth*4)
		renderer.Copy(tex, nil, nil)
		renderer.Present()
		// ~60fps
		sdl.Delay(16)
	}
}

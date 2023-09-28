package main

import (
	"fmt"
	"time"
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

func getCenter() pos {
	return pos{float32(winWidth) / 2, float32(winHeight) / 2}
}

// updating the ball
// computes new position, and collision logic
func (ball *ball) update(leftPaddle, rightPaddle *paddle, elapsedTime float32) {
	// log ball current speed
	fmt.Printf("%+v %+v\n", *&ball.xv, *&ball.yv)

	// update ball position
	ball.x += ball.xv * elapsedTime
	ball.y += ball.yv * elapsedTime
	// if ball hits top or bottom boundry then invert y velocity to bounce
	// left side of or is the bottom of the ball
	// right side of or is the top of the ball
	if ball.y-ball.radius < 0 || int(ball.y+ball.radius) > winHeight {
		ball.yv *= -1
	}
	// if ball hits either left and right walls then reset position
	if ball.x < 0 || int(ball.x) > winWidth {
		ball.pos = getCenter()
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
	speed float32
	color color
}

// drawing the paddle
func (paddle *paddle) draw(pixels []byte) {
	startX := int(paddle.x - paddle.w/2)
	startY := int(paddle.y - paddle.h/2)

	for y := 0; y < int(paddle.h); y++ {
		currentY := startY + y
		for x := 0; x < int(paddle.w); x++ {
			setPixel(startX+x, currentY, paddle.color, pixels)
		}
	}
}

// updating the paddle, TODO: implement bounds for moving off the screen
func (paddle *paddle) update(keyState []uint8, elapsedTime float32) {
	fmt.Printf("paddle speed: %+v\n", *&paddle.speed)
	// add or subtract the paddle velocity to the position to move the paddle
	// if the paddles position is at the top of the window, it will no longer go up
	// and the same for the bottom of the window,
	// basically, if the paddle is at the bottom or top of the screen you cannot go any further in that respective direction
	// but you can still go the opposite direction
	if keyState[sdl.SCANCODE_UP] != 0 && !(paddle.y-paddle.h/2 < 0) {
		paddle.y -= paddle.speed * elapsedTime
	} else if keyState[sdl.SCANCODE_DOWN] != 0 && !(paddle.y+paddle.h/2 > float32(winHeight)) {
		paddle.y += paddle.speed * elapsedTime
	}
}

// unbeatable cpu player that cannot lose
func (paddle *paddle) aiUpdate(ball *ball, elapsedTime float32) {
	paddle.y = ball.y
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
	//// INIT Window
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
	//// INIT Window

	// a byte array that will store our pixels, essentially the screen in variable form before its drawn
	pixels := make([]byte, int(winWidth*winHeight)*4)

	// initializing the game entities
	// PLAYER INIT
	player1 := paddle{pos{50, 300}, 20, 100, 300, color{255, 255, 255}}
	player2 := paddle{pos{float32(winWidth) - 50, 300}, 20, 100, 300, color{255, 255, 255}}
	ball := ball{getCenter(), 20, 0, 0, color{255, 255, 255}}

	keyState := sdl.GetKeyboardState()

	var frameStart time.Time
	var elapsedTime float32

	//// Main game loop
	running := true
	for running {
		frameStart = time.Now()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
			}
		}

		// Update object positions
		clear(pixels)

		// Log current positions
		// fmt.Println("Player 1 position:", player1.pos)
		// fmt.Println("Player 2 position:", player2.pos)
		// fmt.Println("Ball position:", ball.pos)

		//updates
		player1.update(keyState, elapsedTime)
		player2.aiUpdate(&ball, elapsedTime)
		ball.update(&player1, &player2, elapsedTime)

		// Draw the objects
		player1.draw(pixels)
		player2.draw(pixels)
		ball.draw(pixels)

		// Update the texture and renderer
		tex.Update(nil, unsafe.Pointer(&pixels[0]), winWidth*4)
		renderer.Copy(tex, nil, nil)
		renderer.Present()
		// ~60fps

		elapsedTime = float32(time.Since(frameStart).Seconds())
		// frame smoothing?
		if elapsedTime < .005 {
			sdl.Delay(5 - uint32(elapsedTime/1000.0))
			elapsedTime = float32(time.Since(frameStart).Seconds())
		}

	}
	//// main game loop
}

// KNOWN BUGS:
// collision error where ball and cpu player meet and both go off screen
// collision error where ball gets stuck behind player and enters a "caught" state and bounces between the player and their respective goal and doesnt reset the ball position
// collision error where ball phases through cpu players paddle
// collision error where ball phases through ceiling and floor

## Sept. 11th '23
    - Added comments
    - Added changelog
    - Updated most struct ints to use float32s bc it reduces cast amounts

## Sept. 12th '23
    - Added time based mechanics so the framerate is not locked
    - Added paddle speed
    - Changed more stuff to float32 
    - Broke the game, currently does not do anything when running, no controls or ball movement. F for pong.

## Sept. 14th '23
    - Updated pong draw function to save a math equation computation for every item in the draw loop so it only computes the y value once for each row

## Sept. 27th '23
    - Added more time dependent functions
    - Fixed several values to be more aligned with the tutorial, still no pong. F
    - Fixed the issue with the game not updating, it had to do with the timing at the bottom of the file.
    - Added frame smoothing?

## Sept. 28th '23
    - Added window collisions for paddles
    - Rebuilt pong for update
    - Added notes 
    - Updated speed values so the ball is faster
    - Rebuilt pong for update
    - Fixed one instance of the collision bug but i still want to improve this
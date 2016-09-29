package chip8

import "github.com/veandco/go-sdl2/sdl"

const(
  scale = 10
  w = 64
  h = 32
)

type Graphics struct {
  window *sdl.Window
  surface *sdl.Surface
  clearRect sdl.Rect
  pixels [w][h]bool
  pixelRects [w][h]sdl.Rect
}

func NewGraphics() (g *Graphics) {
  //Init SDL and create screen
	sdl.Init(sdl.INIT_VIDEO)
  window, _ := sdl.CreateWindow("test",
    sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
    w * scale, h * scale, 0)
  surface, _ := window.GetSurface()

  //Create graphics struct
  g = &Graphics{window: window, surface: surface,
    clearRect: sdl.Rect{X: 0, Y:0, W: w*scale, H: h*scale}}

  //Set default values of pixels and create rects
  for x:=0; x<w; x++ {
    for y:=0; y<h; y++{
      g.pixels[x][y] = false
      g.pixelRects[x][y] = sdl.Rect{
        X: int32(scale*x), Y: int32(scale*y), W: scale, H: scale}
    }
  }

  g.ClearScreen()
	return
}

//Set the screen to black
func (g *Graphics)ClearScreen(){
  g.surface.FillRect(&(g.clearRect), 0x0)
  g.window.UpdateSurface()

  //Clear pixels
  for x, _ := range g.pixels{
    for y, _ := range g.pixels[x] {
      g.pixels[x][y] = false
    }
  }
}

func (g *Graphics)DrawScreen(){
  for x:=0; x<w; x++{
    for y:=0; y<h; y++{
      color := uint32(0)
      if g.pixels[x][y] {
        color = 0xFFFFFF
      }
      g.surface.FillRect(&(g.pixelRects[x][y]), color);
    }
  }

  g.window.UpdateSurface()
}

func (g *Graphics)DrawSprite(x, y uint8, sprite [8][]bool) uint8{
  value := uint8(0)
  for i:= 0; i< len(sprite); i++ {
    for j, pixel := range sprite[i]{
      if pixel {
        g.pixels[(7-i) + int(x)][j + int(y)] =
          !g.pixels[(7-i) + int(x)][j + int(y)]
        value = 1
      }
    }
  }

  return value
}

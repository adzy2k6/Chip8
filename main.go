package main
import(
  "github.com/adzy2k6/Chip8/chip8"
  "io/ioutil"
  "fmt"
  "time"
)

func main(){
  rom, err := ioutil.ReadFile("INVADERS")
  if err != nil{
    fmt.Println(err)
  }

  g := chip8.NewGraphics()
  c := chip8.NewChip8(rom, g)
  for ;;{
    err := c.Tick()
    g.DrawScreen()
    //time.Sleep(100 * time.Microsecond)
    if err != nil{
      fmt.Println("Error")
      fmt.Println(err)
      break
    }
  }

  time.Sleep(5 * time.Second)
}

package chip8

import (
	"math/rand"
	"errors"
	"strconv"
)

type Chip8 struct {
	memory    [4096]uint8
	registers [16]uint8
	stack     [256]uint16
	sp        uint8
	pc        uint16
	i         uint16

	g *Graphics
}

func NewChip8(rom []uint8, g *Graphics) (c *Chip8) {
	c = &Chip8{pc: 0x200, sp: 0, g: g}
	if len(rom) > 0xFFF-0x200 {
		panic("Rom too large")
	}

	for i:=0; i< len(rom); i++{
		c.memory[0x200 + i] = rom[i]
	}
	return c
}

//Process Instructions
func (c *Chip8) Tick() error{
	opcode := uint16(c.memory[c.pc])<<8 | uint16(c.memory[c.pc+1])
	firstHalfByte := c.memory[c.pc] >> 4

	switch firstHalfByte {
	case 0:
		if opcode == 0x00E0 {
			c.g.ClearScreen()
			c.pc += 2
		} else if opcode == 0xEE {
			//Return from subroutine
			c.sp--
			c.pc = c.stack[c.sp]
			c.pc += 2
		} else {
			return undefinedOpcode(opcode, c.pc)
		}
	case 1:
		//Jump to address
		c.pc = opcode & 0xFFF
	case 2:
		//Call Subroutine
		c.stack[c.sp] = c.pc
		c.sp++
		c.pc = parseNumber(opcode)
	case 3:
		if c.compareRegister(opcode) {
			c.pc += 4
		} else {
			c.pc += 2
		}
	case 4:
		if !(c.compareRegister(opcode)) {
			c.pc += 2
		} else {
			c.pc += 4
		}
	case 5:
		reg1, reg2 := parseRegisters(opcode)
		if c.registers[reg1] == c.registers[reg2] {
			c.pc += 4
		} else {
			c.pc += 2
		}
	case 6:
		reg, number := parseRegAndNumber(opcode)
		c.registers[reg] = number
		c.pc += 2
	case 7:
		reg, number := parseRegAndNumber(opcode)
		c.registers[reg] += number
		c.pc += 2
	case 8:
		return c.processOpcode8(opcode)
	case 9:
		reg1, reg2 := parseRegisters(opcode)
		if c.registers[reg1] == c.registers[reg2] {
			c.pc += 2
		} else {
			c.pc += 4
		}
	case 0xA:
		c.i = opcode & 0xFFF
		c.pc += 2
	case 0xB:
		c.pc = uint16(c.registers[0]) + (opcode & 0xFFF)
	case 0xC:
		reg, number := parseRegAndNumber(opcode)
		c.registers[reg] = number & uint8(rand.Uint32())
		c.pc += 2
	case 0xD:
		reg1, reg2 := parseRegisters(opcode)
		number := uint8(opcode & 0xF)
		c.registers[15] = c.drawSprite(reg1, reg2, number)
		c.pc += 2
	case 0xE:
		_, lsb := parseRegAndNumber(opcode)
		if lsb == 0x9E{
			c.pc += 2
		}else if lsb == 0xA1{
			c.pc += 4
		}else{
			undefinedOpcode(opcode, c.pc)
		}

	case 0xF:
		return c.processOpcodeF(opcode);
	default:
		return undefinedOpcode(opcode, c.pc)
	}
	return nil
}

/*
	Compares the number to the register encoded in the opcode
*/
func (c *Chip8) compareRegister(opcode uint16) bool {
	reg, number := parseRegAndNumber(opcode)
	return (number == c.registers[reg])
}

/*
	Returns the register and number encoded in the opcode
*/
func parseRegAndNumber(opcode uint16) (reg, number uint8) {
	reg = uint8(opcode>>8) & 0xF
	number = uint8(opcode & 0xFF)
	return
}

/*
	Return the lower three half bytes from the code
*/
func parseNumber(opcode uint16) (number uint16) {
	number = uint16(opcode & 0xFFF)
	return
}

func parseRegisters(opcode uint16) (reg1, reg2 uint8) {
	reg1 = uint8(opcode>>8) & 0xF
	reg2 = uint8(opcode>>4) & 0xF
	return
}

/*
  Handle all opcodes beginning with 8
*/
func (c *Chip8) processOpcode8(opcode uint16) error{
	lowestByte := opcode & 0xF
	reg1, reg2 := (opcode&0xF00)>>8, (opcode&0xF0)>>4
	switch lowestByte {
	case 0:
		c.registers[reg1] = c.registers[reg2]
	case 1:
		c.registers[reg1] = c.registers[reg1] | c.registers[reg2]
	case 2:
		c.registers[reg1] = c.registers[reg1] & c.registers[reg2]
	case 3:
		c.registers[reg1] = c.registers[reg1] ^ c.registers[reg2]
	case 4:
		c.registers[reg1] = c.registers[reg1] + c.registers[reg2]
		if c.registers[reg2] > c.registers[reg1] {
			c.registers[15] = 1
		} else {
			c.registers[15] = 0
		}
	case 5:
		temp := c.registers[reg1] - c.registers[reg2]
		if c.registers[reg1] < temp {
			c.registers[15] = 0
		} else {
			c.registers[15] = 1
		}
		c.registers[reg1] = temp
	case 6:
		c.registers[15] = c.registers[reg1] & 1
		c.registers[reg1] = c.registers[reg1] >> 1
	case 7:
		c.registers[reg1] = c.registers[reg2] - c.registers[reg1]
		if c.registers[reg1] > c.registers[reg2] {
			c.registers[15] = 0
		} else {
			c.registers[15] = 1
		}
	case 0xE:
		c.registers[15] = c.registers[reg1] & uint8(0x8000>>15)
		c.registers[reg1] = c.registers[reg1] << 1
	default:
		return undefinedOpcode(opcode, c.pc)
	}
	c.pc += 2
	return nil
}

func (c *Chip8) processOpcodeF(opcode uint16)error{
	reg1, lsb := parseRegAndNumber(opcode)
	switch lsb{
	case 0x07:
		c.registers[reg1] = 0
	case 0x0A:
		c.registers[reg1] = 1
	case 0x1E:
		c.i += uint16(c.registers[reg1])
	case 0x65:
		reg1, _ = parseRegisters(opcode)
		index := c.i
		for i:=uint8(0); i<=reg1; i++{
			c.registers[i] = c.memory[index]
			index ++
		}
	}
	c.pc += 2
	return nil
}

func undefinedOpcode(opcode uint16, pc uint16)error{
	error := errors.New(
		"Undefined Opcode: 0x" +
		strconv.FormatUint(uint64(opcode), 16) + " " +
		strconv.FormatUint(uint64(pc), 10),
	)
	return error
}

func (c *Chip8)drawSprite(reg1, reg2, number uint8)uint8{
	x := c.registers[reg1]
	y := c.registers[reg2]
	var sprite [8][]bool

	//Create sprite slices
	for i:=0; i<8; i++ {
		sprite[i] = make([]bool, number)
	}

	c.registers[15] = 0
	//Create sprite arrays
	for i:=uint8(0); i<number; i++{
		row := c.memory[c.i + uint16(i)]
		for j:=uint8(0); j<8; j++ {
			//Set Value
			if (row & (1<<j)) != 0 {
				sprite[j][i] = true
			}else{
				c.registers[15] = 1
				sprite[j][i] = false
			}
		}
	}

	return c.g.DrawSprite(x, y, sprite)
}

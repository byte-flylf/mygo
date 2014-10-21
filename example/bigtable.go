package main

import (
    "fmt"
)

type BitTable struct {
    b []uint64
}

func (t *BitTable) bitToWord(x int) (word, bit uint) {
    return uint(x/64), uint(x%64)
}

func (t *BitTable) Bit(x int) bool {
    w, b := t.bitToWord(x)
    return t.b[w]&(1<<b) != 0
}

func (t *BitTable) SetBit(x int) {
    w, b := t.bitToWord(x)
    t.b[w] |= (1 << b)
}

func (t *BitTable) ClearBit(x int) {
    w, b := t.bitToWord(x)
    t.b[w] = t.b[w] & ^(1 << b)
}

func NewBitTable(size uint) *BitTable {
    if size % 64 != 0 {
        panic("Invalid size")
    }
    s2 := size / 64 * size
    return &BitTable{make([]uint64, s2)}
}

func main() {
    t := NewBitTable(128)
    t.SetBit(45)
    t.SetBit(44)
    fmt.Println("bit 45: ", t.Bit(45))
    fmt.Println("bit 44: ", t.Bit(44))
    t.ClearBit(44)
    fmt.Println("bit 45: ", t.Bit(45))
    fmt.Println("bit 44: ", t.Bit(44))
}

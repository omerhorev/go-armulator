package goarmulator

import (
	"github.com/pkg/errors"
)

// Interface for interacting with registers
// Aarch64 defines 31 general-purpose registers (W0..W30), SP register and PC.
// Each register can be accessed as X (64bit), or W (32bit, lsb).
type RegistersBank struct {
	x  []uint64
	pc uint64
	sp uint64
}

func NewAarch64RegistersBank() *RegistersBank {
	return &RegistersBank{
		x:  make([]uint64, 30),
		pc: 0,
		sp: 0,
	}
}

// Access X register. Returns error if invalid register was provided
func (b *RegistersBank) X(id int) (uint64, error) {
	if id >= len(b.x) {
		return 0, errors.Errorf("unknown register X%d", id)
	}

	return b.x[id], nil
}

// Access W register. Returns error if invalid register was provided
func (b *RegistersBank) W(id int) (uint32, error) {
	if id >= len(b.x) {
		return 0, errors.Errorf("unknown register W%d", id)
	}

	return uint32(b.x[id] & 0xffffffff), nil
}

// Access the PC register
func (b *RegistersBank) PC() uint64 {
	return b.pc
}

// Access the SP register
func (b *RegistersBank) SP() uint64 {
	return b.sp
}

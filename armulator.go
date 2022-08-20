package goarmulator

import (
	// "elf"
	// "elf"
	"context"
	"debug/elf"
	"io"
	"sync"

	"github.com/pkg/errors"
)

// Interface for interacting with registers
type RegistersBank interface {
	// Returns a register by it's number
	Get(id int) (error, Register)

	//Returns a register by it's name
	GetName(name string) (error, Register)

	// Writes a value to a register
	Write(id int, value uint64) error
}

// Interface for
type Register interface {
	// Writes an uint64 to the register
	Write(uint64)

	// Writes an int to the reguster. The int is converted to uint64 then written
	WriteInt(int)
}

type Armulator struct {
	RegistersBank RegistersBank
	File          *elf.File
	finishWg      sync.WaitGroup
	err           error
}

func NewArmulator(elfFile io.ReaderAt) (*Armulator, error) {
	e, err := elf.NewFile(elfFile)
	if err != nil {
		return nil, errors.Wrap(err, "elf")
	}

	a := Armulator{
		File: e,
	}

	return &a, nil
}

func (a *Armulator) Run() error {
	if err := a.Start(); err != nil {
		return err
	}

	a.Wait()
	return nil
}

func (a *Armulator) Start() error {
	return a.StartContext(context.Background())
}

func (a *Armulator) StartContext(ctx context.Context) error {
	return nil
}

func (a *Armulator) Wait() {
	a.finishWg.Wait()
}

func (a *Armulator) Close() error {
	return a.File.Close()
}

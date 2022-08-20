package goarmulator

import (
	// "elf"
	// "elf"
	"context"
	"debug/elf"
	"io"
	"log"
	"sync"

	"github.com/omerhorev/goarmulator/mem"
	"github.com/omerhorev/goarmulator/utils"
	"github.com/pkg/errors"
)

// Armulator is the main structure of this project. It is used to emulate arm programs
// and run them in a platform idependant way
type Armulator struct {
	RegistersBank *RegistersBank
	Memory        *mem.Mem
	ELF           *elf.File
	finishWg      sync.WaitGroup
	Log           *log.Logger
}

// Creates a new armulator object from an ELF file
func NewArmulator(elfFile io.ReaderAt) (*Armulator, error) {
	e, err := elf.NewFile(elfFile)
	if err != nil {
		return nil, errors.Wrap(err, "elf")
	}

	a := Armulator{
		ELF:           e,
		RegistersBank: NewAarch64RegistersBank(),
		Memory:        mem.NewMemFromMemory(),
		Log:           log.New(io.Discard, "", 0),
	}

	return &a, nil
}

// Starts the execution of the program and wait for it to finish.
// Allocate the memory regions specified in the ELF file, set PC to the entrypoint
// and start execution.
func (a *Armulator) Run() error {
	if err := a.Start(); err != nil {
		return err
	}

	a.Wait()
	return nil
}

// Starts the execution of the program.
// Allocate the memory regions specified in the ELF file, set PC to the entrypoint
// and start execution.
func (a *Armulator) Start() error {
	return a.StartContext(context.Background())
}

// Starts the execution of the program with a context.
// Allocate the memory regions specified in the ELF file, set PC to the entrypoint
// and start execution.
func (a *Armulator) StartContext(ctx context.Context) error {
	if err := a.allocateELF(); err != nil {
		return errors.Wrap(err, "elf")
	}

	if err := a.initializeRegisters(); err != nil {
		return errors.Wrap(err, "init registers")
	}

	return nil
}

// Wait for the execution of the program to finish.
func (a *Armulator) Wait() {
	a.finishWg.Wait()
}

// Closes the emulator, thus freeing all resources allocated
func (a *Armulator) Close() error {
	return a.ELF.Close()
}

func (a *Armulator) allocateELF() error {
	if a.ELF.Machine != elf.EM_AARCH64 {
		return errors.Errorf("machine is not Aarch64 (is 0x%x)", int64(a.ELF.Machine))
	}

	for i, prog := range a.ELF.Progs {
		// a.Log.Println("loading prog #%d: %s", prog.Filesz)
		if prog.Type != elf.PT_LOAD {
			continue
		}

		if err := a.createSegmentInMemory(prog); err != nil {
			return errors.Wrapf(err, "segment %d", i)
		}
	}

	return nil
}

func (a *Armulator) createSegmentInMemory(prog *elf.Prog) error {
	var perm mem.Permissions = 0
	if prog.Flags&elf.PF_R != 0 {
		perm |= mem.PermRead
	}
	if prog.Flags&elf.PF_W != 0 {
		perm |= mem.PermWrite
	}
	if prog.Flags&elf.PF_X != 0 {
		perm |= mem.PermExecute
	}

	a.Log.Printf("alloc: addr=0x%x size=0x%x perm=%s", prog.Memsz, prog.Vaddr, perm)
	if err := a.Memory.Alloc(prog.Vaddr, int(prog.Memsz), perm); err != nil {
		return errors.Wrap(err, "segment alloc")
	}

	w := utils.NewOffsetWriter(a.Memory.Raw(), int64(prog.Vaddr))
	r := prog.Open()

	if _, err := io.CopyN(w, r, int64(prog.Filesz)); err != nil {
		return errors.Wrap(err, "segment copy")
	}

	return nil
}

func (a *Armulator) initializeRegisters() error {
	log.Printf("initialize registers: pc=0x%x", a.ELF.Entry)

	a.RegistersBank.pc = a.ELF.Entry

	return nil
}

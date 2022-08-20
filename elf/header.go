package elf

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/pkg/errors"
)

var (
	elfHeaderMagic = []byte{0x7f, 0x45, 0x4c, 0x46}
)

type Header struct {
	ByteOrder                 binary.ByteOrder // Byte Order (little/big endian)
	OsAbi                     OsAbi            // OS ABI
	OsAbiVersion              OsAbiVersion     // OS ABI Version
	Type                      Type             // Identifies object file type
	Machine                   Machine          // Specifies target instruction set architecture
	Entrypoint                uint64           // Mmemory address of the entry point from where the process starts executing
	ProgramHeaderOffset       uint64           // Points to the start of the program header table
	ProgramHeaderEntriesCount int              // Contains the number of entries in the program header table
	SectionHeaderOffset       uint64           // Points to the start of the section header table
	SectionHeaderEntriesCount int              // Contains the number of entries in the section header table
	Flags                     uint32           // Interpretation of this field depends on the target architecture
	SectionHeaderStringIndex  int              // Contains index of the section header table entry that contains the section names
}

// Os ABI as specified in the ELF ABI field
type OsAbi byte

const (
	OsAbiSystemV                    OsAbi = 0x00
	OsAbiHPUX                       OsAbi = 0x01
	OsAbiNetBSD                     OsAbi = 0x02
	OsAbiLinux                      OsAbi = 0x03
	OsAbiGNUHurd                    OsAbi = 0x04
	OsAbiSolaris                    OsAbi = 0x06
	OsAbiAIXMonterey                OsAbi = 0x07
	OsAbiIRIX                       OsAbi = 0x08
	OsAbiFreeBSD                    OsAbi = 0x09
	OsAbiTru64                      OsAbi = 0x0A
	OsAbiNovellModesto              OsAbi = 0x0B
	OsAbiOpenBSD                    OsAbi = 0x0C
	OsAbiOpenVMS                    OsAbi = 0x0D
	OsAbiNonStopKernel              OsAbi = 0x0E
	OsAbiAROS                       OsAbi = 0x0F
	OsAbiFenixOS                    OsAbi = 0x10
	OsAbiNuxiCloudABI               OsAbi = 0x11
	OsAbiStratusTechnologiesOpenVOS OsAbi = 0x12
)

type Type uint16

const (
	TypeNone   Type = 0x00
	TypeRel    Type = 0x01
	TypeExec   Type = 0x02
	TypeDyn    Type = 0x03
	TypeCore   Type = 0x04
	TypeLoOs   Type = 0xFE00
	TypeHiOs   Type = 0xFEFF
	TypeLoProc Type = 0xFF00
	TypeHIProc Type = 0xFFFF
)

// Os ABI version as specified in the ELF ABI version field
type OsAbiVersion byte

type Machine uint16

const (
	Unknown                          Machine = 0x00
	ATandTWE32100                    Machine = 0x01
	SPARC                            Machine = 0x02
	x86                              Machine = 0x03
	Motorola68000M68k                Machine = 0x04
	Motorola88000M88k                Machine = 0x05
	IntelMCU                         Machine = 0x06
	Intel80860                       Machine = 0x07
	MIPS                             Machine = 0x08
	IBMSystem370                     Machine = 0x09
	MIPSRS3000LittleEndian           Machine = 0x0A
	HewlettPackardPARISC             Machine = 0x0E
	Intel80960                       Machine = 0x13
	PowerPC                          Machine = 0x14
	PowerPC64bit                     Machine = 0x15
	S390includingS390x               Machine = 0x16
	IBMSPUSPC                        Machine = 0x17
	NECV800                          Machine = 0x24
	FujitsuFR20                      Machine = 0x25
	TRWRH32                          Machine = 0x26
	MotorolaRCE                      Machine = 0x27
	ARMv7Aarch32                     Machine = 0x28
	DigitalAlpha                     Machine = 0x29
	SuperH                           Machine = 0x2A
	SPARCVersion9                    Machine = 0x2B
	SiemensTriCore                   Machine = 0x2C
	ArgonautRISCCore                 Machine = 0x2D
	HitachiH8300                     Machine = 0x2E
	HitachiH8300H                    Machine = 0x2F
	HitachiH8S                       Machine = 0x30
	HitachiH8500                     Machine = 0x31
	IA64                             Machine = 0x32
	StanfordMIPSX                    Machine = 0x33
	MotorolaColdFire                 Machine = 0x34
	MotorolaM68HC12                  Machine = 0x35
	FujitsuMMA                       Machine = 0x36
	SiemensPCP                       Machine = 0x37
	SonyNCPU                         Machine = 0x38
	DensoNDR1                        Machine = 0x39
	MotorolaStar                     Machine = 0x3A
	ToyotaME16                       Machine = 0x3B
	STMicroelectronicsST100          Machine = 0x3C
	AdvancedLogicCorpTinyJ           Machine = 0x3D
	AMDx8664                         Machine = 0x3E
	Sony                             Machine = 0x3F
	DigitalEquipmentCorpPDP10        Machine = 0x40
	DigitalEquipmentCorpPDP11        Machine = 0x41
	SiemensFX66                      Machine = 0x42
	STMicroelectronicsST9Plus8_16bit Machine = 0x43
	STMicroelectronicsST78bit        Machine = 0x44
	MotorolaMC68HC16                 Machine = 0x45
	MotorolaMC68HC11                 Machine = 0x46
	MotorolaMC68HC08                 Machine = 0x47
	MotorolaMC68HC05                 Machine = 0x48
	SiliconGraphicsSVx               Machine = 0x49
	STMicroelectronicsST19_8bit      Machine = 0x4A
	DigitalVAX                       Machine = 0x4B
	AxisCommunications32bit          Machine = 0x4C
	InfineonTechnologies32bit        Machine = 0x4D
	Element14_64bit                  Machine = 0x4E
	LSILogic16bit                    Machine = 0x4F
	TMS320C6000Family                Machine = 0x8C
	MCSTElbrusE2k                    Machine = 0xAF
	ARM64bitsARMv8Aarch64            Machine = 0xB7
	ZilogZ80                         Machine = 0xDC
	RISCV                            Machine = 0xF3
	BerkeleyPacketFilter             Machine = 0xF7
	WDC65C816                        Machine = 0x101
)

func (e *Elf) parseHeader() error {
	if _, err := e.r.Seek(0, io.SeekStart); err != nil {
		return err
	}

	magic := make([]byte, 4)
	if _, err := io.ReadFull(e.r, magic); err != nil {
		return errors.Wrap(err, "read magic")
	}

	if !bytes.Equal(magic, elfHeaderMagic) {
		return errors.Errorf("invalid magic %x", magic)
	}

	if bitness, err := readByte(e.r); err != nil {
		return errors.Wrap(err, "read bitness")
	} else if bitness != 2 {
		return errors.Errorf("not a 64bit elf")
	}

	if endianness, err := readByte(e.r); err != nil {
		return errors.Wrap(err, "read endianness")
	} else if endianness == 1 {
		e.Header.ByteOrder = binary.LittleEndian
	} else if endianness == 2 {
		e.Header.ByteOrder = binary.BigEndian
	} else {
		return errors.Errorf("unknown endianness %d", endianness)
	}

	if abi, err := readByte(e.r); err != nil {
		return errors.Wrap(err, "read abi")
	} else {
		e.Header.OsAbi = OsAbi(abi)
	}

	if abiversion, err := readByte(e.r); err != nil {
		return errors.Wrap(err, "read abi version")
	} else {
		e.Header.OsAbiVersion = OsAbiVersion(abiversion)
	}

	if _, err := e.r.Seek(7, io.SeekCurrent); err != nil {
		return err
	}

	if err := binary.Read(e.r, e.Header.ByteOrder, &e.Header.Type); err != nil {
		return errors.Wrap(err, "read type")
	}

	if err := binary.Read(e.r, e.Header.ByteOrder, &e.Header.Machine); err != nil {
		return errors.Wrap(err, "read machine")
	}

	var version uint32 = 0
	if err := binary.Read(e.r, e.Header.ByteOrder, &version); err != nil {
		return errors.Wrap(err, "read version")
	}

	if version != 0x1 {
		return errors.Errorf("invalid ELF version %d", version)
	}

	var entrypoint uint64 = 0
	if err := binary.Read(e.r, e.Header.ByteOrder, &entrypoint); err != nil {
		return errors.Wrap(err, "read entrypoint")
	}

	e.Header.Entrypoint = uint64(entrypoint)

	if err := binary.Read(e.r, e.Header.ByteOrder, &e.Header.ProgramHeaderOffset); err != nil {
		return errors.Wrap(err, "read phoff")
	}

	if err := binary.Read(e.r, e.Header.ByteOrder, &e.Header.SectionHeaderOffset); err != nil {
		return errors.Wrap(err, "read shoff")
	}

	if err := binary.Read(e.r, e.Header.ByteOrder, &e.Header.Flags); err != nil {
		return errors.Wrap(err, "read flags")
	}

	if _, err := e.r.Seek(4, io.SeekCurrent); err != nil {
		return err
	}

	var phnum uint16 = 0
	if err := binary.Read(e.r, e.Header.ByteOrder, &phnum); err != nil {
		return errors.Wrap(err, "read phnum")
	}

	e.Header.ProgramHeaderEntriesCount = int(phnum)

	if _, err := e.r.Seek(2, io.SeekCurrent); err != nil {
		return err
	}

	var shnum uint16 = 0
	if err := binary.Read(e.r, e.Header.ByteOrder, &shnum); err != nil {
		return errors.Wrap(err, "read sh")
	}

	e.Header.ProgramHeaderEntriesCount = int(shnum)

	var phshstrndx uint16 = 0
	if err := binary.Read(e.r, e.Header.ByteOrder, &phshstrndx); err != nil {
		return errors.Wrap(err, "read phshstrndx")
	}

	e.Header.SectionHeaderStringIndex = int(phshstrndx)

	return nil
}

func readByte(r io.Reader) (byte, error) {
	b := []byte{0}
	_, err := io.ReadFull(r, b)

	return b[0], err
}

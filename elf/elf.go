package elf

import (
	"io"

	"github.com/pkg/errors"
)

type Elf struct {
	r      io.ReadSeeker
	Header Header
}

func ReadElf(r io.ReadSeeker) (*Elf, error) {
	e := Elf{
		r: r,
	}

	if err := e.parseHeader(); err != nil {
		return nil, errors.Wrap(err, "header")
	}

	return &e, nil
}

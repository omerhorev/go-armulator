package goarmulator

import "fmt"

type ErrorUnknownRegister struct {
	Id int
}

func (e *ErrorUnknownRegister) Error() string {
	return fmt.Sprintf("armulator: unknown interface %d", e.Id)
}

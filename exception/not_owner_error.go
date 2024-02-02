package exception

import "fmt"

type NotOwnerError struct {
	Error string
}

func NewNotOwnerError(error string) NotOwnerError {
	return NotOwnerError{Error: error}
}

func PanicIfNotOwner(err error, s string) {
	if err != nil {
		fmt.Println(s)
		panic(NewNotOwnerError(err.Error()))
	}
}

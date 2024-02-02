package exception

import "fmt"

type NotFoundError struct {
	Error string
}

func NewNotFoundError(error string) NotFoundError {
	return NotFoundError{Error: error}
}

func PanicIfNotFound(err error, s string) {
	if err != nil {
		fmt.Println(s)
		panic(NewNotFoundError(err.Error()))
	}
}

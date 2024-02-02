package helper

import (
	"fmt"
)

func PanicIfError(err error, s string) {
	fmt.Println(err, "fff", err != nil, "rr", err == nil)
	if err != nil {
		fmt.Println(s)
		panic(err)
	}
}

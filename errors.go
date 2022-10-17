package main

import (
	"fmt"
	"os"
)

func errorPrintln(message string, err error) {
	fmt.Println(fmt.Sprintf("msg: (%s) got error: (%s) ", message, err.Error()))
}

func errorIsNotExist(message string, err error) bool {
	if os.IsNotExist(err) {
		errorPrintln(message, err)
		return os.IsNotExist(err)
	}

	return false
}

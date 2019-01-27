package main

import "fmt"

type Error struct {
	Code    int
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprint(e.Code, " ",e.Message)
}

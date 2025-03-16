package main

import "fmt"

type AuraError struct {
	Code    string
	Message string
}

func (ae AuraError) Error() string {
	return fmt.Sprintf("%s - %s", ae.Code, ae.Message)
}

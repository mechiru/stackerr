package main

import (
	"encoding/json"
	"os"

	"github.com/mechiru/stackerr"
)

func main() {
	err := inner1()
	e := logEntry{
		Severity: "ERROR",
		Message:  err.Error(),
	}
	json.NewEncoder(os.Stdout).Encode(e)
}

func inner1() error {
	return inner2()
}

func inner2() error {
	return inner3()
}

func inner3() error {
	return stackerr.New("new error\n")
}

type logEntry struct {
	Severity string `json:"severity"`
	Message  string `json:"message"`
}

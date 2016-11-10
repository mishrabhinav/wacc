package main

import (
	"fmt"
)

type SyntaxError struct {
	file         string
	line, column int
	msg          string
}

func (m *SyntaxError) Error() string {
	return fmt.Sprintf("%s:%d:%d:error: %s", m.file, m.line, m.column, m.msg)
}

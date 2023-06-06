package logging

import (
	"time"

	"github.com/go-stack/stack"
)

type Entry struct {
	Level          LogLevel
	Time           time.Time
	Module         string
	Call           stack.Call
	Message        string
	KeyValues      []interface{}
	FormatSelector string
}

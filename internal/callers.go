package internal

import (
	"runtime"
	"strconv"

	"github.com/pkg/errors"
)

// StackTrace represents the function call stack, with the last function to be
// called first, followed by the function that called it, and so on.
type StackTrace []uintptr

// StackTrace implements stackTrace from github.com/pkg/errors.
func (s StackTrace) StackTrace() errors.StackTrace {
	frames := make(errors.StackTrace, len(s))
	for i, pc := range s {
		frames[i] = errors.Frame(pc)
	}
	return frames
}

// Callers is a wrapper for runtime.Callers that automatically allocates a
// large enough slice to fit the entire stack trace. A skip value of 0 starts
// the trace with Callers, a skip value of 1 starts the trace with the caller
// of Callers, and so on.
func Callers(skip int) StackTrace {
	pc := make([]uintptr, 16)
	for {
		n := runtime.Callers(skip+1, pc)
		if n < len(pc)-1 {
			return pc[:n]
		}

		pc = make([]uintptr, len(pc)*2)
	}
}

// AppendTo appends a textual representation of the stack trace to the end of
// the byte slice.
func (s StackTrace) AppendTo(buf []byte) []byte {
	frames := runtime.CallersFrames([]uintptr(s))

	for i := 0; ; i++ {
		frame, more := frames.Next()

		buf = append(buf, "\n("...)
		buf = strconv.AppendInt(buf, int64(i), 10)
		buf = append(buf, ") "...)
		if frame.Function != "" {
			buf = append(buf, frame.Function...)
		} else {
			buf = append(buf, "(unknown function)"...)
		}
		buf = append(buf, "+0x"...)
		buf = strconv.AppendUint(buf, uint64(frame.PC-frame.Entry), 16)
		if frame.File != "" {
			buf = append(buf, " @ "...)
			buf = append(buf, frame.File...)
			buf = append(buf, ":"...)
			buf = strconv.AppendInt(buf, int64(frame.Line), 10)
		}

		if !more {
			return buf
		}
	}
}

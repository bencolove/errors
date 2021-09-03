package errors

import "runtime"

const (
	MAX_CALL_STACK_DEPTH = 10
)

type callframe struct {
	file string
	line int
	fn   string
}

func extractFrame(frame *runtime.Frame) callframe {
	var funcName string
	funcForPc := runtime.FuncForPC(frame.PC)
	if funcForPc != nil {
		funcName = funcForPc.Name()
	}
	return callframe{
		file: frame.File,
		line: frame.Line,
		fn:   funcName,
	}
}

func GetStackFrame(skipdepth int) *runtime.Frame {
	// skipdepth:=0 identified as the caller of runtime.Caller()
	pc, file, line, ok := runtime.Caller(skipdepth)

	if !ok {
		return nil
	}

	frame := &runtime.Frame{
		PC:   pc,
		File: file,
		Line: line,
	}

	funcForPc := runtime.FuncForPC(pc)
	if funcForPc != nil {
		frame.Func = funcForPc
		frame.Function = funcForPc.Name()
		frame.Entry = funcForPc.Entry()
	}

	return frame
}

func GetCallframes(skip int) []callframe {
	pcs := make([]uintptr, MAX_CALL_STACK_DEPTH)
	runtime.Callers(skip, pcs)

	// collect
	frames := runtime.CallersFrames(pcs)

	// max depth
	var stacks []callframe
	for frame, more := frames.Next(); more; frame, more = frames.Next() {
		callstack := extractFrame(&frame)
		stacks = append(stacks, callstack)
	}

	return stacks
}

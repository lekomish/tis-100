package loader

import (
	"fmt"

	"github.com/yuin/gopher-lua"
)

// wrapWriterError provides a consistent error message for write failures.
func wrapWriterError(ctx, filePath string, err error) error {
	return fmt.Errorf("failed to write %s to file %s: %w", ctx, filePath, err)
}

// mustString checks that the Lua value is a string and returns it.
func mustString(val lua.LValue, field string) (string, error) {
	str, ok := val.(lua.LString)
	if !ok {
		return "", fmt.Errorf("%s: expected string, got %s", field, val.Type())
	}
	return str.String(), nil
}

// mustNumber checks that the Lua value is a number and returns it.
func mustNumber(val lua.LValue, field string) (lua.LNumber, error) {
	num, ok := val.(lua.LNumber)
	if !ok {
		return 0, fmt.Errorf("%s: expected number, got %s", field, val.Type())
	}
	return num, nil
}

// mustTable checks that the Lua value is a table and returns it.
func mustTable(val lua.LValue, field string) (*lua.LTable, error) {
	tab, ok := val.(*lua.LTable)
	if !ok {
		return nil, fmt.Errorf("%s: expected table, got %s", field, val.Type())
	}
	return tab, nil
}

// runLuaFunction calls a global Lua function by name and returns its single result.
func runLuaFunction(lState *lua.LState, functionName string) (lua.LValue, error) {
	fn := lState.GetGlobal(functionName)
	if fn.Type() == lua.LTNil {
		return nil, fmt.Errorf("lua function %q not found", functionName)
	}

	if err := lState.CallByParam(lua.P{
		Fn:      fn,
		NRet:    1,
		Protect: true,
	}); err != nil {
		return nil, fmt.Errorf("error while calling %q function: %w", functionName, err)
	}
	result := lState.Get(-1)
	lState.Pop(1)
	return result, nil
}

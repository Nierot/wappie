// Code generated by "stringer -type=Command"; DO NOT EDIT.

package main

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[CommandDebug-0]
	_ = x[CommandHelp-1]
}

const _Command_name = "CommandDebugCommandHelp"

var _Command_index = [...]uint8{0, 12, 23}

func (i Command) String() string {
	if i < 0 || i >= Command(len(_Command_index)-1) {
		return "Command(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Command_name[_Command_index[i]:_Command_index[i+1]]
}

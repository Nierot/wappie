package main

type Command int

//go:generate stringer -type=Command
const (
	CommandDebug Command = iota
	CommandHelp
)
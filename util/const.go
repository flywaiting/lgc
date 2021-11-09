package util

const (
	Ready = iota
	Run
	Done
	Cancel = 10 + iota
	Wrong
)

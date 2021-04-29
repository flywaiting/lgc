package com

import "runtime"

const (
	Ready = iota
	Running
	Succ
	// 主动停止
	Stop
	// 异常中断
	Interrupt
)

const OS = runtime.GOOS

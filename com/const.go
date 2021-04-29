package com

import "runtime"

const (
	Ready = iota
	Running
	Succ
	// 移除
	Remove
	// 异常中断
	Interrupt
	// 主动停止
	Kill
)

const OS = runtime.GOOS

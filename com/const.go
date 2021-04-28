package com

const (
	Ready = iota
	Running
	Succ
	// 异常中断
	Interrupt
	// 主动停止
	Stop
)
